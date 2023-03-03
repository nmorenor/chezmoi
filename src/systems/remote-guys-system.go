package systems

import (
	"sync"

	"github.com/nmorenor/chezmoi/entities"
	"github.com/nmorenor/chezmoi/net"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/ByteArena/box2d"
	"github.com/Noofbiz/engoBox2dSystem"
)

type RemoteGuysSystem struct {
	entities    []*entities.ControlEntity
	world       *ecs.World
	spriteSheet *common.Spritesheet
	client      *net.RemoteClient
	mutex       *sync.Mutex
}

func NewRemoteGuysSystem(client *net.RemoteClient) *RemoteGuysSystem {
	system := &RemoteGuysSystem{
		entities:    []*entities.ControlEntity{},
		world:       nil,
		spriteSheet: nil,
		client:      client,
		mutex:       &sync.Mutex{},
	}
	client.OnRemoteUpdate = system.onRemoteUpdate
	client.OnSessionJoin = system.onSessionJoin
	client.OnSessionLeave = system.onSessionLeave
	return system
}

func (system *RemoteGuysSystem) onRemoteUpdate(client *net.RemoteClient, from *string, msg net.Message) {
	system.mutex.Lock()
	defer system.mutex.Unlock()
	for _, next := range system.entities {
		if next.Guy != nil && next.Guy.RemoteId != nil && *next.Guy.RemoteId == *from {
			next.Guy.RemoteVector = &msg.Point
			next.Guy.RemoteAnimation = &msg.Animation
			next.Guy.RemotePosition = &msg.Position
		}
	}
}

func (system *RemoteGuysSystem) onSessionJoin(client *net.RemoteClient, target *string) {
	system.mutex.Lock()
	defer system.mutex.Unlock()
	hero := system.CreateHero(system.world, engo.Point{X: engo.GameWidth() / 2, Y: (engo.GameHeight() / 2) + 128}, system.spriteSheet, *target)
	system.Add(
		hero,
		&hero.AnimationComponent,
		&hero.ControlComponent,
		&hero.SpaceComponent,
	)
}

func (system *RemoteGuysSystem) onSessionLeave(client *net.RemoteClient, target *string) {
	if target == nil {
		return
	}
	system.mutex.Lock()
	defer system.mutex.Unlock()
	for _, e := range system.entities {
		if e.Guy.RemoteId != nil && *e.Guy.RemoteId == *target {
			system.Remove(e.BasicEntity)
		}
	}
}

func (system *RemoteGuysSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range system.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		system.entities = append(system.entities[:delete], system.entities[delete+1:]...)
	}
}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (system *RemoteGuysSystem) Update(dt float32) {
	system.mutex.Lock()
	defer system.mutex.Unlock()
	for _, e := range system.entities {
		if e.Guy != nil && e.Guy.RemoteId != nil {
			targetVector := e.Guy.RemoteVector
			targetAnim := e.Guy.RemoteAnimation

			if targetVector != nil && targetAnim != nil {
				anim := ActionsByKey[*targetAnim]
				if anim != nil {
					e.AnimationComponent.SelectAnimationByAction(anim)
				}
				tVector := engo.Point{X: targetVector.X, Y: targetVector.Y}
				speed := dt * entities.SPEED_SCALE
				vector, _ := tVector.Normalize()
				vector.MultiplyScalar(speed)

				e.Body.SetLinearVelocity(box2d.B2Vec2{X: float64(vector.X) * float64(e.SpaceComponent.Width), Y: float64(vector.Y) * float64(e.SpaceComponent.Width)})
				if e.Guy.RemotePosition != nil && targetVector.X == 0 && targetVector.Y == 0 {
					e.SpaceComponent.Position = engo.Point{X: e.Guy.RemotePosition.X, Y: e.Guy.RemotePosition.Y}
				}
			}
			e.Guy.RemoteAnimation = nil
			e.Guy.RemoteVector = nil
			e.Guy.RemotePosition = nil
		}
	}
}

func (c *RemoteGuysSystem) Add(basic *entities.Guy, anim *common.AnimationComponent, control *entities.ControlComponent, space *common.SpaceComponent) {
	c.entities = append(c.entities, &entities.ControlEntity{Guy: basic, AnimationComponent: anim, ControlComponent: control, SpaceComponent: space})
}

func (system *RemoteGuysSystem) New(world *ecs.World) {
	system.world = world
	system.spriteSheet = common.NewSpritesheetFromFile("textures/run_horizontal_32x32_2.png", 32, 64)
	system.client.Initialize()
}

func (system *RemoteGuysSystem) CreateHero(world *ecs.World, point engo.Point, spriteSheet *common.Spritesheet, remoteId string) *entities.Guy {
	hero := &entities.Guy{BasicEntity: ecs.NewBasic(), RemoteId: &remoteId}

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
		case *RemoteGuysSystem:
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

	return hero
}
