package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

// Portal représente une porte/portail entre zones
type Portal struct {
	X, Y              float64
	Width, Height     float64
	DestinationZoneID string // ID de la zone de destination
	DestinationSpawn  string // Nom du spawn point dans la destination
	IsPlayerNearby    bool
	InteractionRange  float64
	Label             string // Texte affiché (ex: "Forêt Maudite")

	// Animation
	AnimationTimer int
	AnimationFrame int
}

// NewPortal crée un nouveau portail
func NewPortal(x, y float64, destZoneID, destSpawn, label string) *Portal {
	return &Portal{
		X:                 x,
		Y:                 y,
		Width:             48,
		Height:            64,
		DestinationZoneID: destZoneID,
		DestinationSpawn:  destSpawn,
		Label:             label,
		InteractionRange:  80.0,
		IsPlayerNearby:    false,
		AnimationTimer:    0,
		AnimationFrame:    0,
	}
}

// Update met à jour le portail
func (p *Portal) Update(playerX, playerY float64) {
	// Calculer la distance avec le joueur
	dx := playerX - p.X
	dy := playerY - p.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Vérifier si le joueur est à portée d'interaction
	p.IsPlayerNearby = distance < p.InteractionRange

	// Animation du portail
	p.AnimationTimer++
	if p.AnimationTimer >= 10 {
		p.AnimationTimer = 0
		p.AnimationFrame = (p.AnimationFrame + 1) % 8
	}
}

// CanInteract vérifie si le joueur peut interagir
func (p *Portal) CanInteract() bool {
	return p.IsPlayerNearby
}

// Draw dessine le portail
func (p *Portal) Draw(screen *ebiten.Image, camera *Camera) {
	// Convertir les coordonnées du monde en coordonnées écran
	screenX, screenY := camera.WorldToScreen(p.X, p.Y)

	// Dessiner la base du portail (rectangle sombre)
	baseColor := color.RGBA{40, 40, 60, 255}
	vector.DrawFilledRect(
		screen,
		float32(screenX-p.Width/2),
		float32(screenY-p.Height),
		float32(p.Width),
		float32(p.Height),
		baseColor,
		false,
	)

	// Bordure du portail (effet magique)
	borderColors := []color.RGBA{
		{100, 100, 255, 255}, // Bleu
		{150, 100, 255, 255}, // Violet
		{100, 150, 255, 255}, // Bleu clair
	}
	borderColor := borderColors[p.AnimationFrame%3]
	vector.StrokeRect(
		screen,
		float32(screenX-p.Width/2),
		float32(screenY-p.Height),
		float32(p.Width),
		float32(p.Height),
		3,
		borderColor,
		false,
	)

	// Effet de particules/énergie au centre
	for i := 0; i < 5; i++ {
		offsetY := float32(i*12) + float32((p.AnimationFrame*2)%60) - 30
		particleColor := color.RGBA{
			R: 100 + uint8(p.AnimationFrame*10),
			G: 100 + uint8(p.AnimationFrame*15),
			B: 200,
			A: uint8(150 - i*20),
		}
		vector.DrawFilledCircle(
			screen,
			float32(screenX),
			float32(screenY-p.Height/2)+offsetY,
			3+float32(i),
			particleColor,
			true,
		)
	}

	// Si le joueur est proche, afficher le prompt
	if p.IsPlayerNearby {
		promptText := "E - " + p.Label
		promptX := int(screenX) - len(promptText)*3
		promptY := int(screenY - p.Height - 20)

		// Fond pour le texte
		textBg := color.RGBA{0, 0, 0, 180}
		vector.DrawFilledRect(
			screen,
			float32(promptX-5),
			float32(promptY-5),
			float32(len(promptText)*6+10),
			20,
			textBg,
			false,
		)

		ebitenutil.DebugPrintAt(screen, promptText, promptX, promptY)
	}
}
