package systems

import (
	"sync"

	"github.com/nmorenor/chezmoi/entities"
	"github.com/nmorenor/chezmoi/net"
	"github.com/nmorenor/chezmoi/options"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/ByteArena/box2d"
	"github.com/Noofbiz/engoBox2dSystem"
)

const (
	SessionEndType = "SessionEndType"
)

type SessionEndMessage struct {
	Session string
}

// Type implements the engo.Message interface
func (SessionEndMessage) Type() string {
	return SessionEndType
}

type RemoteGuysSystem struct {
	entities    []*entities.ControlEntity
	labelFont   *common.Font
	world       *ecs.World
	spriteSheet *common.Spritesheet
	client      *net.RemoteClient
	mutex       *sync.Mutex
}

func NewRemoteGuysSystem(client *net.RemoteClient, labelFont *common.Font) *RemoteGuysSystem {
	system := &RemoteGuysSystem{
		entities:    []*entities.ControlEntity{},
		labelFont:   labelFont,
		world:       nil,
		spriteSheet: nil,
		client:      client,
		mutex:       &sync.Mutex{},
	}
	client.OnRemoteUpdate = system.onRemoteUpdate
	client.OnSessionJoin = system.onSessionJoin
	client.OnSessionLeave = system.onSessionLeave
	client.OnSessionEnd = system.onSessionEnd
	return system
}

func (system *RemoteGuysSystem) onSessionEnd() {
	engo.Mailbox.Dispatch(SessionEndMessage{
		Session: *options.SessionInfo.Session,
	})
	options.Reset()
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

func (system *RemoteGuysSystem) onSessionJoin(client *net.RemoteClient, target *string, targetPosition *net.Point, anim *string) {
	system.mutex.Lock()
	defer system.mutex.Unlock()
	var position engo.Point
	if targetPosition != nil && targetPosition.X != 0 && targetPosition.Y != 0 {
		position = engo.Point{X: targetPosition.X, Y: targetPosition.Y}
	} else {
		position = engo.Point{X: engo.GameWidth() / 2, Y: (engo.GameHeight() / 2) + 128}
	}

	hero := system.CreateHero(system.world, position, system.spriteSheet, *target, *client.Participants[*target])
	if anim != nil {
		targetAnim := ActionsByKey[*anim]
		if targetAnim != nil {
			hero.AnimationComponent.SelectAnimationByAction(targetAnim)
		}
	}
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
			for _, system := range system.world.Systems() {
				switch sys := system.(type) {
				case *common.RenderSystem:
					sys.Remove(e.BasicEntity)
					sys.Remove(e.Label.BasicEntity)
				case *common.AnimationSystem:
					sys.Remove(e.BasicEntity)
				case *RemoteGuysSystem:
					sys.Remove(e.BasicEntity)
				case *engoBox2dSystem.PhysicsSystem:
					sys.Remove(e.BasicEntity)
				case *engoBox2dSystem.CollisionSystem:
					sys.Remove(e.BasicEntity)
				}
			}
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

	if options.SessionInfo.Client == nil {
		engo.SetSceneByName("Start", true)
		return
	}

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
			mid := e.SpaceComponent.Position.X + (e.SpaceComponent.Width / 2) - 4
			target := mid - (e.Label.RenderComponent.Drawable.Width() / 2)
			e.Label.SpaceComponent.Position = engo.Point{X: target, Y: e.SpaceComponent.Position.Y}
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

func (guysSystem *RemoteGuysSystem) CreateHero(world *ecs.World, point engo.Point, spriteSheet *common.Spritesheet, remoteId string, remoteName string) *entities.Guy {
	hero := &entities.Guy{BasicEntity: ecs.NewBasic(), RemoteId: &remoteId, Name: &remoteName}

	hostLabel := entities.Label{BasicEntity: ecs.NewBasic()}
	hostLabel.RenderComponent.Drawable = common.Text{
		Font: guysSystem.labelFont,
		Text: *hero.Name,
	}
	hostLabel.SpaceComponent.Position = point

	hostLabel.RenderComponent.SetZIndex(6)
	hero.Label = &hostLabel

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
			sys.Add(
				&hero.Label.BasicEntity,
				&hero.Label.RenderComponent,
				&hero.Label.SpaceComponent,
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
