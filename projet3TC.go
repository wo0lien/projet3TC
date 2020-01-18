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

	for i := range slices {
		grayscale.GrayFilter(slices[i][0])
	}

	img, err := imagetools.Open("assets/test.png")

	if err != nil {
		log.Printf("failed opening file: %s", err)
		panic(err.Error())
	}

	t := imagetools.Cut(img, 4)

	for i := range t {
		for j := 0; j < len(t[i]); j++ {
			t[i][j] = edge.FSobel(t[i][j])
		}
	}

	result := imagetools.Rebuild(t)

	imagetools.Export(result, "result.png")
}
