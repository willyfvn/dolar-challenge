package main

import (
	"fmt"

	"github.com/willyfvn/dolar-challenge.git/client"
	"github.com/willyfvn/dolar-challenge.git/server"
)

func main() {

	go server.InitializeServer()

	cotacao := client.FetchCotacao()
	fmt.Println(cotacao)

}
