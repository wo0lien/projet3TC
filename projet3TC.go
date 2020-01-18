package main

import (
	"fmt"
	"github.com/oliamb/cutter"
	"github.com/wo0lien/projet3TC/filters/edge"
	"github.com/wo0lien/projet3TC/filters/grayscale"
	"image"
	"image/png"
	"log"
	"math"
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
	var _ = cutter.Crop

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

	slices := concurrentFilter(imgSrc)

	for _, slice := range slices {
		grayscale.GrayFilter(slice)
	}
}

func concurrentFilter(img image.Image) []image.Image {

	//traitements sur l'image
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	var numCPU = runtime.NumCPU()
	slice := int(math.Floor(float64(h) / float64(numCPU)))

	slices := make([]image.Image, numCPU+1)

	cpt := 0
	for y := 0; y < h; y = y + slice {

		ymax := min(y+slice, h)
		// create an image copy of the slice
		imgSliced, err := cutter.Crop(img, cutter.Config{
			Width:   w,
			Height:  ymax - y,
			Anchor:  image.Point{0, y},
			Options: cutter.Copy,
		})

		if err != nil {
			log.Printf("Error while slicing the image")
			panic(err.Error())
		}
		//z := strconv.Itoa(y)
		//name := "slice" + z
		//exportImage(imgSliced, name)
		slices[cpt] = imgSliced
		cpt++
	}

	return slices
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func exportImage(img image.Image, name string) {
	// Encode the grayscale image to the new file
	newFileName := name
	newfile, err := os.Create(newFileName)

	if err != nil {
		log.Printf("failed creating png output: %s", err)
		panic(err.Error())
	}
	defer newfile.Close()
	png.Encode(newfile, img)
}
