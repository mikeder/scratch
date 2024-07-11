package main

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed sun.kage
var sunShader []byte

//go:embed water.kage
var waterShader []byte

func main() {
	// compile the shader(s)
	sun, err := ebiten.NewShader(sunShader)
	if err != nil {
		panic(err)
	}
	water, err := ebiten.NewShader(waterShader)
	if err != nil {
		panic(err)
	}

	// create game struct
	game := &Game{sun: sun, water: water}

	// configure window and run game
	ebiten.SetWindowTitle("Retro Sun")
	ebiten.SetWindowSize(1024, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err = ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	sun   *ebiten.Shader
	water *ebiten.Shader
	time  int
}

// Assume a fixed layout.
func (g *Game) Layout(_, _ int) (int, int) {
	return 1024, 720
}

// No logic to update.
func (g *Game) Update() error {
	g.time++
	return nil
}

// Core drawing function from where we call DrawTrianglesShader.
func (g *Game) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]any{
		"Resolution": []float32{float32(w), float32(h)},
		"Time":       float32(g.time) / 120,
	}

	// draw sun
	screen.DrawRectShader(w, h, g.sun, op)

	// capture current screen
	cap := ebiten.NewImageFromImage(screen)
	op.Images[0] = cap

	// draw water
	screen.DrawRectShader(w, h, g.water, op)
}
