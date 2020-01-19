package edge

//sobel algo from https://stackoverflow.com/questions/17815687/image-processing-implementing-sobel-filter
import (
	"image"
	"image/color"
	_ "image/jpeg" //for test
	_ "image/png"  //for test
	"math"

	"github.com/wo0lien/projet3TC/imagetools"
)

/*
FSobel used for filter noise
*/
func FSobel(in image.Image) image.Image {
	loadedImage := in

	//Creation of a new image with the same dimensions as the input one
	b := loadedImage.Bounds()
	minx, miny := b.Min.X, b.Min.Y
	maxx, maxy := b.Max.X, b.Max.Y
	w, h := b.Dx(), b.Dy()
	myImage := image.NewRGBA(loadedImage.Bounds())

	//convertion in greyscale with the BW algo
	gris := make([][]int16, h)
	for i := range gris {
		gris[i] = make([]int16, w)
	}
	for cpt := miny; cpt < maxy; cpt++ {
		for cpt2 := minx; cpt2 < maxx; cpt2++ {
			red, gr, blue, _ := loadedImage.At(cpt2, cpt).RGBA()
			gris[cpt][cpt2] = int16(0.2125*float32(red*255/65535) + 0.7154*float32(gr*255/65535) + 0.0721*float32(blue*255/65535))

		}
	}
	//Edge-detection algorithm applied to each pixel
	var maxG float64 = 0 //we save the highest value of gradient for mapping the values
	gradient := make([][]float64, h)
	for i := range gradient {
		gradient[i] = make([]float64, w)
	}
	for cpt := miny + 1; cpt < maxy-2; cpt++ {
		for cpt2 := minx + 1; cpt2 < maxx-2; cpt2++ {
			var gx float64
			gx = float64(-1*gris[cpt-1][cpt2-1] + 1*gris[cpt+1][cpt2-1] + -2*gris[cpt-1][cpt2] + 2*gris[cpt+1][cpt2] - 1*gris[cpt-1][cpt2+1] + 1*gris[cpt+1][cpt2+1])
			gy := float64(-1*gris[cpt-1][cpt2-1] - 2*gris[cpt][cpt2-1] - 1*gris[cpt+1][cpt2-1] + 1*gris[cpt-1][cpt2+1] + 2*gris[cpt][cpt2+1] + 1*gris[cpt+1][cpt2+1])
			gradient[cpt][cpt2] = math.Sqrt(gx*gx + gy*gy)
			if gradient[cpt][cpt2] > maxG {
				maxG = gradient[cpt][cpt2]
			}

		}
	}
	for cpt := miny + 1; cpt < maxy-2; cpt++ {
		for cpt2 := minx + 1; cpt2 < maxx-2; cpt2++ {
			var valsobel uint8
			if gradient[cpt][cpt2] > 255 {
				valsobel = 255
			}
			valsobel = uint8(gradient[cpt][cpt2] * 255 / maxG)
			myImage.Set(cpt2, cpt, color.RGBA{valsobel, valsobel, valsobel, 255})
		}
	}
	return myImage

}

//------ Beginning of the concurrent part

type portion struct {
	id  int
	img image.Image
}

/*
ConcurrentEdgeFilter Return the grayscale of an image and compute concurrently
*/
func ConcurrentEdgeFilter(imgSrc image.Image) image.Image {

	out := make(chan portion)
	slices := imagetools.CropChevauchement(imgSrc, 4, 2) //on laisse 2 pixels de chevauchement

	for i := 0; i < 4; i++ {
		go edgWorker(i, out, slices[i][0])
	}

	for i := 0; i < 4; i++ {
		slice := <-out
		slices[slice.id][0] = slice.img
	}

	imgEnd := imagetools.RebuildChevauchement(slices, 2)

	return imgEnd

}

func edgWorker(id int, out chan portion, img image.Image) {
	imgOut := FSobel(img)
	var ret portion
	ret.img = imgOut
	ret.id = id
	out <- ret
}
