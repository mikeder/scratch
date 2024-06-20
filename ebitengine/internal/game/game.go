package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kyroy/kdtree"
	"github.com/unitoftime/ecs"
)

const (
	GameStateMenu = iota
	GameStatePlaying
	GameStateOver

	ScreenWidth  = 1920
	ScreenHeight = 1080
)

type input struct {
	up     bool
	down   bool
	left   bool
	right  bool
	fire   bool
	cursor Vec2

	enter bool
	exit  bool
}

type window struct {
	width  int
	height int
}

type Game struct {
	center Vec2
	dt     time.Duration
	input  *input
	op     *ebiten.DrawImageOptions
	window window
	world  *ecs.World

	// ui ebi

	// state
	gameState   uint
	playerAdded bool
	tree        *kdtree.KDTree
}

var _ ebiten.Game = (*Game)(nil)

func NewGame() *Game {
	g := &Game{
		input:     new(input),
		gameState: GameStateMenu,
		op:        new(ebiten.DrawImageOptions),
		tree:      kdtree.New(nil),
		window:    window{ScreenWidth, ScreenHeight},
		world:     ecs.NewWorld(),
	}
	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	start := time.Now()
	if !g.playerAdded && g.center != Vec2Zero {
		gid := g.world.NewId()
		player := NewGopher(gid, g.center)
		g.world.Write(gid, ecs.C(player))
		g.playerAdded = true
	}

	ReadInputs(g.input)

	switch g.gameState {
	case GameStateMenu:
		// ShowMenu()
		if g.input.enter || g.input.fire {
			g.gameState = GameStatePlaying
		}
	case GameStatePlaying:
		SpawnCrabs(g.center, g.world)
		MoveGopher(g.input, g.world)
		SpawnBullets(g.center, g.input, g.world)
		MoveBullets(g.world)
		ExpireBullets(g.world)
		MoveCrabs(g.world)
		KillCrabs(g.tree, g.world)
		UpdateKDTree(g.tree, g.world)
		if g.input.exit {
			g.gameState = GameStateMenu
		}
	default:
		// do stuff
	}

	g.dt = time.Since(start)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	b := screen.Bounds()
	x := b.Dx() / 2
	y := b.Dy() / 2
	g.center = Vec2{X: float64(x), Y: float64(y)}

	switch g.gameState {
	case GameStateMenu:
		DrawWorld(screen, g.op)

		PrintDebugText(screen, g.input, g.world)
		// ShowMenu()
	case GameStatePlaying:
		DrawWorld(screen, g.op)

		PrintDebugText(screen, g.input, g.world)
		DrawCrabs(screen, g.op, g.world)
		DrawGopher(screen, g.op, g.world)
		DrawBullets(screen, g.op, g.world)

	default:
		// do stuff
	}
}
