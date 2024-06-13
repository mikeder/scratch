package game

import (
	"bytes"
	"image"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hongshibao/go-kdtree"
	"github.com/unitoftime/ecs"
)

const (
	crabSpeed float64 = 1.5
)

var (
	crabImage1      *ebiten.Image
	crabImage2      *ebiten.Image
	crabImage3      *ebiten.Image
	crabSpawnTicker *time.Ticker
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

	crabSpawnTicker = time.NewTicker(time.Millisecond * 20)
}

var _ kdtree.Point = (Crab{})

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

// Return the total number of dimensions
func (c Crab) Dim() int {
	return 2
}

// Return the value X_{dim}, dim is started from 0
func (c Crab) GetValue(dim int) float64 {
	switch dim {
	case 0:
		return c.pos.X
	case 1:
		return c.pos.Y
	default:
		panic("unsupported dimension")
	}
}

// Return the distance between two points
func (c Crab) Distance(other kdtree.Point) float64 {
	// d = √((x2-x1)2 + (y2-y1)2)
	s0 := math.Pow(other.GetValue(0)-c.GetValue(0), 2)
	s1 := math.Pow(other.GetValue(1)-c.GetValue(1), 2)
	sum := s0 + s1
	ret := math.Sqrt(sum)
	return ret
}

// Return the distance between the point and the plane X_{dim}=val
func (c Crab) PlaneDistance(val float64, dim int) float64 {
	tmp := c.GetValue(dim) - val
	return tmp * tmp
}

func SpawnCrabs(center Vec2, world *ecs.World) {
	select {
	case <-crabSpawnTicker.C:
		pos := randomPositionAround(center, 500, 1200)
		id := world.NewId()
		g := NewCrab(id, pos)
		world.Write(id, ecs.C(g))
	default:
		return
	}
}

func MoveCrabs(world *ecs.World) {
	player := ecs.Query1[Gopher](world)
	crabs := ecs.Query1[Crab](world)

	var pp Vec2
	player.MapId(func(_ ecs.Id, a *Gopher) {
		pp = a.pos
	})

	crabs.MapId(func(id ecs.Id, a *Crab) {
		dir := pp.Sub(a.pos).Clamp(Vec2{-180, -180}, Vec2{180, 180})
		a.pos.X += dir.X * crabSpeed * 0.005
		a.pos.Y += dir.Y * crabSpeed * 0.005
	})

}

func KillCrabs(world *ecs.World) {
	bullets := ecs.Query1[Bullet](world)
	crabs := ecs.Query1[Crab](world)

	bullets.MapId(func(bid ecs.Id, b *Bullet) {
		crabs.MapId(func(cid ecs.Id, c *Crab) {
			if b.pos.Distance(c.pos) <= 20 {
				// ecs.Delete(world, bid)
				ecs.Delete(world, cid)
			}
		})
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
