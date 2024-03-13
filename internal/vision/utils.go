package vision

import (
	"image"

	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
)

func bpToRect(bp *visionpb.BoundingPoly) image.Rectangle {
	vertices := bp.Vertices
	if len(vertices) != 4 {
		return image.Rectangle{}
	}

	min, max := vertices[0], vertices[2]
	return image.Rect(int(min.X), int(min.Y), int(max.X), int(max.Y))
}
