package cli

import (
	"fmt"
	"image"
	"os"
	fp "path/filepath"

	"github.com/RadhiFadlillah/vision-my-pdf/internal/vision"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

func saveDebugImages(pages []vision.Page, outputDir string) error {
	// Prepare font
	roboto := canvas.NewFontFamily("Roboto")
	roboto.MustLoadSystemFont("Roboto", canvas.FontBold)

	for _, page := range pages {
		if err := saveDebugImage(page, roboto, outputDir); err != nil {
			return err
		}
	}

	return nil
}

func saveDebugImage(page vision.Page, fontFamily *canvas.FontFamily, outputDir string) error {
	// Prepare output for this image
	imgName := cleanFileName(page.Image)
	debugOutput := fp.Join(outputDir, imgName) + ".png"

	// Open the image
	f, err := os.Open(page.Image)
	if err != nil {
		return fmt.Errorf("debug open error for \"%s\": %w", imgName, err)
	}
	defer f.Close()

	// Convert image to canvas
	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("debug canvas error for \"%s\": %w", imgName, err)
	}

	// Create blank canvas
	imgRect := img.Bounds()
	c := canvas.New(float64(imgRect.Dx()), float64(imgRect.Dy()))
	ctx := canvas.NewContext(c)
	ctx.DrawImage(0, 0, img, canvas.DPMM(1.0))
	ctx.SetCoordSystem(canvas.CartesianIV)
	ctx.SetStrokeWidth(2)

	// Draw paragraph box
	font := fontFamily.Face(96, canvas.Red)

	for i, p := range page.Paragraphs {
		pRect := p.BoundingBox
		ctx.SetStrokeColor(canvas.Red)
		drawRect(ctx, pRect, 20)

		text := canvas.NewTextLine(font, fmt.Sprintf("P-%d", i), canvas.Left)
		textX := float64(pRect.Min.X - 20)
		textY := float64(pRect.Min.Y - 20 - 10)
		ctx.DrawText(textX, textY, text)

		// Draw each line
		for _, l := range p.Lines {
			ctx.SetStrokeColor(canvas.Green)
			drawRect(ctx, l.BoundingBox, 10)

			// Draw each word
			for _, w := range l.Words {
				ctx.SetStrokeColor(canvas.Blue)
				drawRect(ctx, w.BoundingBox, 0)
			}
		}
	}

	// Save the debug image
	if err = renderers.Write(debugOutput, c); err != nil {
		return fmt.Errorf("debug save error for \"%s\": %w", imgName, err)
	}

	logrus.Printf("saved debug for \"%s\"", imgName)
	return nil
}

func drawRect(ctx *canvas.Context, rect image.Rectangle, padding int) {
	// Adjust padding
	padding += int(ctx.StrokeWidth)

	// Create rectangle
	ctx.MoveTo(float64(rect.Min.X-padding), float64(rect.Min.Y-padding))
	ctx.LineTo(float64(rect.Max.X+padding), float64(rect.Min.Y-padding))
	ctx.LineTo(float64(rect.Max.X+padding), float64(rect.Max.Y+padding))
	ctx.LineTo(float64(rect.Min.X-padding), float64(rect.Max.Y+padding))
	ctx.LineTo(float64(rect.Min.X-padding), float64(rect.Min.Y-padding))
	ctx.Stroke()
}
