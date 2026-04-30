package ui

import (
	"e-soulslike/entities"
	"e-soulslike/settings"
	"e-soulslike/stats"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// MenuState représente l'état actuel du menu
type MenuState int

const (
	MenuStateMain MenuState = iota
	MenuStateStats
	MenuStateOptions
	MenuStateTutorial
	MenuStateLevelUp
)

// PauseMenu représente le menu pause complet
type PauseMenu struct {
	screenWidth  int
	screenHeight int
	currentState MenuState

	// Boutons du menu principal
	mainButtons []*Button

	// Sous-menus
	statsDisplay *StatsDisplay
	optionsMenu  *OptionsMenu
	tutorial     *Tutorial
	levelUpMenu  *LevelUpMenu

	// Références
	gameStats   *stats.GameStats
	settings    *settings.Settings
	playerStats *entities.PlayerStats

	// Callbacks
	OnResume     func()
	OnQuickSave  func()
	OnQuitToMenu func()
}

// NewPauseMenu crée un nouveau menu pause
func NewPauseMenu(screenWidth, screenHeight int, gameStats *stats.GameStats, gameSettings *settings.Settings, playerStats *entities.PlayerStats) *PauseMenu {
	pm := &PauseMenu{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		currentState: MenuStateMain,
		gameStats:    gameStats,
		settings:     gameSettings,
		playerStats:  playerStats,
		mainButtons:  make([]*Button, 0),
	}

	// Position du menu
	centerX := float32(screenWidth / 2)
	centerY := float32(screenHeight / 2)
	menuWidth := float32(400)

	// Créer les boutons du menu principal
	buttonY := centerY - 80
	buttonWidth := float32(250)
	buttonHeight := float32(40)

	// Bouton: Reprendre
	pm.mainButtons = append(pm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY, buttonWidth, buttonHeight,
		"Reprendre (ECHAP)",
		func() {
			if pm.OnResume != nil {
				pm.OnResume()
			}
		},
	))

	// Bouton: Statistiques
	pm.mainButtons = append(pm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+50, buttonWidth, buttonHeight,
		"Statistiques",
		func() {
			pm.currentState = MenuStateStats
		},
	))

	// Bouton "Montée de niveau" retiré - disponible uniquement aux feux de camp

	// Bouton: Options
	pm.mainButtons = append(pm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+100, buttonWidth, buttonHeight,
		"Options",
		func() {
			pm.currentState = MenuStateOptions
		},
	))

	// Bouton: Tutoriel
	pm.mainButtons = append(pm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+150, buttonWidth, buttonHeight,
		"Tutoriel",
		func() {
			pm.currentState = MenuStateTutorial
		},
	))

	// Bouton: Sauvegarde rapide
	pm.mainButtons = append(pm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+200, buttonWidth, buttonHeight,
		"Sauvegarde Rapide",
		func() {
			if pm.OnQuickSave != nil {
				pm.OnQuickSave()
			}
		},
	))

	// Bouton: Quitter vers Menu
	pm.mainButtons = append(pm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+250, buttonWidth, buttonHeight,
		"Quitter vers Menu",
		func() {
			if pm.OnQuitToMenu != nil {
				pm.OnQuitToMenu()
			}
		},
	))

	// Créer les sous-menus
	menuX := centerX - menuWidth/2
	menuY := float32(50)
	menuHeight := float32(screenHeight - 100)

	pm.statsDisplay = NewStatsDisplay(menuX, menuY, menuWidth, menuHeight, gameStats)

	pm.optionsMenu = NewOptionsMenu(menuX, menuY, menuWidth, menuHeight, gameSettings, func() {
		pm.currentState = MenuStateMain
	})

	pm.tutorial = NewTutorial(menuX, menuY, menuWidth, menuHeight, func() {
		pm.currentState = MenuStateMain
	})

	pm.levelUpMenu = NewLevelUpMenu(menuX, menuY, menuWidth, menuHeight, playerStats, func() {
		pm.currentState = MenuStateMain
	})

	return pm
}

