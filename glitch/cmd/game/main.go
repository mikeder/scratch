package main

import "github.com/mikeder/scratchygo/glitch/internal/game"

func main() {
	g, err := game.NewGame()
	if err != nil {
		panic(err)
	}
	g.Run()
}
