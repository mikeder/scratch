package game

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/unitoftime/ecs"
)

var (
	crabImage *ebiten.Image
	ticker    *time.Ticker
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(CrabWalk_png))
	if err != nil {
		log.Fatal(err)
	}
	tmp := ebiten.NewImageFromImage(img)

	s := tmp.Bounds().Size()
	crabImage = ebiten.NewImage(s.X, s.Y)

	crabImage.DrawImage(tmp, nil)

	ticker = time.NewTicker(time.Millisecond * 20)
}

type Crab struct {
	id       ecs.Id
	image    *ebiten.Image
	pos      Vec2
	velocity Vec2
}

func NewCrab(id ecs.Id, pos Vec2) Crab {
	return Crab{
		id:       id,
		image:    crabImage,
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

		fmt.Printf("spawn crab: %+v\n", g)
	default:
		return
	}
}

func DespawnCrabs(world *ecs.World) {
	q := ecs.Query1[Crab](world)

	q.MapId(func(id ecs.Id, a *Crab) {
		if ok := ecs.Delete(world, id); ok {
			fmt.Print("deleted: ", id)
		}
	})
}

func DrawCrabs(screen *ebiten.Image, op *ebiten.DrawImageOptions, world *ecs.World) {
	q := ecs.Query1[Crab](world)

	q.MapId(func(id ecs.Id, g *Crab) {
		op.GeoM.Reset()
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

// fn random_position_around(pos: Vec2, min: f32, max: f32) -> (f32, f32) {
//     let mut rng = rand::thread_rng();
//     let angle = rng.gen_range(0.0..PI * 2.0);
//     let dist = rng.gen_range(min..max);

//     let offset_x = angle.cos() * dist;
//     let offset_y = angle.sin() * dist;

//     let random_x = pos.x + offset_x;
//     let random_y = pos.y + offset_y;

//     (random_x, random_y)
// }
