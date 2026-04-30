package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// Button représente un bouton cliquable
type Button struct {
	X, Y          float32
	Width, Height float32
	Text          string
	OnClick       func()
	IsHovered     bool
	IsPressed     bool
	Enabled       bool
}

// NewButton crée un nouveau bouton
func NewButton(x, y, width, height float32, text string, onClick func()) *Button {
	return &Button{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		Text:    text,
		OnClick: onClick,
		Enabled: true,
	}
}

// Update met à jour l'état du bouton
func (b *Button) Update(mouseX, mouseY int, mousePressed, mouseJustPressed bool) {
	if !b.Enabled {
		b.IsHovered = false
		b.IsPressed = false
		return
	}

	// Vérifier si la souris est sur le bouton
	mx, my := float32(mouseX), float32(mouseY)
	b.IsHovered = mx >= b.X && mx <= b.X+b.Width &&
		my >= b.Y && my <= b.Y+b.Height

	// Si on clique sur le bouton
	if b.IsHovered && mouseJustPressed {
		b.IsPressed = true
		if b.OnClick != nil {
			b.OnClick()
		}
	} else if !mousePressed {
		b.IsPressed = false
	}
}

// Draw dessine le bouton
func (b *Button) Draw(screen *ebiten.Image) {
	var bgColor, borderColor color.RGBA

	if !b.Enabled {
		// Bouton désactivé
		bgColor = color.RGBA{60, 60, 60, 255}
		borderColor = color.RGBA{80, 80, 80, 255}
	} else if b.IsPressed {
		// Bouton pressé
		bgColor = color.RGBA{80, 120, 200, 255}
		borderColor = color.RGBA{100, 150, 255, 255}
	} else if b.IsHovered {
		// Bouton survolé
		bgColor = color.RGBA{70, 100, 180, 255}
		borderColor = color.RGBA{100, 150, 255, 255}
	} else {
		// Bouton normal
		bgColor = color.RGBA{50, 80, 150, 255}
		borderColor = color.RGBA{80, 120, 200, 255}
	}

	// Dessiner le fond du bouton
	vector.DrawFilledRect(screen, b.X, b.Y, b.Width, b.Height, bgColor, false)

	// Dessiner la bordure
	vector.StrokeRect(screen, b.X, b.Y, b.Width, b.Height, 2, borderColor, false)

	// Dessiner le texte (centré) avec ebitenutil
	textX := int(b.X + b.Width/2 - float32(len(b.Text)*6)/2)
	textY := int(b.Y + b.Height/2 - 4)

	// Utiliser DebugPrintAt pour afficher le texte
	// Note: La couleur sera toujours blanche avec DebugPrint,
	// mais c'est lisible contrairement aux rectangles
	ebitenutil.DebugPrintAt(screen, b.Text, textX, textY)
}

// SetEnabled active ou désactive le bouton
func (b *Button) SetEnabled(enabled bool) {
	b.Enabled = enabled
}

// Contains vérifie si un point est dans le bouton
func (b *Button) Contains(x, y int) bool {
	fx, fy := float32(x), float32(y)
	return fx >= b.X && fx <= b.X+b.Width &&
		fy >= b.Y && fy <= b.Y+b.Height
}
