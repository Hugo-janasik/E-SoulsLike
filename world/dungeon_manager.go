package world

// DungeonManager gère tous les donjons du jeu
type DungeonManager struct {
	dungeons        map[string]*Dungeon
	currentDungeon  *Dungeon
	isInDungeon     bool
}

// NewDungeonManager crée un nouveau gestionnaire de donjons
func NewDungeonManager() *DungeonManager {
	dm := &DungeonManager{
		dungeons:    make(map[string]*Dungeon),
		isInDungeon: false,
	}

	// Créer les 3 donjons
	dm.createDungeons()

	return dm
}

// createDungeons crée les 3 donjons du jeu
func (dm *DungeonManager) createDungeons() {
	// Donjon Facile : 5 étages
	easyDungeon := NewDungeon(
		"forest_dungeon",
		"Forêt des Âmes Perdues",
		"Un donjon mystérieux rempli d'esprits perdus",
		DifficultyEasy,
		5,
	)
	easyDungeon.CreateFloors()
	dm.AddDungeon(easyDungeon)

	// Donjon Intermédiaire : 10 étages
	mediumDungeon := NewDungeon(
		"ruins_dungeon",
		"Ruines de l'Ancien Royaume",
		"Des ruines anciennes cachant des secrets dangereux",
		DifficultyMedium,
		10,
	)
	mediumDungeon.CreateFloors()
	dm.AddDungeon(mediumDungeon)

	// Donjon Difficile : 20 étages
	hardDungeon := NewDungeon(
		"abyss_dungeon",
		"Abîme des Ténèbres",
		"Les profondeurs les plus sombres et dangereuses",
		DifficultyHard,
		20,
	)
	hardDungeon.CreateFloors()
	dm.AddDungeon(hardDungeon)
}

// AddDungeon ajoute un donjon au gestionnaire
func (dm *DungeonManager) AddDungeon(dungeon *Dungeon) {
	dm.dungeons[dungeon.ID] = dungeon
}

// GetDungeon retourne un donjon par son ID
func (dm *DungeonManager) GetDungeon(id string) *Dungeon {
	return dm.dungeons[id]
}

// GetCurrentDungeon retourne le donjon actuel
func (dm *DungeonManager) GetCurrentDungeon() *Dungeon {
	return dm.currentDungeon
}

// IsInDungeon retourne true si le joueur est dans un donjon
func (dm *DungeonManager) IsInDungeon() bool {
	return dm.isInDungeon
}

// EnterDungeon fait entrer le joueur dans un donjon
func (dm *DungeonManager) EnterDungeon(dungeonID string) (*Floor, bool) {
	dungeon, exists := dm.dungeons[dungeonID]
	if !exists {
		return nil, false
	}

	dm.currentDungeon = dungeon
	dm.isInDungeon = true
	dm.currentDungeon.CurrentFloor = 0 // Commencer au premier étage

	return dm.currentDungeon.GetCurrentFloor(), true
}

// ExitDungeon fait sortir le joueur du donjon actuel
func (dm *DungeonManager) ExitDungeon() {
	dm.currentDungeon = nil
	dm.isInDungeon = false
}

// GetCurrentFloor retourne l'étage actuel du donjon actuel
func (dm *DungeonManager) GetCurrentFloor() *Floor {
	if dm.currentDungeon == nil {
		return nil
	}
	return dm.currentDungeon.GetCurrentFloor()
}

// ChangeFloor change d'étage dans le donjon actuel
func (dm *DungeonManager) ChangeFloor(floorNumber int) (*Floor, bool) {
	if dm.currentDungeon == nil {
		return nil, false
	}

	if dm.currentDungeon.SetCurrentFloor(floorNumber) {
		return dm.currentDungeon.GetCurrentFloor(), true
	}

	return nil, false
}

// GetCurrentZone retourne la zone de l'étage actuel
func (dm *DungeonManager) GetCurrentZone() *Zone {
	floor := dm.GetCurrentFloor()
	if floor == nil {
		return nil
	}
	return floor.Zone
}

// Update met à jour le donjon actuel
func (dm *DungeonManager) Update(playerX, playerY float64) {
	if !dm.isInDungeon || dm.currentDungeon == nil {
		return
	}

	zone := dm.GetCurrentZone()
	if zone == nil {
		return
	}

	// Mettre à jour les feux de camp
	for _, campfire := range zone.Campfires {
		if cf, ok := campfire.(interface{ Update(float64, float64) }); ok {
			cf.Update(playerX, playerY)
		}
	}

	// Mettre à jour les portails
	for _, portal := range zone.Portals {
		portal.Update(playerX, playerY)
	}

	// Mettre à jour les ennemis
	for _, enemy := range zone.Enemies {
		if en, ok := enemy.(interface{ Update(float64, float64) }); ok {
			en.Update(playerX, playerY)
		}
	}
}

// CheckPortalInteraction vérifie si le joueur peut interagir avec un portail
func (dm *DungeonManager) CheckPortalInteraction() *Portal {
	zone := dm.GetCurrentZone()
	if zone == nil {
		return nil
	}

	for _, portal := range zone.Portals {
		if portal.CanInteract() {
			return portal
		}
	}

	return nil
}

// CheckCampfireInteraction vérifie si le joueur peut interagir avec un feu de camp
func (dm *DungeonManager) CheckCampfireInteraction() interface{} {
	zone := dm.GetCurrentZone()
	if zone == nil {
		return nil
	}

	for _, campfire := range zone.Campfires {
		if cf, ok := campfire.(interface{ CanInteract() bool }); ok {
			if cf.CanInteract() {
				return campfire
			}
		}
	}

	return nil
}

// GetAllDungeons retourne la liste de tous les donjons
func (dm *DungeonManager) GetAllDungeons() []*Dungeon {
	dungeons := make([]*Dungeon, 0, len(dm.dungeons))
	for _, dungeon := range dm.dungeons {
		dungeons = append(dungeons, dungeon)
	}
	return dungeons
}

// GetDungeonInfo retourne les informations d'un donjon pour l'affichage
func (dm *DungeonManager) GetDungeonInfo(dungeonID string) (name string, floors int, difficulty string, exists bool) {
	dungeon, exists := dm.dungeons[dungeonID]
	if !exists {
		return "", 0, "", false
	}

	difficultyStr := ""
	switch dungeon.Difficulty {
	case DifficultyEasy:
		difficultyStr = "Facile"
	case DifficultyMedium:
		difficultyStr = "Intermédiaire"
	case DifficultyHard:
		difficultyStr = "Difficile"
	}

	return dungeon.Name, dungeon.TotalFloors, difficultyStr, true
}
