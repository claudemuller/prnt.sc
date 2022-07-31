package pkg

import (
	"bytes"
	"image"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowImage(img image.Image) {
	a := app.New()
	win := a.NewWindow("prnt.sc")
	maxRetries := 3

	var (
		content   *fyne.Container
		canvasImg *canvas.Image
		text1     *widget.Button
	)

	goBtn := func() {
		if img == nil {
			var err error

			img, err = GetNewImage("https://prnt.sc/", 6, &maxRetries)
			if err != nil {
				log.Printf("error getting new image: %v", err)
			}
		}

		imgWidth := float32(img.Bounds().Size().X)
		imgHeight := float32(img.Bounds().Size().Y + 20)

		canvasImg = canvas.NewImageFromImage(img)
		img = nil

		canvasImg.SetMinSize(fyne.Size{Width: imgWidth, Height: imgHeight})
		content = container.New(layout.NewVBoxLayout(), canvasImg, text1)

		win.SetContent(content)
		win.Resize(fyne.NewSize(imgWidth, imgHeight))
		win.CenterOnScreen()
	}

	text1 = widget.NewButton("get a new pic", goBtn)

	goBtn()

	win.ShowAndRun()
}

func GetNewImage(prntscURL string, idLen int, maxRetries *int) (image.Image, error) {
	var imgURL string

	retries := 0
	for imgURL == "" {
		if retries == *maxRetries {
			log.Fatalf("retries exhausted")
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

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	return img, nil
}
