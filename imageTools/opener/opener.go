package opener

import (
	"image"
	"log"
	"os"
)

func open(name string) image.Image {

	infile, err := os.Open(name)

	if err != nil {
		log.Printf("failed opening file: %s", err)
		panic(err.Error())
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		panic(err.Error())
	}
	return img
}
