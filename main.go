// Convert an image to a photomosaic
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func main() {

	pictureFile, er := os.Open("images/cover4.jpg")
	if er != nil {
		fmt.Print(er.Error())
	}
	defer pictureFile.Close()

	si, _ := jpeg.Decode(pictureFile)

	width := 1200
	height := 675
	origin := image.Point{0, 0}
	end := image.Point{width, height}

	newImage := image.NewRGBA(image.Rectangle{origin, end})

	blue := color.RGBA{0, 0, 255, 0xff}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				newImage.Set(x, y, blue)
			case x >= width/2 && y >= height/2: // lower right quadrant
				newImage.Set(x, y, color.White)
			default:
				newImage.Set(x, y, si.At(x, y))
			}
		}
	}

	oi, _ := os.Create("images/output.jpg")
	err := jpeg.Encode(oi, newImage, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	oi.Close()
}
