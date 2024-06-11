package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func PrintDebugText(screen *ebiten.Image) {
	fps := ebiten.ActualFPS()
	tps := ebiten.ActualTPS()
	txt := fmt.Sprintf("FPS: %0.2f\nTPS: %0.2f\n", fps, tps)

	ebitenutil.DebugPrint(screen, txt)
}
