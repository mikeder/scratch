package game

import (
	"bytes"
	"image"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/unitoftime/ecs"
)

var (
	crabImage1 *ebiten.Image
	crabImage2 *ebiten.Image
	crabImage3 *ebiten.Image
	ticker     *time.Ticker
)

func init() {
	c1png, _, err := image.Decode(bytes.NewReader(Crab1_png))
	if err != nil {
		log.Fatal(err)
	}
	c2png, _, err := image.Decode(bytes.NewReader(Crab2_png))
	if err != nil {
		log.Fatal(err)
	}
	c3png, _, err := image.Decode(bytes.NewReader(Crab3_png))
	if err != nil {
		log.Fatal(err)
	}
	crabImage1 = ebiten.NewImageFromImage(c1png)
	crabImage2 = ebiten.NewImageFromImage(c2png)
	crabImage3 = ebiten.NewImageFromImage(c3png)

	ticker = time.NewTicker(time.Millisecond * 20)
}

type Crab struct {
	id       ecs.Id
	image    *ebiten.Image
	pos      Vec2
	velocity Vec2
}

func NewCrab(id ecs.Id, pos Vec2) Crab {
	img := func() *ebiten.Image {
		switch rand.Intn(3) {
		case 0:
			return crabImage1
		case 1:
			return crabImage2
		default:
			return crabImage3
		}
	}()

	return Crab{
		id:       id,
		image:    img,
		pos:      pos,
		velocity: Vec2{rand.Float64(), rand.Float64()},
	}
}

func SpawnCrabs(center Vec2, world *ecs.World) {
	select {
	case <-ticker.C:
		pos := randomPositionAround(center, 400, 800)
		id := world.NewId()
		g := NewCrab(id, pos)
		world.Write(id, ecs.C(g))
	default:
		return
	}
}

func DespawnCrabs(world *ecs.World) {
	q := ecs.Query1[Crab](world)

	q.MapId(func(id ecs.Id, a *Crab) {
		if ok := ecs.Delete(world, id); ok {
			// fmt.Print("deleted: ", id)
		}
	})
}

func DrawCrabs(screen *ebiten.Image, op *ebiten.DrawImageOptions, world *ecs.World) {
	q := ecs.Query1[Crab](world)

	q.MapId(func(id ecs.Id, g *Crab) {
		op.GeoM.Reset()
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(float64(g.pos.X), float64(g.pos.Y))
		screen.DrawImage(g.image, op)
	})
}

func randomPositionAround(pos Vec2, min, max float32) Vec2 {
	angle := 0 + rand.Float64()*(math.Pi*2-0)
	dist := min + rand.Float32()*(max-min)
	offsetX := math.Cos(angle) * float64(dist)
	offsetY := math.Sin(angle) * float64(dist)
	return Vec2{pos.X + (offsetX), pos.Y + (offsetY)}
}
