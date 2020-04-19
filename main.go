// Convert an image to a photomosaic
package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
)

const (
	pieceWidth  = 1
	pieceHeight = 1
)

type piece struct {
	origin        [2]int
	averageColour [3]uint32
}

func main() {
	// Load the source image
	src, _ := imaging.Open("images/sky.jpg")

	// Determine new dimentions to be a multiple of the respective piece dimentions
	newx := (src.Bounds().Max.X / pieceWidth) * pieceWidth
	newy := (src.Bounds().Max.Y / pieceHeight) * pieceHeight

	// Crop the image around the edges
	croppedImage := imaging.CropCenter(src, newx, newy)
	// fmt.Printf("%T\n", croppedImage.At(600, 400))
	// s, q, f, k := croppedImage.At(600, 400).RGBA()
	// fmt.Printf("%v\t%v\t%v\t%v\n", (s)>>8, (q)>>8, (f)>>8, k)

	numPiecesX := newx / pieceWidth
	numPiecesY := newy / pieceHeight
	// totalPieces := numPiecesX * numPiecesY

	var pieces [][]piece
	originx, originy := 0, 0
	for j := 0; j < numPiecesY; j++ {
		originy = j*pieceHeight + 1

		var row []piece
		for i := 0; i < numPiecesX; i++ {
			originx = i*pieceWidth + 1
			p := piece{origin: [2]int{originx, originy}}
			row = append(row, p)
		}
		pieces = append(pieces, row)
	}
	// width := 1200
	// height := 650
	orig := image.Point{0, 0}
	end := image.Point{newx, newy}
	// cyan := color.RGBA{100, 200, 200, 0xff}

	newImage := image.NewRGBA(image.Rectangle{orig, end})
	// fmt.Println(pieces)
	// For each piece
	for u, m := range pieces {
		for v, n := range m {
			var sum [3]uint32

			// For each pixel within that piece
			for i := n.origin[0]; i < n.origin[0]+pieceWidth; i++ {
				for j := n.origin[1]; j < n.origin[1]+pieceHeight; j++ {
					r, g, b, a := croppedImage.At(i, j).RGBA()
					_ = a
					// if (r + g + b) != 0 {
					// 	fmt.Println(r>>8, g>>8, b>>8)
					// }
					sum[0] += (r)
					sum[1] += (g)
					sum[2] += (b)
					// newImage.Set(i, j, croppedImage.At(i, j))
				}
			}
			sum[0] = sum[0] / (pieceWidth * pieceHeight)
			sum[1] = sum[1] / (pieceWidth * pieceHeight)
			sum[2] = sum[2] / (pieceWidth * pieceHeight)
			pieces[u][v].averageColour = sum
			for i := n.origin[0]; i < n.origin[0]+50; i++ {
				for j := n.origin[1]; j < n.origin[1]+50; j++ {
					newImage.Set(i, j, color.RGBA{uint8(sum[0]), uint8(sum[1]), uint8(sum[2]), 0xff})
				}
			}
		}
	}

	// for _, g := range pieces {
	// 	for _, h := range g {
	// 		fmt.Println(h.averageColour)
	// 	}
	// }

	saveImage(newImage, "images/comeon.jpg")
}

func saveImage(i *image.RGBA, p string) {
	oi, _ := os.Create(p)
	defer oi.Close()
	err := jpeg.Encode(oi, i, nil)
	if err != nil {
		panic(err.Error())
	}
	// func newPiece() *piece {
	// 	return &piece{width: pieceWidth, height: pieceHeight}
	// }
}
