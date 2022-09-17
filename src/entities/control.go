package entities

import (
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
