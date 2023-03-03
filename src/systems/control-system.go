package systems

import (
	"github.com/ByteArena/box2d"
	"github.com/nmorenor/chezmoi/entities"
	"github.com/nmorenor/chezmoi/net"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var (
	WalkUpAction    *common.Animation
	WalkDownAction  *common.Animation
	WalkLeftAction  *common.Animation
	WalkRightAction *common.Animation
	StopUpAction    *common.Animation
	StopDownAction  *common.Animation
	StopLeftAction  *common.Animation
	StopRightAction *common.Animation
	SkillAction     *common.Animation
	Actions         []*common.Animation
	ActionsByKey    map[string]*common.Animation

	UpButton    = "up"
	DownButton  = "down"
	LeftButton  = "left"
	RightButton = "right"
	PauseButton = "pause"
)

type ControlSystem struct {
	entities []entities.ControlEntity
	client   *net.RemoteClient
}

func NewControlSystem(client *net.RemoteClient) *ControlSystem {
	return &ControlSystem{
		entities: make([]entities.ControlEntity, 0),
		client:   client,
	}
}

func (c *ControlSystem) Add(basic *entities.Guy, anim *common.AnimationComponent, control *entities.ControlComponent, space *common.SpaceComponent) {
	c.entities = append(c.entities, entities.ControlEntity{Guy: basic, AnimationComponent: anim, ControlComponent: control, SpaceComponent: space})
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func setAnimation(e entities.ControlEntity) *common.Animation {
	var result *common.Animation
	if engo.Input.Button(UpButton).JustPressed() {
		e.AnimationComponent.SelectAnimationByAction(WalkUpAction)
		result = WalkUpAction
	} else if engo.Input.Button(DownButton).JustPressed() {
		e.AnimationComponent.SelectAnimationByAction(WalkDownAction)
		result = WalkDownAction
	} else if engo.Input.Button(LeftButton).JustPressed() {
		e.AnimationComponent.SelectAnimationByAction(WalkLeftAction)
		result = WalkLeftAction
	} else if engo.Input.Button(RightButton).JustPressed() {
		e.AnimationComponent.SelectAnimationByAction(WalkRightAction)
		result = WalkRightAction
	} else if engo.Input.Button(UpButton).JustReleased() {
		e.AnimationComponent.SelectAnimationByAction(StopUpAction)
		result = StopUpAction
		if engo.Input.Button(LeftButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkLeftAction)
			result = WalkLeftAction
		} else if engo.Input.Button(RightButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkRightAction)
			result = WalkRightAction
		} else if engo.Input.Button(UpButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkUpAction)
			result = WalkUpAction
		} else if engo.Input.Button(DownButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkDownAction)
			result = WalkDownAction
		}
	} else if engo.Input.Button(DownButton).JustReleased() {
		e.AnimationComponent.SelectAnimationByAction(StopDownAction)
		result = StopDownAction
		if engo.Input.Button(LeftButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkLeftAction)
			result = WalkLeftAction
		} else if engo.Input.Button(RightButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkRightAction)
			result = WalkRightAction
		} else if engo.Input.Button(UpButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkUpAction)
			result = WalkUpAction
		} else if engo.Input.Button(DownButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkDownAction)
			result = WalkDownAction
		}
	} else if engo.Input.Button(LeftButton).JustReleased() {
		e.AnimationComponent.SelectAnimationByAction(StopLeftAction)
		result = StopLeftAction
		if engo.Input.Button(LeftButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkLeftAction)
			result = WalkLeftAction
		} else if engo.Input.Button(RightButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkRightAction)
			result = WalkRightAction
		} else if engo.Input.Button(UpButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkUpAction)
			result = WalkUpAction
		} else if engo.Input.Button(DownButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkDownAction)
			result = WalkDownAction
		}
	} else if engo.Input.Button(RightButton).JustReleased() {
		e.AnimationComponent.SelectAnimationByAction(StopRightAction)
		result = StopRightAction
		if engo.Input.Button(LeftButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkLeftAction)
			result = WalkLeftAction
		} else if engo.Input.Button(RightButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkRightAction)
			result = WalkRightAction
		} else if engo.Input.Button(UpButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkUpAction)
			result = WalkUpAction
		} else if engo.Input.Button(DownButton).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkDownAction)
			result = WalkDownAction
		}
	}
	return result
}

func getSpeed(e entities.ControlEntity) (p engo.Point, changed bool) {
	p.X = engo.Input.Axis(e.ControlComponent.SchemeHoriz).Value()
	p.Y = engo.Input.Axis(e.ControlComponent.SchemeVert).Value()
	origX, origY := p.X, p.Y

	if engo.Input.Button(UpButton).JustPressed() {
		p.Y = -1
	} else if engo.Input.Button(DownButton).JustPressed() {
		p.Y = 1
	}
	if engo.Input.Button(LeftButton).JustPressed() {
		p.X = -1
	} else if engo.Input.Button(RightButton).JustPressed() {
		p.X = 1
	}

	if engo.Input.Button(UpButton).JustReleased() || engo.Input.Button(DownButton).JustReleased() {
		p.Y = 0
		changed = true
		if engo.Input.Button(UpButton).Down() {
			p.Y = -1
		} else if engo.Input.Button(DownButton).Down() {
			p.Y = 1
		} else if engo.Input.Button(LeftButton).Down() {
			p.X = -1
		} else if engo.Input.Button(RightButton).Down() {
			p.X = 1
		}
	}
	if engo.Input.Button(LeftButton).JustReleased() || engo.Input.Button(RightButton).JustReleased() {
		p.X = 0
		changed = true
		if engo.Input.Button(LeftButton).Down() {
			p.X = -1
		} else if engo.Input.Button(RightButton).Down() {
			p.X = 1
		} else if engo.Input.Button(UpButton).Down() {
			p.Y = -1
		} else if engo.Input.Button(DownButton).Down() {
			p.Y = 1
		}
	}
	changed = changed || p.X != origX || p.Y != origY
	return
}

func (c *ControlSystem) Update(dt float32) {
	for _, e := range c.entities {
		anim := setAnimation(e)

		if vector, changed := getSpeed(e); changed {
			speed := dt * entities.SPEED_SCALE
			vector, _ = vector.Normalize()
			vector.MultiplyScalar(speed)
			e.Body.SetLinearVelocity(box2d.B2Vec2{X: float64(vector.X) * float64(e.SpaceComponent.Width), Y: float64(vector.Y) * float64(e.SpaceComponent.Width)})
			if anim != nil {
				animName := anim.Name
				go c.client.SendMessage(net.Point{X: vector.X, Y: vector.Y}, net.Point{X: e.SpaceComponent.Position.X, Y: e.SpaceComponent.Position.Y}, animName, c.client.Client.Id)
			}
		}
	}
}
