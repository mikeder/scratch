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
	q3 := ecs.Query1[Projectile](world)

	var bullets int
	q3.MapId(func(_ ecs.Id, _ *Projectile) {
		bullets++
	})

	var crabs int
	q1.MapId(func(_ ecs.Id, _ *Crab) {
		crabs++
	})

	var gophers int
	var dir Vec2
	q2.MapId(func(_ ecs.Id, a *Gopher) {
		dir = input.cursor.Sub(a.pos)
		gophers++
	})

	fps := ebiten.ActualFPS()
	tps := ebiten.ActualTPS()

	txt := fmt.Sprintf(`
Bullets: %d
Crabs: %d
Gophers: %d
Cursor: %v
Dir: %v
FPS: %0.2f
TPS: %0.2f
	`, bullets, crabs, gophers, input.cursor, dir, fps, tps)
	ebitenutil.DebugPrint(screen, txt)
}
