package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
)

const (
	InstagramWidth  int    = 1080
	InstagramHeight int    = 1350
	ResultDirectory string = "carousel"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Missing image path argument")
		os.Exit(1)
	}
	imageFileName := args[0]
	file, err := os.Open(imageFileName)
	if err != nil {
		fmt.Println("Unable to open image")
		os.Exit(1)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("Unable to decode image")
		os.Exit(1)
	}
	err = os.MkdirAll(ResultDirectory, 0755)
	if err != nil {
		fmt.Println("Unable to create carousel folder")
		os.Exit(1)
	}
	carousel := splitCarousel(img)
	for idx, block := range carousel {
		destFile := filepath.Join(ResultDirectory, fmt.Sprintf("%d.jpg", idx+1))
		err := saveImage(block, destFile)
		if err != nil {
			fmt.Println(fmt.Sprintf("Unable to save image number %d", idx+1))
		}
	}

	fmt.Println("Carousel saved successfully")
}

func splitCarousel(img image.Image) []image.Image {
	numberOfImages := img.Bounds().Max.X / InstagramWidth
	carousel := make([]image.Image, numberOfImages)

	for x := 0; x < numberOfImages; x++ {
		xStart := x * InstagramWidth
		xEnd := xStart + InstagramWidth
		carousel[x] = img.(SubImager).SubImage(image.Rect(xStart, 0, xEnd, InstagramHeight))
	}

	return carousel
}

func saveImage(img image.Image, filename string) error {
	croppedImageFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer croppedImageFile.Close()
	if err := jpeg.Encode(croppedImageFile, img, nil); err != nil {
		return err
	}
	return nil
}
