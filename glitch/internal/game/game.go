package game

import (
	"time"

	"github.com/mikeder/scratchygo/glitch/internal/game/render"
	"github.com/unitoftime/ecs"
	"github.com/unitoftime/glitch"
	"github.com/unitoftime/glitch/shaders"
)

const FixedTimeStep time.Duration = 16 * time.Millisecond

const (
	GameStateMenu = iota
	GameStatePlaying
)

type Game struct {
	camera    *glitch.CameraOrtho
	shader    *glitch.Shader
	scheduler *ecs.Scheduler
	textAtlas *glitch.Atlas
	window    *glitch.Window
	world     *ecs.World
}

func NewGame() (*Game, error) {
	// setup window
	window, err := glitch.NewWindow(0, 0, "Game", glitch.WindowConfig{
		Vsync: true,
	})
	if err != nil {
		return nil, err
	}

	// setup shader
	shader, err := glitch.NewShader(shaders.SpriteShader)
	if err != nil {
		return nil, err
	}

	// setup scheduler
	scheduler := ecs.NewScheduler()
	scheduler.SetFixedTimeStep(FixedTimeStep)

	// setup camera
	camera := glitch.NewCameraOrtho()
	camera.DepthRange = glitch.Vec2{-1, 127}
	camera.SetOrtho2D(window.Bounds())
	camera.SetView2D(0, 0, .95, .95)

	// setup world
	world := ecs.NewWorld()

	// setup text atlas
	textAtlas, err := glitch.DefaultAtlas()
	if err != nil {
		return nil, err
	}

	return &Game{
		camera,
		shader,
		scheduler,
		textAtlas,
		window,
		world,
	}, nil
}

func (g *Game) Run() {
	pass := glitch.NewRenderPass(g.shader)

	input := []ecs.System{}
	physics := []ecs.System{}
	render := []ecs.System{
		render.ClearSystem(g.window, g.camera, pass),
		render.MenuSystem(g.window, g.camera, g.textAtlas, pass),
		render.DrawSystem(g.window, g.camera, pass),
	}

	g.scheduler.AppendInput(input...)
	g.scheduler.AppendPhysics(physics...)
	g.scheduler.AppendRender(render...)

	g.scheduler.Run()
}
