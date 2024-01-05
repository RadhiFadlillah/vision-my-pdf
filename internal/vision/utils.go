package vision

import (
	"image"
	"io"

	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
)

func validateImage(f io.Reader) (bool, error) {
	// Decode image
	i, _, err := image.Decode(f)
	if err != nil {
		return false, err
	}

	// Make sure image not empty
	bounds := i.Bounds().Size()
	if valid := bounds.X > 1 && bounds.Y > 1; !valid {
		return false, nil
	}

	return true, nil
}

func bpToRect(bp *visionpb.BoundingPoly) image.Rectangle {
	vertices := bp.Vertices
	if len(vertices) != 4 {
		return image.Rectangle{}
	}

	min, max := vertices[0], vertices[2]
	return image.Rect(int(min.X), int(min.Y), int(max.X), int(max.Y))
}
