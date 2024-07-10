package game

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/unitoftime/ecs"
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

type HealthPickup struct {
	id    ecs.Id
	img   *ebiten.Image
	pos   Vec2
	heals int
	uses  int
}

func NewHealthPickup(id ecs.Id, pos Vec2) HealthPickup {
	return HealthPickup{
		id:    id,
		img:   healthPickup,
		pos:   pos,
		heals: 25,
		uses:  1,
	}
}

func SpawnHealth(ticker *time.Ticker, world *ecs.World) {
	select {
	case <-ticker.C:
		q := ecs.Query1[Gopher](world)

		var pos Vec2
		q.MapId(func(id ecs.Id, a *Gopher) {
			pos = a.pos
		})

		id := world.NewId()
		world.Write(id, ecs.C(NewHealthPickup(id, randomPositionAround(pos, 200, 400))))
	default:
		return
	}
}

func DrawHealth(screen *ebiten.Image, op *ebiten.DrawImageOptions, world *ecs.World) {
	q := ecs.Query1[HealthPickup](world)

	q.MapId(func(id ecs.Id, h *HealthPickup) {
		op.GeoM.Reset()
		op.ColorScale.Reset()

		op.GeoM.Scale(1.5, 1.5)
		op.GeoM.Translate(float64(h.pos.X), float64(h.pos.Y))
		op.ColorScale.SetA(1)

		screen.DrawImage(h.img, op)
	})

}

type HealthText struct {
	amount int
	addedT time.Time
	pos    Vec2
}

func NewHealthText(amount int, pos Vec2) HealthText {
	return HealthText{
		amount: amount,
		addedT: time.Now(),
		pos:    pos,
	}
}

func DrawHealthText(screen *ebiten.Image, world *ecs.World) {
	q := ecs.Query1[HealthText](world)

	q.MapId(func(id ecs.Id, a *HealthText) {
		if time.Since(a.addedT) > time.Millisecond*300 {
			ecs.Delete(world, id)
		}

		var txt string
		op := &text.DrawOptions{}
		if a.amount > 0 {
			txt = fmt.Sprintf("+%d", a.amount)
			op.ColorScale.ScaleWithColor(color.NRGBA{50, 255, 50, 255})
		} else {
			txt = fmt.Sprintf("%d", a.amount)
			op.ColorScale.ScaleWithColor(color.NRGBA{255, 50, 50, 255})
		}

		op.GeoM.Translate(a.pos.X, a.pos.Y)
		op.LineSpacing = 24
		op.PrimaryAlign = text.AlignEnd
		text.Draw(screen, txt, &text.GoTextFace{
			Source: textFaceSource,
			Size:   fontSize,
		}, op)
	})
}
