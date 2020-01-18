package edge

//sobel algo from https://stackoverflow.com/questions/17815687/image-processing-implementing-sobel-filter
import (
	"image"
	"image/color"
	_ "image/jpeg" //for test
	_ "image/png"  //for test
	"math"
)

/*
FSobel used for filter noise
*/
func FSobel(in image.Image) image.Image {
	loadedImage := in

	//Creation of a new image with the same dimensions as the input one
	b := loadedImage.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y
	myImage := image.NewRGBA(loadedImage.Bounds())

	//convertion in greyscale with the BW algo
	gris := make([][]int16, imgHeight)
	for i := range gris {
		gris[i] = make([]int16, imgWidth)
	}
	for cpt := 0; cpt < imgHeight; cpt++ {
		for cpt2 := 0; cpt2 < imgWidth; cpt2++ {
			red, gr, blue, _ := loadedImage.At(cpt2, cpt).RGBA()
			gris[cpt][cpt2] = int16(0.2125*float32(red*255/65535) + 0.7154*float32(gr*255/65535) + 0.0721*float32(blue*255/65535))

		}
	}
	//Edge-detection algorithm applied to each pixel
	var maxG float64 = 0 //we save the highest value of gradient for mapping the values
	gradient := make([][]float64, imgHeight)
	for i := range gradient {
		gradient[i] = make([]float64, imgWidth)
	}
	for cpt := 1; cpt < imgHeight-2; cpt++ {
		for cpt2 := 1; cpt2 < imgWidth-2; cpt2++ {
			var gx float64
			gx = float64(-1*gris[cpt-1][cpt2-1] + 1*gris[cpt+1][cpt2-1] + -2*gris[cpt-1][cpt2] + 2*gris[cpt+1][cpt2] - 1*gris[cpt-1][cpt2+1] + 1*gris[cpt+1][cpt2+1])
			gy := float64(-1*gris[cpt-1][cpt2-1] - 2*gris[cpt][cpt2-1] - 1*gris[cpt+1][cpt2-1] + 1*gris[cpt-1][cpt2+1] + 2*gris[cpt][cpt2+1] + 1*gris[cpt+1][cpt2+1])
			gradient[cpt][cpt2] = math.Sqrt(gx*gx + gy*gy)
			if gradient[cpt][cpt2] > maxG {
				maxG = gradient[cpt][cpt2]
			}

		}
	}
	for cpt := 1; cpt < imgHeight-2; cpt++ {
		for cpt2 := 1; cpt2 < imgWidth-2; cpt2++ {
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
