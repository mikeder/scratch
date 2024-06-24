package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	keyPress     = inpututil.IsKeyJustPressed
	keyRelease   = inpututil.IsKeyJustReleased
	mousePress   = inpututil.IsMouseButtonJustPressed
	mouseRelease = inpututil.IsMouseButtonJustReleased
)

func ReadInputs(input *input) {
	if keyPress(ebiten.KeyW) || keyPress(ebiten.KeyArrowUp) {
		input.up = true
	}
	if keyRelease(ebiten.KeyW) || keyRelease(ebiten.KeyArrowUp) {
		input.up = false
	}

	if keyPress(ebiten.KeyS) || keyPress(ebiten.KeyArrowDown) {
		input.down = true
	}
	if keyRelease(ebiten.KeyS) || keyRelease(ebiten.KeyArrowDown) {
		input.down = false
	}

	if keyPress(ebiten.KeyA) || keyPress(ebiten.KeyArrowLeft) {
		input.left = true
	}
	if keyRelease(ebiten.KeyA) || keyRelease(ebiten.KeyArrowLeft) {
		input.left = false
	}

	if keyPress(ebiten.KeyD) || keyPress(ebiten.KeyArrowRight) {
		input.right = true
	}
	if keyRelease(ebiten.KeyD) || keyRelease(ebiten.KeyArrowRight) {
		input.right = false
	}

	if keyPress(ebiten.KeySpace) || mousePress(ebiten.MouseButtonLeft) {
		input.fire = true
	}
	if keyRelease(ebiten.KeySpace) || mouseRelease(ebiten.MouseButtonLeft) {
		input.fire = false
	}

	// ENTER
	if keyPress(ebiten.KeyEnter) {
		input.enter = true
	}
	if keyRelease(ebiten.KeyEnter) {
		input.enter = false
	}

	// EXIT
	if keyPress(ebiten.KeyEscape) || keyPress(ebiten.KeyDelete) {
		input.exit = true
	}
	if keyRelease(ebiten.KeyEscape) || keyRelease(ebiten.KeyDelete) {
		input.exit = false
	}

	// if inpututil.AppendJustPressedTouchIDs()

	cx, cy := ebiten.CursorPosition()
	input.cursor = Vec2{float64(cx), float64(cy)}
}
