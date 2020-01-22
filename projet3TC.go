package main

import (
	"flag"
	"fmt"

	"github.com/wo0lien/projet3TC/tcp/serveur"
)

func main() {

	port := flag.Int("port", 8080, "Port du serveur")
	flag.Parse()

	fmt.Println("Starting the server")

	//launching the server with concurrent handlings
	serveur.StartServer(*port)

}
