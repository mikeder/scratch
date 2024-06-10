package render

import (
	"fmt"
	"image/color"
	"time"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/glitch"
	"github.com/unitoftime/glitch/ui"
)

func MenuSystem(win *glitch.Window, camera *glitch.CameraOrtho, atlas *glitch.Atlas, pass *glitch.RenderPass) ecs.System {
	fmt.Println("MenuSystem")

	group := ui.NewGroup(win, camera, atlas, pass)
	group.SetLayer(1)
	// group.Debug = true

	menuRect := win.Bounds().SliceHorizontal(800).SliceVertical(600)
	rect1 := glitch.NewGeomDraw().Rectangle(menuRect, 2)
	rect1.SetColor(glitch.RGBA{R: 0, G: 1, B: 1, A: 1})

	// playButtonRect := menuRect.CutBottom(50)

	normal := NewImage(200, 100, color.RGBA{1, 0, 1, 1})
	texture := glitch.NewTexture(normal, false)
	buttonSprite := glitch.NewSprite(texture, texture.Bounds())

	hover := NewImage(200, 100, color.RGBA{0, 1, 1, 0})
	texture2 := glitch.NewTexture(hover, false)
	buttonSprite2 := glitch.NewSprite(texture2, texture2.Bounds())

	pressed := NewImage(200, 100, color.RGBA{0, 1, 1, 1})
	texture3 := glitch.NewTexture(pressed, false)
	buttonSprite3 := glitch.NewSprite(texture3, texture3.Bounds())

	style := ui.Style{
		Normal:  ui.NewSpriteStyle(buttonSprite, glitch.White),
		Hovered: ui.NewSpriteStyle(buttonSprite2, glitch.Black),
		Pressed: ui.NewSpriteStyle(buttonSprite3, glitch.Greyscale(1)),
		Text:    ui.NewTextStyle(),
	}
	style.Text = style.Text.Scale(0.6)

	return ecs.NewSystem(func(dt time.Duration) {
		ui.Clear()
		group.Clear()

		group.Text("Play Game", menuRect, ui.NewTextStyle())
		// group.Button("Play", playButtonRect, style)

		rect1.Draw(pass, glitch.Mat4Ident)
	})
}
