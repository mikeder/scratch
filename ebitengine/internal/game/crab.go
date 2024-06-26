package game

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
	"github.com/unitoftime/ecs"
)

type Dead struct{}

var (
	crabImage1 *ebiten.Image
	crabImage2 *ebiten.Image
	crabImage3 *ebiten.Image
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
}

type Crab struct {
	id        ecs.Id
	health    Health
	image     *ebiten.Image
	pos       Vec2
	spawnedAt time.Time
	killedAt  time.Time
	speed     float64
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
	}

	speed := func() float64 {
		switch rand.Intn(5) {
		case 0:
			return .7
		case 1:
			return .9
		case 2:
			return 1.10
		case 3:
			return 1.15
		case 4:
			return 1.2
		default:
			return 5.0
		}
	}

	return Crab{
		id:        id,
		image:     img(),
		pos:       pos,
		spawnedAt: time.Now(),
		killedAt:  time.Time{},
		speed:     speed(),
		health:    crabDefaultHealth,
	}
}

func (c *Crab) IsDead() bool {
	return c.health <= 0
}

func SpawnCrabs(ticker *time.Ticker, world *ecs.World) {
	select {
	case <-ticker.C:
		q := ecs.Query1[Gopher](world)
		var pos Vec2
		q.MapId(func(id ecs.Id, a *Gopher) {
			pos = a.pos
		})

		for range 2 {
			id := world.NewId()
			world.Write(id, ecs.C(NewCrab(id, randomPositionAround(pos, 500, 1200))))
		}
	default:
		return
	}
}

func MoveCrabs(world *ecs.World) {
	player := ecs.Query1[Gopher](world)
	crabs := ecs.Query1[Crab](world, ecs.Without(Dead{}))

	var pp Vec2
	player.MapId(func(_ ecs.Id, g *Gopher) {
		pp = g.pos
	})

	crabs.MapId(func(id ecs.Id, c *Crab) {
		dir := pp.Sub(c.pos).Clamp(Vec2{-180, -180}, Vec2{180, 180})
		c.pos.X += dir.X * c.speed * 0.005
		c.pos.Y += dir.Y * c.speed * 0.005
	})
}

func BulletHitsCrab(tree *kdtree.KDTree, world *ecs.World) {
	bullets := ecs.Query1[Projectile](world)

	bullets.MapId(func(bid ecs.Id, b *Projectile) {
		nn := tree.KNN(&points.Point2D{X: b.pos.X, Y: b.pos.Y}, 1)
		for i := range nn {
			c := nn[i].(*points.Point).Data.(*Crab)
			if b.team == teamCrab || c.IsDead() {
				continue
			}
			if b.pos.Distance(c.pos) < float64(c.image.Bounds().Dx()/4) {
				c.health.Remove(5)
			}
		}
	})
}

func KillCrabs(counter *uint, world *ecs.World) {
	crabs := ecs.Query1[Crab](world)

	crabs.MapId(func(bid ecs.Id, c *Crab) {
		if c.health <= 0 && c.killedAt.IsZero() {
			*counter += 1
			c.killedAt = time.Now()
			world.Write(c.id, ecs.C(Dead{}))
		}
	})
}

func DeleteCrabs(world *ecs.World) {
	crabs := ecs.Query1[Crab](world, ecs.With(Dead{}))

	crabs.MapId(func(bid ecs.Id, c *Crab) {
		if time.Since(c.killedAt) > 750*time.Millisecond {
			ecs.Delete(world, c.id)
		}
	})
}

func DrawCrabs(screen *ebiten.Image, op *ebiten.DrawImageOptions, world *ecs.World) {
	q := ecs.Query1[Crab](world)

	q.MapId(func(id ecs.Id, c *Crab) {
		alpha := func() float32 {
			v := float32(1)
			t := time.Since(c.spawnedAt).Seconds()
			if t > 1 && c.killedAt.IsZero() {
				return v
			}
			if c.IsDead() {
				return float32(1 - time.Since(c.killedAt).Seconds())
			}

			return float32(0.1 + t)
		}()

		op.GeoM.Reset()
		op.ColorScale.Reset()

		if !c.killedAt.IsZero() {
			op.GeoM.Scale(0.5, -0.5)
			op.GeoM.Translate(0, float64(c.image.Bounds().Dy()/2))
			op.ColorScale.ScaleWithColor(color.Black)
		} else {
			op.GeoM.Scale(0.5, 0.5)
		}
		op.GeoM.Translate(float64(c.pos.X), float64(c.pos.Y))
		op.ColorScale.SetA(alpha)

		screen.DrawImage(c.image, op)
	})

}

func CrabShoots(tree *kdtree.KDTree, ticker *time.Ticker, world *ecs.World) {
	select {
	case <-ticker.C:
		q2 := ecs.Query1[Gopher](world)
		var target Vec2
		q2.MapId(func(id ecs.Id, a *Gopher) {
			target = a.pos
		})

		// get crab nearest to gopher
		nn := tree.KNN(&points.Point2D{X: target.X, Y: target.Y}, 1)
		for i := range nn {
			attacker := nn[i].(*points.Point).Data.(*Crab)

			dir := attacker.pos.Sub(target).Clamp(Vec2{-360, -360}, Vec2{360, 360})

			bid := world.NewId()
			world.Write(bid, ecs.C(NewCrabBullet(attacker.id, bid, attacker.pos, dir)), ecs.C(teamCrab))
		}
	default:
		return
	}
}
