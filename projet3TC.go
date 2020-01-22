package main

import (
	"flag"
	"fmt"
	"github.com/wo0lien/projet3TC/tcp/serveur"
)

func main() {

	port := flag.Int("p", 8080, "Port du serveur")
	host := flag.String("h", "localhost", "ip du serveur si on veut le lancer sur un reseau")
	concurrent := flag.Bool("c", false, "Choisit si les filtres doivent etre concurrents ou pas")
	flag.Parse()

	fmt.Println("Starting the server")
	if *concurrent {
		fmt.Println("concurrent")
	} else {
		fmt.Println("non-concurrent")
	}

	//launching the server with concurrent handlings
	serveur.StartServer(*port, *host, *concurrent)

}
