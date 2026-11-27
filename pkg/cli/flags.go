package main

import(
	"fmt"
	"flag"
	"os"
	"strings"
	"path/filepath"
)

type Flags struct {
    c       int
    t       int
    i    	int
    f 		string
}

func GetFlags()(Flags) {

	/**
	-c [int]: Quantidade de pacotes por envio
	-i [int]: Intervalo entre envio de pacotes
	-t [int]: Define um timeout por pacote em milissegundos
	-f [string]: Lê uma lista de domínios ou IP's de um arquivo .txt
	/?: Exibem a ajuda e a lista completa de comandos e sintaxes no Prompt de Comando.
	**/
	
	//Declarando flags
	times := flag.Int("c", 4, "Define quantos pacotes serão enviados")
	rest := flag.Int("i", 1000, "Define o intervalo entre envio de pacotes em milissegundos")
	timeout := flag.Int("t", 4000, "Define o tempo limite de espera por cada resposta em milissegundos")
	fromFile := flag.String("f", "", "Lê uma lista de domínios ou IP's de um arquivo .txt")
	help := flag.Bool("/?", false, "Exibe a lista de comandos")

	flag.Parse()//Interpreta argumentos recebidos

	//Testando flags

	if *times != 4 {
		fmt.Println("flag -c alterada")
	}

	if *rest != 1000 {
		fmt.Println("flag -i alterada")
	}

	if *timeout != 4000 {
		fmt.Println("flag -t alterada")
	}

	if *fromFile != "" {
		fmt.Println("flag -f alterada")
	}

	if *help {
		fmt.Println("flag /? alterada")
	}

	return Flags{
		c:*times,
		i:*rest,
		t:*timeout,
		f:*fromFile,
	}
}

func validateFlags(flags Flags) error {
	if flags.c <= 0 {
		return fmt.Errorf("flag -c deve ser maior que zero")
	}
	if flags.i < 0 || flags.i > 60000 {
		return fmt.Errorf("flag -i deve estar entre 0 e %d milissegundos", 60000)
	}
	if flags.t <= 0 || flags.t > 5000 {
		return fmt.Errorf("flag -t deve estar entre 1 e %d milissegundos", 5000)
	}
	if flags.f != "" {

        if !strings.EqualFold(filepath.Ext(flags.f), ".txt") {
            return fmt.Errorf("flag -f deve apontar para um arquivo .txt")
        }

        info, err := os.Stat(flags.f)

        if err != nil {
            return fmt.Errorf("não foi possível acessar o arquivo %q: %w", flags.f, err)
        }

        if !info.Mode().IsRegular() {
            return fmt.Errorf("o caminho informado não é um arquivo regular")
        }
    }
	return nil
}
