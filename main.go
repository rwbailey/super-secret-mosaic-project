// Convert an image to a photomosaic
package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
)

// Define the dimensions of the pieces
const (
	pieceWidth  = 10
	pieceHeight = 10

	X = 0
	Y = 1

	R = 0
	G = 1
	B = 2
)

type piece struct {
	origin        [2]int
	averageColour [3]uint32
}

func main() {
	// Load the source image
	src, _ := imaging.Open("images/cover.jpg")

	// Determine new dimensions to be a multiple of the respective piece dimensions
	// Integer division helps here
	newx := (src.Bounds().Max.X / pieceWidth) * pieceWidth
	newy := (src.Bounds().Max.Y / pieceHeight) * pieceHeight

	// Crop the image around the edges to fit an exact number of pieces
	croppedImage := imaging.CropCenter(src, newx, newy)

	numPiecesX := newx / pieceWidth
	numPiecesY := newy / pieceHeight
	// totalPieces := numPiecesX * numPiecesY

	// 2D slice to hold the pieces
	var pieces [][]piece

	// The origin of each piece relative to the source image
	originx, originy := 0, 0

	// Populate the pieces slice
	for j := 0; j < numPiecesY; j++ {
		originy = j * pieceHeight

		var row []piece
		for i := 0; i < numPiecesX; i++ {
			originx = i * pieceWidth
			p := piece{origin: [2]int{originx, originy}}
			row = append(row, p)
		}
		pieces = append(pieces, row)
	}

	// Output image created now for debugging purposes
	orig := image.Point{0, 0}
	end := image.Point{newx, newy}
	newImage := image.NewRGBA(image.Rectangle{orig, end})

	// Debug
	yellow := color.RGBA{255, 255, 0, 0xff}
	_ = yellow

	// For each piece
	for u, m := range pieces {
		for v, n := range m {

			// Array to hold the sum of the rbg values of each pixel within the piece
			var sum [3]uint32

			// For each pixel within that piece
			for i := n.origin[X]; i < n.origin[X]+pieceWidth; i++ {
				for j := n.origin[Y]; j < n.origin[Y]+pieceHeight; j++ {
					r, g, b, a := croppedImage.At(i, j).RGBA()
					_ = a
					sum[R] += r
					sum[G] += g
					sum[B] += b
				}
			}

			// *** Average over the piece
			sum[R] = (sum[R] / (pieceWidth * pieceHeight)) >> 8
			sum[G] = (sum[G] / (pieceWidth * pieceHeight)) >> 8
			sum[B] = (sum[B] / (pieceWidth * pieceHeight)) >> 8

			pieces[u][v].averageColour = sum

			// Set colours in output image  here for debugging
			for i := n.origin[X]; i < n.origin[X]+pieceWidth; i++ {
				for j := n.origin[Y]; j < n.origin[Y]+pieceHeight; j++ {
					newImage.Set(i, j, color.RGBA{uint8(sum[R]), uint8(sum[G]), uint8(sum[B]), 0xff})
				}
			}
		}
	}

	saveImage(newImage, "images/output.jpg")
}

func saveImage(i *image.RGBA, p string) {
	oi, _ := os.Create(p)
	defer oi.Close()
	err := jpeg.Encode(oi, i, nil)
	if err != nil {
		panic(err.Error())
	}
}
