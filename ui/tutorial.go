package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Tutorial affiche le tutoriel du jeu
type Tutorial struct {
	X, Y          float32
	Width, Height float32
	CurrentPage   int
	MaxPages      int
	Buttons       []*Button
	OnClose       func()
}

// NewTutorial crée un nouveau tutoriel
func NewTutorial(x, y, width, height float32, onClose func()) *Tutorial {
	tut := &Tutorial{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		CurrentPage: 0,
		MaxPages:    4,
		OnClose:     onClose,
		Buttons:     make([]*Button, 0),
	}

	// Boutons de navigation
	buttonY := y + height - 50

	// Bouton Précédent
	tut.Buttons = append(tut.Buttons, NewButton(
		x+20, buttonY, 100, 30,
		"< Precedent",
		func() {
			if tut.CurrentPage > 0 {
				tut.CurrentPage--
			}
		},
	))

	// Bouton Suivant
	tut.Buttons = append(tut.Buttons, NewButton(
		x+130, buttonY, 100, 30,
		"Suivant >",
		func() {
			if tut.CurrentPage < tut.MaxPages-1 {
				tut.CurrentPage++
			}
		},
	))

	// Bouton Fermer
	tut.Buttons = append(tut.Buttons, NewButton(
		x+width-120, buttonY, 100, 30,
		"Fermer",
		onClose,
	))

	return tut
}

// Update met à jour le tutoriel
func (tut *Tutorial) Update(mouseX, mouseY int, mousePressed, mouseJustPressed bool) {
	for _, button := range tut.Buttons {
		button.Update(mouseX, mouseY, mousePressed, mouseJustPressed)
	}

	// Désactiver les boutons selon la page
	tut.Buttons[0].SetEnabled(tut.CurrentPage > 0)
	tut.Buttons[1].SetEnabled(tut.CurrentPage < tut.MaxPages-1)
}

// Draw dessine le tutoriel
func (tut *Tutorial) Draw(screen *ebiten.Image) {
	// Fond
	bgColor := color.RGBA{30, 30, 40, 240}
	vector.DrawFilledRect(screen, tut.X, tut.Y, tut.Width, tut.Height, bgColor, false)

	// Bordure
	borderColor := color.RGBA{100, 150, 255, 255}
	vector.StrokeRect(screen, tut.X, tut.Y, tut.Width, tut.Height, 2, borderColor, false)

	// Contenu selon la page
	contentX := int(tut.X) + 20
	contentY := int(tut.Y) + 20

	switch tut.CurrentPage {
	case 0:
		tut.drawPage1(screen, contentX, contentY)
	case 1:
		tut.drawPage2(screen, contentX, contentY)
	case 2:
		tut.drawPage3(screen, contentX, contentY)
	case 3:
		tut.drawPage4(screen, contentX, contentY)
	}

	// Indicateur de page
	pageText := fmt.Sprintf("Page %d/%d", tut.CurrentPage+1, tut.MaxPages)
	pageX := int(tut.X + tut.Width/2 - 30)
	pageY := int(tut.Y + tut.Height - 60)
	ebitenutil.DebugPrintAt(screen, pageText, pageX, pageY)

	// Dessiner les boutons
	for _, button := range tut.Buttons {
		button.Draw(screen)
	}
}

func (tut *Tutorial) drawPage1(screen *ebiten.Image, x, y int) {
	lines := []string{
		"=== BIENVENUE DANS E-SOULSLIKE ===",
		"",
		"E-SoulsLike est un jeu d'action-aventure 2D",
		"inspiré de Zelda et Dark Souls.",
		"",
		"Vous incarnez un héros qui doit explorer",
		"un monde dangereux rempli d'ennemis.",
		"",
		"Utilisez vos compétences de combat et",
		"votre agilité pour survivre !",
		"",
		"Appuyez sur 'Suivant' pour continuer...",
	}

	for i, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, x, y+i*20)
	}
}

func (tut *Tutorial) drawPage2(screen *ebiten.Image, x, y int) {
	lines := []string{
		"=== CONTROLES DE BASE ===",
		"",
		"DEPLACEMENTS:",
		"  ZQSD ou WASD    - Se deplacer",
		"  Fleches         - Se deplacer (alternatif)",
		"",
		"ACTIONS:",
		"  ESPACE          - Attaquer (cout: 20 stamina)",
		"  SHIFT           - Sprint (consomme stamina)",
		"  SHIFT + Direction - Esquive/Roll",
		"                    (cout: 30 stamina)",
		"",
		"SYSTEME:",
		"  ECHAP           - Menu pause",
	}

	for i, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, x, y+i*20)
	}
}

func (tut *Tutorial) drawPage3(screen *ebiten.Image, x, y int) {
	lines := []string{
		"=== SYSTEMES DE JEU ===",
		"",
		"STAMINA:",
		"  - Utilisée pour attaquer, sprinter, esquiver",
		"  - Se régénère automatiquement",
		"  - Gérez bien votre stamina en combat!",
		"",
		"COMBAT:",
		"  - Attaquez les ennemis dans votre champ",
		"  - Esquivez pour éviter les dégâts",
		"  - Les ennemis sont étourdis quand touchés",
		"",
		"ENNEMIS:",
		"  - Basic: Équilibré",
		"  - Fast: Rapide mais fragile",
		"  - Tank: Lent mais résistant",
	}

	for i, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, x, y+i*20)
	}
}

func (tut *Tutorial) drawPage4(screen *ebiten.Image, x, y int) {
	lines := []string{
		"=== CONSEILS ===",
		"",
		"1. Gardez toujours un oeil sur votre stamina",
		"   (barre verte sous vos points de vie)",
		"",
		"2. Utilisez l'esquive pour éviter les dégâts",
		"   et repositionner stratégiquement",
		"",
		"3. Les ennemis ont des rayons de détection,",
		"   approchez-vous pour les attirer",
		"",
		"4. N'oubliez pas de sauvegarder régulièrement",
		"   via le menu pause!",
		"",
		"Bonne chance, aventurier!",
	}

	for i, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, x, y+i*20)
	}
}

// Ajouter l'import fmt en haut du fichier
