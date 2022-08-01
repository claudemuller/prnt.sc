package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
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

	var nextBtn, saveBtn widget.Clickable

	for {
		e := <-state.Win.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			resizeWin := func(m unit.Metric, c *app.Config) {
				s := img.Image.Bounds().Size()
				s.Y += 40
				c.Size = s
			}

			for nextBtn.Clicked() {
				img, _ = GetNewImage(prntscURL, idLen, state.MaxRetries)

				state.Win.Option(resizeWin)
			}

			for saveBtn.Clicked() {
				saveLoc := os.Getenv("HOME")

				f, err := os.Create(fmt.Sprintf("%s/temp/%s.%s", saveLoc, img.Filename["fn"], img.Filename["ext"]))
				if err != nil {
					log.Printf("could not create file: %v", err)
				}
				defer f.Close()

				switch img.Filename["ext"] {
				case "png":
					if err = png.Encode(f, img.Image); err != nil {
						log.Printf("failed to save .png: %v", err)
					}
				case "jpeg":
					fallthrough
				case "jpg":
					if err = jpeg.Encode(f, img.Image, nil); err != nil {
						log.Printf("failed to save .jpg: %v", err)
					}
				}

				state.Win.Option(resizeWin)
			}

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					pngImageOp := paint.NewImageOp(img.Image)
					pngImageOp.Add(gtx.Ops)

					imgWidget := widget.Image{
						Src:   pngImageOp,
						Scale: 1,
					}

					return imgWidget.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &nextBtn, "get another pic")
							btn.CornerRadius = 0

							return btn.Layout(gtx)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &saveBtn, "save pic")
							btn.CornerRadius = 0

							return btn.Layout(gtx)
						}),
					)
				}),
			)

			state.Win.Option(resizeWin)

			e.Frame(gtx.Ops)
		}
	}
}

type Filename map[string]string

type Img struct {
	Filename Filename
	Image    image.Image
}

func GetNewImage(prntscURL string, idLen int, maxRetries *int) (*Img, error) {
	var imgURL string

	retries := 0
	for imgURL == "" {
		if retries > *maxRetries {
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

	i, err := decodeImg(imgURL, imgData)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	return i, nil
}

func decodeImg(imgURL string, imgData []byte) (*Img, error) {
	urlPieces := strings.Split(imgURL, "/")
	filename := strings.Split(urlPieces[len(urlPieces)-1], ".")
	ext := filename[1]
	fn := filename[0]

	var err error

	var im image.Image

	switch ext {
	case "png":
		im, err = png.Decode(bytes.NewReader(imgData))
	case "jpeg":
		fallthrough
	case "jpg":
		im, err = jpeg.Decode(bytes.NewReader(imgData))
	case "gif":
		im, err = gif.Decode(bytes.NewReader(imgData))
	default:
		err = errors.New("decoding failed")
	}

	return &Img{
		Filename: Filename{
			"fn":  fn,
			"ext": ext,
		},
		Image: im,
	}, err
}
