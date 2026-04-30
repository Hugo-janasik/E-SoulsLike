package world

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// dungeonPalette définit la palette de couleurs d'un thème de donjon
type dungeonPalette struct {
	floorBase   color.RGBA // couleur de base du sol
	floorJoint  color.RGBA // joints entre les dalles
	wallTop     color.RGBA // dessus du mur (face visible depuis le haut)
	wallFace    color.RGBA // face avant du mur (visible depuis le sud)
	wallShadow  color.RGBA // ombre pied/côté gauche du mur
	wallInner   color.RGBA // mur interne / plafond (loin des salles)
	floorShadow color.RGBA // ombre au sol projetée par le mur nord
	void        color.RGBA // vide absolu (hors de tout couloir)
}

var dungeonPalettes = map[ZoneTheme]dungeonPalette{
	ThemeDefault: {
		floorBase:   color.RGBA{68, 62, 55, 255},
		floorJoint:  color.RGBA{42, 38, 32, 255},
		wallTop:     color.RGBA{95, 88, 80, 255},
		wallFace:    color.RGBA{55, 48, 42, 255},
		wallShadow:  color.RGBA{18, 15, 12, 255},
		wallInner:   color.RGBA{22, 20, 17, 255},
		floorShadow: color.RGBA{0, 0, 0, 90},
		void:        color.RGBA{6, 5, 4, 255},
	},
	ThemeForest: {
		floorBase:   color.RGBA{38, 52, 28, 255},
		floorJoint:  color.RGBA{24, 34, 16, 255},
		wallTop:     color.RGBA{58, 78, 45, 255},
		wallFace:    color.RGBA{28, 40, 18, 255},
		wallShadow:  color.RGBA{8, 12, 5, 255},
		wallInner:   color.RGBA{12, 18, 8, 255},
		floorShadow: color.RGBA{0, 5, 0, 80},
		void:        color.RGBA{3, 5, 2, 255},
	},
	ThemeCatacombs: {
		floorBase:   color.RGBA{40, 34, 30, 255},
		floorJoint:  color.RGBA{24, 19, 16, 255},
		wallTop:     color.RGBA{68, 56, 50, 255},
		wallFace:    color.RGBA{30, 24, 20, 255},
		wallShadow:  color.RGBA{6, 4, 3, 255},
		wallInner:   color.RGBA{10, 8, 6, 255},
		floorShadow: color.RGBA{10, 0, 0, 100},
		void:        color.RGBA{2, 1, 1, 255},
	},
	ThemeRuins: {
		floorBase:   color.RGBA{60, 52, 40, 255},
		floorJoint:  color.RGBA{38, 32, 22, 255},
		wallTop:     color.RGBA{88, 74, 58, 255},
		wallFace:    color.RGBA{48, 38, 28, 255},
		wallShadow:  color.RGBA{14, 10, 7, 255},
		wallInner:   color.RGBA{20, 16, 11, 255},
		floorShadow: color.RGBA{5, 3, 0, 85},
		void:        color.RGBA{4, 3, 2, 255},
	},
}

