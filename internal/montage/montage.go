package montage

import (
	"image"
	"image/draw"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
)

type Montage struct {
	Image   image.Image
	Paths   []string
	YLimits []int
	Bounds  []image.Rectangle
}

func Create(imagePaths ...string) (Montage, error) {
	// Prepare variables
	var empty Montage

	// Convert image paths so it's standalone
	imagePaths = append([]string{}, imagePaths...)

	// If there is only one image, use it as it is
	if len(imagePaths) == 1 {
		img, err := imgio.Open(imagePaths[0])
		if err != nil {
			return empty, err
		}
		imgRect := img.Bounds()

		return Montage{
			Image:   img,
			Paths:   []string{imagePaths[0]},
			YLimits: []int{imgRect.Dy()},
			Bounds:  []image.Rectangle{imgRect},
		}, nil
	}

	// Open and calculate each image bound
	var maxWidth int
	var totalHeight int
	images := make([]image.Image, len(imagePaths))

	for i, imgPath := range imagePaths {
		// Open the image
		img, err := imgio.Open(imgPath)
		if err != nil {
			return empty, err
		}
		images[i] = img

		// Get the image bounds
		rect := img.Bounds()
		totalHeight += rect.Dy()
		if width := rect.Dx(); width > maxWidth {
			maxWidth = width
		}
	}

	// Create an empty canvas
	canvas := image.NewRGBA(image.Rect(0, 0, maxWidth, totalHeight))
	draw.Draw(canvas, canvas.Rect, image.White, image.Pt(0, 0), draw.Over)

	// Put each image in the canvas
	var yOffset int
	var yLimits []int
	var imageBounds []image.Rectangle
	for _, img := range images {
		// Get the image bound
		imgRect := img.Bounds()
		imgWidth, imgHeight := imgRect.Dx(), imgRect.Dy()

		// Draw the image
		drawArea := image.Rect(0, yOffset, imgWidth, yOffset+imgHeight)
		draw.Draw(canvas, drawArea, img, img.Bounds().Min, draw.Over)

		// Increase the Y offset
		yOffset += imgHeight

		// Save the limit and boundaries
		yLimits = append(yLimits, yOffset)
		imageBounds = append(imageBounds, drawArea)
	}

	// Invert image, since Google vision seems to yield better performance
	// with white text on black background.
	canvas = effect.Invert(canvas)

	// Return the montage
	return Montage{
		Image:   canvas,
		Paths:   imagePaths,
		YLimits: yLimits,
		Bounds:  imageBounds,
	}, nil
}

func (m Montage) Name() string {
	var names []string
	for _, imgPath := range m.Paths {
		imgName := filepath.Base(imgPath)
		imgName = strings.TrimSuffix(imgName, "_ocr.png")
		names = append(names, imgName)
	}
	return strings.Join(names, "-") + ".png"
}
