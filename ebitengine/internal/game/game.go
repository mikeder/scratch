package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/unitoftime/ecs"
)

const (
	GameStateMenu = iota
	GameStatePlaying

	ScreenWidth  = 1024
	ScreenHeight = 768
)

type Game struct {
	count int
	keys  []ebiten.Key
	world *ecs.World
}

var _ ebiten.Game = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		world: ecs.NewWorld(),
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	// https://github.com/mikeder/larpa/blob/main/src/server/mod.rs#L276

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	g.count++
	g.count %= 200

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	cf := float32(g.count)
	vector.DrawFilledRect(screen, 10+cf, 10+cf, 100+cf, 100+cf, color.RGBA{0x90, 0x80, 0x80, 5}, true)
	vector.DrawFilledRect(screen, 20+cf, 10+cf, 100+cf, 100+cf, color.RGBA{0x00, 0x80, 0x80, 5}, true)
	vector.DrawFilledRect(screen, 30+cf, 10+cf, 100+cf, 100+cf, color.RGBA{0x90, 0x80, 0x80, 5}, true)
	vector.DrawFilledRect(screen, 40+cf, 10+cf, 100+cf, 100+cf, color.RGBA{0x20, 0x80, 0x80, 5}, true)
	vector.DrawFilledRect(screen, 50+cf, 50+cf, 100+cf, 100+cf, color.RGBA{0x90, 0x80, 0x80, 5}, true)

	PrintDebugText(screen, g.keys)
}
