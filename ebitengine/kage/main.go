package main

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

var Resolutions struct {
	low  [2]int
	med  [2]int
	high [2]int
}

func init() {
	Resolutions.low = [2]int{960, 540}
	Resolutions.med = [2]int{1920, 1080}
	Resolutions.high = [2]int{3840, 2160}
}

//go:embed sun.kage
var sunShader []byte

//go:embed stars.kage
var starsShader []byte

//go:embed water.kage
var waterShader []byte

func main() {
	// compile the shader(s)
	sun, err := ebiten.NewShader(sunShader)
	if err != nil {
		panic(err)
	}
	stars, err := ebiten.NewShader(starsShader)
	if err != nil {
		panic(err)
	}
	water, err := ebiten.NewShader(waterShader)
	if err != nil {
		panic(err)
	}

	// create game struct
	res := Resolutions.low
	game := &Game{res: res, stars: stars, sun: sun, water: water}

	// configure window and run game
	ebiten.SetWindowTitle("Retro Sun")
	ebiten.SetWindowSize(res[0], res[1])
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err = ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	res   [2]int
	stars *ebiten.Shader
	sun   *ebiten.Shader
	water *ebiten.Shader
	time  int
}

// Assume a fixed layout.
func (g *Game) Layout(_, _ int) (int, int) {
	return g.res[0], g.res[1]
}

// No logic to update.
func (g *Game) Update() error {
	g.time++
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		g.res = Resolutions.low
	}
	if ebiten.IsKeyPressed(ebiten.KeyM) {
		g.res = Resolutions.med
	}
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		g.res = Resolutions.high
	}
	return nil
}

// Core drawing function from where we call DrawTrianglesShader.
func (g *Game) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]any{
		"Resolution": []float32{float32(w), float32(h)},
		"Time":       float32(g.time) / 240,
	}

	// draw stars
	screen.DrawRectShader(w, h, g.stars, op)

	// draw sun
	screen.DrawRectShader(w, h, g.sun, op)

	// capture current screen
	cap := ebiten.NewImageFromImage(screen)
	op.Images[0] = cap

	// draw water
	screen.DrawRectShader(w, h, g.water, op)
}
