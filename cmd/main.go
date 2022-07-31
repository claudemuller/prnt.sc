package main

import (
	"flag"
	"log"

	"prnt.sc/internal/pkg"
)

func main() {
	const (
		prntscURL = "https://prnt.sc/"
		idLen     = 6
	)

	maxRetries := flag.Int("retries", 3, "the number of retries if an image URL can't be found")

	flag.Parse()
	log.SetPrefix("prnt.sc >> ")

	img, _ := pkg.GetNewImage(prntscURL, idLen, maxRetries)
	pkg.ShowImage(img)
}
