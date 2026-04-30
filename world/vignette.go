package world

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// vignetteConfig paramètre l'effet torche selon le thème du donjon
type vignetteConfig struct {
	lightRadius float64 // rayon intérieur transparent (pixels écran)
	fadeRadius  float64 // rayon auquel l'obscurité est complète
	maxAlpha    float64 // opacité maximale du noir [0-255]
}

var vignetteConfigs = map[ZoneTheme]vignetteConfig{
	ThemeDefault:   {lightRadius: 260, fadeRadius: 520, maxAlpha: 195},
	ThemeForest:    {lightRadius: 280, fadeRadius: 560, maxAlpha: 180},
	ThemeCatacombs: {lightRadius: 175, fadeRadius: 420, maxAlpha: 225},
	ThemeRuins:     {lightRadius: 225, fadeRadius: 500, maxAlpha: 208},
}

// Vignette gère l'effet de lumière de torche centré sur le joueur.
// L'image est pré-calculée une seule fois (2× taille écran) et positionnée
// chaque frame sur le joueur avec un léger scintillement sinusoïdal.
type Vignette struct {
	img     *ebiten.Image
	screenW int
	screenH int
	tickerT float64 // temps interne pour le scintillement
}

// NewVignette crée la vignette pour le thème donné.
// screenW/screenH : dimensions de la fenêtre de jeu.
func NewVignette(screenW, screenH int, theme ZoneTheme) *Vignette {
	cfg, ok := vignetteConfigs[theme]
	if !ok {
		cfg = vignetteConfigs[ThemeDefault]
	}

	v := &Vignette{screenW: screenW, screenH: screenH}
	v.img = buildVignetteImage(screenW*2, screenH*2, cfg)
	return v
}

// buildVignetteImage construit un gradient radial noir-transparent dans une image
// de taille (w×h), centré, avec le rayon et l'opacité max donnés.
func buildVignetteImage(w, h int, cfg vignetteConfig) *ebiten.Image {
	cx := float64(w) / 2
	cy := float64(h) / 2

	pixels := make([]byte, w*h*4)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			dist := math.Sqrt(dx*dx + dy*dy)

			var alpha float64
			if dist <= cfg.lightRadius {
				alpha = 0
			} else if dist >= cfg.fadeRadius {
				alpha = cfg.maxAlpha
			} else {
				t := (dist - cfg.lightRadius) / (cfg.fadeRadius - cfg.lightRadius)
				// Ease-in cubique pour un fondu progressif et naturel
				t = t * t * (3 - 2*t)
				alpha = t * cfg.maxAlpha
			}

			i := (y*w + x) * 4
			pixels[i] = 0
			pixels[i+1] = 0
			pixels[i+2] = 0
			pixels[i+3] = byte(alpha)
		}
	}

	img := ebiten.NewImage(w, h)
	img.WritePixels(pixels)
	return img
}

// Update avance l'animation de scintillement (à appeler chaque frame).
func (v *Vignette) Update() {
	v.tickerT += 0.035
}

// Draw dessine la vignette centrée sur la position écran du joueur.
// playerScreenX, playerScreenY : coordonnées pixel du joueur à l'écran.
func (v *Vignette) Draw(screen *ebiten.Image, playerScreenX, playerScreenY float64) {
	// Scintillement : dérive légère de position (~3px max) via deux sinus déphasés
	flickerX := math.Sin(v.tickerT)*2.2 + math.Sin(v.tickerT*1.61)*1.3
	flickerY := math.Cos(v.tickerT*0.93)*1.8 + math.Cos(v.tickerT*2.1)*0.9

	// Opacité légèrement variable (±4%)
	alphaScale := float32(0.97 + 0.03*math.Sin(v.tickerT*0.7+1.0))

	opts := &ebiten.DrawImageOptions{}
	// Centrer l'image 2×écran sur le joueur
	offX := playerScreenX - float64(v.screenW) + flickerX
	offY := playerScreenY - float64(v.screenH) + flickerY
	opts.GeoM.Translate(offX, offY)
	opts.ColorScale.ScaleAlpha(alphaScale)

	screen.DrawImage(v.img, opts)
}
