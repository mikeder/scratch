package render

import (
	"fmt"
	"time"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/glitch"
)

// Clear the current frame, reset the render pass layer and clear the window.
func ClearSystem(win *glitch.Window, camera *glitch.CameraOrtho, pass *glitch.RenderPass) ecs.System {
	fmt.Println("ClearSystem")
	return ecs.NewSystem(func(dt time.Duration) {
		pass.Clear()
		pass.SetLayer(0)
		glitch.Clear(win, glitch.Black)
	})
}

// Draw the render pass and update the window.
func DrawSystem(win *glitch.Window, camera *glitch.CameraOrtho, pass *glitch.RenderPass) ecs.System {
	fmt.Println("DrawSystem")
	pass.SetCamera2D(camera)

	return ecs.NewSystem(func(dt time.Duration) {
		pass.Draw(win)
		win.Update()
	})
}
