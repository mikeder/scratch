package game

import (
	_ "embed"
)

var (
	//go:embed crab1.png
	Crab1_png []byte
	//go:embed crab2.png
	Crab2_png []byte
	//go:embed crab3.png
	Crab3_png []byte

	//go:embed gopher.png
	Gopher_png []byte
)
