package cutter

import (
	"github.com/oliamb/cutter"
	"image"
	"log"
	"math"
	"runtime"
)

/*
Cut slice image using the computer CPUNumber
return a slice of images
*/
func Cut(img image.Image) []image.Image {

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
