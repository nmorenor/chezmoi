package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/nmorenor/chezmoi/options"
)

type entity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type SetupSessionSystem struct {
	e         entity
	NextScene string
	HostMode  bool
}

func (s *SetupSessionSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	s.e = entity{basic, space}
}

func (s *SetupSessionSystem) Remove(basic ecs.BasicEntity) {}

func (s *SetupSessionSystem) Update(float32) {
	if engo.Input.Mouse.Action == engo.Press {
		x := engo.Input.Mouse.X + engo.ResizeXOffset/(2*engo.GetGlobalScale().X)
		y := engo.Input.Mouse.Y + engo.ResizeYOffset/(2*engo.GetGlobalScale().Y)
		if s.e.SpaceComponent.Contains(engo.Point{X: x, Y: y}) {
			options.SessionInfo.HostMode = s.HostMode
			engo.SetSceneByName(s.NextScene, true)
		}
	}
}
