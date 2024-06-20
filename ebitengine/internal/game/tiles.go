package game

import (
	"bytes"
	"image"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

type TileMatrix struct {
	layers [][]int
}

const (
	tileSize  = 16
	tileScale = 2
)

var (
	tilesImage *ebiten.Image
	tileMatrix *TileMatrix

	worldImage *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)

	cols := (ScreenWidth / tileSize)
	rows := (ScreenHeight / tileSize)

	grassTiles := []int{218, 219, 243, 244}
	flowerTiles := []int{301, 302, 303, 304}

	grass := func() []int {
		var layer []int
		for range cols {
			for range rows {
				layer = append(layer, grassTiles[rand.Intn(len(grassTiles))])
			}
		}
		return layer
	}()

	flowers := func() []int {
		var layer []int
		for range cols {
			for range rows {
				rnd := rand.Float64()
				if rnd > 0.95 {
					layer = append(layer, flowerTiles[rand.Intn(len(flowerTiles))])
				} else {
					layer = append(layer, 0)
				}
			}
		}
		return layer
	}()

	tileMatrix = &TileMatrix{layers: [][]int{grass, flowers}}
	drawWorldImage()
}

func DrawWorld(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	// const xCount = ScreenWidth / tileSize

	// w := tilesImage.Bounds().Dx()
	// tileXCount := w / tileSize

	// for _, layer := range tileMatrix.layers {
	// 	for i, tile := range layer {
	// 		op.GeoM.Reset()

	// 		op.GeoM.Translate(float64((i%xCount)*tileSize), float64((i/xCount)*tileSize))
	// 		op.GeoM.Scale(tileScale, tileScale)

	// 		sx := (tile % tileXCount) * tileSize
	// 		sy := (tile / tileXCount) * tileSize
	// 		screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	// 	}
	// }

	op.GeoM.Reset()
	screen.DrawImage(worldImage, op)
}

func drawWorldImage() {
	const xCount = ScreenWidth / tileSize

	w := tilesImage.Bounds().Dx()
	tileXCount := w / tileSize

	op := new(ebiten.DrawImageOptions)
	img := ebiten.NewImage(ScreenWidth, ScreenHeight)
	for _, layer := range tileMatrix.layers {
		for i, tile := range layer {
			op.GeoM.Reset()

			op.GeoM.Translate(float64((i%xCount)*tileSize), float64((i/xCount)*tileSize))
			op.GeoM.Scale(tileScale, tileScale)

			sx := (tile % tileXCount) * tileSize
			sy := (tile / tileXCount) * tileSize
			img.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}
	worldImage = img
}
