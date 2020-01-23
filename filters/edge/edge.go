package edge

//sobel algo from https://stackoverflow.com/questions/17815687/image-processing-implementing-sobel-filter
import (
	"image"
	"image/color"
	_ "image/jpeg" //for test
	_ "image/png"  //for test
	"math"
	"runtime"

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
			gris[cpt-miny][cpt2-minx] = int16(0.2125*float32(red*255/65535) + 0.7154*float32(gr*255/65535) + 0.0721*float32(blue*255/65535))

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
			gx = float64(-1*gris[cpt-1-miny][cpt2-1-minx] + 1*gris[cpt+1-miny][cpt2-1-minx] + -2*gris[cpt-1-miny][cpt2-minx] + 2*gris[cpt+1-miny][cpt2-minx] - 1*gris[cpt-1-miny][cpt2+1-minx] + 1*gris[cpt+1-miny][cpt2+1-minx])
			gy := float64(-1*gris[cpt-1-miny][cpt2-1-minx] - 2*gris[cpt-miny][cpt2-1-minx] - 1*gris[cpt+1-miny][cpt2-1-minx] + 1*gris[cpt-1-miny][cpt2+1-minx] + 2*gris[cpt-miny][cpt2+1-minx] + 1*gris[cpt+1-miny][cpt2+1-minx])
			gradient[cpt-miny][cpt2-minx] = math.Sqrt(gx*gx + gy*gy)
			if gradient[cpt-miny][cpt2-minx] > maxG {
				maxG = gradient[cpt-miny][cpt2-minx]
			}

		}
	}
	for cpt := miny + 1; cpt < maxy-2; cpt++ {
		for cpt2 := minx + 1; cpt2 < maxx-2; cpt2++ {
			var valsobel uint8
			if gradient[cpt-miny][cpt2-minx] > 255 {
				valsobel = 255
			}
			valsobel = uint8(gradient[cpt-miny][cpt2-minx] * 255 / maxG)
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
	ngo:=runtime.NumCPU()
	out := make(chan portion)
	slices := imagetools.CropChevauchement(imgSrc, ngo, 10) //on laisse 5 pixels de chevauchement

	for i := 0; i < ngo; i++ {
		go edgWorker(i, out, slices[i][0])
	}

	for i := 0; i < ngo; i++ {
		slice := <-out
		slices[slice.id][0] = slice.img
	}

	imgEnd := imagetools.RebuildChevauchement(slices, 10)

	return imgEnd

}

func edgWorker(id int, out chan portion, img image.Image) {
	imgOut := FSobel(img)
	var ret portion
	ret.img = imgOut
	ret.id = id
	out <- ret
}
