package game

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func PrintDebugText(screen *ebiten.Image, keys []ebiten.Key) {
	var s strings.Builder
	for i := range keys {
		s.WriteString(keys[i].String())
	}

	fps := ebiten.ActualFPS()
	tps := ebiten.ActualTPS()
	txt := fmt.Sprintf("Input: %s\nFPS: %0.2f\nTPS: %0.2f\n", s.String(), fps, tps)

	ebitenutil.DebugPrint(screen, txt)
}