// Update met à jour le menu pause
func (pm *PauseMenu) Update(mouseX, mouseY int, mousePressed, mouseJustPressed bool) {
	switch pm.currentState {
	case MenuStateMain:
		for _, button := range pm.mainButtons {
			button.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
		}

	case MenuStateStats:
		// Pas de boutons dans l'affichage des stats, juste Échap pour sortir
		// (géré dans game.go)

	case MenuStateOptions:
		pm.optionsMenu.Update(mouseX, mouseY, mousePressed, mouseJustPressed)

	case MenuStateTutorial:
		pm.tutorial.Update(mouseX, mouseY, mousePressed, mouseJustPressed)

	case MenuStateLevelUp:
		pm.levelUpMenu.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
	}
}

// Draw dessine le menu pause
func (pm *PauseMenu) Draw(screen *ebiten.Image) {
	// Overlay semi-transparent (fond sombre)
	overlay := color.RGBA{0, 0, 0, 180}
	vector.DrawFilledRect(screen, 0, 0, float32(pm.screenWidth), float32(pm.screenHeight), overlay, false)

	// Dessiner selon l'état actuel
	switch pm.currentState {
	case MenuStateMain:
		pm.drawMainMenu(screen)

	case MenuStateStats:
		pm.statsDisplay.Draw(screen)
		pm.drawBackHint(screen)

	case MenuStateOptions:
		pm.optionsMenu.Draw(screen)

	case MenuStateTutorial:
		pm.tutorial.Draw(screen)

	case MenuStateLevelUp:
		pm.levelUpMenu.Draw(screen)
	}
}

// drawMainMenu dessine le menu principal
func (pm *PauseMenu) drawMainMenu(screen *ebiten.Image) {
	centerX := float32(pm.screenWidth / 2)
	centerY := float32(pm.screenHeight / 2)

	// Panneau du menu
	menuWidth := float32(400)
	menuHeight := float32(450)
	menuX := centerX - menuWidth/2
	menuY := centerY - menuHeight/2

	// Fond du menu
	menuBg := color.RGBA{40, 40, 50, 255}
	vector.DrawFilledRect(screen, menuX, menuY, menuWidth, menuHeight, menuBg, false)

	// Bordure du menu
	borderColor := color.RGBA{100, 150, 255, 255}
	vector.StrokeRect(screen, menuX, menuY, menuWidth, menuHeight, 3, borderColor, false)

	// Titre "PAUSE"
	titleText := "=== PAUSE ==="
	titleX := int(centerX) - len(titleText)*6/2
	titleY := int(menuY) + 30
	ebitenutil.DebugPrintAt(screen, titleText, titleX, titleY)

	// Indicateur visuel "PAUSE" clignotant en haut
	if ebiten.TPS() == 0 || (ebiten.ActualTPS() > 0 && int(ebiten.ActualTPS()/2)%2 == 0) {
		pauseIndicator := "|| PAUSE ||"
		indicatorX := int(centerX) - len(pauseIndicator)*6/2
		indicatorY := 30
		ebitenutil.DebugPrintAt(screen, pauseIndicator, indicatorX, indicatorY)
	}

	// Dessiner les boutons
	for _, button := range pm.mainButtons {
		button.Draw(screen)
	}
}

// drawBackHint affiche un indice pour revenir en arrière
func (pm *PauseMenu) drawBackHint(screen *ebiten.Image) {
	hint := "Appuyez sur ECHAP pour revenir"
	hintX := pm.screenWidth/2 - len(hint)*6/2
	hintY := pm.screenHeight - 30
	ebitenutil.DebugPrintAt(screen, hint, hintX, hintY)
}

// HandleEscape gère la touche Échap dans les sous-menus
func (pm *PauseMenu) HandleEscape() {
	if pm.currentState != MenuStateMain {
		pm.currentState = MenuStateMain
	}
}

// GetCurrentState retourne l'état actuel du menu
func (pm *PauseMenu) GetCurrentState() MenuState {
	return pm.currentState
}

// GetLevelUpMenu retourne le menu de level up
func (pm *PauseMenu) GetLevelUpMenu() *LevelUpMenu {
	return pm.levelUpMenu
}
