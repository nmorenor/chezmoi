package systems

import (
	"math/rand"
	"time"

	"github.com/chezmoi/entities"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type MouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type TilesSystem struct {
	mouseTracker MouseTracker
	world        *ecs.World
	tiles        []*entities.Tile
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
		}

	}
}
