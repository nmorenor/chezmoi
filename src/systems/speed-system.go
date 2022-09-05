package systems

import (
	"github.com/chezmoi/entities"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type SpeedSystem struct {
	entities []entities.SpeedEntity
}

func (s *SpeedSystem) New(*ecs.World) {
	engo.Mailbox.Listen(entities.SPEED_MESSAGE, func(message engo.Message) {
		speed, isSpeed := message.(entities.SpeedMessage)
		if isSpeed {
			// fmt.Printf("%#v\n", speed.Point)
			for _, e := range s.entities {
				if e.ID() == speed.BasicEntity.ID() {
					e.SpeedComponent.Point = speed.Point
				}
			}
		}
	})
}

func (s *SpeedSystem) Add(basic *ecs.BasicEntity, speed *entities.SpeedComponent, space *common.SpaceComponent) {
	s.entities = append(s.entities, entities.SpeedEntity{BasicEntity: basic, SpeedComponent: speed, SpaceComponent: space})
}

func (s *SpeedSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *SpeedSystem) Update(dt float32) {
	for _, e := range s.entities {
		speed := engo.GameWidth() * dt
		e.SpaceComponent.Position.X = e.SpaceComponent.Position.X + speed*e.SpeedComponent.Point.X
		e.SpaceComponent.Position.Y = e.SpaceComponent.Position.Y + speed*e.SpeedComponent.Point.Y

		// Add Game Border Limits
		var heightLimit float32 = entities.LevelHeight - e.SpaceComponent.Height
		if e.SpaceComponent.Position.Y < 0 {
			e.SpaceComponent.Position.Y = 0
		} else if e.SpaceComponent.Position.Y > heightLimit {
			e.SpaceComponent.Position.Y = heightLimit
		}

		var widthLimit float32 = entities.LevelWidth - e.SpaceComponent.Width
		if e.SpaceComponent.Position.X < 0 {
			e.SpaceComponent.Position.X = 0
		} else if e.SpaceComponent.Position.X > widthLimit {
			e.SpaceComponent.Position.X = widthLimit
		}
	}
}
