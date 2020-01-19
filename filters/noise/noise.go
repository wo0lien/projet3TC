package noise

import (
	"image"
	"image/color"
	"math"

	"github.com/wo0lien/projet3TC/filters/uint32slice"
	"github.com/wo0lien/projet3TC/imagetools"
)

// Fmediane utile
func Fmediane(in image.Image, p int) image.Image {
	loadedImage := in

	b := loadedImage.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y
	myImage := image.NewRGBA(image.Rect(0, 0, imgWidth-2*p, imgHeight-2*p))
	t := (2*p + 1) * (2*p + 1)
	var red = make([]uint32, t)
	var green = make([]uint32, t)
	var blue = make([]uint32, t)

	for cpt := p; cpt < imgWidth-p; cpt++ {
		for cpt2 := p; cpt2 < imgHeight-p; cpt2++ {
			i := 0
			for cptwi := -p; cptwi < p+1; cptwi++ {
				for cpthe := -p; cpthe < p+1; cpthe++ {
					red[i], green[i], blue[i], _ = loadedImage.At(cpt+cptwi, cpt2+cpthe).RGBA()
					i++
				}
			}
			uint32slice.SortUint32s(red)
			uint32slice.SortUint32s(green)
			uint32slice.SortUint32s(blue)
			ind := uint(math.Floor(float64(t) / 2))
			valrouge, valvert, valbleu := uint8(red[ind]*255/65535), uint8(green[ind]*255/65535), uint8(blue[ind]*255/65535)
			myImage.Set(cpt-p, cpt2-p, color.RGBA{valrouge, valvert, valbleu, 255})
		}
	}
	return myImage

}

// Fmean utile
func Fmean(img image.Image, p int) image.Image {

	b := img.Bounds()
	minx, miny := b.Min.X, b.Min.Y
	maxx, maxy := b.Max.X, b.Max.Y
	myImage := image.NewRGBA(image.Rect(minx, miny, maxx-2*p, maxy-2*p))
	var valred uint32
	var valgreen uint32
	var valblue uint32

	for cpt := minx + p; cpt < maxx-p; cpt++ {
		for cpt2 := miny + p; cpt2 < maxy-p; cpt2++ {
			i := 0
			valred, valgreen, valblue = 0, 0, 0
			for cptwi := -p; cptwi < p+1; cptwi++ {
				for cpthe := -p; cpthe < p+1; cpthe++ {
					red, green, blue, _ := img.At(cpt+cptwi, cpt2+cpthe).RGBA()
					valred, valgreen, valblue = valred+red, valgreen+green, valblue+blue
					i++
				}
			}

			valrouge, valvert, valbleu := uint8((valred/(uint32(i)+1))*255/65535), uint8((valgreen/(uint32(i)+1))*255/65535), uint8((valblue/(uint32(i)+1))*255/65535)
			myImage.Set(cpt-p, cpt2-p, color.RGBA{valrouge, valvert, valbleu, 255})
		}
	}
	return myImage

}

//-----Beginning fo the concurrent part

type portion struct {
	id  int
	img image.Image
}

/*
ConcurrentFmean Return the image with less noise and compute concurrently
*/
func ConcurrentFmean(imgSrc image.Image, p int) image.Image {

	out := make(chan portion)
	slices := imagetools.Crop(imgSrc, 4)

	for i := 0; i < 4; i++ {
		go meanWorker(i, p, out, slices[i][0])
	}

	for i := 0; i < 4; i++ {
		slice := <-out
		slices[slice.id][0] = slice.img
	}

	imgEnd := imagetools.Rebuild(slices)

	return imgEnd

}

func meanWorker(id int, p int, out chan portion, img image.Image) {
	imgOut := Fmean(img, p)
	var ret portion
	ret.img = imgOut
	ret.id = id
	out <- ret
}
