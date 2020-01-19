package grayscale

import (
	"image"
	"image/color"
	"math"

	"github.com/wo0lien/projet3TC/imagetools"
)

/*
GrayFilter Return the grayscale of an image
*/
func GrayFilter(imgSrc image.Image) image.Image {

	// Create a new grayscale image
	bounds := imgSrc.Bounds()
	wn, hn := bounds.Min.X, bounds.Min.Y
	wx, hx := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{wn, hn}, image.Point{wx, hx}})
	for x := wn; x < wx; x++ {
		for y := hn; y < hx; y++ {
			imageColor := imgSrc.At(x, y)
			rr, gg, bb, _ := imageColor.RGBA()
			r := math.Pow(float64(rr), 2.2)
			g := math.Pow(float64(gg), 2.2)
			b := math.Pow(float64(bb), 2.2)
			m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
			Y := uint16(m + 0.5)
			grayColor := color.Gray{uint8(Y >> 8)}
			grayScale.Set(x, y, grayColor)
		}
	}

	return grayScale
}

//------ Beginning of the concurrent part

type portion struct {
	id  int
	img image.Image
}

/*
ConcurrentGrayFilter Return the grayscale of an image and compute concurrently
*/
func ConcurrentGrayFilter(imgSrc image.Image) image.Image {

	out := make(chan portion)
	slices := imagetools.Crop(imgSrc, 4)

	for i := 0; i < 4; i++ {
		go gsWorker(i, out, slices[i][0])
	}

	for i := 0; i < 4; i++ {
		slice := <-out
		slices[slice.id][0] = slice.img
	}

	imgEnd := imagetools.Rebuild(slices)

	return imgEnd

}

func gsWorker(id int, out chan portion, img image.Image) {
	imgOut := GrayFilter(img)
	var ret portion
	ret.img = imgOut
	ret.id = id
	out <- ret
}
