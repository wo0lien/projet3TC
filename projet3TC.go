package main

import (
	"fmt"
	"github.com/wo0lien/projet3TC/filters/edge"
	"github.com/wo0lien/projet3TC/filters/grayscale"
	"image"
	"log"
	"os"
)

/*
Structure des renvois des workers pour la pool de workers du grayscale
*/

func main() {

	//blanc variables only to avoid import errors
	var _ = fmt.Printf
	var _ = edge.FSobel

	filename := "epice.png"
	infile, err := os.Open(filename)

	if err != nil {
		log.Printf("failed opening %s: %s", filename, err)
		panic(err.Error())
	}
	defer infile.Close()

	imgSrc, _, err := image.Decode(infile)
	if err != nil {
		panic(err.Error())
	}

	grayscale.GrayFilter(imgSrc)

}
