package ui

import (
	"e-soulslike/settings"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// OptionsMenu affiche le menu des options
type OptionsMenu struct {
	X, Y          float32
	Width, Height float32
	Settings      *settings.Settings
	Buttons       []*Button
	OnClose       func()
}

// NewOptionsMenu crée un nouveau menu d'options
func NewOptionsMenu(x, y, width, height float32, gameSettings *settings.Settings, onClose func()) *OptionsMenu {
	om := &OptionsMenu{
		X:        x,
		Y:        y,
		Width:    width,
		Height:   height,
		Settings: gameSettings,
		OnClose:  onClose,
		Buttons:  make([]*Button, 0),
	}

	// Créer les boutons
	buttonY := y + 200

	// Bouton: Toggle Audio
	om.Buttons = append(om.Buttons, NewButton(
		x+20, buttonY, 200, 30,
		"Audio: ON/OFF",
		func() {
			om.Settings.AudioEnabled = !om.Settings.AudioEnabled
		},
	))

	// Bouton: Volume +
	om.Buttons = append(om.Buttons, NewButton(
		x+20, buttonY+40, 95, 30,
		"Volume +",
		func() {
			om.Settings.MasterVolume += 0.1
			if om.Settings.MasterVolume > 1.0 {
				om.Settings.MasterVolume = 1.0
			}
		},
	))

	// Bouton: Volume -
	om.Buttons = append(om.Buttons, NewButton(
		x+125, buttonY+40, 95, 30,
		"Volume -",
		func() {
			om.Settings.MasterVolume -= 0.1
			if om.Settings.MasterVolume < 0.0 {
				om.Settings.MasterVolume = 0.0
			}
		},
	))

	// Bouton: Sauvegarder les options
	om.Buttons = append(om.Buttons, NewButton(
		x+20, buttonY+120, 200, 30,
		"Sauvegarder",
		func() {
			om.Settings.SaveToFile("settings.json")
		},
	))

	// Bouton: Fermer
	om.Buttons = append(om.Buttons, NewButton(
		x+20, buttonY+160, 200, 30,
		"Retour",
		om.OnClose,
	))

	return om
}

// Update met à jour le menu des options
func (om *OptionsMenu) Update(mouseX, mouseY int, mousePressed, mouseJustPressed bool) {
	for _, button := range om.Buttons {
		button.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
	}
}

// Draw dessine le menu des options
func (om *OptionsMenu) Draw(screen *ebiten.Image) {
	// Fond
	bgColor := color.RGBA{30, 30, 40, 240}
	vector.DrawFilledRect(screen, om.X, om.Y, om.Width, om.Height, bgColor, false)

	// Bordure
	borderColor := color.RGBA{100, 150, 255, 255}
	vector.StrokeRect(screen, om.X, om.Y, om.Width, om.Height, 2, borderColor, false)

	// Titre
	titleX := int(om.X) + 10
	titleY := int(om.Y) + 10
	ebitenutil.DebugPrintAt(screen, "=== OPTIONS ===", titleX, titleY)

	// Afficher les paramètres actuels
	startY := titleY + 40
	lines := []string{
		fmt.Sprintf("Audio: %s", map[bool]string{true: "Active", false: "Desactive"}[om.Settings.AudioEnabled]),
		fmt.Sprintf("Volume maitre: %.0f%%", om.Settings.MasterVolume*100),
		fmt.Sprintf("Volume musique: %.0f%%", om.Settings.MusicVolume*100),
		fmt.Sprintf("Volume effets: %.0f%%", om.Settings.SFXVolume*100),
		"",
		"Ajustez les parametres ci-dessous:",
	}

	for i, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, titleX, startY+i*18)
	}

	// Dessiner les boutons
	for _, button := range om.Buttons {
		button.Draw(screen)
	}
}
