package entities

import (
	"e-soulslike/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

// Campfire représente un feu de camp (point de sauvegarde)
type Campfire struct {
	X, Y              float64
	InteractionRange  float64
	IsPlayerNearby    bool
	AnimationTimer    int
	AnimationFrame    int
}

// NewCampfire crée un nouveau feu de camp
func NewCampfire(x, y float64) *Campfire {
	return &Campfire{
		X:                x,
		Y:                y,
		InteractionRange: 80.0, // Distance pour interagir
		AnimationTimer:   0,
		AnimationFrame:   0,
	}
}

// Update met à jour le feu de camp
func (c *Campfire) Update(playerX, playerY float64) {
	// Animation du feu
	c.AnimationTimer++
	if c.AnimationTimer >= 10 {
		c.AnimationTimer = 0
		c.AnimationFrame = (c.AnimationFrame + 1) % 4
	}

	// Vérifier si le joueur est à portée
	dx := playerX - c.X
	dy := playerY - c.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	c.IsPlayerNearby = distance < c.InteractionRange
}

// Draw dessine le feu de camp
func (c *Campfire) Draw(screen *ebiten.Image, camera *world.Camera) {
	screenX, screenY := camera.WorldToScreen(c.X, c.Y)

	// Base du feu (cercle gris foncé pour représenter les pierres)
	baseColor := color.RGBA{60, 60, 60, 255}
	vector.DrawFilledCircle(screen, float32(screenX), float32(screenY), 20, baseColor, true)

	// Flammes (triangles qui changent de taille)
	flameSize := float32(15 + c.AnimationFrame*2)
	flameColor1 := color.RGBA{255, 150, 0, 255}  // Orange
	flameColor2 := color.RGBA{255, 100, 0, 255}  // Orange foncé
	flameColor3 := color.RGBA{255, 200, 50, 255} // Jaune

	// Dessiner plusieurs flammes pour l'animation
	for i := 0; i < 3; i++ {
		offset := float32(i * 10)
		flameColor := flameColor1
		if i == 1 {
			flameColor = flameColor2
		} else if i == 2 {
			flameColor = flameColor3
		}

		// Dessiner une flamme (cercle qui monte)
		x := float32(screenX)
		y := float32(screenY) - flameSize - offset

		vector.DrawFilledCircle(screen, x, y, flameSize/2, flameColor, true)
	}

	// Lueur autour du feu
	glowColor := color.RGBA{255, 200, 100, 50}
	vector.DrawFilledCircle(screen, float32(screenX), float32(screenY), 40+float32(c.AnimationFrame*2), glowColor, true)

	// Si le joueur est proche, afficher le prompt
	if c.IsPlayerNearby {
		promptText := "Appuyez sur E pour interagir"
		textX := int(screenX) - len(promptText)*3
		textY := int(screenY) - 50
		ebitenutil.DebugPrintAt(screen, promptText, textX, textY)
	}
}

// CanInteract retourne vrai si le joueur peut interagir avec le feu
func (c *Campfire) CanInteract() bool {
	return c.IsPlayerNearby
}
