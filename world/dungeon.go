package world

import (
	"fmt"
)

// DungeonDifficulty représente la difficulté d'un donjon
type DungeonDifficulty int

const (
	DifficultyEasy DungeonDifficulty = iota
	DifficultyMedium
	DifficultyHard
)

// Floor représente un étage d'un donjon
type Floor struct {
	Number        int       // Numéro de l'étage (1, 2, 3, etc.)
	Zone          *Zone     // La zone (map) de cet étage
	IsRestArea    bool      // True si c'est une salle de repos avec feu de camp
	NextFloorPortal *Portal // Portail vers l'étage suivant (nil si dernier étage)
	PrevFloorPortal *Portal // Portail vers l'étage précédent (nil si premier étage)
	Generator     *DungeonGenerator // Générateur utilisé pour cet étage
}

// Dungeon représente un donjon complet avec plusieurs étages
type Dungeon struct {
	ID          string
	Name        string
	Description string
	Difficulty  DungeonDifficulty
	TotalFloors int
	Floors      []*Floor
	CurrentFloor int // Index de l'étage actuel (0-based)

	// Point d'entrée du donjon (pour revenir au hub)
	EntrancePortal *Portal
}

// NewDungeon crée un nouveau donjon
func NewDungeon(id, name, description string, difficulty DungeonDifficulty, totalFloors int) *Dungeon {
	return &Dungeon{
		ID:          id,
		Name:        name,
		Description: description,
		Difficulty:  difficulty,
		TotalFloors: totalFloors,
		Floors:      make([]*Floor, 0, totalFloors),
		CurrentFloor: 0,
	}
}

// CreateFloors génère tous les étages du donjon
func (d *Dungeon) CreateFloors() {
	var theme ZoneTheme
	switch d.Difficulty {
	case DifficultyEasy:
		theme = ThemeForest
	case DifficultyMedium:
		theme = ThemeRuins
	case DifficultyHard:
		theme = ThemeCatacombs
	}

	for i := 1; i <= d.TotalFloors; i++ {
		// Tous les 2 étages, créer une salle de repos
		isRestArea := (i%2 == 0 && i < d.TotalFloors)

		var floor *Floor
		if isRestArea {
			floor = d.createRestFloor(i, theme)
		} else {
			floor = d.createCombatFloor(i, theme)
		}

		d.Floors = append(d.Floors, floor)
	}

	// Créer les portails entre les étages
	d.linkFloors()
}

// createCombatFloor crée un étage de combat normal avec génération procédurale
func (d *Dungeon) createCombatFloor(floorNumber int, theme ZoneTheme) *Floor {
	// Taille de la map basée sur la difficulté (augmentée)
	var width, height int
	switch d.Difficulty {
	case DifficultyEasy:
		width, height = 80, 80 // Augmenté de 50 à 80
	case DifficultyMedium:
		width, height = 90, 90 // Augmenté de 60 à 90
	case DifficultyHard:
		width, height = 100, 100 // Augmenté de 70 à 100
	}

	zoneName := fmt.Sprintf("%s - Étage %d", d.Name, floorNumber)
	zoneID := fmt.Sprintf("%s_floor_%d", d.ID, floorNumber)

	// Créer le générateur avec seed unique basé sur l'ID du donjon et le numéro d'étage
	seed := int64(hashString(d.ID)) + int64(floorNumber)*1000
	generator := NewDungeonGenerator(width, height, seed)

	// Générer la tilemap procédurale
	tilemap := generator.Generate(false)

	// Appliquer le thème
	tilemap = applyThemeToTilemap(tilemap, theme)

	// Créer la zone avec la tilemap générée
	zone := &Zone{
		ID:          zoneID,
		Name:        zoneName,
		Description: d.Description,
		TileMap:     tilemap,
		SpawnPoints: make(map[string]SpawnPoint),
		Enemies:     make([]interface{}, 0),
		Campfires:   make([]interface{}, 0),
		Portals:     make([]*Portal, 0),
		Theme:       theme,
	}

	// Points de spawn basés sur les positions générées
	entranceX, entranceY := generator.GetEntrancePosition()
	exitX, exitY := generator.GetExitPosition()

	zone.AddSpawnPoint("entrance", float64(entranceX*TileSize), float64(entranceY*TileSize))
	zone.AddSpawnPoint("exit", float64(exitX*TileSize), float64(exitY*TileSize))
	zone.AddSpawnPoint("main", float64(entranceX*TileSize), float64(entranceY*TileSize))

	return &Floor{
		Number:     floorNumber,
		Zone:       zone,
		IsRestArea: false,
		Generator:  generator,
	}
}

