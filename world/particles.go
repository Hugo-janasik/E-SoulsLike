package world

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// particle représente une particule d'ambiance en coordonnées monde.
type particle struct {
	x, y     float64  // position monde
	vx, vy   float64  // vitesse (monde/frame)
	life     float64  // vie restante [0, 1]
	decay    float64  // réduction de vie par frame
	size     float32  // rayon en pixels écran
	r, g, b  uint8    // couleur de base
	maxAlpha uint8    // alpha maximum (à vie pleine)
	wobble   float64  // phase de l'oscillation latérale
	wobbleAmp float64 // amplitude de l'oscillation
}

// particleConfig décrit le comportement des particules pour un thème donné.
type particleConfig struct {
	maxCount int
	types    []particleType
}

type particleType struct {
	r, g, b           uint8
	maxAlpha          uint8
	sizeMin, sizeMax  float32
	vxMin, vxMax      float64
	vyMin, vyMax      float64
	decayMin, decayMax float64
	wobbleAmp         float64
	weight            int // probabilité relative
}

var particleConfigs = map[ZoneTheme]particleConfig{
	ThemeForest: {
		maxCount: 55,
		types: []particleType{
			// Spores verdâtres qui montent doucement
			{r: 160, g: 220, b: 110, maxAlpha: 130,
				sizeMin: 1.0, sizeMax: 2.2,
				vxMin: -0.12, vxMax: 0.12,
				vyMin: -0.45, vyMax: -0.15,
				decayMin: 0.003, decayMax: 0.007,
				wobbleAmp: 0.18, weight: 6},
			// Poussière blanche translucide
			{r: 210, g: 230, b: 200, maxAlpha: 80,
				sizeMin: 0.8, sizeMax: 1.5,
				vxMin: -0.08, vxMax: 0.08,
				vyMin: -0.10, vyMax: 0.05,
				decayMin: 0.004, decayMax: 0.008,
				wobbleAmp: 0.10, weight: 4},
		},
	},
	ThemeCatacombs: {
		maxCount: 65,
		types: []particleType{
			// Braises oranges qui montent vite
			{r: 230, g: 90, b: 20, maxAlpha: 210,
				sizeMin: 1.0, sizeMax: 2.5,
				vxMin: -0.18, vxMax: 0.18,
				vyMin: -1.10, vyMax: -0.55,
				decayMin: 0.010, decayMax: 0.020,
				wobbleAmp: 0.22, weight: 5},
			// Cendres grises qui descendent lentement
			{r: 160, g: 145, b: 135, maxAlpha: 90,
				sizeMin: 0.8, sizeMax: 1.8,
				vxMin: -0.06, vxMax: 0.10,
				vyMin: 0.05, vyMax: 0.22,
				decayMin: 0.004, decayMax: 0.009,
				wobbleAmp: 0.08, weight: 5},
			// Petites étincelles jaunes
			{r: 255, g: 210, b: 60, maxAlpha: 180,
				sizeMin: 0.6, sizeMax: 1.2,
				vxMin: -0.25, vxMax: 0.25,
				vyMin: -1.40, vyMax: -0.70,
				decayMin: 0.018, decayMax: 0.030,
				wobbleAmp: 0.30, weight: 2},
		},
	},
	ThemeRuins: {
		maxCount: 48,
		types: []particleType{
			// Poussière de pierre, dérive légèrement vers la droite
			{r: 155, g: 140, b: 115, maxAlpha: 100,
				sizeMin: 0.8, sizeMax: 1.8,
				vxMin: 0.10, vxMax: 0.40,
				vyMin: -0.05, vyMax: 0.12,
				decayMin: 0.005, decayMax: 0.010,
				wobbleAmp: 0.06, weight: 7},
			// Gravier plus gros qui tombe
			{r: 130, g: 118, b: 98, maxAlpha: 70,
				sizeMin: 1.2, sizeMax: 2.0,
				vxMin: 0.05, vxMax: 0.20,
				vyMin: 0.15, vyMax: 0.40,
				decayMin: 0.007, decayMax: 0.014,
				wobbleAmp: 0.04, weight: 3},
		},
	},
	ThemeDefault: {
		maxCount: 40,
		types: []particleType{
			// Grains de poussière neutres
			{r: 175, g: 165, b: 155, maxAlpha: 85,
				sizeMin: 0.7, sizeMax: 1.6,
				vxMin: -0.10, vxMax: 0.10,
				vyMin: -0.15, vyMax: 0.08,
				decayMin: 0.004, decayMax: 0.009,
				wobbleAmp: 0.10, weight: 10},
		},
	},
}

