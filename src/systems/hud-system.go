package systems

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/f1bonacc1/glippy"
	"github.com/nmorenor/chezmoi/entities"
	"github.com/nmorenor/chezmoi/net"
	"github.com/nmorenor/chezmoi/options"
)

const (
	HUDEventType = "HUDEventType"
	CLOSE        = 1
	SHARE        = 2
)

type HUDEventMessage struct {
	EventType int
}

// Type implements the engo.Message interface
func (HUDEventMessage) Type() string {
	return HUDEventType
}

type HudSystem struct {
	world     *ecs.World
	labelFont *common.Font
	client    *net.RemoteClient
	close     bool
}

func NewHudSystem(labelFont *common.Font, client *net.RemoteClient) *HudSystem {
	return &HudSystem{labelFont: labelFont, client: client}
}

func (cb *HudSystem) New(world *ecs.World) {
	cb.world = world
	rectangle, _ := common.LoadedSprite("rectangle-small.svg")
	rectangleBase := entities.Button{BasicEntity: ecs.NewBasic()}

	rectangleBase.RenderComponent.Drawable = rectangle
	rectangleBase.RenderComponent.SetZIndex(1001)
	rectangleBase.SpaceComponent = common.SpaceComponent{Position: engo.Point{X: 0, Y: (engo.WindowHeight()) - rectangle.Height()}, Width: rectangle.Width(), Height: rectangleBase.Height}
	rectangleBase.SetShader(common.HUDShader)

	closeCircle, _ := common.LoadedSprite("close-circle.svg")
	closeSessionButton := entities.Button{BasicEntity: ecs.NewBasic(), MouseComponent: common.MouseComponent{Track: true}}
	closeSessionButton.RenderComponent.SetZIndex(1002)
	closeSessionButton.RenderComponent.Drawable = closeCircle
	closeSessionButton.SpaceComponent.Position.X = 25
	closeSessionButton.SpaceComponent.Position.Y = (rectangleBase.SpaceComponent.Position.Y + 12)
	closeSessionButton.SpaceComponent.Height = closeCircle.Height()
	closeSessionButton.SpaceComponent.Width = closeCircle.Width()
	closeSessionButton.SetShader(common.HUDShader)

	shareSprite, _ := common.LoadedSprite("share.svg")
	shareButton := entities.Button{BasicEntity: ecs.NewBasic(), MouseComponent: common.MouseComponent{Track: true}}
	shareButton.RenderComponent.SetZIndex(1002)
	shareButton.RenderComponent.Drawable = shareSprite
	shareButton.SpaceComponent.Position.X = closeSessionButton.SpaceComponent.Position.X + closeSessionButton.SpaceComponent.Width + 2
	shareButton.SpaceComponent.Position.Y = closeSessionButton.SpaceComponent.Position.Y
	shareButton.SpaceComponent.Height = shareSprite.Height()
	shareButton.SpaceComponent.Width = shareSprite.Width()
	shareButton.SetShader(common.HUDShader)

	closeSystem := &HudEventSystem{EventType: CLOSE}
	world.AddSystem(closeSystem)
	closeSystem.Add(&closeSessionButton.BasicEntity, &closeSessionButton.SpaceComponent)

	shareSystem := &HudEventSystem{EventType: SHARE}
	world.AddSystem(shareSystem)
	shareSystem.Add(&shareButton.BasicEntity, &shareButton.SpaceComponent)

	cursorSystem := NewCursorSystem()
	world.AddSystem(cursorSystem)
	cursorSystem.Add(&closeSessionButton.BasicEntity, &closeSessionButton.SpaceComponent, &closeSessionButton.RenderComponent, &closeSessionButton.MouseComponent)
	cursorSystem.Add(&shareButton.BasicEntity, &shareButton.SpaceComponent, &shareButton.RenderComponent, &shareButton.MouseComponent)

	engo.Mailbox.Listen(HUDEventType, func(m engo.Message) {
		msg, ok := m.(HUDEventMessage)
		if !ok {
			return
		}
		if msg.EventType == CLOSE {
			cb.close = true
			if cb.client.Client.Host {
				engo.Mailbox.Dispatch(SessionEndMessage{
					Session: *cb.client.Client.Session,
				})
			}
			cb.client.Client.Close()
			options.Reset()
		}
		if msg.EventType == SHARE {
			glippy.Set(*cb.client.Client.Session)
			fmt.Println("Session on clipboard: " + *cb.client.Client.Session)
		}
	})

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(
				&rectangleBase.BasicEntity,
				&rectangleBase.RenderComponent,
				&rectangleBase.SpaceComponent,
			)
			sys.Add(
				&closeSessionButton.BasicEntity,
				&closeSessionButton.RenderComponent,
				&closeSessionButton.SpaceComponent,
			)
			sys.Add(
				&shareButton.BasicEntity,
				&shareButton.RenderComponent,
				&shareButton.SpaceComponent,
			)
		case *common.MouseSystem:
			sys.Add(&closeSessionButton.BasicEntity,
				&closeSessionButton.MouseComponent, nil, nil)
		}
	}
}

func (c *HudSystem) Remove(basic ecs.BasicEntity) {

}

func (c *HudSystem) Update(dt float32) {
	if c.close {
		engo.SetSceneByName("Start", true)
	}
}

type HudEventSystem struct {
	e         entity
	EventType int
}

func (s *HudEventSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	s.e = entity{basic, space}
}

func (s *HudEventSystem) Remove(basic ecs.BasicEntity) {}

func (s *HudEventSystem) Update(float32) {
	if engo.Input.Mouse.Action == engo.Press {
		x := engo.Input.Mouse.X + engo.ResizeXOffset/(2*engo.GetGlobalScale().X)
		y := engo.Input.Mouse.Y + engo.ResizeYOffset/(2*engo.GetGlobalScale().Y)
		if s.e.SpaceComponent.Contains(engo.Point{X: x, Y: y}) {
			engo.Mailbox.Dispatch(HUDEventMessage{
				EventType: s.EventType,
			})
		}
	}
}
