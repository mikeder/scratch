package game

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	crabDefaultHealth   Health = 5
	playerDefaultHealth Health = 100
	playerInjuredHealth Health = 25
	playerDeadHealth    Health = 0
)

type Health int

func (h *Health) Remove(hr int) {
	tmp := *h - Health(hr)
	if tmp <= 0 {
		*h = Health(0)
		return
	}
	*h = tmp
}

func (h *Health) Add(ha int) {
	*h += Health(ha)
}

var (
	healthPickup *ebiten.Image
)

func init() {
	ebitenPng, _, err := image.Decode(bytes.NewReader(Ebiten_png))
	if err != nil {
		log.Fatal(err)
	}
	healthPickup = ebiten.NewImageFromImage(ebitenPng)
}

func SpawnHealth()
