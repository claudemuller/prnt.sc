package main

import (
	"flag"
	"image"
	"log"
	"os"

	"prnt.sc/internal/pkg"

	"gioui.org/app"
	"gioui.org/unit"
)

func main() {
	state := pkg.NewState(flag.Int("retries", 3, "the number of retries if an image URL can't be found"))

	flag.Parse()
	log.SetPrefix("prnt.sc >> ")

	go func() {
		state.Win = app.NewWindow(func(m unit.Metric, c *app.Config) {
			c.Title = "prnt.sc"
			c.Size = image.Point{X: 50, Y: 50}
		})

		err := pkg.Run(state)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()
	app.Main()
}
