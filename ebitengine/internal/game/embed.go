package game

import (
	_ "embed"
)

var (
	//go:embed crab-walk.png
	CrabWalk_png []byte

	//go:embed gopher.png
	Gopher_png []byte
)
