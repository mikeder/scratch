package game

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hongshibao/go-kdtree"
	"github.com/unitoftime/ecs"
)

const (
	GameStateMenu = iota
	GameStatePlaying

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
}

type Game struct {
	center Vec2
	dt     time.Duration
	input  *input
	op     *ebiten.DrawImageOptions
	world  *ecs.World

	// state
	playerAdded bool
	treeMux     *sync.RWMutex
	tree        *kdtree.KDTree
}

var _ ebiten.Game = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		input:   new(input),
		op:      new(ebiten.DrawImageOptions),
		treeMux: &sync.RWMutex{},
		tree:    kdtree.NewKDTree(nil),
		world:   ecs.NewWorld(),
	}
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
	SpawnCrabs(g.center, g.world)
	MoveGopher(g.input, g.world)
	SpawnBullets(g.center, g.input, g.world)
	MoveBullets(g.world)
	ExpireBullets(g.world)
	// UpdateKDTree(g.treeMux, g.tree, g.world)
	MoveCrabs(g.world)
	KillCrabs(g.world)

	g.dt = time.Since(start)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	b := screen.Bounds()
	x := b.Dx() / 2
	y := b.Dy() / 2
	g.center = Vec2{X: float64(x), Y: float64(y)}

	PrintDebugText(screen, g.input, g.world)
	DrawCrabs(screen, g.op, g.world)
	DrawGopher(screen, g.op, g.world)
	DrawBullets(screen, g.op, g.world)
}