// AmbientParticles gère le pool de particules d'ambiance pour un donjon.
type AmbientParticles struct {
	pool    []particle
	maxCount int
	cfg     particleConfig
	rng     *rand.Rand
	totalWeight int
}

// NewAmbientParticles crée le système de particules pour le thème donné.
func NewAmbientParticles(theme ZoneTheme) *AmbientParticles {
	cfg, ok := particleConfigs[theme]
	if !ok {
		cfg = particleConfigs[ThemeDefault]
	}
	tw := 0
	for _, pt := range cfg.types {
		tw += pt.weight
	}
	return &AmbientParticles{
		pool:        make([]particle, 0, cfg.maxCount),
		maxCount:    cfg.maxCount,
		cfg:         cfg,
		rng:         rand.New(rand.NewSource(rand.Int63())),
		totalWeight: tw,
	}
}

// Update déplace les particules existantes et en fait naître de nouvelles
// dans la zone visible de la caméra.
func (ap *AmbientParticles) Update(camera *Camera) {
	// Avancer les particules vivantes
	n := 0
	for i := range ap.pool {
		p := &ap.pool[i]
		p.life -= p.decay
		if p.life <= 0 {
			continue
		}
		// Oscillation latérale (donne du naturel)
		p.wobble += 0.05
		p.x += p.vx + math.Sin(p.wobble)*p.wobbleAmp*0.06
		p.y += p.vy
		ap.pool[n] = *p
		n++
	}
	ap.pool = ap.pool[:n]

	// Faire naître de nouvelles particules jusqu'au maximum
	// On en fait naître plusieurs par frame pour remplir rapidement le pool au démarrage
	spawnBurst := 3
	for i := 0; i < spawnBurst && len(ap.pool) < ap.maxCount; i++ {
		ap.pool = append(ap.pool, ap.spawnParticle(camera))
	}
}

// spawnParticle crée une particule à un endroit aléatoire de la zone visible.
func (ap *AmbientParticles) spawnParticle(camera *Camera) particle {
	// Zone de spawn : vue caméra + 20% de marge de chaque côté
	margin := 0.20
	w := camera.Width / camera.Zoom
	h := camera.Height / camera.Zoom
	x := camera.X + w*(-margin+ap.rng.Float64()*(1+2*margin))
	y := camera.Y + h*(-margin+ap.rng.Float64()*(1+2*margin))

	pt := ap.pickType()
	sz := pt.sizeMin + ap.rng.Float32()*(pt.sizeMax-pt.sizeMin)
	vx := pt.vxMin + ap.rng.Float64()*(pt.vxMax-pt.vxMin)
	vy := pt.vyMin + ap.rng.Float64()*(pt.vyMax-pt.vyMin)
	decay := pt.decayMin + ap.rng.Float64()*(pt.decayMax-pt.decayMin)

	return particle{
		x: x, y: y,
		vx: vx, vy: vy,
		life: 1.0, decay: decay,
		size: sz,
		r: pt.r, g: pt.g, b: pt.b,
		maxAlpha:  pt.maxAlpha,
		wobble:    ap.rng.Float64() * math.Pi * 2,
		wobbleAmp: pt.wobbleAmp,
	}
}

// pickType choisit un type de particule aléatoirement selon les poids.
func (ap *AmbientParticles) pickType() particleType {
	roll := ap.rng.Intn(ap.totalWeight)
	cumul := 0
	for _, pt := range ap.cfg.types {
		cumul += pt.weight
		if roll < cumul {
			return pt
		}
	}
	return ap.cfg.types[len(ap.cfg.types)-1]
}

// Draw dessine toutes les particules vivantes.
func (ap *AmbientParticles) Draw(screen *ebiten.Image, camera *Camera) {
	for _, p := range ap.pool {
		sx, sy := camera.WorldToScreen(p.x, p.y)
		alpha := uint8(float64(p.maxAlpha) * p.life)
		if alpha < 5 {
			continue // trop transparent pour être visible
		}
		c := color.RGBA{p.r, p.g, p.b, alpha}
		vector.DrawFilledCircle(screen, float32(sx), float32(sy), p.size, c, false)
	}
}
