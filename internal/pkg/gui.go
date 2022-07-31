package pkg

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

const (
	prntscURL = "https://prnt.sc/"
	idLen     = 6
)

func Run(state *State) error {
	th := material.NewTheme(gofont.Collection())
	img, _ := GetNewImage(prntscURL, idLen, state.MaxRetries)

	var ops op.Ops

	var button widget.Clickable

	for {
		e := <-state.Win.Events()
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
				img, _ = GetNewImage(prntscURL, idLen, state.MaxRetries)

				state.Win.Option(resizeWin)
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

			state.Win.Option(resizeWin)

			e.Frame(gtx.Ops)
		}
	}
}

func GetNewImage(prntscURL string, idLen int, maxRetries *int) (image.Image, error) {
	var imgURL string

	retries := 0
	for imgURL == "" {
		if retries == *maxRetries {
			log.Fatalf("retries exhausted")
		}

		if retries > 0 {
			log.Printf("failed to retrieve image, trying again: retry[%d]", retries)
		}

		id, err := GenID(idLen)
		if err != nil {
			log.Fatalf("failed to get random ID: %v", err)
		}

		data, err := Fetch(prntscURL + id)
		if err != nil {
			log.Fatalf("failed to retrieve site data: %v", err)
		}

		imgURL, err = ScrapePicURL(string(data))
		if err != nil {
			log.Fatalf("failed to extract picture: %v", err)
		}

		retries++
	}

	imgData, err := Fetch(imgURL)
	if err != nil {
		log.Fatalf("failed to retrieve image: %v", err)
	}

	img, err := decodeImg(imgURL, imgData)
	if err != nil {
		log.Printf("failed to decode image, fetching a new image: %v", err)
	}

	return img, nil
}

func decodeImg(imgURL string, imgData []byte) (image.Image, error) {
	urlPieces := strings.Split(imgURL, ".")
	ext := urlPieces[len(urlPieces)-1]

	switch ext {
	case "png":
		return png.Decode(bytes.NewReader(imgData))
	case "jpeg":
		fallthrough
	case "jpg":
		return jpeg.Decode(bytes.NewReader(imgData))
	}

	return nil, errors.New("decoding failed")
}
