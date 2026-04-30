package game

import (
	"e-soulslike/combat"
	"e-soulslike/entities"
	"e-soulslike/input"
	"e-soulslike/save"
	"e-soulslike/settings"
	"e-soulslike/stats"
	"e-soulslike/ui"
	"e-soulslike/world"
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// GameState représente l'état actuel du jeu
type GameState int

const (
	GameStateMainMenu GameState = iota
	GameStatePlaying
	GameStateDead // Joueur mort — affiche l'écran YOU DIED
)

// Game représente l'état principal du jeu
type Game struct {
	width            int
	height           int
	currentState     GameState
	player           *entities.Player
	zoneManager      *world.ZoneManager
	dungeonManager   *world.DungeonManager
	camera           *world.Camera
	inputMgr         *input.InputManager
	combatSystem     *combat.CombatSystem
	mainMenu         *ui.MainMenu
	pauseMenu        *ui.PauseMenu
	campfireMenu     *ui.CampfireMenu
	inventoryMenu    *ui.InventoryMenu
	hud              *ui.HUD
	isPaused         bool
	isAtCampfire     bool
	isInInventory    bool
	gameStats        *stats.GameStats
	settings         *settings.Settings
	saveMessage      string
	saveMessageTimer int

	// Transition entre zones
	isTransitioning    bool
	transitionTimer    int
	transitionDuration int
	transitionAlpha    float64
	nextZoneID         string
	nextSpawnPoint     string

	// Mort du joueur
	deathScreen       *ui.DeathScreen
	respawnZoneID     string // zone du dernier feu de camp visité
	respawnSpawnPoint string // spawn point associé

	// Effet torche dans les donjons
	vignette  *world.Vignette
	particles *world.AmbientParticles
}

// NewGame crée une nouvelle instance du jeu
func NewGame(width, height int) *Game {
	g := &Game{
		width:              width,
		height:             height,
		currentState:       GameStateMainMenu,
		inputMgr:           input.NewInputManager(),
		transitionDuration: 30, // 0.5 secondes à 60 FPS
	}

	// Charger ou créer les paramètres
	var err error
	g.settings, err = settings.LoadSettingsFromFile("settings.json")
	if err != nil {
		g.settings = settings.NewDefaultSettings()
	}

	// Initialiser les statistiques
	g.gameStats = stats.NewGameStats()

	// Initialiser le HUD
	g.hud = ui.NewHUD(width, height)

	// Initialiser le menu principal
	g.mainMenu = ui.NewMainMenu(width, height, g.settings)

	// Définir les callbacks du menu principal
	g.mainMenu.OnNewGame = func() {
		g.startNewGame()
	}
	g.mainMenu.OnLoadGame = func() {
		g.loadGame()
	}
	g.mainMenu.OnQuit = func() {
		os.Exit(0)
	}

	// Initialiser le menu pause (sera utilisé pendant le jeu)
	emptyStats := entities.NewPlayerStats()
	g.pauseMenu = ui.NewPauseMenu(width, height, g.gameStats, g.settings, emptyStats)
	g.isPaused = false

	// Définir les callbacks du menu pause
	g.pauseMenu.OnResume = func() {
		g.isPaused = false
	}
	g.pauseMenu.OnQuickSave = func() {
		g.quickSave()
	}
	g.pauseMenu.OnQuitToMenu = func() {
		g.quitToMainMenu()
	}

	// Initialiser le menu d'inventaire (sera mis à jour quand le joueur est créé)
	g.inventoryMenu = ui.NewInventoryMenu(width, height, nil, func() {
		g.isInInventory = false
	})
	g.isInInventory = false

	return g
}

// Update met à jour la logique du jeu (appelé 60 fois par seconde)
func (g *Game) Update() error {
	// Mettre à jour l'input manager (toujours nécessaire)
	g.inputMgr.Update()

	// Si on est dans le menu principal
	if g.currentState == GameStateMainMenu {
		// Gérer Échap dans le menu principal
		if g.inputMgr.IsKeyJustPressed(ebiten.KeyEscape) {
			if g.mainMenu.GetCurrentState() != ui.MainMenuStateMain {
				g.mainMenu.HandleEscape()
			}
		}

		// Mettre à jour le menu principal
		mouseX, mouseY := g.inputMgr.GetMousePosition()
		mousePressed := g.inputMgr.IsMousePressed()
		mouseJustPressed := g.inputMgr.IsMouseJustPressed()
		g.mainMenu.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
		return nil
	}

	// Écran de mort : n'importe quelle touche pour ressusciter
	if g.currentState == GameStateDead {
		g.deathScreen.Update()
		if g.deathScreen.ReadyForInput {
			anyKey := g.inputMgr.IsKeyJustPressed(ebiten.KeySpace) ||
				g.inputMgr.IsKeyJustPressed(ebiten.KeyEnter) ||
				g.inputMgr.IsKeyJustPressed(ebiten.KeyEscape) ||
				g.inputMgr.IsMouseJustPressed()
			if anyKey {
				g.respawnPlayer()
			}
		}
		return nil
	}

	// Si on est en jeu
	if g.currentState == GameStatePlaying {
		// Gérer les transitions de zone
		if g.isTransitioning {
			g.updateTransition()
			return nil
		}

		// Gérer la touche Échap pour pause/reprendre (sauf si l'inventaire est ouvert)
		if g.inputMgr.IsKeyJustPressed(ebiten.KeyEscape) && !g.isInInventory {
			if g.isAtCampfire {
				g.campfireMenu.HandleEscape()
			} else if g.isPaused {
				if g.pauseMenu.GetCurrentState() != ui.MenuStateMain {
					g.pauseMenu.HandleEscape()
				} else {
					g.isPaused = false
				}
			} else {
				g.isPaused = true
			}
		}

		// Gérer l'interaction avec les feux de camp et portails (touche E)
		if g.inputMgr.IsKeyJustPressed(ebiten.KeyE) && !g.isPaused && !g.isAtCampfire {
			var portal *world.Portal
			var campfire interface{}

			// Vérifier les interactions selon où on est
			if g.dungeonManager.IsInDungeon() {
				portal = g.dungeonManager.CheckPortalInteraction()
				campfire = g.dungeonManager.CheckCampfireInteraction()
			} else {
				portal = g.zoneManager.CheckPortalInteraction()
				campfire = g.zoneManager.CheckCampfireInteraction()
			}

			// Vérifier interaction avec portail
			if portal != nil {
				g.startZoneTransition(portal.DestinationZoneID, portal.DestinationSpawn)
			} else if campfire != nil {
				// Vérifier interaction avec feu de camp
				g.isAtCampfire = true
			}
		}

		// Si au feu de camp, mettre à jour le menu feu de camp
		if g.isAtCampfire {
			mouseX, mouseY := g.inputMgr.GetMousePosition()
			mousePressed := g.inputMgr.IsMousePressed()
			mouseJustPressed := g.inputMgr.IsMouseJustPressed()
			g.campfireMenu.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
			return nil
		}

		// Si en pause, mettre à jour le menu avec la souris
		if g.isPaused {
			mouseX, mouseY := g.inputMgr.GetMousePosition()
			mousePressed := g.inputMgr.IsMousePressed()
			mouseJustPressed := g.inputMgr.IsMouseJustPressed()
			g.pauseMenu.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
			return nil
		}

		// Gérer l'ouverture/fermeture de l'inventaire avec 'I' ou 'B'
		if (g.inputMgr.IsKeyJustPressed(ebiten.KeyI) || g.inputMgr.IsKeyJustPressed(ebiten.KeyB)) && !g.isPaused && !g.isAtCampfire {
			g.isInInventory = !g.isInInventory
		}

		// Si l'inventaire est ouvert, mettre à jour le menu d'inventaire
		if g.isInInventory {
			g.inventoryMenu.Update(g.inputMgr)
			return nil
		}

		// Mettre à jour le zone manager ou dungeon manager selon où on est
		if g.dungeonManager.IsInDungeon() {
			g.dungeonManager.Update(g.player.X, g.player.Y)
			if g.vignette != nil {
				g.vignette.Update()
			}
			// Incrémenter le tick de la tilemap (animations de torches)
			if currentZone := g.dungeonManager.GetCurrentZone(); currentZone != nil {
				currentZone.TileMap.Tick++
			}
			// Mettre à jour les particules d'ambiance
			if g.particles != nil {
				g.particles.Update(g.camera)
			}
		} else {
			g.zoneManager.Update(g.player.X, g.player.Y)
		}

		// Récupérer la tilemap actuelle pour les collisions
		var currentTilemap *world.TileMap
		if g.dungeonManager.IsInDungeon() {
			currentZone := g.dungeonManager.GetCurrentZone()
			if currentZone != nil {
				currentTilemap = currentZone.TileMap
			}
		} else {
			currentZone := g.zoneManager.GetCurrentZone()
			if currentZone != nil {
				currentTilemap = currentZone.TileMap
			}
		}

		// Tracker l'esquive pour les stats
		wasDodging := g.player.IsDodging

		// Mettre à jour le joueur avec la tilemap pour les collisions
		g.player.Update(g.inputMgr, currentTilemap)

		// Enregistrer l'esquive si elle vient de démarrer
		if !wasDodging && g.player.IsDodging {
			g.gameStats.RecordDodge()
		}

		// Mettre à jour la caméra pour suivre le joueur
		g.camera.Follow(g.player.X, g.player.Y)

		// Mettre à jour le système de combat
		g.combatSystem.Update()

		// Détecter la mort du joueur après la frame de combat
		if !g.player.IsAlive() {
			g.triggerDeath()
		}
	}

	return nil
}

// triggerDeath déclenche la séquence de mort du joueur
func (g *Game) triggerDeath() {
	// Sauvegarder les âmes perdues et les réinitialiser à 0
	soulsLost := g.player.Stats.Souls
	g.player.Stats.Souls = 0

	// Créer l'écran de mort
	g.deathScreen = ui.NewDeathScreen(g.width, g.height, soulsLost)
	g.currentState = GameStateDead
	g.isPaused = false
	g.isAtCampfire = false
	g.isInInventory = false
}

// respawnPlayer téléporte le joueur au dernier feu de camp et réinitialise les ennemis
func (g *Game) respawnPlayer() {
	// Sortir du donjon si nécessaire
	if g.dungeonManager.IsInDungeon() {
		g.dungeonManager.ExitDungeon()
		g.vignette = nil
		g.particles = nil
	}

	// Changer de zone vers le dernier feu de camp
	if g.respawnZoneID != "" {
		g.zoneManager.ChangeZone(g.respawnZoneID)
	}

	// Téléporter le joueur au spawn point du feu de camp
	spawnPoint := g.respawnSpawnPoint
	if spawnPoint == "" {
		spawnPoint = "main"
	}
	spawnX, spawnY, _ := g.zoneManager.GetCurrentZone().GetSpawnPoint(spawnPoint)
	g.player.X = spawnX
	g.player.Y = spawnY

	// Restaurer la vie et la stamina
	g.player.Health = g.player.MaxHealth
	g.player.Stamina = g.player.MaxStamina
	g.player.IsInvincible = false
	g.player.InvincibleTimer = 0

	// Réinitialiser les ennemis de la zone
	g.respawnEnemiesInCurrentZone()

	// Mettre à jour le système de combat avec la nouvelle zone
	currentZone := g.zoneManager.GetCurrentZone()
	enemies := make([]*entities.Enemy, 0)
	for _, e := range currentZone.Enemies {
		if enemy, ok := e.(*entities.Enemy); ok {
			enemies = append(enemies, enemy)
		}
	}
	g.combatSystem = combat.NewCombatSystem(g.player, enemies, g.gameStats)

	// Recentrer la caméra
	g.camera.Follow(g.player.X, g.player.Y)

	// Retourner en jeu
	g.currentState = GameStatePlaying
	g.saveMessage = "Ressuscite au dernier feu de camp"
	g.saveMessageTimer = 120
}

// respawnEnemiesInCurrentZone réinitialise tous les ennemis de la zone courante
func (g *Game) respawnEnemiesInCurrentZone() {
	var currentZone *world.Zone
	if g.dungeonManager.IsInDungeon() {
		currentZone = g.dungeonManager.GetCurrentZone()
	} else {
		currentZone = g.zoneManager.GetCurrentZone()
	}
	if currentZone == nil {
		return
	}
	for _, e := range currentZone.Enemies {
		if enemy, ok := e.(*entities.Enemy); ok {
			enemy.Respawn()
		}
	}
}

// updateTransition gère l'animation de transition entre zones
func (g *Game) updateTransition() {
	g.transitionTimer++

	// Première moitié: fade to black
	if g.transitionTimer < g.transitionDuration/2 {
		g.transitionAlpha = float64(g.transitionTimer) / float64(g.transitionDuration/2)
	} else if g.transitionTimer == g.transitionDuration/2 {
		// Milieu de la transition: changer de zone
		g.changeZone(g.nextZoneID, g.nextSpawnPoint)
		g.transitionAlpha = 1.0
	} else if g.transitionTimer < g.transitionDuration {
		// Deuxième moitié: fade from black
		remaining := g.transitionDuration - g.transitionTimer
		g.transitionAlpha = float64(remaining) / float64(g.transitionDuration/2)
	} else {
		// Fin de la transition
		g.isTransitioning = false
		g.transitionTimer = 0
		g.transitionAlpha = 0
	}
}

// startZoneTransition démarre une transition vers une nouvelle zone
func (g *Game) startZoneTransition(zoneID, spawnPoint string) {
	g.isTransitioning = true
	g.transitionTimer = 0
	g.nextZoneID = zoneID
	g.nextSpawnPoint = spawnPoint
}

// changeZone change la zone actuelle (gère zones normales et donjons)
func (g *Game) changeZone(zoneID, spawnPoint string) {
	// Garde défensive : ne pas continuer si le joueur n'existe pas encore
	if g.player == nil {
		return
	}

	// Vérifier si c'est un donjon
	dungeon := g.dungeonManager.GetDungeon(zoneID)
	if dungeon != nil {
		// Entrer dans le donjon
		floor, success := g.dungeonManager.EnterDungeon(zoneID)
		if success && floor != nil {
			// Téléporter le joueur au spawn point
			spawnX, spawnY, _ := floor.Zone.GetSpawnPoint(spawnPoint)
			g.player.X = spawnX
			g.player.Y = spawnY

			// Mettre à jour le système de combat avec les nouveaux ennemis
			enemies := make([]*entities.Enemy, 0)
			for _, e := range floor.Zone.Enemies {
				if enemy, ok := e.(*entities.Enemy); ok {
					enemies = append(enemies, enemy)
				}
			}
			g.combatSystem = combat.NewCombatSystem(g.player, enemies, g.gameStats)

			// Créer la vignette torche et les particules adaptées au thème du donjon
			g.vignette = world.NewVignette(g.width, g.height, floor.Zone.Theme)
			g.particles = world.NewAmbientParticles(floor.Zone.Theme)

			// Message de transition
			g.saveMessage = fmt.Sprintf(">>> %s - Étage 1/%d <<<", dungeon.Name, dungeon.TotalFloors)
			g.saveMessageTimer = 120 // 2 secondes
		}
		return
	}

	// Vérifier si c'est un étage de donjon
	if g.dungeonManager.IsInDungeon() {
		currentDungeon := g.dungeonManager.GetCurrentDungeon()
		if currentDungeon != nil {
			// Chercher l'étage dans le donjon actuel
			for i, floor := range currentDungeon.Floors {
				if floor.Zone.ID == zoneID {
					// Changer d'étage
					currentDungeon.CurrentFloor = i

					// Téléporter le joueur au spawn point
					spawnX, spawnY, _ := floor.Zone.GetSpawnPoint(spawnPoint)
					g.player.X = spawnX
					g.player.Y = spawnY

					// Mettre à jour le système de combat avec les nouveaux ennemis
					enemies := make([]*entities.Enemy, 0)
					for _, e := range floor.Zone.Enemies {
						if enemy, ok := e.(*entities.Enemy); ok {
							enemies = append(enemies, enemy)
						}
					}
					g.combatSystem = combat.NewCombatSystem(g.player, enemies, g.gameStats)

					// Message de transition
					floorText := ""
					if floor.IsRestArea {
						floorText = "Salle de Repos"
					} else {
						floorText = fmt.Sprintf("Étage %d/%d", floor.Number, currentDungeon.TotalFloors)
					}
					g.saveMessage = fmt.Sprintf(">>> %s <<<", floorText)
					g.saveMessageTimer = 120 // 2 secondes
					return
				}
			}
		}

		// Si on arrive ici et que la destination est "firelink", sortir du donjon
		if zoneID == "firelink" {
			g.dungeonManager.ExitDungeon()
			g.vignette = nil
			g.particles = nil
		}
	}

	// Zone normale (pas un donjon)
	if g.zoneManager.ChangeZone(zoneID) {
		// Téléporter le joueur au spawn point
		spawnX, spawnY, _ := g.zoneManager.GetCurrentZone().GetSpawnPoint(spawnPoint)
		g.player.X = spawnX
		g.player.Y = spawnY

		// Mettre à jour le système de combat avec les nouveaux ennemis
		currentZone := g.zoneManager.GetCurrentZone()
		enemies := make([]*entities.Enemy, 0)
		for _, e := range currentZone.Enemies {
			if enemy, ok := e.(*entities.Enemy); ok {
				enemies = append(enemies, enemy)
			}
		}
		g.combatSystem = combat.NewCombatSystem(g.player, enemies, g.gameStats)

		// Message de transition
		g.saveMessage = fmt.Sprintf(">>> %s <<<", currentZone.Name)
		g.saveMessageTimer = 120 // 2 secondes
	}
}

// Draw dessine le jeu à l'écran (appelé après Update)
func (g *Game) Draw(screen *ebiten.Image) {
	// Si on est dans le menu principal
	if g.currentState == GameStateMainMenu {
		g.mainMenu.Draw(screen)
		return
	}

	// Écran de mort
	if g.currentState == GameStateDead {
		// Fond noir de base puis overlay du death screen par-dessus
		screen.Fill(color.RGBA{0, 0, 0, 255})
		g.deathScreen.Draw(screen)
		return
	}

	// Si on est en jeu
	if g.currentState == GameStatePlaying {
		// Fond noir dans les donjons, gris foncé ailleurs
		if g.dungeonManager.IsInDungeon() {
			screen.Fill(color.RGBA{0, 0, 0, 255})
		} else {
			screen.Fill(color.RGBA{40, 40, 45, 255})
		}

		// Récupérer la zone actuelle (donjon ou zone normale)
		var currentZone *world.Zone
		if g.dungeonManager.IsInDungeon() {
			currentZone = g.dungeonManager.GetCurrentZone()
		} else {
			currentZone = g.zoneManager.GetCurrentZone()
		}

		if currentZone != nil {
			// Dessiner la map avec décalage de caméra
			currentZone.TileMap.Draw(screen, g.camera)

			// Dessiner les portails
			for _, portal := range currentZone.Portals {
				portal.Draw(screen, g.camera)
			}

			// Dessiner les feux de camp
			for _, campfire := range currentZone.Campfires {
				if cf, ok := campfire.(*entities.Campfire); ok {
					cf.Draw(screen, g.camera)
				} else if fg, ok := campfire.(*entities.FireGuardian); ok {
					fg.Draw(screen, g.camera)
				}
			}

			// Dessiner les ennemis
			for _, enemy := range currentZone.Enemies {
				if en, ok := enemy.(*entities.Enemy); ok {
					en.Draw(screen, g.camera)
				}
			}
		}

		// Dessiner le joueur
		g.player.Draw(screen, g.camera)

		// Particules d'ambiance (par-dessus entités, sous la vignette)
		if g.particles != nil {
			g.particles.Draw(screen, g.camera)
		}

		// Effet torche/vignette (par-dessus le monde, sous le HUD)
		if g.vignette != nil {
			playerScreenX, playerScreenY := g.camera.WorldToScreen(g.player.X, g.player.Y)
			g.vignette.Draw(screen, playerScreenX, playerScreenY)
		}

		// Dessiner le HUD (santé, stamina, âmes)
		g.hud.Draw(screen, g.player)

		// Afficher le nom de la zone en haut à gauche
		if currentZone != nil {
			zoneText := currentZone.Name
			ebitenutil.DebugPrintAt(screen, zoneText, 20, 120)
		}

		// Afficher les informations de debug
		debugText := "E-SoulsLike - ZQSD/WASD: Bouger | Espace: Attaquer | Shift: Sprint/Esquive | E: Interagir | Echap: Pause"
		ebitenutil.DebugPrint(screen, debugText)

		// Afficher le message de sauvegarde/transition si actif
		if g.saveMessageTimer > 0 {
			g.saveMessageTimer--
			messageX := g.width/2 - len(g.saveMessage)*6/2
			messageY := g.height / 2
			ebitenutil.DebugPrintAt(screen, g.saveMessage, messageX, messageY)
		}

		// Si au feu de camp, afficher le menu feu de camp
		if g.isAtCampfire {
			g.campfireMenu.Draw(screen)
		}

		// Si le jeu est en pause, afficher le menu pause par-dessus
		if g.isPaused {
			g.pauseMenu.Draw(screen)
		}

		// Si l'inventaire est ouvert, afficher le menu d'inventaire par-dessus
		if g.isInInventory {
			g.inventoryMenu.Draw(screen)
		}

		// Effet de transition (fade to black)
		if g.isTransitioning && g.transitionAlpha > 0 {
			alpha := uint8(g.transitionAlpha * 255)
			overlayColor := color.RGBA{0, 0, 0, alpha}
			ebitenutil.DrawRect(screen, 0, 0, float64(g.width), float64(g.height), overlayColor)
		}
	}
}

// Layout définit la taille logique de l'écran
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

// quickSave crée une sauvegarde rapide
func (g *Game) quickSave() {
	saveData := save.NewSaveData("QuickSave")

	// Remplir les données du joueur
	saveData.PlayerX = g.player.X
	saveData.PlayerY = g.player.Y
	saveData.PlayerHealth = g.player.Health
	saveData.PlayerStamina = g.player.Stamina

	// Sauvegarder les stats et la progression
	saveData.PlayerLevel = g.player.Stats.Level
	saveData.PlayerSouls = g.player.Stats.Souls
	saveData.PlayerForce = g.player.Stats.Force
	saveData.PlayerStaminaStat = g.player.Stats.Stamina
	saveData.PlayerVie = g.player.Stats.Vie

	// Sauvegarder l'équipement
	if g.player.Equipment.Weapon != nil {
		saveData.EquippedWeaponID = g.player.Equipment.Weapon.ID
	}
	if g.player.Equipment.Chest != nil {
		saveData.EquippedChestID = g.player.Equipment.Chest.ID
	}
	if g.player.Equipment.Helmet != nil {
		saveData.EquippedHelmetID = g.player.Equipment.Helmet.ID
	}
	if g.player.Equipment.Boots != nil {
		saveData.EquippedBootsID = g.player.Equipment.Boots.ID
	}
	if g.player.Equipment.Ring != nil {
		saveData.EquippedRingID = g.player.Equipment.Ring.ID
	}
	if g.player.Equipment.Amulet != nil {
		saveData.EquippedAmuletID = g.player.Equipment.Amulet.ID
	}
	if g.player.Equipment.Belt != nil {
		saveData.EquippedBeltID = g.player.Equipment.Belt.ID
	}

	// Sauvegarder l'inventaire
	saveData.InventoryItems = make(map[string]int)
	for _, slot := range g.player.Inventory.Slots {
		if slot != nil && slot.Item != nil {
			// Extraire l'Item de base pour obtenir l'ID
			var itemID string
			switch item := slot.Item.(type) {
			case *entities.Weapon:
				itemID = item.ID
			case *entities.Armor:
				itemID = item.ID
			case *entities.Accessory:
				itemID = item.ID
			case *entities.Item:
				itemID = item.ID
			}
			if itemID != "" {
				saveData.InventoryItems[itemID] = slot.Quantity
			}
		}
	}

	// Remplir les statistiques
	saveData.EnemiesKilled = g.gameStats.EnemiesKilled
	saveData.PlayTime = g.gameStats.PlayTime
	saveData.TotalDamageDealt = g.gameStats.TotalDamageDealt
	saveData.TotalDamageTaken = g.gameStats.TotalDamageTaken

	// Sauvegarder l'état du donjon
	saveData.IsInDungeon = g.dungeonManager.IsInDungeon()
	if saveData.IsInDungeon {
		currentDungeon := g.dungeonManager.GetCurrentDungeon()
		if currentDungeon != nil {
			saveData.CurrentDungeonID = currentDungeon.ID
			saveData.CurrentFloor = currentDungeon.CurrentFloor + 1 // +1 car CurrentFloor est 0-based
		}
	}

	// Sauvegarder
	err := save.QuickSave(saveData)
	if err != nil {
		g.saveMessage = fmt.Sprintf("Erreur sauvegarde: %v", err)
	} else {
		g.saveMessage = "Sauvegarde reussie!"
	}
	g.saveMessageTimer = 180 // 3 secondes à 60 FPS
}

// populateZones ajoute les entités dans les zones
func (g *Game) populateZones() {
	// Place du Feu - Hub central
	firelink := g.zoneManager.GetZone("firelink")
	firelink.AddCampfire(entities.NewFireGuardian(1600, 1740)) // Décalé vers le bas du centre (1600, 1600)
}

// isFarEnoughFromSpawn vérifie si une position est assez loin du spawn point
func isFarEnoughFromSpawn(x, y, spawnX, spawnY float64, minDistance float64) bool {
	dx := x - spawnX
	dy := y - spawnY
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance >= minDistance
}

// populateDungeons ajoute les entités dans les donjons de manière procédurale
func (g *Game) populateDungeons() {
	const minSpawnDistance = 400.0 // Distance minimale en pixels entre le spawn et les ennemis

	// Donjon Facile (5 étages)
	easyDungeon := g.dungeonManager.GetDungeon("forest_dungeon")
	if easyDungeon != nil {
		for _, floor := range easyDungeon.Floors {
			if floor.IsRestArea && floor.Generator != nil {
				// Salle de repos : ajouter un Fire Guardian à une position aléatoire dans la salle
				tileX, tileY := floor.Generator.GetRandomFloorPosition()
				x := float64(tileX * world.TileSize)
				y := float64(tileY * world.TileSize)
				floor.Zone.AddCampfire(entities.NewFireGuardian(x, y))
			} else if floor.Generator != nil {
				// Étage de combat : placer des ennemis aléatoirement dans les salles
				// Obtenir la position de spawn du joueur pour cet étage
				spawnX, spawnY, _ := floor.Zone.GetSpawnPoint("main")

				numEnemies := 2 + floor.Number // Plus d'ennemis à mesure qu'on monte
				spawned := 0
				maxAttempts := numEnemies * 10 // Éviter une boucle infinie

				for attempts := 0; attempts < maxAttempts && spawned < numEnemies; attempts++ {
					tileX, tileY := floor.Generator.GetRandomFloorPosition()
					x := float64(tileX * world.TileSize)
					y := float64(tileY * world.TileSize)

					// Vérifier que l'ennemi n'est pas trop proche du spawn
					if isFarEnoughFromSpawn(x, y, spawnX, spawnY, minSpawnDistance) {
						floor.Zone.AddEnemy(entities.NewEnemy(x, y, entities.EnemyTypeBasic))
						spawned++
					}
				}
			}
		}
	}

	// Donjon Intermédiaire (10 étages)
	mediumDungeon := g.dungeonManager.GetDungeon("ruins_dungeon")
	if mediumDungeon != nil {
		for _, floor := range mediumDungeon.Floors {
			if floor.IsRestArea && floor.Generator != nil {
				// Salle de repos : ajouter un Fire Guardian à une position aléatoire dans la salle
				tileX, tileY := floor.Generator.GetRandomFloorPosition()
				x := float64(tileX * world.TileSize)
				y := float64(tileY * world.TileSize)
				floor.Zone.AddCampfire(entities.NewFireGuardian(x, y))
			} else if floor.Generator != nil {
				// Étage de combat : placer des ennemis aléatoirement
				spawnX, spawnY, _ := floor.Zone.GetSpawnPoint("main")

				numEnemies := 3 + floor.Number
				spawned := 0
				maxAttempts := numEnemies * 10

				for attempts := 0; attempts < maxAttempts && spawned < numEnemies; attempts++ {
					tileX, tileY := floor.Generator.GetRandomFloorPosition()
					x := float64(tileX * world.TileSize)
					y := float64(tileY * world.TileSize)

					if isFarEnoughFromSpawn(x, y, spawnX, spawnY, minSpawnDistance) {
						floor.Zone.AddEnemy(entities.NewEnemy(x, y, entities.EnemyTypeBasic))
						spawned++
					}
				}
			}
		}
	}

	// Donjon Difficile (20 étages)
	hardDungeon := g.dungeonManager.GetDungeon("abyss_dungeon")
	if hardDungeon != nil {
		for _, floor := range hardDungeon.Floors {
			if floor.IsRestArea && floor.Generator != nil {
				// Salle de repos : ajouter un Fire Guardian à une position aléatoire dans la salle
				tileX, tileY := floor.Generator.GetRandomFloorPosition()
				x := float64(tileX * world.TileSize)
				y := float64(tileY * world.TileSize)
				floor.Zone.AddCampfire(entities.NewFireGuardian(x, y))
			} else if floor.Generator != nil {
				// Étage de combat : placer beaucoup d'ennemis aléatoirement
				spawnX, spawnY, _ := floor.Zone.GetSpawnPoint("main")

				numEnemies := 4 + floor.Number
				spawned := 0
				maxAttempts := numEnemies * 10

				for attempts := 0; attempts < maxAttempts && spawned < numEnemies; attempts++ {
					tileX, tileY := floor.Generator.GetRandomFloorPosition()
					x := float64(tileX * world.TileSize)
					y := float64(tileY * world.TileSize)

					if isFarEnoughFromSpawn(x, y, spawnX, spawnY, minSpawnDistance) {
						floor.Zone.AddEnemy(entities.NewEnemy(x, y, entities.EnemyTypeBasic))
						spawned++
					}
				}
			}
		}
	}
}

// startNewGame démarre une nouvelle partie
func (g *Game) startNewGame() {
	// Initialiser la caméra
	g.camera = world.NewCamera(float64(g.width), float64(g.height))

	// Créer le gestionnaire de zones
	g.zoneManager = world.NewZoneManager()

	// Créer le gestionnaire de donjons
	g.dungeonManager = world.NewDungeonManager()

	// Peupler les zones avec les entités
	g.populateZones()

	// Peupler les donjons avec les entités
	g.populateDungeons()

	// Créer le joueur au spawn point de la zone de départ
	currentZone := g.zoneManager.GetCurrentZone()
	spawnX, spawnY, _ := currentZone.GetSpawnPoint("main")
	g.player = entities.NewPlayer(spawnX, spawnY)

	// Initialiser le système de combat (convertir []interface{} en []*entities.Enemy)
	enemies := make([]*entities.Enemy, 0)
	for _, e := range currentZone.Enemies {
		if enemy, ok := e.(*entities.Enemy); ok {
			enemies = append(enemies, enemy)
		}
	}
	g.combatSystem = combat.NewCombatSystem(g.player, enemies, g.gameStats)

	// Réinitialiser les statistiques
	g.gameStats = stats.NewGameStats()

	// Mettre à jour le menu pause avec les nouvelles stats
	g.pauseMenu = ui.NewPauseMenu(g.width, g.height, g.gameStats, g.settings, g.player.Stats)
	g.pauseMenu.OnResume = func() {
		g.isPaused = false
	}
	g.pauseMenu.OnQuickSave = func() {
		g.quickSave()
	}
	g.pauseMenu.OnQuitToMenu = func() {
		g.quitToMainMenu()
	}
	if g.pauseMenu.GetLevelUpMenu() != nil {
		g.pauseMenu.GetLevelUpMenu().OnStatsChanged = func() {
			g.player.ApplyStatsToPlayer()
		}
	}

	// Initialiser le menu feu de camp
	g.campfireMenu = ui.NewCampfireMenu(g.width, g.height, g.gameStats, g.player)
	g.campfireMenu.OnClose = func() {
		g.isAtCampfire = false
	}
	g.campfireMenu.OnRest = func() {
		g.player.Health = g.player.MaxHealth
		g.player.Stamina = g.player.MaxStamina
		// Mémoriser ce feu de camp comme point de respawn
		if g.dungeonManager.IsInDungeon() {
			if d := g.dungeonManager.GetCurrentDungeon(); d != nil {
				g.respawnZoneID = d.Floors[d.CurrentFloor].Zone.ID
			}
		} else {
			if z := g.zoneManager.GetCurrentZone(); z != nil {
				g.respawnZoneID = z.ID
			}
		}
		g.respawnSpawnPoint = "main"
		// Réinitialiser les ennemis (repos au feu de camp = respawn ennemis)
		g.respawnEnemiesInCurrentZone()
	}
	g.campfireMenu.OnSave = func() {
		g.quickSave()
	}
	if g.campfireMenu.GetLevelUpMenu() != nil {
		g.campfireMenu.GetLevelUpMenu().OnStatsChanged = func() {
			g.player.ApplyStatsToPlayer()
		}
	}

	// Initialiser le menu d'inventaire avec le joueur
	g.inventoryMenu = ui.NewInventoryMenu(g.width, g.height, g.player, func() {
		g.isInInventory = false
	})

	// Zone de respawn par défaut = zone de départ
	if startZone := g.zoneManager.GetCurrentZone(); startZone != nil {
		g.respawnZoneID = startZone.ID
	}
	g.respawnSpawnPoint = "main"

	// Passer en mode jeu
	g.currentState = GameStatePlaying
	g.isPaused = false
	g.isAtCampfire = false
	g.isInInventory = false

	// Message de bienvenue
	g.saveMessage = "Bienvenue dans E-SoulsLike!"
	g.saveMessageTimer = 120

	// Sauvegarder immédiatement
	g.quickSave()
}

// loadGame charge la partie sauvegardée
func (g *Game) loadGame() {
	// Charger les données de sauvegarde
	saveData, err := save.LoadQuickSave()
	if err != nil {
		g.startNewGame()
		g.saveMessage = "Aucune sauvegarde trouvee, nouvelle partie!"
		g.saveMessageTimer = 180
		return
	}

	// Initialiser la caméra
	g.camera = world.NewCamera(float64(g.width), float64(g.height))

	// Créer le gestionnaire de zones
	g.zoneManager = world.NewZoneManager()

	// Créer le gestionnaire de donjons
	g.dungeonManager = world.NewDungeonManager()

	// Peupler les zones avec les entités
	g.populateZones()

	// Peupler les donjons avec les entités
	g.populateDungeons()

	// Restaurer l'état du donjon si le joueur était dans un donjon
	if saveData.IsInDungeon && saveData.CurrentDungeonID != "" {
		g.dungeonManager.EnterDungeon(saveData.CurrentDungeonID)
		if saveData.CurrentFloor > 0 {
			g.dungeonManager.ChangeFloor(saveData.CurrentFloor)
		}
	}

	// Créer le joueur avec les données sauvegardées
	g.player = entities.NewPlayer(saveData.PlayerX, saveData.PlayerY)
	g.player.Health = saveData.PlayerHealth
	g.player.Stamina = saveData.PlayerStamina

	// Restaurer les stats et la progression
	if saveData.PlayerLevel > 0 { // Vérifier que les stats sont présentes dans la sauvegarde
		g.player.Stats.Level = saveData.PlayerLevel
		g.player.Stats.Souls = saveData.PlayerSouls
		g.player.Stats.Force = saveData.PlayerForce
		g.player.Stats.Stamina = saveData.PlayerStaminaStat
		g.player.Stats.Vie = saveData.PlayerVie
		g.player.Stats.RecalculateStats()
		// Mettre à jour les max du joueur
		g.player.MaxHealth = int(g.player.Stats.MaxHealth)
		g.player.MaxStamina = int(g.player.Stats.MaxStamina)
	}

	// Restaurer l'équipement
	if saveData.EquippedWeaponID != "" {
		if saveData.EquippedWeaponID == "sword_basic" {
			weapon := entities.NewStartingSword()
			g.player.Equipment.Weapon = weapon
		}
	}
	if saveData.EquippedChestID != "" {
		if saveData.EquippedChestID == "chest_beginner" {
			g.player.Equipment.Chest = entities.NewBeginnerChestplate()
		}
	}
	if saveData.EquippedHelmetID != "" {
		if saveData.EquippedHelmetID == "helmet_beginner" {
			g.player.Equipment.Helmet = entities.NewBeginnerHelmet()
		}
	}
	if saveData.EquippedBootsID != "" {
		if saveData.EquippedBootsID == "boots_beginner" {
			g.player.Equipment.Boots = entities.NewBeginnerBoots()
		}
	}
	if saveData.EquippedRingID != "" {
		if saveData.EquippedRingID == "ring_beginner" {
			g.player.Equipment.Ring = entities.NewBeginnerRing()
		}
	}
	if saveData.EquippedAmuletID != "" {
		if saveData.EquippedAmuletID == "amulet_beginner" {
			g.player.Equipment.Amulet = entities.NewBeginnerAmulet()
		}
	}
	if saveData.EquippedBeltID != "" {
		if saveData.EquippedBeltID == "leggings_beginner" {
			g.player.Equipment.Belt = entities.NewBeginnerLeggings()
		}
	}

	// Recalculer les stats après avoir équipé les items
	g.player.UpdateDerivedStats()

	// Restaurer l'inventaire
	for itemID, quantity := range saveData.InventoryItems {
		// Charger les items selon leur ID
		switch itemID {
		case "sword_basic":
			weapon := entities.NewStartingSword()
			g.player.Inventory.AddItem(weapon, quantity)
		case "helmet_beginner":
			g.player.Inventory.AddItem(entities.NewBeginnerHelmet(), quantity)
		case "chest_beginner":
			g.player.Inventory.AddItem(entities.NewBeginnerChestplate(), quantity)
		case "boots_beginner":
			g.player.Inventory.AddItem(entities.NewBeginnerBoots(), quantity)
		case "amulet_beginner":
			g.player.Inventory.AddItem(entities.NewBeginnerAmulet(), quantity)
		case "ring_beginner":
			g.player.Inventory.AddItem(entities.NewBeginnerRing(), quantity)
		case "leggings_beginner":
			g.player.Inventory.AddItem(entities.NewBeginnerLeggings(), quantity)
		}
	}

	// Initialiser le système de combat
	currentZone := g.zoneManager.GetCurrentZone()
	enemies := make([]*entities.Enemy, 0)
	for _, e := range currentZone.Enemies {
		if enemy, ok := e.(*entities.Enemy); ok {
			enemies = append(enemies, enemy)
		}
	}
	g.combatSystem = combat.NewCombatSystem(g.player, enemies, g.gameStats)

	// Charger les statistiques
	g.gameStats = stats.NewGameStats()
	g.gameStats.EnemiesKilled = saveData.EnemiesKilled
	g.gameStats.PlayTime = saveData.PlayTime
	g.gameStats.TotalDamageDealt = saveData.TotalDamageDealt
	g.gameStats.TotalDamageTaken = saveData.TotalDamageTaken

	// Mettre à jour le menu pause avec les stats chargées
	g.pauseMenu = ui.NewPauseMenu(g.width, g.height, g.gameStats, g.settings, g.player.Stats)
	g.pauseMenu.OnResume = func() {
		g.isPaused = false
	}
	g.pauseMenu.OnQuickSave = func() {
		g.quickSave()
	}
	g.pauseMenu.OnQuitToMenu = func() {
		g.quitToMainMenu()
	}

	// Initialiser le menu feu de camp
	g.campfireMenu = ui.NewCampfireMenu(g.width, g.height, g.gameStats, g.player)
	g.campfireMenu.OnClose = func() {
		g.isAtCampfire = false
	}
	g.campfireMenu.OnRest = func() {
		g.player.Health = g.player.MaxHealth
		g.player.Stamina = g.player.MaxStamina
		// Mémoriser ce feu de camp comme point de respawn
		if g.dungeonManager.IsInDungeon() {
			if d := g.dungeonManager.GetCurrentDungeon(); d != nil {
				g.respawnZoneID = d.Floors[d.CurrentFloor].Zone.ID
			}
		} else {
			if z := g.zoneManager.GetCurrentZone(); z != nil {
				g.respawnZoneID = z.ID
			}
		}
		g.respawnSpawnPoint = "main"
		// Réinitialiser les ennemis (repos au feu de camp = respawn ennemis)
		g.respawnEnemiesInCurrentZone()
	}
	g.campfireMenu.OnSave = func() {
		g.quickSave()
	}

	// Initialiser le menu d'inventaire avec le joueur
	g.inventoryMenu = ui.NewInventoryMenu(g.width, g.height, g.player, func() {
		g.isInInventory = false
	})

	// Zone de respawn par défaut = zone courante au chargement
	if loadedZone := g.zoneManager.GetCurrentZone(); loadedZone != nil {
		g.respawnZoneID = loadedZone.ID
	}
	g.respawnSpawnPoint = "main"

	// Passer en mode jeu
	g.currentState = GameStatePlaying
	g.isPaused = false
	g.isInInventory = false

	// Afficher un message de confirmation
	g.saveMessage = "Partie chargee avec succes!"
	g.saveMessageTimer = 180
}

// quitToMainMenu retourne au menu principal
func (g *Game) quitToMainMenu() {
	// Sauvegarder avant de quitter
	g.quickSave()

	// Retourner au menu principal
	g.currentState = GameStateMainMenu
	g.isPaused = false
}
