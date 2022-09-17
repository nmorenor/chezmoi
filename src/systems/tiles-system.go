package systems

import (
	"math/rand"
	"time"

	"github.com/chezmoi/entities"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/ByteArena/box2d"
	"github.com/Noofbiz/engoBox2dSystem"
)

type MouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type TilesSystem struct {
	mouseTracker MouseTracker
	world        *ecs.World
	tiles        []*entities.Tile
	objects      []*entities.CollisionObject
	elapsed      float32
	buildTime    float32
	built        int
}

// Remove is called whenever an Entity is removed from the World, in order to remove it from this sytem as well
func (cb *TilesSystem) Remove(ecs.BasicEntity) {
	// do nothing
}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (cb *TilesSystem) Update(dt float32) {
	renderSystem := new(common.RenderSystem)
	for _, system := range cb.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			renderSystem = sys
		}
	}
	if renderSystem == nil {
		return
	}
	cb.elapsed += dt
	if cb.elapsed >= cb.buildTime {
		cb.elapsed = 0
		cb.built++
	}
	cb.addTiles(renderSystem)
}

func (cb *TilesSystem) addTiles(renderSystem *common.RenderSystem) {
	for _, v := range cb.tiles {
		renderSystem.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
	}
}

func (cb *TilesSystem) addObjectsToCollision(collisionSystem *engoBox2dSystem.CollisionSystem) {
	for _, v := range cb.objects {
		collisionSystem.Add(&v.BasicEntity, &v.SpaceComponent, &v.Box2dComponent)
	}
}

func (cb *TilesSystem) addObjectsToPhysics(physicsSystem *engoBox2dSystem.PhysicsSystem) {
	for _, v := range cb.objects {
		physicsSystem.Add(&v.BasicEntity, &v.SpaceComponent, &v.Box2dComponent)
	}
}

func (cb *TilesSystem) initTiles() {
	resource, err := engo.Files.Resource("levels/interiors.tmx")
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(common.TMXResource)
	levelData := tmxResource.Level

	entities.LevelWidth = levelData.Bounds().Max.X
	entities.LevelHeight = levelData.Bounds().Max.Y
	entities.LevelBounds = levelData.Bounds()

	tiles := make([]*entities.Tile, 0)
	layer := float32(2)
	for _, tileLayer := range levelData.TileLayers {
		for _, tileElement := range tileLayer.Tiles {
			if tileElement.Image != nil {
				tile := &entities.Tile{BasicEntity: ecs.NewBasic()}
				center := tile.Center()
				tile.Rotation = tileElement.Rotation
				tile.SetCenter(center)
				tile.RenderComponent = common.RenderComponent{
					Drawable:    tileElement.Image,
					Scale:       engo.Point{X: 1, Y: 1},
					StartZIndex: layer,
				}
				tile.SpaceComponent = common.SpaceComponent{
					Position: tileElement.Point,
					Width:    32,
					Height:   32,
				}
				tiles = append(tiles, tile)
			}
		}
		layer++
	}
	cb.tiles = tiles

	objects := make([]*entities.CollisionObject, 0)
	for _, objectLayer := range levelData.ObjectLayers {
		for _, objectElement := range objectLayer.Objects {
			if objectElement.Type == "room" {
				continue
			}
			collisionObject := entities.CollisionObject{
				BasicEntity: ecs.NewBasic(),
				SpaceComponent: common.SpaceComponent{
					Position: engo.Point{X: objectElement.X, Y: objectElement.Y},
					Width:    objectElement.Width,
					Height:   objectElement.Height,
				},
			}

			objectBodyDef := box2d.NewB2BodyDef()
			objectBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(collisionObject.SpaceComponent.Center())
			objectBodyDef.Angle = engoBox2dSystem.Conv.DegToRad(collisionObject.SpaceComponent.Rotation)
			if objectElement.Lines != nil && len(objectElement.Lines) > 0 {
				objectBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
				collisionObject.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(objectBodyDef)
				var starBodyShape box2d.B2PolygonShape
				var vertices []box2d.B2Vec2

				for _, nextLines := range objectElement.Lines {
					for _, nextLine := range nextLines.Lines {
						vertices = append(vertices, box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(nextLine.P1.X), Y: engoBox2dSystem.Conv.PxToMeters(nextLine.P1.Y)})
						vertices = append(vertices, box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(nextLine.P2.X), Y: engoBox2dSystem.Conv.PxToMeters(nextLine.P2.Y)})
					}
				}
				starBodyShape.Set(vertices, len(vertices))
				starFixtureDef := box2d.B2FixtureDef{Shape: &starBodyShape}
				collisionObject.Box2dComponent.Body.CreateFixtureFromDef(&starFixtureDef)
			} else {
				collisionObject.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(objectBodyDef)
				var objectBodyShape box2d.B2PolygonShape
				objectBodyShape.SetAsBox(engoBox2dSystem.Conv.PxToMeters(collisionObject.SpaceComponent.Width/2),
					engoBox2dSystem.Conv.PxToMeters(collisionObject.SpaceComponent.Height/2))
				objectFixtureDef := box2d.B2FixtureDef{Shape: &objectBodyShape}
				collisionObject.Box2dComponent.Body.CreateFixtureFromDef(&objectFixtureDef)
			}
			objects = append(objects, &collisionObject)
		}
	}
	cb.objects = objects
	common.CameraBounds = levelData.Bounds()
}

func (cb *TilesSystem) New(world *ecs.World) {
	cb.mouseTracker.BasicEntity = ecs.NewBasic()
	cb.mouseTracker.MouseComponent = common.MouseComponent{Track: true}
	cb.world = world
	cb.initTiles()
	rand.Seed(time.Now().UnixNano())

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&cb.mouseTracker.BasicEntity, &cb.mouseTracker.MouseComponent, nil, nil)
		case *common.RenderSystem:
			cb.addTiles(sys)
		case *engoBox2dSystem.PhysicsSystem:
			cb.addObjectsToPhysics(sys)
		case *engoBox2dSystem.CollisionSystem:
			cb.addObjectsToCollision(sys)
		}

	}
}
