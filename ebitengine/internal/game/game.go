package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kyroy/kdtree"
	"github.com/unitoftime/ecs"
)

type GameState int

const (
	GameStateMenu GameState = iota
	GameStatePlaying
	GameStateOver

	ScreenWidth  = 1920
	ScreenHeight = 1080

	bulletTickStart     = 200 * time.Millisecond
	bulletTickEnd       = 20 * time.Millisecond
	crabTickStart       = 500 * time.Millisecond
	crabTickEnd         = 50 * time.Millisecond
	crabBulletTickStart = 5 * time.Second
	crabBulletTickEnd   = 50 * time.Millisecond
)

type input struct {
	up     bool
	down   bool
	left   bool
	right  bool
	fire   bool
	cursor Vec2

	enter bool
	reset bool
	exit  bool
}

type window struct {
	width  int
	height int
}

type Game struct {
	// window
	center Vec2
	// dt     time.Duration
	input  *input
	op     *ebiten.DrawImageOptions
	window window

	// state
	gameState   GameState
	playerAdded bool
	score       uint
	nextWave    uint
	waveNum     uint
	tree        *kdtree.KDTree
	world       *ecs.World

	// tickers
	bulletTicker           *time.Ticker
	currentBulletTickD     time.Duration
	crabTicker             *time.Ticker
	currentCrabTickD       time.Duration
	crabBulletTicker       *time.Ticker
	currentCrabBulletTickD time.Duration
	treeTicker             *time.Ticker
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

		score:    0,
		nextWave: 100,
		waveNum:  1,

		bulletTicker:           time.NewTicker(bulletTickStart),
		currentBulletTickD:     bulletTickStart,
		crabTicker:             time.NewTicker(crabTickStart),
		currentCrabTickD:       crabTickStart,
		crabBulletTicker:       time.NewTicker(crabBulletTickStart),
		currentCrabBulletTickD: crabBulletTickStart,
		treeTicker:             time.NewTicker(time.Millisecond * 40),
	}
	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	if !g.playerAdded && g.center != Vec2Zero {
		gid := g.world.NewId()
		player := NewGopher(gid, g.center)
		g.world.Write(gid, ecs.C(player))
		g.playerAdded = true
		time.Sleep(time.Millisecond * 200)
	}

	ReadInputs(g.input)

	switch g.gameState {
	case GameStateMenu:
		g.Reset()
		if g.input.enter || g.input.fire {
			g.gameState = GameStatePlaying
		}
	case GameStatePlaying:
		SpawnCrabs(g.crabTicker, g.world)
		MoveGopher(g.input, g.world)
		SpawnBullets(g.center, g.bulletTicker, g.input, g.world)
		MoveBullets(g.world)
		ExpireBullets(g.world)
		MoveCrabs(g.world)
		BulletHitsCrab(g.tree, g.world)
		BulletHitsGopher(g.world)
		KillCrabs(&g.score, g.world)
		DeleteCrabs(g.world)
		CrabShoots(g.tree, g.crabBulletTicker, g.world)
		KillGopher(&g.gameState, g.tree, g.world)
		UpdateKDTree(g.tree, g.treeTicker, g.world)
		g.UpdateWave()
		if g.input.exit {
			g.gameState = GameStateMenu
		}
	case GameStateOver:
		if g.input.reset || g.input.exit {
			g.gameState = GameStateMenu
		}
	default:
		// do stuff
	}

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
		StartMenu(screen)
	case GameStatePlaying:
		DrawWorld(screen, g.op)
		DrawGopher(screen, g.op, g.world)
		DrawCrabs(screen, g.op, g.world)
		DrawBullets(screen, g.op, g.world)
		PlayMenu(g.score, g.waveNum, g.world, screen)
	case GameStateOver:
		DrawWorld(screen, g.op)
		DrawGopher(screen, g.op, g.world)
		DrawCrabs(screen, g.op, g.world)
		DrawBullets(screen, g.op, g.world)
		OverMenu(g.score, g.waveNum, screen)
	default:
		// do stuff
	}

	PrintDebugText(screen, g.input, g.world)
}

func (g *Game) Reset() {
	q1 := ecs.Query1[Projectile](g.world)
	q2 := ecs.Query1[Crab](g.world)
	q3 := ecs.Query1[Gopher](g.world)

	q1.MapId(func(id ecs.Id, b *Projectile) {
		ecs.Delete(g.world, id)
	})
	q2.MapId(func(id ecs.Id, c *Crab) {
		ecs.Delete(g.world, id)
	})
	q3.MapId(func(id ecs.Id, p *Gopher) {
		ecs.Delete(g.world, id)
	})

	g.playerAdded = false
	g.waveNum = 1
	g.nextWave = 100
	g.bulletTicker.Reset(bulletTickStart)
	g.currentBulletTickD = bulletTickStart

	g.score = 0
	g.crabTicker.Reset(crabTickStart)
	g.currentCrabTickD = crabTickStart

	g.crabBulletTicker.Reset(crabBulletTickStart)
	g.currentCrabBulletTickD = crabBulletTickStart
}

func (g *Game) UpdateWave() {
	current := g.nextWave
	if g.score >= current {
		g.nextWave = current * 2
		g.waveNum += 1

		// bullet ticker
		qd := g.currentBulletTickD.Nanoseconds() / 4
		bd := time.Duration(g.currentBulletTickD.Nanoseconds() - qd)
		if bd <= bulletTickEnd {
			bd = bulletTickEnd
		}
		g.currentBulletTickD = bd
		g.bulletTicker.Reset(bd)

		// crab ticker
		cd := time.Duration(g.currentCrabTickD.Nanoseconds() / 2)
		if cd < crabTickEnd {
			cd = crabTickStart
		}
		g.currentCrabTickD = cd
		g.crabTicker.Reset(cd)

		// crab bullet ticker
		cbd := time.Duration(g.currentCrabBulletTickD.Nanoseconds() / 2)
		if cbd < crabBulletTickEnd {
			cbd = crabBulletTickStart
		}
		g.currentCrabBulletTickD = cbd
		g.crabBulletTicker.Reset(cbd)
	}
}
