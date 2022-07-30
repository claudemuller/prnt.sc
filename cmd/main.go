package main

import (
	"bytes"
	"flag"
	"image"
	"log"

	"prnt.sc/internal"
)

func main() {
	const (
		prntscURL = "https://prnt.sc/"
		idLen     = 6
	)

	maxRetries := flag.Int("retries", 3, "the number of retries if an image URL can't be found")

	flag.Parse()
	log.SetPrefix("prnt.sc >> ")

	var imgURL string

	retries := 0
	for imgURL == "" {
		if retries == *maxRetries {
			log.Fatalf("retries exhausted")
		}

		id, err := internal.GenID(idLen)
		if err != nil {
			log.Fatalf("failed to get random ID: %v", err)
		}

		data, err := internal.Fetch(prntscURL + id)
		if err != nil {
			log.Fatalf("failed to retrieve site data: %v", err)
		}

		imgURL, err = internal.ScrapePicURL(string(data))
		if err != nil {
			log.Fatalf("failed to extract picture: %v", err)
		}

		retries++
	}

	imgData, err := internal.Fetch(imgURL)
	if err != nil {
		log.Fatalf("failed to retrieve image: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	internal.ShowImage(img)
}
