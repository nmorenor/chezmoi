package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type ControlComponent struct {
	SchemeVert  string
	SchemeHoriz string
}

type ControlEntity struct {
	*ecs.BasicEntity
	*common.AnimationComponent
	*ControlComponent
	*common.SpaceComponent
}
