package grayscale

import (
	"github.com/wo0lien/projet3TC/imagetools"
	"testing"
)

func TestGrayFilter(t *testing.T) {
	img, err := imagetools.Open("../../assets/epice.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := GrayFilter(img)

	imagetools.Export(result, "grayfilterexport.png")
}

func TestConcurrentGrayFilter(t *testing.T) {
	img, err := imagetools.Open("../../assets/hubble.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := ConcurrentGrayFilter(img)

	imagetools.Export(result, "concugrayfilterexport.png")
}
