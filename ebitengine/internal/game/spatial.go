package game

import (
	"math"
	"sync"
	"time"

	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
	"github.com/unitoftime/ecs"
)

type Vec2 struct {
	X float64
	Y float64
}

var Vec2Zero = Vec2{0, 0}

func (v Vec2) Add(rhs Vec2) Vec2 {
	return Vec2{X: v.X + rhs.X, Y: v.Y + rhs.Y}
}

func (v Vec2) Sub(rhs Vec2) Vec2 {
	return Vec2{X: v.X - rhs.X, Y: v.Y - rhs.Y}
}

func (v Vec2) Min(rhs Vec2) Vec2 {
	return Vec2{X: math.Min(v.X, rhs.X), Y: math.Min(v.Y, rhs.Y)}
}

func (v Vec2) Max(rhs Vec2) Vec2 {
	return Vec2{X: math.Max(v.X, rhs.X), Y: math.Max(v.Y, rhs.Y)}
}

func (v Vec2) Clamp(min, max Vec2) Vec2 {
	if min.X > max.X || min.Y > max.Y {
		panic("Clamp: expected min <= max")
	}
	return v.Max(min).Min(max)
}

// Return the distance between two points
func (v Vec2) Distance(other Vec2) float64 {
	// d = √((x2-x1)2 + (y2-y1)2)
	s0 := math.Pow(other.X-v.X, 2)
	s1 := math.Pow(other.Y-v.Y, 2)
	sum := s0 + s1
	d := math.Sqrt(sum)
	return d
}

var updateTicker = time.NewTicker(time.Millisecond * 20)

func UpdateKDTree(mut *sync.RWMutex, tree *kdtree.KDTree, world *ecs.World) {
	select {
	case <-updateTicker.C:
		q := ecs.Query1[Crab](world)

		var kp []kdtree.Point
		q.MapId(func(id ecs.Id, a *Crab) {
			kp = append(kp, points.NewPoint([]float64{a.pos.X, a.pos.Y}, a))
		})
		mut.Lock()
		*tree = *kdtree.New(kp)
		mut.Unlock()
	default:
		return
	}
}
