package main

import (
	_ "image/png"
	"os"

	"github.com/RadhiFadlillah/vision-my-pdf/internal/cli"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cli.NewApp().Run(os.Args); err != nil {
		logrus.Fatalln(err)
	}
}
