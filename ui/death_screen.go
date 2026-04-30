package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	deathFadeInDuration  = 90  // 1.5s pour le fondu au noir
	deathHoldDuration    = 60  // 1s d'attente avant d'afficher "appuyez sur une touche"
	deathReadyAfter      = deathFadeInDuration + deathHoldDuration
	deathTextScale       = 4   // Facteur d'agrandissement du texte "YOU DIED"
)

// DeathScreen gère l'écran de mort (style Dark Souls)
type DeathScreen struct {
	screenWidth  int
	screenHeight int
	timer        int     // frames depuis l'apparition
	alpha        float64 // opacité du fondu (0.0 → 1.0)
	soulsLost    int
	ReadyForInput bool // vrai quand le joueur peut appuyer pour ressusciter
}

// NewDeathScreen crée un écran de mort
func NewDeathScreen(screenWidth, screenHeight, soulsLost int) *DeathScreen {
	return &DeathScreen{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		soulsLost:    soulsLost,
	}
}

// Update fait avancer l'animation du death screen
func (ds *DeathScreen) Update() {
	ds.timer++

	// Fondu progressif vers le noir
	if ds.timer <= deathFadeInDuration {
		ds.alpha = float64(ds.timer) / float64(deathFadeInDuration)
	} else {
		ds.alpha = 1.0
	}

	// Débloquer la reprise après le délai d'attente
	if ds.timer >= deathReadyAfter {
		ds.ReadyForInput = true
	}
}

// Draw dessine l'écran de mort
func (ds *DeathScreen) Draw(screen *ebiten.Image) {
	// Overlay noir progressif
	overlayAlpha := uint8(ds.alpha * 210)
	vector.DrawFilledRect(screen,
		0, 0,
		float32(ds.screenWidth), float32(ds.screenHeight),
		color.RGBA{0, 0, 0, overlayAlpha},
		false,
	)

	// Ne rien afficher pendant les premières frames du fondu
	if ds.alpha < 0.5 {
		return
	}

	textAlpha := (ds.alpha - 0.5) / 0.5 // 0→1 dans la 2e moitié du fondu

	centerX := float64(ds.screenWidth) / 2
	centerY := float64(ds.screenHeight) / 2

	// --- Titre "YOU DIED" agrandi ---
	// On dessine dans une petite image puis on la scale x4
	titleStr := "YOU DIED"
	titleImg := ebiten.NewImage(len(titleStr)*7+4, 14)
	ebitenutil.DebugPrintAt(titleImg, titleStr, 2, 1)

	titleOp := &ebiten.DrawImageOptions{}
	titleOp.GeoM.Scale(deathTextScale, deathTextScale)
	scaledW := float64(len(titleStr)*7+4) * deathTextScale
	scaledH := float64(14) * deathTextScale
	titleOp.GeoM.Translate(centerX-scaledW/2, centerY-scaledH/2-20)

	// Teinte rouge + alpha
	titleOp.ColorScale.Scale(1.8, 0.15, 0.15, float32(textAlpha))

	screen.DrawImage(titleImg, titleOp)

	// --- Barre décorative rouge ---
	barY := float32(centerY) + float32(scaledH)/2 - 10
	barAlpha := uint8(textAlpha * 180)
	vector.DrawFilledRect(screen,
		float32(centerX)-140, barY,
		280, 2,
		color.RGBA{180, 30, 30, barAlpha},
		false,
	)

	// --- Âmes perdues ---
	if ds.soulsLost > 0 {
		soulsStr := fmt.Sprintf("Ames perdues : %d", ds.soulsLost)
		soulsX := int(centerX) - len(soulsStr)*3
		soulsY := int(barY) + 12
		ebitenutil.DebugPrintAt(screen, soulsStr, soulsX, soulsY)
	}

	// --- "Appuyer pour continuer" (apparaît après le délai) ---
	if ds.ReadyForInput {
		hint := "Appuyez sur une touche pour ressusciter au dernier feu de camp"
		hintX := int(centerX) - len(hint)*3
		hintY := int(centerY) + 80
		ebitenutil.DebugPrintAt(screen, hint, hintX, hintY)
	}
}
