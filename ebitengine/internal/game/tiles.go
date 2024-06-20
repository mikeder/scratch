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
	tileSize = 16
)

var (
	tilesImage *ebiten.Image
	tileMatrix *TileMatrix
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)

	cols := ScreenWidth / tileSize
	rows := ScreenHeight / tileSize

	grass := []int{218, 219, 243, 244}

	layer1 := func() []int {
		var layer []int
		for range cols {
			for range rows {
				layer = append(layer, grass[rand.Intn(len(grass))])
			}
		}
		return layer
	}()

	tileMatrix = &TileMatrix{layers: [][]int{layer1}}
}

func DrawWorld(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	const xCount = ScreenWidth / tileSize

	w := tilesImage.Bounds().Dx()
	tileXCount := w / tileSize

	for _, layer := range tileMatrix.layers {
		for i, tile := range layer {
			op.GeoM.Reset()

			op.GeoM.Translate(float64((i%xCount)*tileSize), float64((i/xCount)*tileSize))
			op.GeoM.Scale(3, 3)

			sx := (tile % tileXCount) * tileSize
			sy := (tile / tileXCount) * tileSize
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}

}
