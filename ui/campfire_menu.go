package ui

import (
	"e-soulslike/entities"
	"e-soulslike/stats"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// CampfireMenuState représente l'état du menu feu de camp
type CampfireMenuState int

const (
	CampfireMenuStateMain CampfireMenuState = iota
	CampfireMenuStateLevelUp
	CampfireMenuStateShop
)

// CampfireMenu représente le menu du feu de camp
type CampfireMenu struct {
	screenWidth  int
	screenHeight int
	currentState CampfireMenuState

	// Boutons du menu principal
	mainButtons []*Button

	// Sous-menus
	levelUpMenu *LevelUpMenu
	shopMenu    *ShopMenu

	// Références
	gameStats   *stats.GameStats
	playerStats *entities.PlayerStats
	player      *entities.Player // Ajouté pour le shop

	// Callbacks
	OnClose func()
	OnRest  func()
	OnSave  func()
}

// NewCampfireMenu crée un nouveau menu de feu de camp
func NewCampfireMenu(screenWidth, screenHeight int, gameStats *stats.GameStats, player *entities.Player) *CampfireMenu {
	cm := &CampfireMenu{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		currentState: CampfireMenuStateMain,
		gameStats:    gameStats,
		playerStats:  player.Stats,
		player:       player,
		mainButtons:  make([]*Button, 0),
	}

	// Position du menu
	centerX := float32(screenWidth / 2)
	centerY := float32(screenHeight / 2)
	menuWidth := float32(400)

	// Créer les boutons du menu principal
	buttonY := centerY - 60
	buttonWidth := float32(250)
	buttonHeight := float32(40)

	// Bouton: Se reposer (restaure santé/stamina)
	cm.mainButtons = append(cm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY, buttonWidth, buttonHeight,
		"Se reposer",
		func() {
			// Restaurer santé/stamina puis fermer
			if cm.OnRest != nil {
				cm.OnRest()
			}
			if cm.OnClose != nil {
				cm.OnClose()
			}
		},
	))

	// Bouton: Monter de niveau
	cm.mainButtons = append(cm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+50, buttonWidth, buttonHeight,
		"Montee de Niveau",
		func() {
			cm.currentState = CampfireMenuStateLevelUp
		},
	))

	// Bouton: Magasin
	cm.mainButtons = append(cm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+100, buttonWidth, buttonHeight,
		"Magasin",
		func() {
			cm.currentState = CampfireMenuStateShop
		},
	))

	// Bouton: Sauvegarder
	cm.mainButtons = append(cm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+150, buttonWidth, buttonHeight,
		"Sauvegarder",
		func() {
			if cm.OnSave != nil {
				cm.OnSave()
			}
		},
	))

	// Bouton: Quitter
	cm.mainButtons = append(cm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+200, buttonWidth, buttonHeight,
		"Quitter",
		func() {
			if cm.OnClose != nil {
				cm.OnClose()
			}
		},
	))

	// Créer le menu de level up
	menuX := centerX - menuWidth/2
	menuY := float32(50)
	menuHeight := float32(screenHeight - 100)

	cm.levelUpMenu = NewLevelUpMenu(menuX, menuY, menuWidth, menuHeight, player.Stats, func() {
		cm.currentState = CampfireMenuStateMain
	})

	// Créer le menu de shop
	cm.shopMenu = NewShopMenu(menuX, menuY, menuWidth, menuHeight, player, func() {
		cm.currentState = CampfireMenuStateMain
	})

	return cm
}

// Update met à jour le menu
func (cm *CampfireMenu) Update(mouseX, mouseY int, mousePressed, mouseJustPressed bool) {
	switch cm.currentState {
	case CampfireMenuStateMain:
		for _, button := range cm.mainButtons {
			button.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
		}

	case CampfireMenuStateLevelUp:
		cm.levelUpMenu.Update(mouseX, mouseY, mousePressed, mouseJustPressed)

	case CampfireMenuStateShop:
		cm.shopMenu.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
	}
}

// Draw dessine le menu
func (cm *CampfireMenu) Draw(screen *ebiten.Image) {
	// Overlay semi-transparent
	overlay := color.RGBA{0, 0, 0, 180}
	vector.DrawFilledRect(screen, 0, 0, float32(cm.screenWidth), float32(cm.screenHeight), overlay, false)

	// Dessiner selon l'état actuel
	switch cm.currentState {
	case CampfireMenuStateMain:
		cm.drawMainMenu(screen)

	case CampfireMenuStateLevelUp:
		cm.levelUpMenu.Draw(screen)

	case CampfireMenuStateShop:
		cm.shopMenu.Draw(screen)
	}
}

// drawMainMenu dessine le menu principal du feu de camp
func (cm *CampfireMenu) drawMainMenu(screen *ebiten.Image) {
	centerX := float32(cm.screenWidth / 2)
	centerY := float32(cm.screenHeight / 2)

	// Panneau du menu (augmenté pour le bouton Magasin)
	menuWidth := float32(400)
	menuHeight := float32(400)
	menuX := centerX - menuWidth/2
	menuY := centerY - menuHeight/2

	// Fond du menu
	menuBg := color.RGBA{40, 40, 50, 255}
	vector.DrawFilledRect(screen, menuX, menuY, menuWidth, menuHeight, menuBg, false)

	// Bordure du menu
	borderColor := color.RGBA{255, 150, 50, 255} // Orange pour rappeler le feu
	vector.StrokeRect(screen, menuX, menuY, menuWidth, menuHeight, 3, borderColor, false)

	// Titre "FEU DE CAMP"
	titleText := "=== FEU DE CAMP ==="
	titleX := int(centerX) - len(titleText)*6/2
	titleY := int(menuY) + 30
	ebitenutil.DebugPrintAt(screen, titleText, titleX, titleY)

	// Dessiner les boutons
	for _, button := range cm.mainButtons {
		button.Draw(screen)
	}

	// Indice en bas
	hintText := "Appuyez sur E ou ECHAP pour quitter"
	hintX := int(centerX) - len(hintText)*6/2
	hintY := int(menuY + menuHeight) - 20
	ebitenutil.DebugPrintAt(screen, hintText, hintX, hintY)
}

// HandleEscape gère la touche Échap
func (cm *CampfireMenu) HandleEscape() {
	if cm.currentState != CampfireMenuStateMain {
		cm.currentState = CampfireMenuStateMain
	} else {
		if cm.OnClose != nil {
			cm.OnClose()
		}
	}
}

// GetCurrentState retourne l'état actuel du menu
func (cm *CampfireMenu) GetCurrentState() CampfireMenuState {
	return cm.currentState
}

// GetLevelUpMenu retourne le menu de level up
func (cm *CampfireMenu) GetLevelUpMenu() *LevelUpMenu {
	return cm.levelUpMenu
}
