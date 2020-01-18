package rebuilder

import (
	"image"
	"image/draw"
)

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
