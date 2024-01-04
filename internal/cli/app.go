package cli

import (
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:      "vision-my-pdf",
		Usage:     "generate HOCR using Google Vision API, to be used with OCRmyPDF",
		UsageText: "vision-my-pdf [flags] ocrmypdf-dir",
		Flags:     appFlags,
	}
}
