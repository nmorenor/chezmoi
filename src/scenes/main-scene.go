package scenes

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/ByteArena/box2d"
	"github.com/chezmoi/assets"
	"github.com/chezmoi/systems"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"golang.org/x/image/font/gofont/gosmallcaps"

	"github.com/Noofbiz/engoBox2dSystem"
)

type MainScene struct {
}

func (scene *MainScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (scene *MainScene) Preload() {
	files := []string{
		"levels/Interiors_free_32x32.png",
		"levels/Room_Builder_free_32x32.png",
		"levels/interiors.tmx",
		"textures/run_horizontal_32x32_2.png",
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
func (scene *MainScene) Setup(updater engo.Updater) {
	common.SetBackground(color.Black)

	world, _ := updater.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.AnimationSystem{})
	world.AddSystem(&common.MouseSystem{})

	kbs := common.NewKeyboardScroller(
		400, engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis)
	world.AddSystem(kbs)
	world.AddSystem(&common.MouseZoomer{
		ZoomSpeed: -0.125,
	})

	world.AddSystem(&engoBox2dSystem.PhysicsSystem{VelocityIterations: 1, PositionIterations: 1})
	world.AddSystem(&engoBox2dSystem.CollisionSystem{})

	world.AddSystem(&systems.TilesSystem{})
	world.AddSystem(&systems.ControlSystem{})
	world.AddSystem(&systems.GuysSystem{})

	engoBox2dSystem.World.SetGravity(box2d.B2Vec2{X: 0, Y: 0})

	engo.Input.RegisterAxis(
		"vertical",
		engo.AxisKeyPair{Min: engo.KeyArrowUp, Max: engo.KeyArrowDown},
		engo.AxisKeyPair{Min: engo.KeyW, Max: engo.KeyS},
	)

	engo.Input.RegisterAxis(
		"horizontal",
		engo.AxisKeyPair{Min: engo.KeyArrowLeft, Max: engo.KeyArrowRight},
		engo.AxisKeyPair{Min: engo.KeyA, Max: engo.KeyD},
	)
}

func (scene *MainScene) Start() {
	opts := engo.RunOptions{
		Title:          "Chez moi",
		Width:          800,
		Height:         600,
		StandardInputs: true,
	}
	engo.Run(opts, scene)
}

var ChezMoiInstance = chezMoi()

func chezMoi() *MainScene {
	scene := &MainScene{}
	return scene
}
