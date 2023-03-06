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

type StartOptionsScene struct {
}

func (scene *StartOptionsScene) Type() string { return "Start" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (scene *StartOptionsScene) Preload() {
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

}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (scene *StartOptionsScene) Setup(updater engo.Updater) {
	common.SetBackground(color.Black)
	common.CameraBounds = engo.AABB{}
	w, _ := updater.(*ecs.World)
	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

	headerFont := &common.Font{
		URL:  "fonts/pixelfont.ttf",
		FG:   color.White,
		Size: 32,
	}
	headerFont.CreatePreloaded()
	startLabel := entities.Label{BasicEntity: ecs.NewBasic()}
	startLabel.RenderComponent.Drawable = common.Text{
		Font: headerFont,
		Text: "Start",
	}
	startLabel.RenderComponent.SetZIndex(1)
	startLabel.SpaceComponent.Position = engo.Point{X: 50, Y: 50}
	w.AddEntity(&startLabel)

	hostSessionSystem := &systems.SetupSessionSystem{NextScene: "Username", HostMode: true}
	w.AddSystem(hostSessionSystem)

	sgs, _ := common.LoadedSprite("textures/button.png")
	hostButton := entities.Button{BasicEntity: ecs.NewBasic()}
	hostButton.RenderComponent.Drawable = sgs
	hostButton.RenderComponent.SetZIndex(2)
	hostButton.SpaceComponent.Position = engo.Point{X: (engo.WindowWidth() / 2) - 50, Y: (engo.WindowHeight() / 2) + 50}
	hostButton.SpaceComponent.Width = hostButton.RenderComponent.Drawable.Width()
	hostButton.SpaceComponent.Height = hostButton.RenderComponent.Drawable.Height()

	buttonFont := &common.Font{
		URL:  "fonts/pixelfont.ttf",
		FG:   color.Black,
		Size: 10,
	}
	buttonFont.CreatePreloaded()
	hostLabel := entities.Label{BasicEntity: ecs.NewBasic()}
	hostLabel.RenderComponent.Drawable = common.Text{
		Font: buttonFont,
		Text: "Host",
	}
	hostLabel.RenderComponent.SetZIndex(3)
	hostLabel.SpaceComponent.Position = engo.Point{X: (hostButton.SpaceComponent.Position.X + 7), Y: hostButton.SpaceComponent.Position.Y + 2}
	w.AddEntity(&hostLabel)

	hostSessionSystem.Add(&hostButton.BasicEntity, &hostButton.SpaceComponent)

	joinSessionSystem := &systems.SetupSessionSystem{NextScene: "JoinUsername", HostMode: false}
	w.AddSystem(joinSessionSystem)

	joinButton := entities.Button{BasicEntity: ecs.NewBasic()}
	joinButton.RenderComponent.Drawable = sgs
	joinButton.RenderComponent.SetZIndex(2)
	joinButton.SpaceComponent.Position = engo.Point{X: (hostButton.SpaceComponent.Position.X + hostButton.SpaceComponent.Width + 25), Y: hostButton.SpaceComponent.Position.Y}
	joinButton.SpaceComponent.Width = joinButton.RenderComponent.Drawable.Width()
	joinButton.SpaceComponent.Height = joinButton.RenderComponent.Drawable.Height()

	joinLabel := entities.Label{BasicEntity: ecs.NewBasic()}
	joinLabel.RenderComponent.Drawable = common.Text{
		Font: buttonFont,
		Text: "Join",
	}
	joinLabel.RenderComponent.SetZIndex(3)
	joinLabel.SpaceComponent.Position = engo.Point{X: (joinButton.SpaceComponent.Position.X + 7), Y: joinButton.SpaceComponent.Position.Y + 2}
	w.AddEntity(&joinLabel)

	joinSessionSystem.Add(&joinButton.BasicEntity, &joinButton.SpaceComponent)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(
				&hostButton.BasicEntity,
				&hostButton.RenderComponent,
				&hostButton.SpaceComponent,
			)
			sys.Add(
				&joinButton.BasicEntity,
				&joinButton.RenderComponent,
				&joinButton.SpaceComponent,
			)
		}
	}

}
