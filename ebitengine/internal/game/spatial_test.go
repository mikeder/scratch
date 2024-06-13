package game_test

import (
	"testing"

	"github.com/hongshibao/go-kdtree"
	"github.com/mikeder/scratchygo/ebitengine/internal/game"
)

func TestVec2KDTree(t *testing.T) {
	p0 := &game.Vec2{0, 0}
	p1 := &game.Vec2{1, 1}
	p2 := &game.Vec2{2, 2}
	p3 := &game.Vec2{3, 3}

	points := []kdtree.Point{p0, p1, p2, p3}

	tree := kdtree.NewKDTree(points)

	p4 := &game.Vec2{7, 7}
	nn := tree.KNN(p4, 1)

	if len(nn) != 1 {
		t.Logf("should have 1 neighbor\n")
		t.Fail()
	}
	got := nn[0].Distance(p4)
	want := 5.656854249492381
	if got != want {
		t.Logf("distance - got: %f, want: %f\n", got, want)
		t.Fail()
	}
}
