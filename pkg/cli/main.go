package main

import(
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	colorReset = "\033[0m" //padrão
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
)

func main() {
	programFlags := GetFlags()

	//Valiação das flags
	if err := validateFlags(programFlags); err != nil {
		printError(err)
		return
	}
	// Para utilizar flag.Args() é necessário chamar flag.Parse() presente em Flags()

	hosts := FFlag(programFlags)
	
	//Converte nos respectivos ips
	ips, err := ResolveTargets(hosts)
	if err != nil {
		printError(err)
	}

	//Printando os alvos
	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(table, "ALVO\tENDEREÇOS RESOLVIDOS\tSTATUS")
	for _, host := range hosts {
		resolvedIPs, ok := ips[host]
		if !ok {
			fmt.Fprintf(table, "%s\tNão resolvido\t%s\n", host, colorize("FALHA", colorRed, os.Stdout))
			continue
		}

		addresses := make([]string, len(resolvedIPs))
		for i, ip := range resolvedIPs {
			addresses[i] = ip.String()
		}
		fmt.Fprintf(table, "%s\t%s\t%s\n", host, strings.Join(addresses, ", "), colorize("RESOLVIDO", colorGreen, os.Stdout))
	}

	if err := table.Flush(); err != nil {
		printError(fmt.Errorf("erro ao exibir os alvos: %w", err))
		return
	}

	delay := programFlags.i
	times := programFlags.c

	//Enviando pacotes
	fmt.Println("\nEnvio:")
	for packet := 1; packet <= times; packet++ {
		for _, host := range hosts {
			_, resolved := ips[host]
			if resolved {
				message := fmt.Sprintf("[%s] Enviando pacote %d/%d", host, packet, times)
				fmt.Println(colorize(message, colorCyan, os.Stdout))
			}
		}
		if packet < times && delay > 0 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

}

func FFlag(programFlags Flags)([]string){
	//Resolve de onde vem as flas (do terminal ou de um arquivo)
	var hosts []string
	if programFlags.f == ""{
		//flag.Parse()
		hosts = flag.Args() //Retorna o que não for flag ou argumento
	}else{
		var err error
		
		hosts, err = FileToSlice(programFlags.f)

		if err != nil {
			printError(err)
		}
	}
	return hosts
}

func printError(err error) {
	message := "Erro: " + err.Error()
	fmt.Fprintln(os.Stderr, colorize(message, colorRed, os.Stderr))
}

func colorize(message, color string, output *os.File) string {
	if os.Getenv("NO_COLOR") != "" {
		return message
	}

	info, err := output.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice == 0 {
		return message
	}
	return color + message + colorReset
}
