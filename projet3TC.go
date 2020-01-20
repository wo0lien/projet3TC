package main

import (
	"fmt"
	"log"

	"github.com/wo0lien/projet3TC/filters/edge"
	"github.com/wo0lien/projet3TC/filters/grayscale"
	"github.com/wo0lien/projet3TC/imagetools"
)

/*
Structure des renvois des workers pour la pool de workers du grayscale
*/

func main() {

	//blanc variables only to avoid import errors
	var _ = fmt.Printf
	var _ = edge.FSobel
	var _ = grayscale.GrayFilter

	imgSrc, err := imagetools.Open("assets/hubble.png")
	if err != nil {
		log.Printf("error loading file")
	}
    var _ = grayscale.ConcurrentGrayFilter(imgSrc)
	
    //imagetools.Export(imgOut, "hubblegray.png")

}
