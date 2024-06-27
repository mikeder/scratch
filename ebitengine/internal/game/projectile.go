package game

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/unitoftime/ecs"
)

const (
	teamGopher int = iota
	teamCrab

	bulletSpeed    float64       = 16.0 * 0.005
	bulletLifetime time.Duration = 5 * time.Second
)

var (
	bulletImage1    *ebiten.Image
	bulletImage2    *ebiten.Image
	crabBulletImage *ebiten.Image
)

func init() {
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

	crabBulletImage = newCrabBulletImage()
}

func newCrabBulletImage() *ebiten.Image {
	l1 := ebiten.NewImage(16, 16)
	l1.Fill(color.NRGBA{255, 255, 255, 128})

	l2 := ebiten.NewImage(8, 8)
	l2.Fill(color.NRGBA{255, 25, 25, 255})

	l3 := ebiten.NewImage(4, 4)
	l3.Fill(color.NRGBA{255, 255, 255, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(l1.Bounds().Dx())/4, float64(l1.Bounds().Dy())/4)
	l1.DrawImage(l2, op)
	l1.DrawImage(l3, op)

	return l1
}

type Projectile struct {
	pid     ecs.Id
	oid     ecs.Id
	image   *ebiten.Image
	expires time.Time
	dir     Vec2
	pos     Vec2
	speed   float64
	team    int
}

func NewGoBullet(pid, oid ecs.Id, speed float64, dir, pos Vec2) Projectile {
	img := func() *ebiten.Image {
		switch rand.Intn(2) {
		case 0:
			return bulletImage1
		default:
			return bulletImage2
		}
	}()

	return Projectile{
		pid:     pid,
		oid:     oid,
		image:   img,
		pos:     pos,
		dir:     dir,
		expires: time.Now().Add(bulletLifetime),
		speed:   speed,
		team:    teamGopher,
	}
}

func NewCrabBullet(oid, pid ecs.Id, pos Vec2, dir Vec2) Projectile {
	return Projectile{
		pid:     pid,
		oid:     oid,
		image:   crabBulletImage,
		pos:     pos,
		dir:     dir,
		speed:   bulletSpeed / 6,
		expires: time.Now().Add(bulletLifetime),
		team:    teamCrab,
	}
}

func SpawnBullets(center Vec2, ticker *time.Ticker, input *input, world *ecs.World) {
	select {
	case <-ticker.C:
		if input.fire {
			q := ecs.Query1[Gopher](world)

			var oid ecs.Id
			var pos Vec2
			q.MapId(func(id ecs.Id, a *Gopher) {
				oid = a.id
				pos = a.pos.Add(Vec2{20, 0})
			})

			dir := pos.Sub(input.cursor).Clamp(Vec2{-360, -360}, Vec2{360, 360})

			pid := world.NewId()
			world.Write(pid, ecs.C(NewGoBullet(oid, pid, bulletSpeed, dir, pos)), ecs.C(teamGopher))

			pid = world.NewId()
			world.Write(pid, ecs.C(NewGoBullet(oid, pid, bulletSpeed, dir.Add(Vec2{18, 18}), pos)), ecs.C(teamGopher))

			pid = world.NewId()
			world.Write(pid, ecs.C(NewGoBullet(oid, pid, bulletSpeed, dir.Sub(Vec2{18, 18}), pos)), ecs.C(teamGopher))
		}
	default:
		return
	}
}

func MoveBullets(world *ecs.World) {
	q := ecs.Query1[Projectile](world)

	q.MapId(func(id ecs.Id, a *Projectile) {
		a.pos.X -= a.dir.X * a.speed
		a.pos.Y -= a.dir.Y * a.speed
	})
}

func DrawBullets(screen *ebiten.Image, op *ebiten.DrawImageOptions, world *ecs.World) {
	q := ecs.Query1[Projectile](world)

	q.MapId(func(id ecs.Id, g *Projectile) {
		op.ColorScale.Reset()
		op.GeoM.Reset()
		op.GeoM.Scale(1.5, 1.5)
		op.GeoM.Translate(float64(g.pos.X), float64(g.pos.Y))
		screen.DrawImage(g.image, op)
	})
}

func ExpireBullets(world *ecs.World) {
	q := ecs.Query1[Projectile](world)

	q.MapId(func(id ecs.Id, a *Projectile) {
		if a.expires.Before(time.Now()) {
			ecs.Delete(world, id)
		}
	})
}
