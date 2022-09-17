package systems

import (
	"github.com/chezmoi/entities"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/ByteArena/box2d"
	"github.com/Noofbiz/engoBox2dSystem"
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

	hero.AnimationComponent = common.NewAnimationComponent(spriteSheet.Drawables(), 0.1)

	hero.AnimationComponent.AddAnimations(Actions)
	hero.AnimationComponent.SelectAnimationByName("downstop")

	hero.ControlComponent = entities.ControlComponent{
		SchemeHoriz: "horizontal",
		SchemeVert:  "vertical",
	}

	hero.RenderComponent.SetZIndex(6)

	dudeBodyDef := box2d.NewB2BodyDef()
	dudeBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	dudeBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(point)
	dudeBodyDef.Angle = engoBox2dSystem.Conv.DegToRad(hero.SpaceComponent.Rotation)
	dudeBodyDef.FixedRotation = true
	hero.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(dudeBodyDef)
	var dudeBodyShape box2d.B2PolygonShape
	dudeBodyShape.SetAsBox(engoBox2dSystem.Conv.PxToMeters(hero.SpaceComponent.Height/3),
		engoBox2dSystem.Conv.PxToMeters((hero.SpaceComponent.Height/2)-(hero.SpaceComponent.Height/3)))
	dudeFixtureDef := box2d.B2FixtureDef{
		Shape:    &dudeBodyShape,
		Density:  1.0,
		Friction: 0.1,
	}
	hero.Box2dComponent.Body.CreateFixtureFromDef(&dudeFixtureDef)

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
				hero,
				&hero.AnimationComponent,
				&hero.ControlComponent,
				&hero.SpaceComponent,
			)
		case *engoBox2dSystem.PhysicsSystem:
			sys.Add(
				&hero.BasicEntity,
				&hero.SpaceComponent,
				&hero.Box2dComponent,
			)
		case *engoBox2dSystem.CollisionSystem:
			sys.Add(
				&hero.BasicEntity,
				&hero.SpaceComponent,
				&hero.Box2dComponent,
			)
		}
	}

	world.AddSystem(&common.EntityScroller{
		SpaceComponent: &hero.SpaceComponent,
		TrackingBounds: entities.LevelBounds,
	})

	return hero
}
