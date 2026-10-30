package main

import(
	"fmt"
	"flag"
)

func main() {
	Flags()
	
	//Para utilizar flag.Args() é necessário chamar flag.Parse() presente em Flags()
	//flag.Parse()
	hosts := flag.Args()
	ips, err := ResolveTargets(hosts)

	//Testando sem usar os parâmetros do terminal
	//ips, err := ResolveTargets([]string{"8.8", "google.com", "google", "8.8.4.4"})
	
	//Printando resultado
	if err != nil {
		fmt.Println(err)
	}
	for chave, valor := range ips{
		fmt.Println(chave + ":", valor)
	}
	
}
