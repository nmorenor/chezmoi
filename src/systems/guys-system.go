package systems

import (
	"github.com/chezmoi/entities"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type GuysSystem struct {
	localGuy *entities.Guy
	world    *ecs.World
}

func (cb *GuysSystem) Remove(ecs.BasicEntity) {
	// do nothing
}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (cb *GuysSystem) Update(dt float32) {
}

func (cb *GuysSystem) New(world *ecs.World) {
	cb.world = world
	spriteSheet := common.NewSpritesheetFromFile("textures/run_horizontal_32x32_2.png", 32, 64)

	StopUpAction = &common.Animation{
		Name:   "upstop",
		Frames: []int{11},
	}

	StopDownAction = &common.Animation{
		Name:   "downstop",
		Frames: []int{23},
	}

	StopLeftAction = &common.Animation{
		Name:   "leftstop",
		Frames: []int{17},
	}

	StopRightAction = &common.Animation{
		Name:   "rightstop",
		Frames: []int{5},
	}

	WalkUpAction = &common.Animation{
		Name:   "up",
		Frames: []int{6, 7, 8, 9, 10},
		Loop:   true,
	}

	WalkDownAction = &common.Animation{
		Name:   "down",
		Frames: []int{18, 19, 20, 21, 22},
		Loop:   true,
	}

	WalkLeftAction = &common.Animation{
		Name:   "left",
		Frames: []int{12, 13, 14, 15, 16},
		Loop:   true,
	}

	WalkRightAction = &common.Animation{
		Name:   "right",
		Frames: []int{0, 1, 2, 3, 4},
		Loop:   true,
	}

	Actions = []*common.Animation{
		StopUpAction,
		StopDownAction,
		StopLeftAction,
		StopRightAction,
		WalkUpAction,
		WalkDownAction,
		WalkLeftAction,
		WalkRightAction,
	}

	engo.Input.RegisterButton(UpButton, engo.KeyW, engo.KeyArrowUp)
	engo.Input.RegisterButton(LeftButton, engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton(RightButton, engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton(DownButton, engo.KeyS, engo.KeyArrowDown)
	engo.Input.RegisterButton(PauseButton, engo.KeySpace)

	cb.localGuy = cb.CreateHero(
		world,
		engo.Point{X: engo.GameWidth() / 2, Y: (engo.GameHeight() / 2) + 128},
		spriteSheet,
	)
}

func (*GuysSystem) CreateHero(world *ecs.World, point engo.Point, spriteSheet *common.Spritesheet) *entities.Guy {
	hero := &entities.Guy{BasicEntity: ecs.NewBasic()}

	hero.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    float32(32),
		Height:   float32(64),
	}
	hero.RenderComponent = common.RenderComponent{
		Drawable: spriteSheet.Cell(0),
		Scale:    engo.Point{X: 1, Y: 1},
	}

	hero.SpeedComponent = entities.SpeedComponent{}
	hero.AnimationComponent = common.NewAnimationComponent(spriteSheet.Drawables(), 0.1)

	hero.AnimationComponent.AddAnimations(Actions)
	hero.AnimationComponent.SelectAnimationByName("downstop")

	hero.ControlComponent = entities.ControlComponent{
		SchemeHoriz: "horizontal",
		SchemeVert:  "vertical",
	}

	hero.RenderComponent.SetZIndex(6)

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(
				&hero.BasicEntity,
				&hero.RenderComponent,
				&hero.SpaceComponent,
			)

		case *common.AnimationSystem:
			sys.Add(
				&hero.BasicEntity,
				&hero.AnimationComponent,
				&hero.RenderComponent,
			)

		case *ControlSystem:
			sys.Add(
				&hero.BasicEntity,
				&hero.AnimationComponent,
				&hero.ControlComponent,
				&hero.SpaceComponent,
			)

		case *SpeedSystem:
			sys.Add(
				&hero.BasicEntity,
				&hero.SpeedComponent,
				&hero.SpaceComponent,
			)
		}
	}

	world.AddSystem(&common.EntityScroller{
		SpaceComponent: &hero.SpaceComponent,
		TrackingBounds: entities.LevelBounds,
	})

	return hero
}
