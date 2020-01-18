package main

import (
	"fmt"
	"github.com/wo0lien/projet3TC/filters/edge"
	"github.com/wo0lien/projet3TC/filters/grayscale"
	"github.com/wo0lien/projet3TC/imagetools"
	"image"
	"log"
	"os"
	"runtime"
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
		log.Printf("failed opening file: %s", err)
		panic(err.Error())
	}
	defer infile.Close()

	imgSrc, _, err := image.Decode(infile)
	if err != nil {
		panic(err.Error())
	}

	// grayscale.GrayFilter(imgSrc)

	var numCPU = runtime.NumCPU()

	slices := imagetools.Cut(imgSrc, numCPU)

	for _, slice := range slices {
		grayscale.GrayFilter(slice)
	}

	img, err := imagetools.Open("assets/test.png")

	if err != nil {
		log.Printf("failed opening file: %s", err)
		panic(err.Error())
	}

	t := imagetools.Cut(img, 9)

	for i := range t {
		t[i] = edge.FSobel(t[i])
	}

	result := imagetools.Rebuild(t)

}
