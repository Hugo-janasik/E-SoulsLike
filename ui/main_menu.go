package ui

import (
	"e-soulslike/settings"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// MainMenuState représente l'état du menu principal
type MainMenuState int

const (
	MainMenuStateMain MainMenuState = iota
	MainMenuStateOptions
)

// MainMenu représente le menu principal du jeu
type MainMenu struct {
	screenWidth  int
	screenHeight int
	currentState MainMenuState

	// Boutons du menu principal
	mainButtons []*Button

	// Sous-menus
	optionsMenu *OptionsMenu

	// Références
	settings *settings.Settings

	// Callbacks
	OnNewGame  func()
	OnLoadGame func()
	OnQuit     func()
}

// NewMainMenu crée un nouveau menu principal
func NewMainMenu(screenWidth, screenHeight int, gameSettings *settings.Settings) *MainMenu {
	mm := &MainMenu{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		currentState: MainMenuStateMain,
		settings:     gameSettings,
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

	// Bouton: Nouvelle Partie
	mm.mainButtons = append(mm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY, buttonWidth, buttonHeight,
		"Nouvelle Partie",
		func() {
			if mm.OnNewGame != nil {
				mm.OnNewGame()
			}
		},
	))

	// Bouton: Charger Partie
	mm.mainButtons = append(mm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+50, buttonWidth, buttonHeight,
		"Charger Partie",
		func() {
			if mm.OnLoadGame != nil {
				mm.OnLoadGame()
			}
		},
	))

	// Bouton: Options
	mm.mainButtons = append(mm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+100, buttonWidth, buttonHeight,
		"Options",
		func() {
			mm.currentState = MainMenuStateOptions
		},
	))

	// Bouton: Quitter
	mm.mainButtons = append(mm.mainButtons, NewButton(
		centerX-buttonWidth/2, buttonY+150, buttonWidth, buttonHeight,
		"Quitter",
		func() {
			if mm.OnQuit != nil {
				mm.OnQuit()
			}
		},
	))

	// Créer le menu options
	menuX := centerX - menuWidth/2
	menuY := float32(50)
	menuHeight := float32(screenHeight - 100)

	mm.optionsMenu = NewOptionsMenu(menuX, menuY, menuWidth, menuHeight, gameSettings, func() {
		mm.currentState = MainMenuStateMain
	})

	return mm
}

// Update met à jour le menu principal
func (mm *MainMenu) Update(mouseX, mouseY int, mousePressed, mouseJustPressed bool) {
	switch mm.currentState {
	case MainMenuStateMain:
		for _, button := range mm.mainButtons {
			button.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
		}

	case MainMenuStateOptions:
		mm.optionsMenu.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
	}
}

// Draw dessine le menu principal
func (mm *MainMenu) Draw(screen *ebiten.Image) {
	// Fond dégradé
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Dessiner selon l'état actuel
	switch mm.currentState {
	case MainMenuStateMain:
		mm.drawMainMenu(screen)

	case MainMenuStateOptions:
		mm.optionsMenu.Draw(screen)
		mm.drawBackHint(screen)
	}
}

// drawMainMenu dessine le menu principal
func (mm *MainMenu) drawMainMenu(screen *ebiten.Image) {
	centerX := float32(mm.screenWidth / 2)
	centerY := float32(mm.screenHeight / 2)

	// Panneau du menu
	menuWidth := float32(400)
	menuHeight := float32(400)
	menuX := centerX - menuWidth/2
	menuY := centerY - menuHeight/2

	// Fond du menu
	menuBg := color.RGBA{40, 40, 50, 255}
	vector.DrawFilledRect(screen, menuX, menuY, menuWidth, menuHeight, menuBg, false)

	// Bordure du menu
	borderColor := color.RGBA{100, 150, 255, 255}
	vector.StrokeRect(screen, menuX, menuY, menuWidth, menuHeight, 3, borderColor, false)

	// Titre "E-SOULSLIKE"
	titleText := "=== E-SOULSLIKE ==="
	titleX := int(centerX) - len(titleText)*6/2
	titleY := int(menuY) + 30
	ebitenutil.DebugPrintAt(screen, titleText, titleX, titleY)

	// Sous-titre
	subtitleText := "Action RPG 2D"
	subtitleX := int(centerX) - len(subtitleText)*6/2
	subtitleY := titleY + 20
	ebitenutil.DebugPrintAt(screen, subtitleText, subtitleX, subtitleY)

	// Dessiner les boutons
	for _, button := range mm.mainButtons {
		button.Draw(screen)
	}

	// Crédits en bas
	creditsText := "Appuyez ECHAP pour revenir"
	creditsX := int(centerX) - len(creditsText)*6/2
	creditsY := int(menuY+menuHeight) - 20
	ebitenutil.DebugPrintAt(screen, creditsText, creditsX, creditsY)
}

// drawBackHint affiche un indice pour revenir en arrière
func (mm *MainMenu) drawBackHint(screen *ebiten.Image) {
	hint := "Appuyez sur ECHAP pour revenir"
	hintX := mm.screenWidth/2 - len(hint)*6/2
	hintY := mm.screenHeight - 30
	ebitenutil.DebugPrintAt(screen, hint, hintX, hintY)
}

// HandleEscape gère la touche Échap dans les sous-menus
func (mm *MainMenu) HandleEscape() {
	if mm.currentState != MainMenuStateMain {
		mm.currentState = MainMenuStateMain
	}
}

// GetCurrentState retourne l'état actuel du menu
func (mm *MainMenu) GetCurrentState() MainMenuState {
	return mm.currentState
}
