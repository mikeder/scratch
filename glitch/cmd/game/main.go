package main

import (
	"github.com/unitoftime/ecs"
	"github.com/unitoftime/glitch"
	"github.com/unitoftime/glitch/shaders"
)

type Position glitch.Vec2
type Velocity glitch.Vec2
type Collider struct {
	Radius float64
}

type PlayerBundle struct {
	Collider
	Position
	Velocity
}

func main() {
	glitch.Run(game)
}

func game() {
	world := ecs.NewWorld()

	var players []PlayerBundle
	for range 5 {
		players = append(players, PlayerBundle{Position: Position{1, 1}, Velocity: Velocity{1, 0}, Collider: Collider{16}})
	}

	for i := range players {
		id := world.NewId()
		ecs.Write(world, id, ecs.C(players[i]))
	}

	win, err := glitch.NewWindow(0, 0, "Game", glitch.WindowConfig{
		Fullscreen: true,
		Vsync:      true,
	})
	if err != nil {
		panic(err)
	}

	shader, err := glitch.NewShader(shaders.SpriteShader)
	if err != nil {
		panic(err)
	}

	pass := glitch.NewRenderPass(shader)

	camera := glitch.NewCameraOrtho()
	camera.DepthRange = glitch.Vec2{-1, 127}

	camera.SetOrtho2D(win.Bounds())
	camera.SetView2D(0, 0, .95, .95)

	for !win.Closed() {
		if win.Pressed(glitch.KeyEscape) {
			win.Close()
		}

		pass.Clear()
		pass.SetLayer(0)

		glitch.Clear(win, glitch.RGBA{R: 0.1, G: 0.1, B: 0.1, A: 0.8})

		pass.SetCamera2D(camera)
		pass.Draw(win)
		win.Update()
	}

}
