package world

import (
	"image/color"
	"math/rand"
	"time"
)

// Room représente une salle dans le donjon
type Room struct {
	X, Y          int // Position dans la grille
	Width, Height int // Taille de la salle
	CenterX, CenterY int // Centre de la salle (en tiles)
}

// Corridor représente un couloir entre deux salles
type Corridor struct {
	StartX, StartY int
	EndX, EndY     int
}

// DungeonGenerator génère des donjons procéduraux
type DungeonGenerator struct {
	Width, Height int
	Rooms         []*Room
	Corridors     []*Corridor
	Seed          int64
	rng           *rand.Rand
}

// NewDungeonGenerator crée un nouveau générateur de donjon
func NewDungeonGenerator(width, height int, seed int64) *DungeonGenerator {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	return &DungeonGenerator{
		Width:     width,
		Height:    height,
		Rooms:     make([]*Room, 0),
		Corridors: make([]*Corridor, 0),
		Seed:      seed,
		rng:       rand.New(rand.NewSource(seed)),
	}
}

// Generate génère un donjon complet
func (dg *DungeonGenerator) Generate(isRestArea bool) *TileMap {
	// Générer les salles
	dg.generateRooms(isRestArea)

	// Connecter les salles avec des couloirs
	dg.connectRooms()

	// Créer la tilemap
	tilemap := dg.createTileMap()

	return tilemap
}

// generateRooms génère les salles du donjon
func (dg *DungeonGenerator) generateRooms(isRestArea bool) {
	if isRestArea {
		// Salle de repos : une seule grande salle centrale
		roomWidth := 20 + dg.rng.Intn(8)  // 20-27 tiles
		roomHeight := 20 + dg.rng.Intn(8) // 20-27 tiles
		roomX := (dg.Width - roomWidth) / 2
		roomY := (dg.Height - roomHeight) / 2

		room := &Room{
			X:       roomX,
			Y:       roomY,
			Width:   roomWidth,
			Height:  roomHeight,
			CenterX: roomX + roomWidth/2,
			CenterY: roomY + roomHeight/2,
		}
		dg.Rooms = append(dg.Rooms, room)
	} else {
		// Étage normal : 4-7 salles
		numRooms := 4 + dg.rng.Intn(4)

		for i := 0; i < numRooms; i++ {
			maxAttempts := 50
			for attempt := 0; attempt < maxAttempts; attempt++ {
				// Taille aléatoire
				roomWidth := 8 + dg.rng.Intn(9)  // 8-16 tiles
				roomHeight := 8 + dg.rng.Intn(9) // 8-16 tiles

				// Position aléatoire
				roomX := 2 + dg.rng.Intn(dg.Width-roomWidth-4)
				roomY := 2 + dg.rng.Intn(dg.Height-roomHeight-4)

				room := &Room{
					X:       roomX,
					Y:       roomY,
					Width:   roomWidth,
					Height:  roomHeight,
					CenterX: roomX + roomWidth/2,
					CenterY: roomY + roomHeight/2,
				}

				// Vérifier qu'elle ne chevauche pas les autres salles
				if !dg.roomOverlaps(room) {
					dg.Rooms = append(dg.Rooms, room)
					break
				}
			}
		}
	}
}

// roomOverlaps vérifie si une salle chevauche une autre
func (dg *DungeonGenerator) roomOverlaps(newRoom *Room) bool {
	for _, room := range dg.Rooms {
		// Ajouter une marge de 2 tiles entre les salles
		if newRoom.X < room.X+room.Width+2 &&
			newRoom.X+newRoom.Width+2 > room.X &&
			newRoom.Y < room.Y+room.Height+2 &&
			newRoom.Y+newRoom.Height+2 > room.Y {
			return true
		}
	}
	return false
}

// connectRooms connecte les salles avec des couloirs
func (dg *DungeonGenerator) connectRooms() {
	if len(dg.Rooms) <= 1 {
		return
	}

	// Connecter chaque salle à la suivante
	for i := 0; i < len(dg.Rooms)-1; i++ {
		room1 := dg.Rooms[i]
		room2 := dg.Rooms[i+1]

		// Créer un couloir en L entre les deux salles
		dg.createLCorridor(room1.CenterX, room1.CenterY, room2.CenterX, room2.CenterY)
	}

	// Ajouter quelques connexions supplémentaires aléatoires
	numExtraConnections := dg.rng.Intn(2) + 1
	for i := 0; i < numExtraConnections && len(dg.Rooms) > 2; i++ {
		idx1 := dg.rng.Intn(len(dg.Rooms))
		idx2 := dg.rng.Intn(len(dg.Rooms))
		if idx1 != idx2 {
			room1 := dg.Rooms[idx1]
			room2 := dg.Rooms[idx2]
			dg.createLCorridor(room1.CenterX, room1.CenterY, room2.CenterX, room2.CenterY)
		}
	}
}

