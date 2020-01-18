package imagetools

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/oliamb/cutter"
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

/*
Rebuild function to create a big image with smaller ones
*/
func Rebuild(t [][]image.Image) image.Image {
	xmax := 0
	ymax := 0
	for y := 0; y < len(t); y++ {
		ymax = ymax + t[y][0].Bounds().Dy()
		for x := 0; x < len(t[x]); x++ {
			xmax = xmax + t[y][x].Bounds().Dx()
		}
	}
	//rectangle for the big image
	r := image.Rectangle{image.Point{0, 0}, image.Point{xmax, ymax}}

	rgba := image.NewRGBA(r)

	xi := 0
	yi := 0

	for y := 0; y < len(t); y++ {
		xi = 0
		yi = yi + t[y][0].Bounds().Dy()
		for x := 0; x < len(t[x]); x++ {
			pi := image.Point{xi, yi}
			ri := image.Rectangle{pi, pi.Add(t[x][y].Bounds().Size())}

			draw.Draw(rgba, ri, t[x][y], image.Point{0, 0}, draw.Src)

			xi = xi + t[y][x].Bounds().Dx()
		}
	}

	return rgba
}

/*
Open is a function to open a file as image
*/
func Open(filepath string) image.Image {

	infile, err := os.Open(filepath)

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
