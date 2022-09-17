package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Noofbiz/engoBox2dSystem"
)

var (
	LevelWidth  float32
	LevelHeight float32
	LevelBounds engo.AABB
)

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.CollisionComponent
}

type CollisionObject struct {
	ecs.BasicEntity
	common.SpaceComponent
	engoBox2dSystem.Box2dComponent
}
