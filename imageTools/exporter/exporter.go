package exporter

import (
	"image"
	"image/png"
	"log"
	"os"
)

/*
Export create a png file based on the image file given
*/
func Export(img image.Image, name string) {
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
