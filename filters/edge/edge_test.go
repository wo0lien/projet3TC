package edge

import (
	"testing"

	"github.com/wo0lien/projet3TC/imagetools"
)

func TestEdgeFilter(t *testing.T) {
	img, err := imagetools.Open("../../assets/edges.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := FSobel(img)

	imagetools.Export(result, "edges.png")
}

func TestConcurrentEdgeFilter(t *testing.T) {
	img, err := imagetools.Open("../../assets/hudf.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := ConcurrentEdgeFilter(img)

	imagetools.Export(result, "concugrayfilterexport.png")
}
