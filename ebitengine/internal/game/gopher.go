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
	gopherImage *ebiten.Image
	ticker      *time.Ticker
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(Gopher_png))
	if err != nil {
		log.Fatal(err)
	}
	tmp := ebiten.NewImageFromImage(img)

	s := tmp.Bounds().Size()
	gopherImage = ebiten.NewImage(s.X, s.Y)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	gopherImage.DrawImage(tmp, op)

	ticker = time.NewTicker(time.Millisecond * 20)
}

type Gopher struct {
	id       ecs.Id
	image    *ebiten.Image
	pos      vec2
	velocity vec2
}

func NewGopher(id ecs.Id, pos vec2) Gopher {
	return Gopher{
		id:       id,
		image:    gopherImage,
		pos:      pos,
		velocity: vec2{rand.Float32(), rand.Float32()},
	}
}

type vec2 struct {
	x float32
	y float32
}

func SpawnGophers(center vec2, world *ecs.World) {
	select {
	case <-ticker.C:
		pos := randomPositionAround(center, 200, 400)
		id := world.NewId()
		g := NewGopher(id, pos)
		world.Write(id, ecs.C(g))

		fmt.Printf("spawn gopher: %+v", g)
	default:
		return
	}
}

func DespawnGophers(world *ecs.World) {
	q := ecs.Query1[Gopher](world)

	q.MapId(func(id ecs.Id, a *Gopher) {
		if ok := ecs.Delete(world, id); ok {
			fmt.Print("deleted: ", id)
		}
	})
}

func DrawGophers(screen *ebiten.Image, op *ebiten.DrawImageOptions, world *ecs.World) {
	q := ecs.Query1[Gopher](world)

	q.MapId(func(id ecs.Id, g *Gopher) {
		op.GeoM.Reset()
		op.GeoM.Translate(float64(g.pos.x), float64(g.pos.y))
		screen.DrawImage(g.image, op)
	})
}

func randomPositionAround(pos vec2, min, max float32) vec2 {
	angle := 0 + rand.Float64()*(math.Pi*2-0)
	dist := min + rand.Float32()*(max-min)
	offsetX := math.Cos(angle) * float64(dist)
	offsetY := math.Sin(angle) * float64(dist)
	fmt.Printf("offsetx: %0.2f, offsety: %0.2f", offsetX, offsetY)

	return vec2{pos.x + float32(offsetX), pos.y + float32(offsetY)}
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
