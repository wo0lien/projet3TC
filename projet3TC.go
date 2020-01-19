package main

import (
	"fmt"
	"log"
	"strconv"

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

	useGrayScaleFilter()
	//edgeFilter()

}

func useGrayScaleFilter() {
	imgSrc, err := imagetools.Open("edges.png")

	if err != nil {
		log.Printf("error loading file")
	}

	slices := imagetools.Crop(imgSrc, 4)

	for i := 0; i < 4; i++ {
		log.Printf("one filter")
		slices[i][0] = grayscale.GrayFilter(slices[i][0])

		imagetools.Export(slices[i][0], "e"+strconv.Itoa(i)+".png")
	}
	imgEnd := imagetools.Rebuild(slices)
	imagetools.Export(imgEnd, "export.png")
}

func edgeFilter() {
	imgSrc, err := imagetools.Open("edges.png")

	if err != nil {
		log.Printf("error loading file")
	}

	slices := imagetools.Crop(imgSrc, 4)

	for i := 0; i < 4; i++ {
		log.Printf("one filter")
		slices[i][0] = edge.FSobel(slices[i][0])
	}
	imgEnd := imagetools.Rebuild(slices)
	imagetools.Export(imgEnd, "export.png")
}
