package noise

import (
	"testing"

	"github.com/wo0lien/projet3TC/imagetools"
)
func TestMediane(t *testing.T) {
	img, err := imagetools.Open("../../assets/imgbruit.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := Fmediane(img, 3)

	imagetools.Export(result, "mediane.png")
}

func TestConcurrentMediane(t *testing.T) {
	img, err := imagetools.Open("../../assets/cielbruit.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := ConcurrentFmediane(img, 4)

	imagetools.Export(result, "concurrentmediane.png")
}


func TestMean(t *testing.T) {
	img, err := imagetools.Open("../../assets/cielbruit.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := Fmean(img, 5)

	imagetools.Export(result, "mean.png")
}

func TestConcurrentMean(t *testing.T) {
	img, err := imagetools.Open("../../assets/cielbruit.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := ConcurrentFmean(img, 3)

	imagetools.Export(result, "concurrentmean.png")
}
