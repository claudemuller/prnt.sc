package internal

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func ShowImage(img image.Image) {
	a := app.New()
	win := a.NewWindow("prnt.sc")

	canvasImg := canvas.NewImageFromImage(img)
	win.SetContent(canvasImg)

	imgWidth := img.Bounds().Size().X
	imgHeight := img.Bounds().Size().Y

	win.Resize(fyne.NewSize(float32(imgWidth), float32(imgHeight)))

	win.ShowAndRun()
}