// fmax retourne le plus grand des deux float32.
func fmax(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// tileHash retourne une valeur pseudo-aléatoire déterministe [0,255] pour (x,y).
// Utilisé pour varier subtilement la couleur des tuiles sans stocker de données.
func tileHash(x, y int) int {
	h := x*2654435761 ^ y*2246822519
	h = ((h >> 16) ^ h) & 0x7FFFFFFF
	return h & 0xFF
}

// clampByte restreint une valeur entière dans [0, 255].
func clampByte(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}


// drawWallFace dessine la face visible d'un mur (le mur dont la tuile au sud est du sol).
// Simule une vue top-down : dessus éclairé + face avant sombre + ombres.
// tick pilote l'animation des torches.
func drawWallFace(screen *ebiten.Image, pal dungeonPalette, sx, sy, ts float32, tileX, tileY, tick int) {
	capH := ts * 6 / 32 // hauteur du dessus du mur (~6px à zoom 1)

	// Dessus du mur (surface horizontale, plus claire)
	vector.DrawFilledRect(screen, sx, sy, ts, capH, pal.wallTop, false)

	// Face avant du mur (surface verticale, sombre)
	vector.DrawFilledRect(screen, sx, sy+capH, ts, ts-capH, pal.wallFace, false)

	// Ombre gauche (crée l'illusion de profondeur)
	shadowW := fmax(ts*2/32, 1)
	vector.DrawFilledRect(screen, sx, sy+capH, shadowW, ts-capH, pal.wallShadow, false)

	// Ombre bas (pied du mur)
	shadowH := fmax(ts*3/32, 1)
	vector.DrawFilledRect(screen, sx, sy+ts-shadowH, ts, shadowH, pal.wallShadow, false)

	// Ligne de joint horizontal sur la face (détail de pierre)
	jw := fmax(ts/32, 1)
	crackY := sy + capH + (ts-capH)*0.45
	crackColor := color.RGBA{pal.wallShadow.R, pal.wallShadow.G, pal.wallShadow.B, 160}
	vector.DrawFilledRect(screen, sx+shadowW, crackY, ts-shadowW*2, jw, crackColor, false)

	// Variation aléatoire de couleur du dessus (subtilement différent d'un mur à l'autre)
	v := tileHash(tileX, tileY)
	if v < 60 {
		highlight := color.RGBA{
			R: uint8(clampByte(int(pal.wallTop.R) + 12)),
			G: uint8(clampByte(int(pal.wallTop.G) + 10)),
			B: uint8(clampByte(int(pal.wallTop.B) + 8)),
			A: 255,
		}
		vector.DrawFilledRect(screen, sx+ts*0.25, sy+capH*0.2, ts*0.5, capH*0.5, highlight, false)
	}

	// Torche : placée sur ~10% des faces de murs (sélection déterministe par hash)
	if v < 26 {
		drawTorch(screen, sx, sy, ts, capH, tick)
	}
}

// drawTorch dessine une torche animée centrée horizontalement sur la face du mur.
// capH : hauteur du dessus du mur (la torche est accrochée juste en dessous).
// tick : compteur de frame pour l'animation de la flamme.
func drawTorch(screen *ebiten.Image, sx, sy, ts, capH float32, tick int) {
	// --- Bracket (applique murale) ---
	bw := fmax(ts*5/32, 3) // largeur du support
	bh := fmax(ts*7/32, 4) // hauteur du support
	bx := sx + ts/2 - bw/2
	by := sy + capH + ts*0.12
	vector.DrawFilledRect(screen, bx, by, bw, bh, color.RGBA{22, 14, 8, 255}, false)
	// Reflet métallique (bord gauche légèrement plus clair)
	vector.DrawFilledRect(screen, bx, by, fmax(bw*0.2, 1), bh, color.RGBA{55, 38, 20, 180}, false)

	// --- Flamme animée ---
	t := float64(tick) * 0.12
	tsf := float64(ts) // ts en float64 pour les calculs sin/cos

	// Centre de la flamme (oscille horizontalement)
	fx := sx + ts/2 + float32(math.Sin(t)*tsf*0.07)
	fy := by - ts*0.05

	// Glow chaud sur le sol en dessous (large, très transparent)
	glowR := ts * 0.85
	vector.DrawFilledCircle(screen, fx, fy+ts*0.3, glowR, color.RGBA{200, 100, 20, 18}, false)

	// Corps externe de la flamme (orange profond)
	outerR := fmax(ts*0.19+float32(math.Sin(t*1.4)*tsf*0.025), ts*0.10)
	vector.DrawFilledCircle(screen, fx, fy, outerR, color.RGBA{210, 85, 15, 200}, false)

	// Corps central (orange vif)
	midR := outerR * 0.62
	midFx := fx + float32(math.Sin(t*0.9)*tsf*0.04)
	midFy := fy - outerR*0.15
	vector.DrawFilledCircle(screen, midFx, midFy, midR, color.RGBA{240, 150, 25, 220}, false)

	// Cœur de la flamme (jaune chaud)
	coreR := midR * 0.55
	vector.DrawFilledCircle(screen, midFx, midFy-midR*0.2, coreR, color.RGBA{255, 220, 80, 235}, false)

	// Étincelle (petit point brillant, position aléatoire via sinus)
	sparkX := fx + float32(math.Sin(t*2.3)*tsf*0.09)
	sparkY := fy - outerR*0.9 + float32(math.Cos(t*1.8)*tsf*0.04)
	sparkR := fmax(ts*0.045, 1)
	vector.DrawFilledCircle(screen, sparkX, sparkY, sparkR, color.RGBA{255, 245, 180, 200}, false)
}

// drawWallInner dessine un mur interne (plafond, entouré d'autres murs).
// Très sombre avec une texture minimale.
func drawWallInner(screen *ebiten.Image, pal dungeonPalette, sx, sy, ts float32, tileX, tileY int) {
	vector.DrawFilledRect(screen, sx, sy, ts, ts, pal.wallInner, false)

	// Légère variation pour éviter un aplat uniforme
	v := tileHash(tileX, tileY)
	if v < 35 {
		detail := color.RGBA{
			R: uint8(clampByte(int(pal.wallInner.R) + 8)),
			G: uint8(clampByte(int(pal.wallInner.G) + 8)),
			B: uint8(clampByte(int(pal.wallInner.B) + 8)),
			A: 255,
		}
		jw := fmax(ts/32, 1)
		lineY := sy + ts*float32(v%3+1)/4
		vector.DrawFilledRect(screen, sx+ts*0.15, lineY, ts*0.7, jw, detail, false)
	}
}

// drawDungeonFloor dessine une dalle de sol avec un motif de pierres apparentes.
// Ajoute une ombre portée si un mur est au nord (northWall=true).
func drawDungeonFloor(screen *ebiten.Image, pal dungeonPalette, sx, sy, ts float32, tileX, tileY int, northWall bool) {
	// Variation subtile de la couleur de base par tuile
	v := tileHash(tileX, tileY)
	vf := int(v&7) - 3 // -3 à +4
	base := color.RGBA{
		R: uint8(clampByte(int(pal.floorBase.R) + vf)),
		G: uint8(clampByte(int(pal.floorBase.G) + vf)),
		B: uint8(clampByte(int(pal.floorBase.B) + vf)),
		A: 255,
	}

	// Remplissage de base
	vector.DrawFilledRect(screen, sx, sy, ts, ts, base, false)

	// Motif de dalles de pierre (joints en brique décalée)
	// Conçu pour 32px, mis à l'échelle via s = ts/32
	s := ts / 32.0
	jw := fmax(s, 1) // épaisseur d'un joint

	h1 := sy + 10*s // première ligne horizontale
	h2 := sy + 21*s // deuxième ligne horizontale

	// Joints horizontaux
	vector.DrawFilledRect(screen, sx, h1, ts, jw, pal.floorJoint, false)
	vector.DrawFilledRect(screen, sx, h2, ts, jw, pal.floorJoint, false)

	// Joints verticaux rangée 1 (sy → h1)
	vector.DrawFilledRect(screen, sx+16*s, sy, jw, 10*s, pal.floorJoint, false)
	// Joints verticaux rangée 2 (h1 → h2)
	vector.DrawFilledRect(screen, sx+6*s, h1, jw, 11*s, pal.floorJoint, false)
	vector.DrawFilledRect(screen, sx+22*s, h1, jw, 11*s, pal.floorJoint, false)
	// Joints verticaux rangée 3 (h2 → bas)
	vector.DrawFilledRect(screen, sx+10*s, h2, jw, 11*s+jw, pal.floorJoint, false)
	vector.DrawFilledRect(screen, sx+26*s, h2, jw, 11*s+jw, pal.floorJoint, false)

	// Ombre portée par le mur nord (dégradé sur le haut de la tuile)
	if northWall {
		shadowH := ts * 0.38
		// Bande sombre principale
		vector.DrawFilledRect(screen, sx, sy, ts, shadowH*0.55, pal.floorShadow, false)
		// Bande de transition (moitié moins opaque)
		soft := color.RGBA{pal.floorShadow.R, pal.floorShadow.G, pal.floorShadow.B, pal.floorShadow.A / 2}
		vector.DrawFilledRect(screen, sx, sy+shadowH*0.55, ts, shadowH*0.45, soft, false)
	}
}
