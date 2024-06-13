package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/unitoftime/ecs"
)

func PrintDebugText(screen *ebiten.Image, input *input, world *ecs.World) {
	q1 := ecs.Query1[Crab](world)
	q2 := ecs.Query1[Gopher](world)

	var crabs int
	q1.MapId(func(_ ecs.Id, _ *Crab) {
		crabs++
	})

	var gophers int
	q2.MapId(func(_ ecs.Id, _ *Gopher) {
		gophers++
	})

	fps := ebiten.ActualFPS()
	tps := ebiten.ActualTPS()

	txt := fmt.Sprintf(`
Crabs: %d
Gophers: %d
Input: %v
FPS: %0.2f
TPS: %0.2f
	`, crabs, gophers, input, fps, tps)
	ebitenutil.DebugPrint(screen, txt)
}
