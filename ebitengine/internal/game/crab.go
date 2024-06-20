package game

import (
	"bytes"
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
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

	crabSpawnTicker = time.NewTicker(time.Millisecond * 50)
}

type Crab struct {
	id    ecs.Id
	image *ebiten.Image
	pos   Vec2
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
		id:    id,
		image: img,
		pos:   pos,
	}
}

func SpawnCrabs(center Vec2, world *ecs.World) {
	select {
	case <-crabSpawnTicker.C:
		for range 2 {
			id := world.NewId()
			world.Write(id, ecs.C(NewCrab(id, randomPositionAround(center, 500, 1200))))
		}
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

func KillCrabs(tree *kdtree.KDTree, world *ecs.World) {
	bullets := ecs.Query1[Bullet](world)

	bullets.MapId(func(bid ecs.Id, b *Bullet) {
		nn := tree.KNN(&points.Point2D{X: b.pos.X, Y: b.pos.Y}, 1)
		for i := range nn {
			c := nn[i].(*points.Point).Data.(*Crab)
			if b.pos.Distance(c.pos) <= float64(c.image.Bounds().Dx()/2) {
				ecs.Delete(world, c.id)
			}
		}
	})
}

func DrawCrabs(screen *ebiten.Image, op *ebiten.DrawImageOptions, world *ecs.World) {
	q := ecs.Query1[Crab](world)

	q.MapId(func(id ecs.Id, c *Crab) {
		op.GeoM.Reset()
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(float64(c.pos.X), float64(c.pos.Y))
		screen.DrawImage(c.image, op)
	})

}
