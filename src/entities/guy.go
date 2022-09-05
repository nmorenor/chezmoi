package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type Guy struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
	ControlComponent
	SpeedComponent

	Left  bool
	Right bool
	Up    bool
	Down  bool
}
