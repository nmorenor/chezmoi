package scenes

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/nmorenor/chezmoi/assets"
	"github.com/nmorenor/chezmoi/entities"
	"github.com/nmorenor/chezmoi/systems"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"golang.org/x/image/font/gofont/gosmallcaps"
)

type SetUserNameScene struct {
}

func (scene *SetUserNameScene) Type() string { return "Username" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (scene *SetUserNameScene) Preload() {
	files := []string{
		"textures/button.png",
		"fonts/pixelfont.ttf",
	}
	for _, file := range files {
		data, err := assets.Asset(file)
		if err != nil {
			fmt.Printf("Unable to locate asset with URL: %v\n", file)
		}
		err = engo.Files.LoadReaderData(file, bytes.NewReader(data))
		if err != nil {
			fmt.Printf("Unable to load asset with URL: %v\n", file)
		}
	}
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
	engo.Input.RegisterButton("backspace", engo.KeyBackspace)
	engo.Input.RegisterButton("paste", engo.KeyF1)
	engo.Input.RegisterButton("enter", engo.KeyEnter)
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (scene *SetUserNameScene) Setup(updater engo.Updater) {
	common.SetBackground(color.Black)

	w, _ := updater.(*ecs.World)
	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

	tfnt := &common.Font{
		URL:  "fonts/pixelfont.ttf",
		FG:   color.White,
		Size: 24,
	}
	tfnt.CreatePreloaded()
	t := entities.Label{BasicEntity: ecs.NewBasic()}
	t.RenderComponent.Drawable = common.Text{
		Font: tfnt,
		Text: "Username:",
	}
	t.RenderComponent.SetZIndex(1)
	t.SpaceComponent.Position = engo.Point{X: (engo.WindowWidth() / 2) - 70, Y: (engo.WindowHeight() / 2) - 50}
	w.AddEntity(&t)

	w.AddSystem(&systems.TypingSystem{})
}
