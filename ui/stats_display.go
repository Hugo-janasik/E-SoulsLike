package ui

import (
	"e-soulslike/stats"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// StatsDisplay affiche les statistiques du jeu
type StatsDisplay struct {
	X, Y          float32
	Width, Height float32
	Stats         *stats.GameStats
}

// NewStatsDisplay crée un nouvel affichage de stats
func NewStatsDisplay(x, y, width, height float32, gameStats *stats.GameStats) *StatsDisplay {
	return &StatsDisplay{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Stats:  gameStats,
	}
}

// Draw dessine l'affichage des statistiques
func (sd *StatsDisplay) Draw(screen *ebiten.Image) {
	// Fond
	bgColor := color.RGBA{30, 30, 40, 240}
	vector.DrawFilledRect(screen, sd.X, sd.Y, sd.Width, sd.Height, bgColor, false)

	// Bordure
	borderColor := color.RGBA{100, 150, 255, 255}
	vector.StrokeRect(screen, sd.X, sd.Y, sd.Width, sd.Height, 2, borderColor, false)

	// Titre
	titleX := int(sd.X) + 10
	titleY := int(sd.Y) + 10
	ebitenutil.DebugPrintAt(screen, "=== STATISTIQUES ===", titleX, titleY)

	// Mettre à jour le temps de jeu
	sd.Stats.UpdatePlayTime()

	// Statistiques
	startY := titleY + 30
	lineHeight := 18

	lines := []string{
		fmt.Sprintf("Temps de jeu: %dm %ds",
			int(sd.Stats.PlayTime.Minutes()),
			int(sd.Stats.PlayTime.Seconds())%60),
		"",
		"=== COMBAT ===",
		fmt.Sprintf("Ennemis tues: %d", sd.Stats.EnemiesKilled),
		fmt.Sprintf("Degats infliges: %d", sd.Stats.TotalDamageDealt),
		fmt.Sprintf("Degats recus: %d", sd.Stats.TotalDamageTaken),
		fmt.Sprintf("Precision: %.1f%%", sd.Stats.GetAccuracy()),
		fmt.Sprintf("Degats moy/coup: %.1f", sd.Stats.GetAverageDamagePerHit()),
		"",
		"=== ACTIONS ===",
		fmt.Sprintf("Attaques reussies: %d", sd.Stats.AttacksLanded),
		fmt.Sprintf("Attaques ratees: %d", sd.Stats.AttacksMissed),
		fmt.Sprintf("Esquives: %d", sd.Stats.DodgesPerformed),
		"",
		"=== MOUVEMENT ===",
		fmt.Sprintf("Distance: %.0fm", sd.Stats.DistanceTraveled),
		fmt.Sprintf("Stamina utilisee: %d", sd.Stats.StaminaUsed),
	}

	for i, line := range lines {
		y := startY + i*lineHeight
		ebitenutil.DebugPrintAt(screen, line, titleX, y)
	}
}
