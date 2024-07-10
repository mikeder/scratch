package main

import (
	_ "embed"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shader.kage
var shaderProgram []byte

func main() {
	// compile the shader
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil {
		panic(err)
	}

	// create game struct
	game := &Game{shader: shader, start: time.Now()}

	// configure window and run game
	ebiten.SetWindowTitle("Retro Sun")
	ebiten.SetWindowSize(1024, 720)
	err = ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	shader   *ebiten.Shader
	vertices [4]ebiten.Vertex
	degrees  int
	start    time.Time
}

// Assume a fixed layout.
func (g *Game) Layout(_, _ int) (int, int) {
	return 1024, 720
}

// No logic to update.
func (g *Game) Update() error {
	g.degrees += 1
	if g.degrees >= 360 {
		g.degrees = 0
	}
	return nil
}

// Core drawing function from where we call DrawTrianglesShader.
func (g *Game) Draw(screen *ebiten.Image) {
	// map the vertices to the target image
	bounds := screen.Bounds()
	g.vertices[0].DstX = float32(bounds.Min.X) // top-left
	g.vertices[0].DstY = float32(bounds.Min.Y) // top-left
	g.vertices[1].DstX = float32(bounds.Max.X) // top-right
	g.vertices[1].DstY = float32(bounds.Min.Y) // top-right
	g.vertices[2].DstX = float32(bounds.Min.X) // bottom-left
	g.vertices[2].DstY = float32(bounds.Max.Y) // bottom-left
	g.vertices[3].DstX = float32(bounds.Max.X) // bottom-right
	g.vertices[3].DstY = float32(bounds.Max.Y) // bottom-right
	// [VERTEX-NOTE]
	// Other properties will be set on later examples. The full
	// configuration is quite verbose, but you will typically create
	// your own helper functions to do the heavy lifting, and in
	// some cases you can optimize and omit some settings on
	// successive passes.

	// triangle shader options
	var shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderOpts.Uniforms = make(map[string]interface{})
	shaderOpts.Uniforms["Center"] = []float32{
		float32(screen.Bounds().Dx()) / 2,
		float32(screen.Bounds().Dy()) / 2,
	}
	shaderOpts.Uniforms["Time"] = float32(time.Since(g.start).Seconds() / 2.5)
	shaderOpts.Uniforms["Resolution"] = []float32{float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy())}
	shaderOpts.Uniforms["Radius"] = float32(80 + 30*math.Sin(float64(g.degrees)*math.Pi/180.0))

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(g.vertices[:], indices, g.shader, &shaderOpts)
}
