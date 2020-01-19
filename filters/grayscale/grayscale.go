package grayscale

import (
	"image"
	"image/color"
	"math"
)

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
