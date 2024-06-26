package game

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
	"github.com/unitoftime/ecs"
)

const (
	gopherSpeed float64 = 8.0
)

var (
	gopherImage *ebiten.Image
)

func init() {
	gopherPng, _, err := image.Decode(bytes.NewReader(Gopher_png))
	if err != nil {
		log.Fatal(err)
	}
	gopherImage = ebiten.NewImageFromImage(gopherPng)
}

type Gopher struct {
	health   *Health
	id       ecs.Id
	image    *ebiten.Image
	pos      Vec2
	velocity Vec2
}

func NewGopher(id ecs.Id, pos Vec2) Gopher {
	dh := playerDefaultHealth
	return Gopher{
		id:       id,
		health:   &dh,
		image:    gopherImage,
		pos:      pos,
		velocity: Vec2{0, 0},
	}
}

func MoveGopher(input *input, world *ecs.World) {
	q := ecs.Query1[Gopher](world)

	q.MapId(func(id ecs.Id, a *Gopher) {
		if input.up {
			a.pos.Y -= 1 * gopherSpeed
		}
		if input.down {
			a.pos.Y += 1 * gopherSpeed
		}
		if input.left {
			a.pos.X -= 1 * gopherSpeed
		}
		if input.right {
			a.pos.X += 1 * gopherSpeed
		}
	})
}

func DrawGopher(screen *ebiten.Image, op *ebiten.DrawImageOptions, world *ecs.World) {
	q := ecs.Query1[Gopher](world)

	q.MapId(func(id ecs.Id, g *Gopher) {
		op.ColorScale.Reset()
		op.GeoM.Reset()
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(float64(g.pos.X), float64(g.pos.Y))
		screen.DrawImage(g.image, op)
	})
}

func BulletHitsGopher(world *ecs.World) {
	bullets := ecs.Query1[Projectile](world)
	gopher := ecs.Query1[Gopher](world)

	var player *Gopher
	gopher.MapId(func(id ecs.Id, a *Gopher) {
		player = a
	})

	bullets.MapId(func(bid ecs.Id, b *Projectile) {
		if b.team == teamGopher { // TODO: use marker component
			return
		}
		if b.pos.Distance(player.pos) < float64(player.image.Bounds().Dx()/4) {
			player.health.Remove(25)
			ecs.Delete(world, b.pid)
		}
	})
}

func KillGopher(gs *GameState, tree *kdtree.KDTree, world *ecs.World) {
	q := ecs.Query1[Gopher](world)

	q.MapId(func(id ecs.Id, g *Gopher) {

		// get nearest crab
		nn := tree.KNN(&points.Point2D{X: g.pos.X, Y: g.pos.Y}, 1)
		for i := range nn {
			c := nn[i].(*points.Point).Data.(*Crab)
			dis := g.pos.Distance(c.pos)

			// if crab touches gopher hurts
			if dis < float64(c.image.Bounds().Dx()/4) || dis < float64(c.image.Bounds().Dy()/4) {
				// unless crab is already dead
				if c.IsDead() {
					continue
				}
				c.health.Remove(50)
				g.health.Remove(50)
			}

			if *g.health <= playerDeadHealth {
				// game over
				*gs = GameStateOver
			}
		}
	})
}
