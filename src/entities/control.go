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
	*Guy
	*common.AnimationComponent
	*ControlComponent
	*common.SpaceComponent
}

type Label struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type Button struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	common.AudioComponent
}
