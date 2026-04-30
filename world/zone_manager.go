package world

// ZoneManager gère toutes les zones du jeu
type ZoneManager struct {
	zones       map[string]*Zone
	currentZone *Zone
}

// NewZoneManager crée un nouveau gestionnaire de zones
func NewZoneManager() *ZoneManager {
	zm := &ZoneManager{
		zones: make(map[string]*Zone),
	}

	// Créer les zones vides (entités ajoutées depuis game.go)
	zm.createEmptyZones()

	return zm
}

// createEmptyZones crée la structure des zones sans entités
func (zm *ZoneManager) createEmptyZones() {
	// ===== ZONE 1: Place du Feu (Hub central) =====
	firelink := NewZone(
		"firelink",
		"Place du Feu",
		"Un refuge sûr au cœur des ténèbres",
		100, // Augmenté de 60 à 100
		100, // Augmenté de 60 à 100
		ThemeDefault,
	)

	// Points de spawn (centre de la nouvelle map 100x100)
	centerX := float64(100 * 32 / 2) // 1600
	centerY := float64(100 * 32 / 2) // 1600
	firelink.AddSpawnPoint("main", centerX, centerY)
	firelink.AddSpawnPoint("from_dungeon", centerX, centerY) // Retour de donjon

	// Portails vers les 3 donjons (repositionnés pour la nouvelle taille)
	firelink.AddPortal(NewPortal(centerX-600, centerY, "forest_dungeon", "main", "Donjon Facile"))
	firelink.AddPortal(NewPortal(centerX, centerY-600, "ruins_dungeon", "main", "Donjon Intermediaire"))
	firelink.AddPortal(NewPortal(centerX+600, centerY, "abyss_dungeon", "main", "Donjon Difficile"))

	zm.AddZone(firelink)

	// Définir la zone de départ
	zm.currentZone = firelink
}

// AddZone ajoute une zone au gestionnaire
func (zm *ZoneManager) AddZone(zone *Zone) {
	zm.zones[zone.ID] = zone
}

// GetZone retourne une zone par son ID
func (zm *ZoneManager) GetZone(id string) *Zone {
	return zm.zones[id]
}

// GetCurrentZone retourne la zone actuelle
func (zm *ZoneManager) GetCurrentZone() *Zone {
	return zm.currentZone
}

// ChangeZone change la zone actuelle
func (zm *ZoneManager) ChangeZone(zoneID string) bool {
	if zone, exists := zm.zones[zoneID]; exists {
		zm.currentZone = zone
		return true
	}
	return false
}

// Update met à jour la zone actuelle
func (zm *ZoneManager) Update(playerX, playerY float64) {
	if zm.currentZone == nil {
		return
	}

	// Mettre à jour les feux de camp
	for _, campfire := range zm.currentZone.Campfires {
		// Type assertion pour appeler Update
		if cf, ok := campfire.(interface{ Update(float64, float64) }); ok {
			cf.Update(playerX, playerY)
		}
	}

	// Mettre à jour les portails
	for _, portal := range zm.currentZone.Portals {
		portal.Update(playerX, playerY)
	}

	// Mettre à jour les ennemis
	for _, enemy := range zm.currentZone.Enemies {
		// Type assertion pour appeler Update
		if en, ok := enemy.(interface{ Update(float64, float64) }); ok {
			en.Update(playerX, playerY)
		}
	}
}

// CheckPortalInteraction vérifie si le joueur peut interagir avec un portail
func (zm *ZoneManager) CheckPortalInteraction() *Portal {
	if zm.currentZone == nil {
		return nil
	}

	for _, portal := range zm.currentZone.Portals {
		if portal.CanInteract() {
			return portal
		}
	}

	return nil
}

// CheckCampfireInteraction vérifie si le joueur peut interagir avec un feu de camp
func (zm *ZoneManager) CheckCampfireInteraction() interface{} {
	if zm.currentZone == nil {
		return nil
	}

	for _, campfire := range zm.currentZone.Campfires {
		// Type assertion pour appeler CanInteract
		if cf, ok := campfire.(interface{ CanInteract() bool }); ok {
			if cf.CanInteract() {
				return campfire
			}
		}
	}

	return nil
}
