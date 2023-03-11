package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type HeadsUpDisplay struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
}
