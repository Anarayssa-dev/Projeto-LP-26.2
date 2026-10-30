package main

import(
	"fmt"
	"flag"
	"os"
	"net"
	"errors"
)

func main() {
	Flags()
	ips, err := ResolveTargets([]string{"8.8", "google.com", "google", "8.8.4.4"})
	if err != nil {
		fmt.Println(err)
	}
	for chave, valor := range ips{
		fmt.Println(chave + ":", valor)
	}
	
}

func ResolveTargets(hosts []string) (map[string][]net.IP, error){

	//Slice que vai guardar os ips resultantes
	ips := map[string][]net.IP{}
	failedTargets := []string{}
	for _,target := range hosts {
		//Se o alvo for um ip direto
		ip := net.ParseIP(target)

		if ip != nil{
			ips[target] = []net.IP{ip}
		}

		//Se o alvo for um domínio
		domainIPs, err := net.LookupIP(target)//net.LookupIP retorna um slice com os ips do domínio

		if err == nil{
			ips[target] = domainIPs
		}

		//Se ambos falharem, o alvo não foi resolvido
		if ip == nil && err != nil{
			failedTargets = append(failedTargets, target)
		}
	}

	if len(failedTargets)!=0{
		return ips, errors.New("Não possível resolver os seguintes targets: " + fmt.Sprint(failedTargets))
	}
	
	return ips, nil
}

func Flags() {
	/**
	• -t: Envia pacotes de eco continuamente para o destino até ser interrompido com as teclas Ctrl + C.
	• -a: Resolve endereços IP para nomes de host.
	• -n [número]: Define a quantidade de pacotes de solicitação de eco a serem enviados (o padrão no Windows é 4).
	• -l [tamanho]: Especifica o tamanho do buffer de dados em bytes (o padrão é 32 bytes; o máximo é 65500 bytes).
	• -f: Aplica o sinalizador "Não Fragmentar" (Don't Fragment) em pacotes IPv4.
	• -i [TTL]: Define o valor do TTL (Time to Live) no cabeçalho do pacote (o máximo é 255).
	• -w [tempo]: Define o tempo limite (timeout) de espera por cada resposta em milissegundos (o padrão é 4000 ms, ou seja, 4 segundos).
	• -4: Força o uso exclusivo do protocolo IPv4.
	• -6: Força o uso exclusivo do protocolo IPv6.
	• /?: Exibem a ajuda e a lista completa de comandos e sintaxes no Prompt de Comando.
	**/
	
	//Declarando flags
	continuos := flag.Bool("t", false, "Envia pacotes continuamente para o destino até ser interrompido com Ctrl + C")
	hostName := flag.Bool("a", false, "Resolve endereços IP para nomes de host")
	times := flag.Int("n", 4, "Define quantos pacotes serão enviados")
	bufferLen := flag.Int("l", 32, "Especifica o tamanho do buffer de dados em bytes (máximo de 65500)")
	dontFragment := flag.Bool("f", false, "Aplica o sinalizador 'Não Fragmentar' em pacotes IPv4")
	ttl := flag.Int("i", 128, "Define o valor do TTL (máximo de 255)")
	timeout := flag.Int("w", 4000, "Define o tempo limite de espera por cada resposta em milissegundos")
	ipv4 := flag.Bool("4", false, "Força o uso exclusivo do protocolo IPv4")
	ipv6 := flag.Bool("6", false, "Força o uso exclusivo do protocolo IPv6")
	help := flag.Bool("/?", false, "Exibe a lista de comandos")

	//Validando máximo em flags com limite
	//A biblioteca flag retorna o ponteiro para o valor
	if *bufferLen > 65500 {
		fmt.Fprintln(os.Stderr, "Ocorreu um erro! Tamanho do buffer não pode ser maior que 65500") //os.Stderr: saída de erros
	}
	if *ttl > 255 {
		fmt.Fprintln(os.Stderr, "Ocorreu um erro! TTL não pode ser maior que 255")
	}
	//Declarar mínimo também?

	flag.Parse()//Interpreta argumentos recebidos
	//Testando flags
	if *continuos {
	fmt.Println("flag -t ativada")
	}

	if *hostName {
		fmt.Println("flag -a ativada")
	}

	if *times != 4 {
		fmt.Println("flag -n ativada")
	}

	if *bufferLen != 32 {
		fmt.Println("flag -l ativada")
	}

	if *dontFragment {
		fmt.Println("flag -f ativada")
	}

	if *ttl != 128 {
		fmt.Println("flag -i ativada")
	}

	if *timeout != 4000 {
		fmt.Println("flag -w ativada")
	}

	if *ipv4 {
		fmt.Println("flag -4 ativada")
	}

	if *ipv6 {
		fmt.Println("flag -6 ativada")
	}

	if *help {
		fmt.Println("flag /? ativada")
	}
}