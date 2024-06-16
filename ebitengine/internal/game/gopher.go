package game

import (
	"bytes"
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/unitoftime/ecs"
)

const (
	bulletSpeed    float64       = 16.0
	bulletLifetime time.Duration = 3 * time.Second
	gopherSpeed    float64       = 8.0
)

var (
	bulletImage1 *ebiten.Image
	bulletImage2 *ebiten.Image
	gopherImage  *ebiten.Image
)

func init() {
	gopherPng, _, err := image.Decode(bytes.NewReader(Gopher_png))
	if err != nil {
		log.Fatal(err)
	}
	gopherImage = ebiten.NewImageFromImage(gopherPng)

	bulletPng1, _, err := image.Decode(bytes.NewReader(GoBullet1_png))
	if err != nil {
		log.Fatal(err)
	}
	bulletImage1 = ebiten.NewImageFromImage(bulletPng1)

	bulletPng2, _, err := image.Decode(bytes.NewReader(GoBullet2_png))
	if err != nil {
		log.Fatal(err)
	}
	bulletImage2 = ebiten.NewImageFromImage(bulletPng2)
}

type Gopher struct {
	id       ecs.Id
	image    *ebiten.Image
	pos      Vec2
	velocity Vec2
}

func NewGopher(id ecs.Id, pos Vec2) Gopher {
	return Gopher{
		id:       id,
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
		op.GeoM.Reset()
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(float64(g.pos.X), float64(g.pos.Y))
		screen.DrawImage(g.image, op)
	})
}

type Bullet struct {
	id    ecs.Id
	image *ebiten.Image

	expires time.Time
	dir     Vec2
	pos     Vec2
}

func NewBullet(id ecs.Id, pos Vec2, dir Vec2) Bullet {
	img := func() *ebiten.Image {
		switch rand.Intn(2) {
		case 0:
			return bulletImage1
		default:
			return bulletImage2
		}
	}()

	return Bullet{
		id:      id,
		image:   img,
		pos:     pos,
		dir:     dir,
		expires: time.Now().Add(bulletLifetime),
	}
}

func SpawnBullets(center Vec2, input *input, world *ecs.World) {

	if input.fire {
		q := ecs.Query1[Gopher](world)

		var pos Vec2
		q.MapId(func(id ecs.Id, a *Gopher) {
			pos = a.pos.Add(Vec2{20, 0})
		})

		dir := pos.Sub(input.cursor).Clamp(Vec2{-360, -360}, Vec2{360, 360})
		for i := range 1 {
			f := float64(i) * 10
			bid := world.NewId()
			world.Write(bid, ecs.C(NewBullet(bid, pos, dir.Add(Vec2{f, f}))))
		}
	}
}

func MoveBullets(world *ecs.World) {
	q := ecs.Query1[Bullet](world)
	q.MapId(func(id ecs.Id, a *Bullet) {
		a.pos.X -= a.dir.X * bulletSpeed * 0.005
		a.pos.Y -= a.dir.Y * bulletSpeed * 0.005
	})
}

func ExpireBullets(world *ecs.World) {
	q := ecs.Query1[Bullet](world)

	q.MapId(func(id ecs.Id, a *Bullet) {
		if a.expires.Before(time.Now()) {
			ecs.Delete(world, id)
		}
	})
}

func DrawBullets(screen *ebiten.Image, op *ebiten.DrawImageOptions, world *ecs.World) {
	q := ecs.Query1[Bullet](world)

	q.MapId(func(id ecs.Id, g *Bullet) {
		op.GeoM.Reset()
		op.GeoM.Scale(1.5, 1.5)
		op.GeoM.Translate(float64(g.pos.X), float64(g.pos.Y))
		screen.DrawImage(g.image, op)
	})
}
