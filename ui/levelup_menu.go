package ui

import (
	"e-soulslike/entities"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// LevelUpMenu affiche le menu de montée de niveau avec système de validation
type LevelUpMenu struct {
	X, Y          float32
	Width, Height float32
	PlayerStats   *entities.PlayerStats

	// Stats temporaires (en attente de validation)
	pendingForce   int
	pendingStamina int
	pendingVie     int

	// Stats de base (valeurs minimales)
	baseForce   int
	baseStamina int
	baseVie     int

	Buttons        []*Button
	OnClose        func()
	OnStatsChanged func() // Appelé quand les stats sont validées
}

// NewLevelUpMenu crée un nouveau menu de level up
func NewLevelUpMenu(x, y, width, height float32, playerStats *entities.PlayerStats, onClose func()) *LevelUpMenu {
	lm := &LevelUpMenu{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		PlayerStats: playerStats,
		OnClose:     onClose,
		Buttons:     make([]*Button, 0),
	}

	// Initialiser les stats de base et pending
	lm.resetPendingStats()

	// Créer les boutons +/- pour chaque stat
	lm.createButtons()

	return lm
}

// resetPendingStats réinitialise les stats en attente aux valeurs actuelles
func (lm *LevelUpMenu) resetPendingStats() {
	lm.baseForce = lm.PlayerStats.Force
	lm.baseStamina = lm.PlayerStats.Stamina
	lm.baseVie = lm.PlayerStats.Vie

	lm.pendingForce = lm.PlayerStats.Force
	lm.pendingStamina = lm.PlayerStats.Stamina
	lm.pendingVie = lm.PlayerStats.Vie
}

// createButtons crée tous les boutons du menu
func (lm *LevelUpMenu) createButtons() {
	lm.Buttons = make([]*Button, 0)

	buttonWidth := float32(40)
	buttonHeight := float32(30)
	startY := lm.Y + 220

	// Boutons Force: - et +
	lm.Buttons = append(lm.Buttons, NewButton(
		lm.X+200, startY, buttonWidth, buttonHeight,
		"-",
		func() {
			if lm.pendingForce > lm.baseForce {
				lm.pendingForce--
			}
		},
	))
	lm.Buttons = append(lm.Buttons, NewButton(
		lm.X+250, startY, buttonWidth, buttonHeight,
		"+",
		func() {
			cost := lm.getTotalCost() + lm.PlayerStats.GetSoulsForStatUpgrade()
			if lm.PlayerStats.Souls >= cost {
				lm.pendingForce++
			}
		},
	))

	// Boutons Stamina: - et +
	lm.Buttons = append(lm.Buttons, NewButton(
		lm.X+200, startY+40, buttonWidth, buttonHeight,
		"-",
		func() {
			if lm.pendingStamina > lm.baseStamina {
				lm.pendingStamina--
			}
		},
	))
	lm.Buttons = append(lm.Buttons, NewButton(
		lm.X+250, startY+40, buttonWidth, buttonHeight,
		"+",
		func() {
			cost := lm.getTotalCost() + lm.PlayerStats.GetSoulsForStatUpgrade()
			if lm.PlayerStats.Souls >= cost {
				lm.pendingStamina++
			}
		},
	))

	// Boutons Vie: - et +
	lm.Buttons = append(lm.Buttons, NewButton(
		lm.X+200, startY+80, buttonWidth, buttonHeight,
		"-",
		func() {
			if lm.pendingVie > lm.baseVie {
				lm.pendingVie--
			}
		},
	))
	lm.Buttons = append(lm.Buttons, NewButton(
		lm.X+250, startY+80, buttonWidth, buttonHeight,
		"+",
		func() {
			cost := lm.getTotalCost() + lm.PlayerStats.GetSoulsForStatUpgrade()
			if lm.PlayerStats.Souls >= cost {
				lm.pendingVie++
			}
		},
	))

	// Bouton Valider
	lm.Buttons = append(lm.Buttons, NewButton(
		lm.X+20, startY+140, 160, 35,
		"Valider",
		func() {
			lm.validateStats()
		},
	))

	// Bouton Annuler
	lm.Buttons = append(lm.Buttons, NewButton(
		lm.X+200, startY+140, 160, 35,
		"Annuler",
		func() {
			lm.resetPendingStats()
		},
	))

	// Bouton Retour
	lm.Buttons = append(lm.Buttons, NewButton(
		lm.X+20, startY+190, 340, 35,
		"Retour",
		func() {
			lm.resetPendingStats()
			if lm.OnClose != nil {
				lm.OnClose()
			}
		},
	))
}

// getTotalCost calcule le coût total des améliorations en attente
func (lm *LevelUpMenu) getTotalCost() int {
	forceIncrease := lm.pendingForce - lm.baseForce
	staminaIncrease := lm.pendingStamina - lm.baseStamina
	vieIncrease := lm.pendingVie - lm.baseVie

	totalUpgrades := forceIncrease + staminaIncrease + vieIncrease
	costPerUpgrade := lm.PlayerStats.GetSoulsForStatUpgrade()

	return totalUpgrades * costPerUpgrade
}

// validateStats applique les changements de stats si possible
func (lm *LevelUpMenu) validateStats() {
	totalCost := lm.getTotalCost()

	// Vérifier si le joueur a assez d'âmes
	if lm.PlayerStats.Souls < totalCost {
		return // Pas assez d'âmes
	}

	// Appliquer les changements
	lm.PlayerStats.Force = lm.pendingForce
	lm.PlayerStats.Stamina = lm.pendingStamina
	lm.PlayerStats.Vie = lm.pendingVie

	// Recalculer les stats dérivées
	lm.PlayerStats.RecalculateStats()

	// Déduire les âmes
	lm.PlayerStats.Souls -= totalCost

	// Mettre à jour les stats de base pour les prochaines modifications
	lm.baseForce = lm.pendingForce
	lm.baseStamina = lm.pendingStamina
	lm.baseVie = lm.pendingVie

	// Notifier le changement
	if lm.OnStatsChanged != nil {
		lm.OnStatsChanged()
	}
}

// Update met à jour le menu
func (lm *LevelUpMenu) Update(mouseX, mouseY int, mousePressed, mouseJustPressed bool) {
	for _, button := range lm.Buttons {
		button.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
	}

	// Activer/désactiver le bouton Valider selon les âmes disponibles
	totalCost := lm.getTotalCost()
	hasChanges := lm.pendingForce != lm.baseForce || lm.pendingStamina != lm.baseStamina || lm.pendingVie != lm.baseVie

	// Index 6 = bouton Valider
	if len(lm.Buttons) > 6 {
		lm.Buttons[6].SetEnabled(hasChanges && lm.PlayerStats.Souls >= totalCost)
	}
}

// Draw dessine le menu
func (lm *LevelUpMenu) Draw(screen *ebiten.Image) {
	// Fond
	bgColor := color.RGBA{30, 30, 40, 240}
	vector.DrawFilledRect(screen, lm.X, lm.Y, lm.Width, lm.Height, bgColor, false)

	// Bordure
	borderColor := color.RGBA{100, 150, 255, 255}
	vector.StrokeRect(screen, lm.X, lm.Y, lm.Width, lm.Height, 2, borderColor, false)

	// Titre
	titleX := int(lm.X) + 10
	titleY := int(lm.Y) + 10
	ebitenutil.DebugPrintAt(screen, "=== MONTEE DE NIVEAU ===", titleX, titleY)

	// Informations du joueur
	startY := titleY + 30
	lineHeight := 18

	lines := []string{
		fmt.Sprintf("Niveau: %d", lm.PlayerStats.Level),
		fmt.Sprintf("Ames disponibles: %d", lm.PlayerStats.Souls),
		"",
		"=== CARACTERISTIQUES ===",
		fmt.Sprintf("Force: %d -> %d (Degats: %.0f)", lm.baseForce, lm.pendingForce, float64(lm.pendingForce)*2),
		fmt.Sprintf("Stamina: %d -> %d (Max: %.0f)", lm.baseStamina, lm.pendingStamina, 100.0+float64(lm.pendingStamina)*5),
		fmt.Sprintf("Vie: %d -> %d (Max HP: %.0f)", lm.baseVie, lm.pendingVie, 100.0+float64(lm.pendingVie)*10),
		"",
		fmt.Sprintf("Cout par amelioration: %d ames", lm.PlayerStats.GetSoulsForStatUpgrade()),
		fmt.Sprintf("Cout total: %d ames", lm.getTotalCost()),
	}

	for i, line := range lines {
		y := startY + i*lineHeight
		ebitenutil.DebugPrintAt(screen, line, titleX, y)
	}

	// Dessiner les boutons
	for _, button := range lm.Buttons {
		button.Draw(screen)
	}

	// Afficher les valeurs à côté des boutons +/-
	buttonStartY := int(lm.Y + 220)

	// Force
	forceText := fmt.Sprintf("Force: %d", lm.pendingForce)
	ebitenutil.DebugPrintAt(screen, forceText, titleX, buttonStartY+5)

	// Stamina
	staminaText := fmt.Sprintf("Stamina: %d", lm.pendingStamina)
	ebitenutil.DebugPrintAt(screen, staminaText, titleX, buttonStartY+45)

	// Vie
	vieText := fmt.Sprintf("Vie: %d", lm.pendingVie)
	ebitenutil.DebugPrintAt(screen, vieText, titleX, buttonStartY+85)
}
