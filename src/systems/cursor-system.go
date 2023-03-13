package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/nmorenor/chezmoi/entities"
)

const (
	CursorEventType = "CursorEventType"
)

type CursorEventMessage struct {
	Element *entities.UIElement
}

// Type implements the engo.Message interface
func (CursorEventMessage) Type() string {
	return CursorEventType
}

type CursorSystem struct {
	entities []*entities.UIElement
	world    *ecs.World
}

func NewCursorSystem() *CursorSystem {
	return &CursorSystem{
		entities: make([]*entities.UIElement, 0),
	}
}

func (cb *CursorSystem) Add(entity *ecs.BasicEntity, spaceComponent *common.SpaceComponent, renderComponent *common.RenderComponent, mouseComponent *common.MouseComponent) {
	cb.entities = append(cb.entities, &entities.UIElement{
		BasicEntity:     entity,
		SpaceComponent:  spaceComponent,
		RenderComponent: renderComponent,
		MouseComponent:  mouseComponent,
	})
}

func (cb *CursorSystem) New(world *ecs.World) {
	cb.world = world
}

func (c *CursorSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.ID() == basic.ID() {
			delete = index
			break
		}
		for _, system := range c.world.Systems() {
			switch sys := system.(type) {
			case *common.MouseSystem:
				sys.Remove(*e.BasicEntity)
			}
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func (c *CursorSystem) Update(dt float32) {
	hand := false
	for _, e := range c.entities {
		e.MouseComponent.Hovered = false
		x := engo.Input.Mouse.X + engo.ResizeXOffset/(2*engo.GetGlobalScale().X)
		y := engo.Input.Mouse.Y + engo.ResizeYOffset/(2*engo.GetGlobalScale().Y)
		if e.SpaceComponent.Contains(engo.Point{X: x, Y: y}) {
			engo.SetCursor(engo.CursorHand)
			e.MouseComponent.Hovered = true
			hand = true
		}
	}
	if !hand {
		engo.SetCursor(engo.CursorArrow)
	}
}