// createRestFloor crée une salle de repos avec feu de camp
func (d *Dungeon) createRestFloor(floorNumber int, theme ZoneTheme) *Floor {
	// Salle pour les repos (augmentée)
	width, height := 60, 60 // Augmenté de 40 à 60

	zoneName := fmt.Sprintf("%s - Salle de Repos (Étage %d)", d.Name, floorNumber)
	zoneID := fmt.Sprintf("%s_rest_%d", d.ID, floorNumber)

	// Créer le générateur pour la salle de repos
	seed := int64(hashString(d.ID)) + int64(floorNumber)*1000 + 500
	generator := NewDungeonGenerator(width, height, seed)

	// Générer une grande salle centrale
	tilemap := generator.Generate(true)

	// Thème par défaut pour les salles de repos
	tilemap = applyThemeToTilemap(tilemap, ThemeDefault)

	// Créer la zone
	zone := &Zone{
		ID:          zoneID,
		Name:        zoneName,
		Description: "Une salle sûre pour se reposer",
		TileMap:     tilemap,
		SpawnPoints: make(map[string]SpawnPoint),
		Enemies:     make([]interface{}, 0),
		Campfires:   make([]interface{}, 0),
		Portals:     make([]*Portal, 0),
		Theme:       ThemeDefault,
	}

	// Point de spawn au centre
	centerX := float64(width * TileSize / 2)
	centerY := float64(height * TileSize / 2)
	zone.AddSpawnPoint("entrance", centerX, centerY-150)
	zone.AddSpawnPoint("exit", centerX, centerY+150)
	zone.AddSpawnPoint("main", centerX, centerY)

	return &Floor{
		Number:     floorNumber,
		Zone:       zone,
		IsRestArea: true,
		Generator:  generator,
	}
}

// linkFloors crée les portails entre les étages
func (d *Dungeon) linkFloors() {
	for i := 0; i < len(d.Floors); i++ {
		floor := d.Floors[i]

		// Utiliser le générateur pour obtenir des positions valides
		var entranceX, entranceY, exitX, exitY int
		if floor.Generator != nil {
			entranceX, entranceY = floor.Generator.GetEntrancePosition()
			exitX, exitY = floor.Generator.GetExitPosition()
		} else {
			// Fallback si pas de générateur
			entranceX = floor.Zone.TileMap.Width / 2
			entranceY = floor.Zone.TileMap.Height / 2
			exitX = entranceX
			exitY = entranceY
		}

		// Portail vers l'étage suivant
		if i < len(d.Floors)-1 {
			nextFloorID := d.Floors[i+1].Zone.ID
			portalX := float64(exitX * TileSize)
			portalY := float64(exitY * TileSize)
			portal := NewPortal(portalX, portalY, nextFloorID, "entrance", fmt.Sprintf("Étage %d", i+2))
			floor.NextFloorPortal = portal
			floor.Zone.AddPortal(portal)
		}

		// Portail vers l'étage précédent
		if i > 0 {
			prevFloorID := d.Floors[i-1].Zone.ID
			portalX := float64(entranceX * TileSize)
			portalY := float64(entranceY * TileSize)
			portal := NewPortal(portalX, portalY, prevFloorID, "exit", fmt.Sprintf("Étage %d", i))
			floor.PrevFloorPortal = portal
			floor.Zone.AddPortal(portal)
		}

		// Portail de sortie vers le hub UNIQUEMENT dans les salles de repos
		if floor.IsRestArea {
			// Dans une salle de repos, placer le portail de sortie au centre
			centerX := float64(floor.Zone.TileMap.Width * TileSize / 2)
			centerY := float64(floor.Zone.TileMap.Height * TileSize / 2)
			exitPortal := NewPortal(centerX, centerY+100, "firelink", "from_dungeon", "Sortir du donjon")
			floor.Zone.AddPortal(exitPortal)
		}
	}
}

