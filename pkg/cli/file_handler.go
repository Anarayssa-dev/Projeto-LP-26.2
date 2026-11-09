package main

import(
	"os"
	"fmt"
	"bufio"
)

func FileToSlice(fileName string) ([]string, error){
	//Arquivo deve estar na raiz do projeto

	//Lendo o arquivo
	file, err := os.Open("../../"+fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)//criando um novo leitor para o arquivo em questão

	var ipsSlice []string

	//Adicionando ao slice linha por linha
	for scanner.Scan() {	//scanner.Scan() retorna True se tiver outra linha para ler
		line := scanner.Text()
		ipsSlice = append(ipsSlice, line)
	}
	
	err = scanner.Err()
	if err != nil {
		fmt.Println("Erro durante a leitura:", err)
		return nil, err
	}

	return ipsSlice,nil
}