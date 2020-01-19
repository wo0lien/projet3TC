package imagetools

import (
	"image"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestOpen(t *testing.T) {

	ret, err := Open("../assets/test.png")

	if err != nil {
		t.Error("raise error when loaded with a proper file")
	}
	t.Log(reflect.TypeOf(ret))
}

func TestExport(t *testing.T) {

	os.Remove("testExport.png")
	// create a fake image
	upLeft := image.Point{0, 0}
	lowRight := image.Point{100, 100}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	Export(img, "testExport.png")

	err := os.Remove("testExport.png")
	if err != nil {
		t.Error("file not created")
	}

}

func TestRebuild(t *testing.T) {
	img, err := Open("../assets/test.png")

	wrepeat := 1
	hrepeat := 4

	imgbounds := img.Bounds()
	wimg, himg := imgbounds.Max.X, imgbounds.Max.Y

	if err != nil {
		t.Error("open does not work")
	}

	//4 en x 3 en y
	matrix := make([][]image.Image, hrepeat)
	for y := range matrix {
		matrix[y] = make([]image.Image, wrepeat)
	}

	//populate the matrix
	for x := 0; x < hrepeat; x++ {
		for y := 0; y < wrepeat; y++ {
			matrix[x][y] = img
		}
	}

	composed := Rebuild(matrix)

	bounds := composed.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	if w != wimg*wrepeat {
		t.Error("Pas la bonne largeur")
	}

	if h != himg*hrepeat {
		t.Error("Pas la bonne hauteur")
	}

	Export(composed, "composed.png")
}

func TestCut(t *testing.T) {
	img, _ := Open("../assets/test.png")
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	slices := Cut(img, 3)

	for i := range slices {
		Export(slices[i][0], "slice"+strconv.Itoa(i)+".png")
	}

	lastSlice := slices[len(slices)-1][0]
	b := lastSlice.Bounds()

	if w != b.Max.X {
		t.Error("Pas la bonne largeur")
	}
	if h != b.Max.Y {
		t.Error("Pas la bonne hauteur")
	}
}

func TestCrop(t *testing.T) {
	img, _ := Open("../assets/test.png")

	slices := Crop(img, 3)

	for i := range slices {
		Export(slices[i][0], "slice"+strconv.Itoa(i)+".png")
	}
}
