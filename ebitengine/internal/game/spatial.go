package game

import (
	"math"

	"github.com/hongshibao/go-kdtree"
)

type Vec2 struct {
	X float64
	Y float64
}

var Vec2Zero = Vec2{0, 0}

// implement kdtree.Point for Vec2
var _ kdtree.Point = (*Vec2)(nil)

// Return the total number of dimensions
func (v *Vec2) Dim() int {
	return 2
}

// Return the value X_{dim}, dim is started from 0
func (v *Vec2) GetValue(dim int) float64 {
	switch dim {
	case 0:
		return v.X
	case 1:
		return v.Y
	default:
		panic("unsupported dimension")
	}
}

// Return the distance between two points
func (v *Vec2) Distance(other kdtree.Point) float64 {
	// d = √((x2-x1)2 + (y2-y1)2)
	s0 := math.Pow(other.GetValue(0)-v.GetValue(0), 2)
	s1 := math.Pow(other.GetValue(1)-v.GetValue(1), 2)
	sum := s0 + s1
	ret := math.Sqrt(sum)
	return ret
}

// Return the distance between the point and the plane X_{dim}=val
func (v *Vec2) PlaneDistance(val float64, dim int) float64 {
	tmp := v.GetValue(dim) - val
	return tmp * tmp
}
