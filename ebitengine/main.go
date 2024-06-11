package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mikeder/scratchygo/ebitengine/internal/game"
)

func main() {
	g := game.NewGame()

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Game (Ebitengine Demo)")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
