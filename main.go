// Convert an image to a photomosaic
package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
)

// Define the dimentions of the pieces
const (
	pieceWidth  = 50
	pieceHeight = 50
)

type piece struct {
	origin        [2]int
	averageColour [3]uint32
}

func main() {
	// Load the source image
	src, _ := imaging.Open("images/cover4.jpg")

	// Determine new dimentions to be a multiple of the respective piece dimentions
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

			// newImage.Set(n.origin[0], n.origin[1], yellow)

			// For each pixel within that piece
			for i := n.origin[0]; i < n.origin[0]+pieceWidth; i++ {
				for j := n.origin[1]; j < n.origin[1]+pieceHeight; j++ {
					r, g, b, a := croppedImage.At(i, j).RGBA()
					_ = a
					sum[0] += r
					sum[1] += g
					sum[2] += b
					// newImage.Set(i, j, croppedImage.At(i, j))
					// if i == n.origin[0] || j == n.origin[1] {
					// 	newImage.Set(i, j, yellow)
					// }
				}
			}

			// fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", sum[0], sum[1], sum[2], sum[0]/(pieceWidth*pieceHeight), sum[1]/(pieceWidth*pieceHeight), sum[2]/(pieceWidth*pieceHeight), uint8(sum[0]/(pieceWidth*pieceHeight)), uint8(sum[1]/(pieceWidth*pieceHeight)), uint8(sum[2]/(pieceWidth*pieceHeight)))

			// *** Average over the piece
			sum[0] = sum[0] / (pieceWidth * pieceHeight)
			sum[1] = sum[1] / (pieceWidth * pieceHeight)
			sum[2] = sum[2] / (pieceWidth * pieceHeight)

			pieces[u][v].averageColour = sum

			// Set colours in output image  here for debugging
			for i := n.origin[0]; i < n.origin[0]+pieceWidth; i++ {
				for j := n.origin[1]; j < n.origin[1]+pieceHeight; j++ {
					newImage.Set(i, j, color.RGBA{uint8(sum[0] >> 8), uint8(sum[1] >> 8), uint8(sum[2] >> 8), 0xff})
				}
			}
		}
	}

	// for _, g := range pieces {
	// 	for _, h := range g {
	// 		fmt.Println(h.averageColour)
	// 	}
	// }

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
