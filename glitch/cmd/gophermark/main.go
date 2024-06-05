package main

// Try: https://www.shadertoy.com/view/csX3RH

import (
	"embed"
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/unitoftime/glitch"
	"github.com/unitoftime/glitch/shaders"
)

const (
	width  = 1512 * 4
	height = 982 * 4
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

//go:embed gopher.png
var f embed.FS

func loadImage(path string) (*image.NRGBA, error) {
	file, err := f.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	nrgba := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(nrgba, nrgba.Bounds(), img, bounds.Min, draw.Src)
	return nrgba, nil
}

func main() {
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

	glitch.Run(runGame)

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
}

func runGame() {
	win, err := glitch.NewWindow(width, height, "Glitch - Gophermark", glitch.WindowConfig{
		Vsync: true,
	})
	if err != nil {
		panic(err)
	}

	shader, err := glitch.NewShader(shaders.SpriteShader)
	if err != nil {
		panic(err)
	}

	pass := glitch.NewRenderPass(shader)
	pass.DepthTest = true

	manImage, err := loadImage("gopher.png")
	if err != nil {
		panic(err)
	}
	texture := glitch.NewTexture(manImage, false)
	manSprite := glitch.NewSprite(texture, texture.Bounds())

	length := 5000
	man := make([]Man, length)
	for i := range man {
		man[i] = NewMan()
	}

	// Text
	atlas, err := glitch.DefaultAtlas()
	if err != nil {
		panic(err)
	}

	text := atlas.Text("", 1)

	min := time.Duration(0)
	max := time.Duration(0)

	counter := 0
	camera := glitch.NewCameraOrtho()
	camera.DepthRange = glitch.Vec2{-127, 127}

	start := time.Now()
	var dt time.Duration

	geom := glitch.NewGeomDraw()
	geomRect := glitch.R(-16, -16, 16, 16)
	// geomMesh := glitch.NewQuadMesh(geomRect, geomRect)

	updatePositions := func() {
		for i := range man {
			man[i].position[0] += man[i].velocity[0]
			man[i].position[1] += man[i].velocity[1]

			if man[i].position[0] <= 0 || (man[i].position[0]+width/4) >= width {
				man[i].velocity[0] = -man[i].velocity[0]
				man[i].color.G = rand.Float64()
			}
			if man[i].position[1] <= 0 || (man[i].position[1]+height/4) >= height {
				man[i].velocity[1] = -man[i].velocity[1]
				man[i].color.B = rand.Float64()
			}
		}
	}

	mat := glitch.Mat4Ident
	for !win.Closed() {
		if win.Pressed(glitch.KeyEscape) {
			win.Close()
		}
		start = time.Now()

		updatePositions()

		pass.Clear()
		pass.SetLayer(0)

		camera.SetOrtho2D(win.Bounds())
		camera.SetView2D((-width/100)/2, (-height/100)/2, 1.0, 1.0)

		counter = (counter + 1) % 60
		if counter == 0 {
			text.Clear()
			text.Set(fmt.Sprintf("%2.2f (%2.2f, %2.2f) ms",
				1000*dt.Seconds(),
				1000*min.Seconds(),
				1000*max.Seconds()))
			min = 100000000000
			max = 0
		}
		text.DrawColorMask(pass, glitch.Mat4Ident, glitch.White)

		pass.SetLayer(1)
		for i := range man {
			mat = glitch.Mat4Ident
			mat.Scale(0.8, 0.8, 1.0).Translate(man[i].position[0], man[i].position[1], 0)
			manSprite.DrawColorMask(pass, mat, man[i].color)

			// geom.DrawRect(pass, geomRect, mat, man[i].color)
			// geomMesh.DrawColorMask(pass, mat, man[i].color)
			// geom.DrawRect(pass, geomRect, mat, man[i].color)
		}
		geom.FillRect(geomRect)
		// geomMesh.Draw(pass, glitch.Mat4Ident)

		glitch.Clear(win, glitch.RGBA{0.1, 0.2, 0.3, 1.0})

		pass.SetCamera2D(camera)
		pass.Draw(win)

		win.Update()

		dt = time.Since(start)
		if dt > max {
			max = dt
		}
		if dt < min {
			min = dt
		}
		// fmt.Println(dt.Seconds() * 1000)
	}
}

type Man struct {
	position, velocity glitch.Vec2
	color              glitch.RGBA
	layer              uint8
}

func NewMan() Man {
	colors := []glitch.RGBA{
		glitch.RGBA{R: 1, G: 1, B: 1, A: rand.ExpFloat64()},
		glitch.RGBA{R: 0, G: 0, B: rand.ExpFloat64(), A: rand.ExpFloat64()},
		glitch.RGBA{R: 0, G: rand.ExpFloat64(), B: rand.ExpFloat64(), A: rand.ExpFloat64()},
		glitch.RGBA{R: rand.ExpFloat64(), G: rand.ExpFloat64(), B: rand.ExpFloat64(), A: rand.ExpFloat64()},
	}
	randIndex := rand.Intn(len(colors))
	vScale := 1.0
	return Man{
		position: glitch.Vec2{width / 4, height / 4},
		velocity: glitch.Vec2{float64(2 * vScale * (rand.Float64() - 0.5)),
			float64(2 * vScale * (rand.Float64() - 0.5))},
		color: colors[randIndex],
		layer: uint8(randIndex) + 1,
	}
}
