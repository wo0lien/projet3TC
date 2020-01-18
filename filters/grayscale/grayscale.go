package grayscale

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

type grayPoint struct {
	x  int
	y  int
	gc color.Gray
}

/*
GrayFilter export the grayscaled version of an image
*/
func GrayFilter(imgSrc image.Image) {

	//on va split l'image en autant de CPU qu'il y a en hauteur
	// var numCPU = runtime.NumCPU()

	//traitements sur l'image
	bounds := imgSrc.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	//imgSliceSize := h / numCPU

	//channel pour permettre de récuperer les points en gris
	out := make(chan grayPoint)

	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})

	// on creer un certain nombre de go routines
	// for splits := 1; splits < numCPU; splits++ {
	// 	min, max := imgSliceSize*(splits-1), imgSliceSize*splits
	// 	go filterer(min, max, imgSrc, out)
	// }
	var slice = 100

	for y := 0; y < h; y += slice {
		ymax := min(y+slice, h)
		go filterer(y, ymax, imgSrc, out)
	}

	// on recupere ce que l'on a envoyer
	for cpt := 0; cpt < h*w; cpt++ {
		gp := <-out
		grayScale.Set(gp.x, gp.y, gp.gc)
	}

	// Encode the grayscale image to the new file
	newFileName := "grayscale.png"
	newfile, err := os.Create(newFileName)
	if err != nil {
		log.Printf("failed creating png output: %s", err)
		panic(err.Error())
	}
	defer newfile.Close()
	png.Encode(newfile, grayScale)
}

func filterer(ymin int, ymax int, img image.Image, out chan grayPoint) {

	//on recupere la largeur
	bounds := img.Bounds()
	w := bounds.Max.X

	for x := 0; x < w; x++ {
		for y := ymin; y < ymax; y++ {
			imageColor := img.At(x, y)
			rr, gg, bb, _ := imageColor.RGBA()
			r := math.Pow(float64(rr), 2.2)
			g := math.Pow(float64(gg), 2.2)
			b := math.Pow(float64(bb), 2.2)
			m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
			Y := uint16(m + 0.5)
			grayColor := color.Gray{uint8(Y >> 8)}
			var gp grayPoint
			gp.x = x
			gp.y = y
			gp.gc = grayColor
			// on renvoie dans le channel la valeur de retour de ce point
			out <- gp
		}
	}

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
