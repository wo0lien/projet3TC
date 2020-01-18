package edge

import (
	"testing"

	"github.com/wo0lien/projet3TC/imagetools"
)

func TestEdge(t *testing.T) {
	img, err := imagetools.Open("../../imagetools/composed.png")

	if err != nil {
		t.Error("open does not work")
	}

	result := FSobel(img)

	imagetools.Export(result, "edge.png")
}
