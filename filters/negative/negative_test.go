package negative

import (
	"github.com/wo0lien/projet3TC/imagetools"
	"testing"
)

func TestNegFilter(t *testing.T) {
	img, err := imagetools.Open("../../assets/lignes.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := NegativeFilter(img)

	imagetools.Export(result, "negfilterexport.png")
}

func TestConcurrentNegFilter(t *testing.T) {
	img, err := imagetools.Open("../../assets/hudf.png")

	if err != nil {
		t.Error("open does not work")
	}


	result := ConcurrentNegFilter(img)

	imagetools.Export(result, "concunegfilterexport.png")
}