// GetCurrentFloor retourne l'étage actuel
func (d *Dungeon) GetCurrentFloor() *Floor {
	if d.CurrentFloor >= 0 && d.CurrentFloor < len(d.Floors) {
		return d.Floors[d.CurrentFloor]
	}
	return nil
}

// GetFloor retourne un étage par son numéro (1-based)
func (d *Dungeon) GetFloor(floorNumber int) *Floor {
	if floorNumber >= 1 && floorNumber <= len(d.Floors) {
		return d.Floors[floorNumber-1]
	}
	return nil
}

// SetCurrentFloor définit l'étage actuel
func (d *Dungeon) SetCurrentFloor(floorNumber int) bool {
	if floorNumber >= 1 && floorNumber <= len(d.Floors) {
		d.CurrentFloor = floorNumber - 1
		return true
	}
	return false
}

// GetProgress retourne la progression dans le donjon (pourcentage)
func (d *Dungeon) GetProgress() float64 {
	if d.TotalFloors == 0 {
		return 0
	}
	return float64(d.CurrentFloor+1) / float64(d.TotalFloors) * 100
}

// hashString crée un hash simple d'une chaîne
func hashString(s string) int {
	hash := 0
	for i := 0; i < len(s); i++ {
		hash = hash*31 + int(s[i])
	}
	return hash
}

// applyThemeToTilemap applique un thème visuel à une tilemap
func applyThemeToTilemap(tm *TileMap, theme ZoneTheme) *TileMap {
	switch theme {
	case ThemeForest:
		// Forêt plus sombre et verdoyante
		for y := 0; y < tm.Height; y++ {
			for x := 0; x < tm.Width; x++ {
				tile := &tm.Tiles[y][x]
				if tile.Type == TileGrass {
					tile.Color.R = tile.Color.R / 2
					tile.Color.G = uint8(float64(tile.Color.G) * 0.8)
					tile.Color.B = tile.Color.B / 2
				}
			}
		}

	case ThemeCastle:
		// Plus de pierres, tons gris
		for y := 0; y < tm.Height; y++ {
			for x := 0; x < tm.Width; x++ {
				tile := &tm.Tiles[y][x]
				avg := (uint16(tile.Color.R) + uint16(tile.Color.G) + uint16(tile.Color.B)) / 3
				tile.Color.R = uint8(avg)
				tile.Color.G = uint8(avg)
				tile.Color.B = uint8(avg) + 20
			}
		}

	case ThemeCatacombs:
		// Sombre, tons marrons/noirs
		for y := 0; y < tm.Height; y++ {
			for x := 0; x < tm.Width; x++ {
				tile := &tm.Tiles[y][x]
				tile.Color.R = tile.Color.R / 3
				tile.Color.G = tile.Color.G / 3
				tile.Color.B = tile.Color.B / 3
			}
		}

	case ThemeRuins:
		// Mélange de pierre et végétation envahissante
		for y := 0; y < tm.Height; y++ {
			for x := 0; x < tm.Width; x++ {
				tile := &tm.Tiles[y][x]
				tile.Color.R = uint8(float64(tile.Color.R) * 0.7)
				tile.Color.G = uint8(float64(tile.Color.G) * 0.9)
				tile.Color.B = uint8(float64(tile.Color.B) * 0.7)
			}
		}
	}

	return tm
}