// createLCorridor crée un couloir en forme de L
func (dg *DungeonGenerator) createLCorridor(x1, y1, x2, y2 int) {
	// Choisir aléatoirement horizontal puis vertical, ou vertical puis horizontal
	if dg.rng.Intn(2) == 0 {
		// Horizontal puis vertical
		dg.Corridors = append(dg.Corridors, &Corridor{x1, y1, x2, y1})
		dg.Corridors = append(dg.Corridors, &Corridor{x2, y1, x2, y2})
	} else {
		// Vertical puis horizontal
		dg.Corridors = append(dg.Corridors, &Corridor{x1, y1, x1, y2})
		dg.Corridors = append(dg.Corridors, &Corridor{x1, y2, x2, y2})
	}
}

// createTileMap crée la tilemap à partir des salles et couloirs
func (dg *DungeonGenerator) createTileMap() *TileMap {
	tilemap := NewTileMap(dg.Width, dg.Height)

	// Remplir tout de murs
	for y := 0; y < dg.Height; y++ {
		for x := 0; x < dg.Width; x++ {
			tilemap.Tiles[y][x].Type = TileWall
			tilemap.Tiles[y][x].Color = color.RGBA{
				R: uint8(60 + dg.rng.Intn(20)),
				G: uint8(60 + dg.rng.Intn(20)),
				B: uint8(60 + dg.rng.Intn(20)),
				A: 255,
			}
		}
	}

	// Dessiner les salles
	for _, room := range dg.Rooms {
		for y := room.Y; y < room.Y+room.Height; y++ {
			for x := room.X; x < room.X+room.Width; x++ {
				if x >= 0 && x < dg.Width && y >= 0 && y < dg.Height {
					tilemap.Tiles[y][x].Type = TileGrass
					// Sol beaucoup moins coloré (gris-brun foncé)
					tilemap.Tiles[y][x].Color = color.RGBA{
						R: uint8(50 + dg.rng.Intn(20)),
						G: uint8(55 + dg.rng.Intn(20)),
						B: uint8(45 + dg.rng.Intn(20)),
						A: 255,
					}
				}
			}
		}
	}

	// Dessiner les couloirs
	for _, corridor := range dg.Corridors {
		x1, y1 := corridor.StartX, corridor.StartY
		x2, y2 := corridor.EndX, corridor.EndY

		// Couloir horizontal (largeur de 2 tiles)
		if y1 == y2 {
			startX := min(x1, x2)
			endX := max(x1, x2)
			for x := startX; x <= endX; x++ {
				for dy := 0; dy <= 1; dy++ {
					yPos := y1 + dy
					if x >= 0 && x < dg.Width && yPos >= 0 && yPos < dg.Height {
						tilemap.Tiles[yPos][x].Type = TileGrass
						// Sol beaucoup moins coloré (gris-brun foncé)
						tilemap.Tiles[yPos][x].Color = color.RGBA{
							R: uint8(50 + dg.rng.Intn(20)),
							G: uint8(55 + dg.rng.Intn(20)),
							B: uint8(45 + dg.rng.Intn(20)),
							A: 255,
						}
					}
				}
			}
		} else {
			// Couloir vertical (largeur de 2 tiles)
			startY := min(y1, y2)
			endY := max(y1, y2)
			for y := startY; y <= endY; y++ {
				for dx := 0; dx <= 1; dx++ {
					xPos := x1 + dx
					if xPos >= 0 && xPos < dg.Width && y >= 0 && y < dg.Height {
						tilemap.Tiles[y][xPos].Type = TileGrass
						// Sol beaucoup moins coloré (gris-brun foncé)
						tilemap.Tiles[y][xPos].Color = color.RGBA{
							R: uint8(50 + dg.rng.Intn(20)),
							G: uint8(55 + dg.rng.Intn(20)),
							B: uint8(45 + dg.rng.Intn(20)),
							A: 255,
						}
					}
				}
			}
		}
	}

	return tilemap
}

// GetRandomFloorPosition retourne une position aléatoire sur le sol
func (dg *DungeonGenerator) GetRandomFloorPosition() (int, int) {
	if len(dg.Rooms) == 0 {
		return dg.Width / 2, dg.Height / 2
	}

	// Choisir une salle aléatoire
	room := dg.Rooms[dg.rng.Intn(len(dg.Rooms))]

	// Position aléatoire dans la salle (avec marge)
	x := room.X + 2 + dg.rng.Intn(max(1, room.Width-4))
	y := room.Y + 2 + dg.rng.Intn(max(1, room.Height-4))

	return x, y
}

// GetEntrancePosition retourne une position pour l'entrée (première salle)
func (dg *DungeonGenerator) GetEntrancePosition() (int, int) {
	if len(dg.Rooms) == 0 {
		return dg.Width / 2, dg.Height / 2
	}
	room := dg.Rooms[0]
	return room.CenterX, room.CenterY
}

// GetExitPosition retourne une position pour la sortie (dernière salle)
func (dg *DungeonGenerator) GetExitPosition() (int, int) {
	if len(dg.Rooms) == 0 {
		return dg.Width / 2, dg.Height / 2
	}
	room := dg.Rooms[len(dg.Rooms)-1]
	return room.CenterX, room.CenterY
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
