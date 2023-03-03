package entities

import (
	"github.com/Noofbiz/engoBox2dSystem"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
	"github.com/nmorenor/chezmoi/net"
)

const (
	SPEED_MESSAGE = "SpeedMessage"
	SPEED_SCALE   = 16
)

type Guy struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
	engoBox2dSystem.Box2dComponent
	ControlComponent

	RemoteId        *string
	RemoteVector    *net.Point
	RemotePosition  *net.Point
	RemoteAnimation *string
	Left            bool
	Right           bool
	Up              bool
	Down            bool
}
