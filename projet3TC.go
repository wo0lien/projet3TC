package main

import (
	"fmt"
	"log"

	"github.com/wo0lien/projet3TC/filters/edge"
	"github.com/wo0lien/projet3TC/imagetools"
)

/*
Structure des renvois des workers pour la pool de workers du grayscale
*/

func main() {

	//blanc variables only to avoid import errors
	var _ = fmt.Printf
	var _ = edge.FSobel

	testEdge2()

}

func testEdge2() {
	imgSrc, err := imagetools.Open("epice.png")

	if err != nil {
		log.Printf("error loading file")
	}

	slices := imagetools.Crop(imgSrc, 4)

	for i := range slices {
		log.Printf("one filter")
		slices[i][0] = edge.FSobel(slices[i][0])
		log.Printf("one export")
	}
	imgEnd := imagetools.Rebuild(slices)
	imagetools.Export(imgEnd, "export.png")
}

func testEdge() {
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
