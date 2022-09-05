package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type HeadsUpDisplay struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type HudText struct {
	Text1, Text2, Text3, Text4 string
	BasicEntity                *ecs.BasicEntity
	MouseComponent             *common.MouseComponent
	SpaceComponent             *common.SpaceComponent
}
