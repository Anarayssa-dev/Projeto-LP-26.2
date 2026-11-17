package main

import(
	"fmt"
	"flag"
)

func main() {
	programFlags := GetFlags()
	
	//Para utilizar flag.Args() é necessário chamar flag.Parse() presente em Flags()

	var hosts []string

	if programFlags.f == ""{
		//flag.Parse()
		hosts = flag.Args() //Retorna o que não for flag ou argumento
	}else{
		var err error
		
		hosts, err = FileToSlice(programFlags.f)

		if err != nil{
			fmt.Println(err)
		}
	}
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
