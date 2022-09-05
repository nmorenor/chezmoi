package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const (
	SPEED_MESSAGE = "SpeedMessage"
	SPEED_SCALE   = 16
)

var (
	LevelWidth  float32
	LevelHeight float32
	LevelBounds engo.AABB
)

type SpeedMessage struct {
	*ecs.BasicEntity
	engo.Point
}

func (SpeedMessage) Type() string {
	return SPEED_MESSAGE
}

type SpeedComponent struct {
	engo.Point
}

type SpeedEntity struct {
	*ecs.BasicEntity
	*SpeedComponent
	*common.SpaceComponent
}
