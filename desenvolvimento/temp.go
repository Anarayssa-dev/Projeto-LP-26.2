package main

import "fmt"

func main() {
	name := "Pedro"
	var message string
	message = fmt.Sprintf("Hi, %v. Welcome!", name)
	fmt.Println(message)
}