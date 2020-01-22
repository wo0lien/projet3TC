package negative

import (
	"github.com/wo0lien/projet3TC/imagetools"
	"image"
	"image/color"
)

/*
GrayFilter Return the grayscale of an image
*/
func NegativeFilter(imgSrc image.Image) image.Image {

	// Create a new grayscale image
	bounds := imgSrc.Bounds()
	wn, hn := bounds.Min.X, bounds.Min.Y
	wx, hx := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewRGBA(image.Rectangle{image.Point{wn, hn}, image.Point{wx, hx}})
	for x := wn; x < wx; x++ {
		for y := hn; y < hx; y++ {
			imageColor := imgSrc.At(x, y)
			rr, gg, bb, _ := imageColor.RGBA()
			negColor := color.RGBA{uint8(255 - (rr * 255 / 65535)), uint8(255 - (gg * 255 / 65535)), uint8(255 - (bb * 255 / 65535)), 255}
			grayScale.Set(x, y, negColor)
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
func ConcurrentNegFilter(imgSrc image.Image) image.Image {

	out := make(chan portion)
	slices := imagetools.Crop(imgSrc, 4)

	for i := 0; i < 4; i++ {
		go negWorker(i, out, slices[i][0])
	}

	for i := 0; i < 4; i++ {
		slice := <-out
		slices[slice.id][0] = slice.img
	}

	imgEnd := imagetools.Rebuild(slices)

	return imgEnd

}

func negWorker(id int, out chan portion, img image.Image) {
	imgOut := NegativeFilter(img)
	var ret portion
	ret.img = imgOut
	ret.id = id
	out <- ret
}