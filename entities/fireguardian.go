package entities

import (
	"e-soulslike/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

// FireGuardian représente le gardien du feu (PNJ interactif dans les salles de repos)
type FireGuardian struct {
	X, Y              float64
	InteractionRange  float64
	IsPlayerNearby    bool
	AnimationTimer    int
	AnimationFrame    int
	Sprite            *ebiten.Image
}

// NewFireGuardian crée un nouveau gardien du feu
func NewFireGuardian(x, y float64) *FireGuardian {
	fg := &FireGuardian{
		X:                x,
		Y:                y,
		InteractionRange: 80.0, // Distance pour interagir
		AnimationTimer:   0,
		AnimationFrame:   0,
	}

	// Charger le sprite (si disponible)
	sprite, _, err := ebitenutil.NewImageFromFile("assets/npcs/fireguardian/idle.png")
	if err == nil {
		fg.Sprite = sprite
	}

	return fg
}

// Update met à jour le gardien du feu
func (fg *FireGuardian) Update(playerX, playerY float64) {
	// Animation
	fg.AnimationTimer++
	if fg.AnimationTimer >= 15 {
		fg.AnimationTimer = 0
		fg.AnimationFrame = (fg.AnimationFrame + 1) % 4
	}

	// Vérifier si le joueur est à portée
	dx := playerX - fg.X
	dy := playerY - fg.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	fg.IsPlayerNearby = distance < fg.InteractionRange
}

// Draw dessine le gardien du feu
func (fg *FireGuardian) Draw(screen *ebiten.Image, camera *world.Camera) {
	screenX, screenY := camera.WorldToScreen(fg.X, fg.Y)

	// Si le sprite est chargé, l'afficher
	if fg.Sprite != nil {
		op := &ebiten.DrawImageOptions{}

		// Centrer le sprite
		w, h := fg.Sprite.Size()
		op.GeoM.Translate(-float64(w)/2, -float64(h))
		op.GeoM.Translate(screenX, screenY)

		screen.DrawImage(fg.Sprite, op)
	} else {
		// Sinon, dessiner un placeholder (forme de PNJ)
		fg.drawPlaceholder(screen, screenX, screenY)
	}

	// Si le joueur est proche, afficher le prompt
	if fg.IsPlayerNearby {
		promptText := "Appuyez sur E pour interagir"
		textX := int(screenX) - len(promptText)*3
		textY := int(screenY) - 80

		// Fond pour le texte
		textBg := color.RGBA{0, 0, 0, 180}
		vector.DrawFilledRect(
			screen,
			float32(textX-5),
			float32(textY-5),
			float32(len(promptText)*6+10),
			20,
			textBg,
			false,
		)

		ebitenutil.DebugPrintAt(screen, promptText, textX, textY)
	}
}

// drawPlaceholder dessine un placeholder pour le Fire Guardian
func (fg *FireGuardian) drawPlaceholder(screen *ebiten.Image, screenX, screenY float64) {
	// Corps (robe)
	bodyColor := color.RGBA{120, 60, 180, 255} // Violet
	vector.DrawFilledCircle(screen, float32(screenX), float32(screenY-10), 25, bodyColor, true)

	// Tête
	headColor := color.RGBA{255, 220, 180, 255} // Couleur peau
	vector.DrawFilledCircle(screen, float32(screenX), float32(screenY-40), 15, headColor, true)

	// Flamme au-dessus de la tête (symbole du feu)
	flameSize := float32(15 + fg.AnimationFrame*2)
	flameColor1 := color.RGBA{255, 150, 0, 255}  // Orange
	flameColor2 := color.RGBA{255, 200, 50, 255} // Jaune

	vector.DrawFilledCircle(screen, float32(screenX), float32(screenY-60), flameSize/2, flameColor1, true)
	vector.DrawFilledCircle(screen, float32(screenX), float32(screenY-70), flameSize/3, flameColor2, true)

	// Base (pieds)
	baseColor := color.RGBA{80, 40, 120, 255}
	vector.DrawFilledRect(screen, float32(screenX-20), float32(screenY+10), 40, 10, baseColor, false)
}

// CanInteract retourne vrai si le joueur peut interagir avec le gardien
func (fg *FireGuardian) CanInteract() bool {
	return fg.IsPlayerNearby
}
