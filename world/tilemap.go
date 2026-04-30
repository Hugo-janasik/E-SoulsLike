package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
)

// TileType définit les différents types de tuiles
type TileType int

const (
	TileGrass TileType = iota
	TileDirt
	TileStone
	TileWater
	TileWall
)

const (
	TileSize = 32
)

// Tile représente une tuile du monde
type Tile struct {
	Type  TileType
	Color color.RGBA
}

// TileMap représente la carte du monde
type TileMap struct {
	Width        int
	Height       int
	Tiles        [][]Tile
	WalkableWater bool // Si true, l'eau est traversable (pour la zone hub)
}

// NewTileMap crée une nouvelle carte
func NewTileMap(width, height int) *TileMap {
	tm := &TileMap{
		Width:  width,
		Height: height,
		Tiles:  make([][]Tile, height),
	}

	// Générer une carte procédurale simple
	for y := 0; y < height; y++ {
		tm.Tiles[y] = make([]Tile, width)
		for x := 0; x < width; x++ {
			// Génération procédurale basique
			noise := rand.Float64()
			var tileType TileType
			var tileColor color.RGBA

			if noise < 0.6 {
				tileType = TileGrass
				// Sol beaucoup moins coloré (gris-brun foncé)
				tileColor = color.RGBA{
					R: uint8(50 + rand.Intn(20)),
					G: uint8(55 + rand.Intn(20)),
					B: uint8(45 + rand.Intn(20)),
					A: 255,
				}
			} else if noise < 0.85 {
				tileType = TileDirt
				// Terre plus sombre
				tileColor = color.RGBA{
					R: uint8(70 + rand.Intn(20)),
					G: uint8(60 + rand.Intn(15)),
					B: uint8(45 + rand.Intn(15)),
					A: 255,
				}
			} else if noise < 0.95 {
				tileType = TileStone
				// Pierre plus foncée
				tileColor = color.RGBA{
					R: uint8(70 + rand.Intn(15)),
					G: uint8(70 + rand.Intn(15)),
					B: uint8(75 + rand.Intn(15)),
					A: 255,
				}
			} else {
				tileType = TileWater
				// Eau beaucoup plus sombre et désaturée
				tileColor = color.RGBA{
					R: uint8(35 + rand.Intn(15)),
					G: uint8(50 + rand.Intn(20)),
					B: uint8(80 + rand.Intn(30)),
					A: 255,
				}
			}

			tm.Tiles[y][x] = Tile{
				Type:  tileType,
				Color: tileColor,
			}
		}
	}

	return tm
}

// Draw dessine la carte
func (tm *TileMap) Draw(screen *ebiten.Image, camera *Camera) {
	// Calculer la zone visible
	startX := int(camera.X / TileSize)
	startY := int(camera.Y / TileSize)
	endX := startX + int(camera.Width/TileSize) + 2
	endY := startY + int(camera.Height/TileSize) + 2

	// Limiter aux bornes de la carte
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}
	if endX >= tm.Width {
		endX = tm.Width - 1
	}
	if endY >= tm.Height {
		endY = tm.Height - 1
	}

	// Dessiner uniquement les tuiles visibles
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			tile := tm.Tiles[y][x]

			worldX := float64(x * TileSize)
			worldY := float64(y * TileSize)

			screenX, screenY := camera.WorldToScreen(worldX, worldY)

			// Dessiner la tuile
			vector.DrawFilledRect(
				screen,
				float32(screenX),
				float32(screenY),
				TileSize*float32(camera.Zoom),
				TileSize*float32(camera.Zoom),
				tile.Color,
				false,
			)

			// Dessiner une grille légère
			vector.StrokeRect(
				screen,
				float32(screenX),
				float32(screenY),
				TileSize*float32(camera.Zoom),
				TileSize*float32(camera.Zoom),
				1,
				color.RGBA{0, 0, 0, 20},
				false,
			)
		}
	}
}

// GetTile retourne la tuile à la position donnée
func (tm *TileMap) GetTile(x, y int) *Tile {
	if x < 0 || x >= tm.Width || y < 0 || y >= tm.Height {
		return nil
	}
	return &tm.Tiles[y][x]
}

// IsWalkable vérifie si une tuile est traversable
func (tm *TileMap) IsWalkable(x, y int) bool {
	tile := tm.GetTile(x, y)
	if tile == nil {
		return false
	}
	// Les murs ne sont jamais traversables
	if tile.Type == TileWall {
		return false
	}
	// L'eau est traversable si WalkableWater est activé (zone hub)
	if tile.Type == TileWater {
		return tm.WalkableWater
	}
	return true
}
