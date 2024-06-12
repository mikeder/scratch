package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/unitoftime/ecs"
)

const (
	GameStateMenu = iota
	GameStatePlaying

	ScreenWidth  = 1024
	ScreenHeight = 768
)

type Game struct {
	center vec2

	keys  []ebiten.Key
	op    *ebiten.DrawImageOptions
	world *ecs.World
}

var _ ebiten.Game = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		op:    new(ebiten.DrawImageOptions),
		world: ecs.NewWorld(),
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	for i := range g.keys {
		if g.keys[i] == ebiten.KeyArrowUp {
			SpawnGophers(g.center, g.world)
		}
		if g.keys[i] == ebiten.KeyArrowDown {
			DespawnGophers(g.world)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	b := screen.Bounds()
	x := b.Dx() / 2
	y := b.Dy() / 2
	g.center = vec2{x: float32(x), y: float32(y)}

	PrintDebugText(screen, g.keys, g.world)
	DrawGophers(screen, g.op, g.world)
}
