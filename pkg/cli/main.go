package main

import(
	"fmt"
	"flag"
	"time"
)

func main() {
	programFlags := GetFlags()
	
	//Para utilizar flag.Args() é necessário chamar flag.Parse() presente em Flags()

	hosts := FFlag(programFlags)
	
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

	delay := programFlags.i
	times := programFlags.c
	

	for i := range times{
		fmt.Println("Enviando pacote")

		//Enviar pacote aqui

		if i != times - 1{
			time.Sleep(time.Duration(delay)*time.Millisecond)
		}
	}
	
}

func FFlag(programFlags Flags)([]string){
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
	return hosts
}

