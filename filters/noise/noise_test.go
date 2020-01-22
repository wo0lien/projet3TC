package noise

import (
	"testing"

	"github.com/wo0lien/projet3TC/imagetools"
)

func TestMean(t *testing.T) {
	img, err := imagetools.Open("../../assets/bruit2.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := Fmean(img, 3)

	imagetools.Export(result, "mean.png")
}

func TestConcurrentMean(t *testing.T) {
	img, err := imagetools.Open("../../assets/lines.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := ConcurrentFmean(img, 3)

	imagetools.Export(result, "concurrentmean.png")
}
