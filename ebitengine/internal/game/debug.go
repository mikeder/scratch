package game

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/unitoftime/ecs"
)

func PrintDebugText(screen *ebiten.Image, keys []ebiten.Key, world *ecs.World) {
	q := ecs.Query1[Gopher](world)

	var gophers int
	q.MapId(func(_ ecs.Id, _ *Gopher) {
		gophers++
	})

	var s strings.Builder
	for i := range keys {
		s.WriteString(keys[i].String())
	}

	fps := ebiten.ActualFPS()
	tps := ebiten.ActualTPS()

	txt := fmt.Sprintf(`
Gophers: %d
Input: %s
FPS: %0.2f
TPS: %0.2f
	`, gophers, s.String(), fps, tps)
	ebitenutil.DebugPrint(screen, txt)
}
