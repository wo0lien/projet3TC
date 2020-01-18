package imagetools

import (
	"github.com/oliamb/cutter"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
)

/*
Cut slice image using the computer CPUNumber
return a slice of images
*/
func Cut(img image.Image, nbSplit int) [][]image.Image {

	//traitements sur l'image
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	slice := int(math.Floor(float64(h) / float64(nbSplit)))

	slices := make([][]image.Image, nbSplit+1)
	for y := range slices {
		slices[y] = make([]image.Image, 1)
	}

	cpt := 0
	for y := 0; y < h; y = y + slice {
		// create an image copy of the slice
		imgSliced, err := cutter.Crop(img, cutter.Config{
			Width:   w,
			Height:  min(slice, h-y),
			Anchor:  image.Point{0, y},
			Options: cutter.Copy,
		})

		if err != nil {
			log.Printf("Error while slicing the image")
			panic(err.Error())
		}

		slices[cpt][0] = imgSliced
		cpt++
	}

	return slices
}

/*
Rebuild function to create a big image with smaller ones
*/
func Rebuild(t [][]image.Image) image.Image {
	xmax := 0
	ymax := 0
	for y := 0; y < len(t); y++ {
		ymax = ymax + t[y][0].Bounds().Dy()
	}
	for x := 0; x < len(t[0]); x++ {
		xmax = xmax + t[0][x].Bounds().Dx()
	}
	//rectangle for the big image
	r := image.Rectangle{image.Point{0, 0}, image.Point{xmax, ymax}}

	rgba := image.NewRGBA(r)

	xi := 0
	yi := 0

	for y := 0; y < len(t); y++ {
		xi = 0
		for x := 0; x < len(t[y]); x++ {
			pi := image.Point{xi, yi}
			ri := image.Rectangle{pi, pi.Add(t[y][x].Bounds().Size())}

			draw.Draw(rgba, ri, t[y][x], image.Point{0, 0}, draw.Src)

			xi = xi + t[y][x].Bounds().Dx()
		}
		yi = yi + t[y][0].Bounds().Dy()
	}

	return rgba
}

/*
Open is a function to open a file as image
*/
func Open(filepath string) (image.Image, error) {

	infile, err := os.Open(filepath)

	if err != nil {
		log.Printf("failed opening file: %s", err)
		return nil, err
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}
	return img, nil
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

/*
Crop l'image en une multitude de fonctions
*/
func Crop(img image.Image, nbSplit int) [][]image.Image {
	//traitements sur l'image
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	slice := int(math.Floor(float64(h) / float64(nbSplit)))

	//Si on est dans le cas ou la division n'est pas entiere
	if h%nbSplit != 0 {
		slice++
	}

	//handle case where h = nbSplit * slice
	slices := make([][]image.Image, nbSplit)
	for y := range slices {
		slices[y] = make([]image.Image, 1)
	}

	cpt := 0

	for y := 0; y < h; y = y + slice {

		//create a subImage
		rect := image.Rect(0, y, w, min(y+slice, h))
		imgSliced, _ := cropImage(img, rect)

		slices[cpt][0] = imgSliced
		cpt++
	}

	return slices
}

func cropImage(img image.Image, cropRect image.Rectangle) (cropImg image.Image, newImg bool) {
	//Interface for asserting whether `img`
	//implements SubImage or not.
	//This can be defined globally.
	type CropableImage interface {
		image.Image
		SubImage(r image.Rectangle) image.Image
	}

	if p, ok := img.(CropableImage); ok {
		// Call SubImage. This should be fast,
		// since SubImage (usually) shares underlying pixel.
		cropImg = p.SubImage(cropRect)
	} else if cropRect = cropRect.Intersect(img.Bounds()); !cropRect.Empty() {
		// If `img` does not implement `SubImage`,
		// copy (and silently convert) the image portion to RGBA image.
		rgbaImg := image.NewRGBA(cropRect)
		for y := cropRect.Min.Y; y < cropRect.Max.Y; y++ {
			for x := cropRect.Min.X; x < cropRect.Max.X; x++ {
				rgbaImg.Set(x, y, img.At(x, y))
			}
		}
		cropImg = rgbaImg
		newImg = true
	} else {
		// Return an empty RGBA image
		cropImg = &image.RGBA{}
		newImg = true
	}

	return cropImg, newImg
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
