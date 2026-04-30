package world

// Zone représente une zone du jeu (comme dans Dark Souls)
type Zone struct {
	ID          string
	Name        string
	Description string
	TileMap     *TileMap

	// Points de spawn pour le joueur dans cette zone
	SpawnPoints map[string]SpawnPoint

	// Entités de la zone (interface{} pour éviter cycle d'import)
	Enemies   []interface{}
	Campfires []interface{}
	Portals   []*Portal

	// Style visuel de la zone
	Theme ZoneTheme
}

// SpawnPoint représente un point d'apparition du joueur
type SpawnPoint struct {
	X float64
	Y float64
}

// ZoneTheme définit le style visuel d'une zone
type ZoneTheme int

const (
	ThemeDefault ZoneTheme = iota
	ThemeForest            // Forêt sombre
	ThemeCastle            // Château
	ThemeCatacombs         // Catacombes
	ThemeRuins             // Ruines
	ThemeHub               // Lobby principal (Place du Feu)
)

// NewZone crée une nouvelle zone
func NewZone(id, name, description string, width, height int, theme ZoneTheme) *Zone {
	return &Zone{
		ID:          id,
		Name:        name,
		Description: description,
		TileMap:     NewThemedTileMap(width, height, theme),
		SpawnPoints: make(map[string]SpawnPoint),
		Enemies:     make([]interface{}, 0),
		Campfires:   make([]interface{}, 0),
		Portals:     make([]*Portal, 0),
		Theme:       theme,
	}
}

// AddSpawnPoint ajoute un point de spawn
func (z *Zone) AddSpawnPoint(name string, x, y float64) {
	z.SpawnPoints[name] = SpawnPoint{X: x, Y: y}
}

// GetSpawnPoint retourne un point de spawn
func (z *Zone) GetSpawnPoint(name string) (float64, float64, bool) {
	if sp, exists := z.SpawnPoints[name]; exists {
		return sp.X, sp.Y, true
	}
	// Point de spawn par défaut au centre
	return float64(z.TileMap.Width * TileSize / 2), float64(z.TileMap.Height * TileSize / 2), false
}

// AddEnemy ajoute un ennemi à la zone
func (z *Zone) AddEnemy(enemy interface{}) {
	z.Enemies = append(z.Enemies, enemy)
}

// AddCampfire ajoute un feu de camp à la zone
func (z *Zone) AddCampfire(campfire interface{}) {
	z.Campfires = append(z.Campfires, campfire)
}

// AddPortal ajoute un portail à la zone
func (z *Zone) AddPortal(portal *Portal) {
	z.Portals = append(z.Portals, portal)
}

// NewThemedTileMap crée une tilemap avec un thème spécifique
func NewThemedTileMap(width, height int, theme ZoneTheme) *TileMap {
	// Le hub utilise son propre générateur de tilemap
	if theme == ThemeHub {
		return GenerateHubTilemap(width, height)
	}

	tm := NewTileMap(width, height)

	// Pour le thème par défaut (zone hub), l'eau est traversable
	if theme == ThemeDefault {
		tm.WalkableWater = true
	}

	// Modifier les couleurs selon le thème
	switch theme {
	case ThemeForest:
		// Forêt plus sombre et verdoyante
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				tile := &tm.Tiles[y][x]
				if tile.Type == TileGrass {
					// Herbe plus foncée
					tile.Color.R = tile.Color.R / 2
					tile.Color.G = uint8(float64(tile.Color.G) * 0.8)
					tile.Color.B = tile.Color.B / 2
				}
			}
		}

	case ThemeCastle:
		// Plus de pierres, tons gris
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				tile := &tm.Tiles[y][x]
				// Tout devient plus gris/pierre
				avg := (uint16(tile.Color.R) + uint16(tile.Color.G) + uint16(tile.Color.B)) / 3
				tile.Color.R = uint8(avg)
				tile.Color.G = uint8(avg)
				tile.Color.B = uint8(avg) + 20
			}
		}

	case ThemeCatacombs:
		// Sombre, tons marrons/noirs
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				tile := &tm.Tiles[y][x]
				tile.Color.R = tile.Color.R / 3
				tile.Color.G = tile.Color.G / 3
				tile.Color.B = tile.Color.B / 3
			}
		}

	case ThemeRuins:
		// Mélange de pierre et végétation envahissante
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				tile := &tm.Tiles[y][x]
				tile.Color.R = uint8(float64(tile.Color.R) * 0.7)
				tile.Color.G = uint8(float64(tile.Color.G) * 0.9)
				tile.Color.B = uint8(float64(tile.Color.B) * 0.7)
			}
		}
	}

	return tm
}
