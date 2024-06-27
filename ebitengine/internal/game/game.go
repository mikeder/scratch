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
	healthTick          = 10 * time.Second
	treeUpdateTick      = 40 * time.Millisecond
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
	tickers     *Tickers
	tree        *kdtree.KDTree
	world       *ecs.World
}

type Tickers struct {
	bulletTicker       *time.Ticker
	currentBulletTickD time.Duration

	crabTicker       *time.Ticker
	currentCrabTickD time.Duration

	crabBulletTicker       *time.Ticker
	currentCrabBulletTickD time.Duration

	healthTicker *time.Ticker

	treeUpdateTicker *time.Ticker
}

func NewTickers() *Tickers {
	return &Tickers{
		bulletTicker:           time.NewTicker(bulletTickStart),
		currentBulletTickD:     bulletTickStart,
		crabTicker:             time.NewTicker(crabTickStart),
		currentCrabTickD:       crabTickStart,
		crabBulletTicker:       time.NewTicker(crabBulletTickStart),
		currentCrabBulletTickD: crabBulletTickStart,
		healthTicker:           time.NewTicker(healthTick),
		treeUpdateTicker:       time.NewTicker(treeUpdateTick),
	}
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

		tickers: NewTickers(),
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
		SpawnCrabs(g.tickers.crabTicker, g.world)
		SpawnHealth(g.tickers.healthTicker, g.world)
		MoveGopher(g.input, g.world)
		GopherPickupHealth(g.world)
		SpawnBullets(g.center, g.tickers.bulletTicker, g.input, g.world)
		MoveBullets(g.world)
		ExpireBullets(g.world)
		MoveCrabs(g.world)
		BulletHitsCrab(g.tree, g.world)
		BulletHitsGopher(g.world)
		KillCrabs(&g.score, g.world)
		DeleteCrabs(g.world)
		CrabShoots(g.tree, g.tickers.crabBulletTicker, g.world)
		KillGopher(&g.gameState, g.tree, g.world)
		UpdateKDTree(g.tree, g.tickers.treeUpdateTicker, g.world)
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
		DrawHealth(screen, g.op, g.world)
		DrawGopher(screen, g.op, g.world)
		DrawCrabs(screen, g.op, g.world)
		DrawBullets(screen, g.op, g.world)
		DrawHealthText(screen, g.world)
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
	q1 := ecs.Query1[Crab](g.world)
	q1.MapId(func(id ecs.Id, c *Crab) {
		ecs.Delete(g.world, id)
	})

	q2 := ecs.Query1[Gopher](g.world)
	q2.MapId(func(id ecs.Id, c *Gopher) {
		ecs.Delete(g.world, id)
	})

	q3 := ecs.Query1[HealthPickup](g.world)
	q3.MapId(func(id ecs.Id, p *HealthPickup) {
		ecs.Delete(g.world, id)
	})

	q4 := ecs.Query1[Projectile](g.world)
	q4.MapId(func(id ecs.Id, p *Projectile) {
		ecs.Delete(g.world, id)
	})

	g.playerAdded = false
	g.nextWave = 100
	g.score = 0
	g.tickers = NewTickers()
	g.waveNum = 1
}

func (g *Game) UpdateWave() {
	current := g.nextWave
	if g.score >= current {
		g.nextWave = current * 2
		g.waveNum += 1

		// bullet ticker
		qd := g.tickers.currentBulletTickD.Nanoseconds() / 4
		bd := time.Duration(g.tickers.currentBulletTickD.Nanoseconds() - qd)
		if bd <= bulletTickEnd {
			bd = bulletTickEnd
		}
		g.tickers.currentBulletTickD = bd
		g.tickers.bulletTicker.Reset(bd)

		// crab ticker
		cd := time.Duration(g.tickers.currentCrabTickD.Nanoseconds() / 2)
		if cd < crabTickEnd {
			cd = crabTickStart
		}
		g.tickers.currentCrabTickD = cd
		g.tickers.crabTicker.Reset(cd)

		// crab bullet ticker
		cbd := time.Duration(g.tickers.currentCrabBulletTickD.Nanoseconds() / 2)
		if cbd < crabBulletTickEnd {
			cbd = crabBulletTickStart
		}
		g.tickers.currentCrabBulletTickD = cbd
		g.tickers.crabBulletTicker.Reset(cbd)
	}
}
