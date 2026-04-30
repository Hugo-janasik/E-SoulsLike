package ui

import (
	"e-soulslike/entities"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// HUD affiche les informations du joueur à l'écran
type HUD struct {
	screenWidth  int
	screenHeight int
}

// NewHUD crée un nouveau HUD
func NewHUD(screenWidth, screenHeight int) *HUD {
	return &HUD{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// Draw dessine le HUD
func (h *HUD) Draw(screen *ebiten.Image, player *entities.Player) {
	// Dessiner santé et stamina en haut à gauche
	h.drawHealthStamina(screen, player)

	// Dessiner les âmes en haut à droite
	h.drawSouls(screen, player)
}

// drawHealthStamina dessine les barres de vie et stamina
func (h *HUD) drawHealthStamina(screen *ebiten.Image, player *entities.Player) {
	// Position
	x := float32(20)
	y := float32(50)
	barWidth := float32(200)
	barHeight := float32(20)

	// Fond pour les barres
	bgColor := color.RGBA{30, 30, 30, 200}
	vector.DrawFilledRect(screen, x-5, y-5, barWidth+10, barHeight*2+15, bgColor, false)

	// Barre de santé
	healthLabel := fmt.Sprintf("HP: %d/%d", player.Health, player.MaxHealth)
	ebitenutil.DebugPrintAt(screen, healthLabel, int(x), int(y)-15)

	// Fond de la barre de santé
	vector.DrawFilledRect(screen, x, y, barWidth, barHeight, color.RGBA{50, 50, 50, 255}, false)
	// Vie actuelle
	healthPercent := float32(player.Health) / float32(player.MaxHealth)
	if healthPercent < 0 {
		healthPercent = 0
	}
	healthWidth := barWidth * healthPercent
	healthColor := color.RGBA{220, 50, 50, 255}
	vector.DrawFilledRect(screen, x, y, healthWidth, barHeight, healthColor, false)
	// Bordure
	vector.StrokeRect(screen, x, y, barWidth, barHeight, 2, color.RGBA{100, 100, 100, 255}, false)

	// Barre de stamina
	staminaY := y + barHeight + 5
	staminaLabel := fmt.Sprintf("Stamina: %d/%d", player.Stamina, player.MaxStamina)
	ebitenutil.DebugPrintAt(screen, staminaLabel, int(x), int(staminaY)-15)

	// Fond de la barre de stamina
	vector.DrawFilledRect(screen, x, staminaY, barWidth, barHeight, color.RGBA{50, 50, 50, 255}, false)
	// Stamina actuelle
	staminaPercent := float32(player.Stamina) / float32(player.MaxStamina)
	if staminaPercent < 0 {
		staminaPercent = 0
	}
	staminaWidth := barWidth * staminaPercent
	staminaColor := color.RGBA{100, 220, 100, 255}
	vector.DrawFilledRect(screen, x, staminaY, staminaWidth, barHeight, staminaColor, false)
	// Bordure
	vector.StrokeRect(screen, x, staminaY, barWidth, barHeight, 2, color.RGBA{100, 100, 100, 255}, false)
}

// drawSouls dessine le compteur d'âmes
func (h *HUD) drawSouls(screen *ebiten.Image, player *entities.Player) {
	// Position en haut à droite
	soulsText := fmt.Sprintf("Ames: %d", player.Stats.Souls)
	textWidth := len(soulsText) * 6 // Approximation de la largeur
	x := h.screenWidth - textWidth - 20
	y := 50

	// Fond
	bgColor := color.RGBA{30, 30, 30, 200}
	vector.DrawFilledRect(screen, float32(x-10), float32(y-10), float32(textWidth+20), 30, bgColor, false)

	// Texte
	ebitenutil.DebugPrintAt(screen, soulsText, x, y)

	// Niveau en dessous
	levelText := fmt.Sprintf("Niveau: %d", player.Stats.Level)
	ebitenutil.DebugPrintAt(screen, levelText, x, y+15)
}
