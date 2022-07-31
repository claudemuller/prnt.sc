package main

import (
	"flag"
	"log"
	"os"

	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"prnt.sc/internal/pkg"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

const (
	prntscURL = "https://prnt.sc/"
	idLen     = 6
)

func main() {
	maxRetries := flag.Int("retries", 3, "the number of retries if an image URL can't be found")

	flag.Parse()
	log.SetPrefix("prnt.sc >> ")

	go func() {
		w := app.NewWindow()

		err := run(w, maxRetries)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window, maxRetries *int) error {
	th := material.NewTheme(gofont.Collection())
	img, _ := pkg.GetNewImage(prntscURL, idLen, maxRetries)

	var ops op.Ops

	var button widget.Clickable

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			resizeWin := func(m unit.Metric, c *app.Config) {
				s := img.Bounds().Size()
				s.Y += 40
				c.Size = s
			}

			for button.Clicked() {
				img, _ = pkg.GetNewImage(prntscURL, idLen, maxRetries)

				w.Option(resizeWin)
			}

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					pngImageOp := paint.NewImageOp(img)
					pngImageOp.Add(gtx.Ops)

					imgWidget := widget.Image{
						Src:   pngImageOp,
						Scale: 1,
					}

					return imgWidget.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &button, "get another pic")
						btn.CornerRadius = 0

						return btn.Layout(gtx)
					})
				}),
			)

			w.Option(resizeWin)

			e.Frame(gtx.Ops)
		}
	}
}
