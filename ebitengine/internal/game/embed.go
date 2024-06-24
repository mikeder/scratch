package game

import (
	_ "embed"
)

var (
	//go:embed assets/crab1.png
	Crab1_png []byte
	//go:embed assets/crab2.png
	Crab2_png []byte
	//go:embed assets/crab3.png
	Crab3_png []byte
	//go:embed assets/gopher.png
	Gopher_png []byte
	//go:embed assets/gobullet1.png
	GoBullet1_png []byte
	//go:embed assets/gobullet2.png
	GoBullet2_png []byte
)
