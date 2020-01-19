package noise

import (
	"testing"

	"github.com/wo0lien/projet3TC/imagetools"
)

func TestNoise(t *testing.T) {
	img, err := imagetools.Open("../../assets/cielbruit.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := Fmediane(img, 2)

	imagetools.Export(result, "edge.png")
}
