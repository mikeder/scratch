package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

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

//go:embed sun.kage.go
var sunShader []byte

//go:embed stars.kage.go
var starsShader []byte

//go:embed water.kage.go
var waterShader []byte

func main() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		go func() {
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
		}()
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	// compile the shader(s)
	sun, err := ebiten.NewShader(sunShader)
	if err != nil {
		return err
	}
	stars, err := ebiten.NewShader(starsShader)
	if err != nil {
		return err
	}
	water, err := ebiten.NewShader(waterShader)
	if err != nil {
		return err
	}

	// create game struct
	res := Resolutions.high
	offscreen := ebiten.NewImage(res[0], res[1])
	game := &Game{offscreen: offscreen, res: res, stars: stars, sun: sun, water: water}

	// configure window and run game
	ebiten.SetWindowTitle("Retro Sun")
	ebiten.SetWindowSize(res[0], res[1])
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err = ebiten.RunGame(game)
	if err != nil {
		return err
	}
	return nil
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	res       [2]int
	offscreen *ebiten.Image
	stars     *ebiten.Shader
	sun       *ebiten.Shader
	water     *ebiten.Shader
	time      int
}

// Assume a fixed layout.
func (g *Game) Layout(_, _ int) (int, int) {
	return g.res[0], g.res[1]
}

// No logic to update.
func (g *Game) Update() error {
	g.time++
	return nil
}

// Core drawing function from where we call DrawTrianglesShader.
func (g *Game) Draw(screen *ebiten.Image) {
	w, h := g.offscreen.Bounds().Dx(), g.offscreen.Bounds().Dy()

	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]any{
		"Resolution": []float32{float32(w), float32(h)},
		"Time":       float32(g.time) / 240,
	}

	// draw stars
	g.offscreen.DrawRectShader(w, h, g.stars, op)

	// draw sun
	g.offscreen.DrawRectShader(w, h, g.sun, op)

	// capture current screen
	op.Images[0] = g.offscreen

	// draw water
	screen.DrawRectShader(w, h, g.water, op)
}
