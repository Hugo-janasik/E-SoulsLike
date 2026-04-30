package world

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// hubPalette regroupe toutes les couleurs du hub.
type hubPalette struct {
	plazaBase   color.RGBA // pierre chaude de la plaza
	plazaJoint  color.RGBA // joints entre dalles de la plaza
	innerBase   color.RGBA // cercle intérieur autour du feu
	innerJoint  color.RGBA
	grassBase   color.RGBA // herbe extérieure
	grassLight  color.RGBA // herbe plus claire (variation)
	grassDark   color.RGBA // herbe plus sombre (ombre)
	waterBase   color.RGBA // eau sombre
	waterSheen  color.RGBA // reflet lumineux sur l'eau
	wallTop     color.RGBA // dessus de la ruine
	wallFace    color.RGBA // face avant de la ruine
	wallShadow  color.RGBA // ombre pied/côté de la ruine
	wallInner   color.RGBA // mur interne (non visible)
	floorShadow color.RGBA // ombre portée des ruines sur le sol
}

var hubPal = hubPalette{
	plazaBase:   color.RGBA{118, 105, 88, 255},
	plazaJoint:  color.RGBA{88, 78, 65, 255},
	innerBase:   color.RGBA{98, 82, 62, 255},
	innerJoint:  color.RGBA{72, 60, 46, 255},
	grassBase:   color.RGBA{46, 66, 28, 255},
	grassLight:  color.RGBA{62, 88, 38, 255},
	grassDark:   color.RGBA{32, 48, 18, 255},
	waterBase:   color.RGBA{28, 40, 60, 255},
	waterSheen:  color.RGBA{50, 72, 105, 255},
	wallTop:     color.RGBA{135, 122, 102, 255},
	wallFace:    color.RGBA{82, 72, 58, 255},
	wallShadow:  color.RGBA{38, 33, 26, 255},
	wallInner:   color.RGBA{52, 46, 36, 255},
	floorShadow: color.RGBA{0, 0, 0, 45},
}

// drawHubTile dispatche le rendu d'une tuile du hub selon son type et son contexte.
func drawHubTile(
	screen *ebiten.Image,
	tile *Tile,
	southTile *Tile,
	northWall bool,
	sx, sy, ts float32,
	tileX, tileY, tick int,
) {
	switch tile.Type {
	case TileWall:
		southIsFloor := southTile != nil && southTile.Type != TileWall
		if southIsFloor {
			drawHubWallFace(screen, sx, sy, ts, tileX, tileY)
		} else {
			drawHubWallInner(screen, sx, sy, ts, tileX, tileY)
		}
	case TileWater:
		drawHubWater(screen, sx, sy, ts, tileX, tileY, tick)
	case TileStone:
		drawHubInnerStone(screen, sx, sy, ts, tileX, tileY, northWall)
	case TileDirt:
		drawHubPlaza(screen, sx, sy, ts, tileX, tileY, northWall)
	default: // TileGrass et tout le reste → herbe
		drawHubGrass(screen, sx, sy, ts, tileX, tileY)
	}
}

// ── Herbe ────────────────────────────────────────────────────────────────────

func drawHubGrass(screen *ebiten.Image, sx, sy, ts float32, tileX, tileY int) {
	v := tileHash(tileX, tileY)

	// Couleur de base avec variation par tuile
	vf := int(v&7) - 3
	base := color.RGBA{
		R: uint8(clampByte(int(hubPal.grassBase.R) + vf)),
		G: uint8(clampByte(int(hubPal.grassBase.G) + vf*2)),
		B: uint8(clampByte(int(hubPal.grassBase.B) + vf)),
		A: 255,
	}
	vector.DrawFilledRect(screen, sx, sy, ts, ts, base, false)

	s := ts / 32.0

	// Quelques brins d'herbe dessinés comme de petits rectangles verticaux
	if v < 200 {
		bladeColor := hubPal.grassLight
		numBlades := 2 + (v % 3)
		for i := 0; i < int(numBlades); i++ {
			bv := tileHash(tileX*7+i, tileY*11+i)
			bx := sx + float32(bv%28+2)*s
			bh := float32(3+bv%4) * s
			by := sy + ts - bh
			bw := fmax(s*1.2, 1)
			bColor := color.RGBA{
				R: uint8(clampByte(int(bladeColor.R) + int(bv%10) - 5)),
				G: uint8(clampByte(int(bladeColor.G) + int(bv%14) - 7)),
				B: bladeColor.B,
				A: 200,
			}
			vector.DrawFilledRect(screen, bx, by, bw, bh, bColor, false)
		}
	}

	// Petits cailloux dans ~15% des cases
	if v > 220 {
		stoneColor := color.RGBA{85, 78, 65, 160}
		bv := tileHash(tileX*3, tileY*5)
		rx := sx + float32(bv%24+4)*s
		ry := sy + float32((bv/31)%22+4)*s
		rr := fmax(s*1.8, 1.5)
		vector.DrawFilledCircle(screen, rx, ry, rr, stoneColor, false)
	}
}

// ── Plaza (pierre chaude) ─────────────────────────────────────────────────────

func drawHubPlaza(screen *ebiten.Image, sx, sy, ts float32, tileX, tileY int, northWall bool) {
	v := tileHash(tileX, tileY)
	vf := int(v&5) - 2
	base := color.RGBA{
		R: uint8(clampByte(int(hubPal.plazaBase.R) + vf)),
		G: uint8(clampByte(int(hubPal.plazaBase.G) + vf)),
		B: uint8(clampByte(int(hubPal.plazaBase.B) + vf)),
		A: 255,
	}
	vector.DrawFilledRect(screen, sx, sy, ts, ts, base, false)

	// Motif de grandes dalles : 2 rangées horizontales (plus larges que donjon)
	s := ts / 32.0
	jw := fmax(s, 1)

	h1 := sy + 15*s
	vector.DrawFilledRect(screen, sx, h1, ts, jw, hubPal.plazaJoint, false)

	// Joints verticaux — rangée 1
	vector.DrawFilledRect(screen, sx+18*s, sy, jw, 15*s, hubPal.plazaJoint, false)
	// Joints verticaux — rangée 2
	vector.DrawFilledRect(screen, sx+8*s, h1, jw, 17*s+jw, hubPal.plazaJoint, false)
	vector.DrawFilledRect(screen, sx+24*s, h1, jw, 17*s+jw, hubPal.plazaJoint, false)

	// Ombre portée depuis un mur au nord
	if northWall {
		shadow := hubPal.floorShadow
		vector.DrawFilledRect(screen, sx, sy, ts, ts*0.35, shadow, false)
		softer := color.RGBA{shadow.R, shadow.G, shadow.B, shadow.A / 2}
		vector.DrawFilledRect(screen, sx, sy+ts*0.35, ts, ts*0.18, softer, false)
	}
}

// ── Pierre intérieure (cercle du feu) ─────────────────────────────────────────

func drawHubInnerStone(screen *ebiten.Image, sx, sy, ts float32, tileX, tileY int, northWall bool) {
	v := tileHash(tileX, tileY)
	vf := int(v&7) - 3
	base := color.RGBA{
		R: uint8(clampByte(int(hubPal.innerBase.R) + vf)),
		G: uint8(clampByte(int(hubPal.innerBase.G) + vf)),
		B: uint8(clampByte(int(hubPal.innerBase.B) + vf)),
		A: 255,
	}
	vector.DrawFilledRect(screen, sx, sy, ts, ts, base, false)

	s := ts / 32.0
	jw := fmax(s, 1)

	// Dalles en éventail (joints concentriques approximés par des lignes décalées)
	h1 := sy + 14*s
	vector.DrawFilledRect(screen, sx, h1, ts, jw, hubPal.innerJoint, false)
	vector.DrawFilledRect(screen, sx+12*s, sy, jw, 14*s, hubPal.innerJoint, false)
	vector.DrawFilledRect(screen, sx+20*s, h1, jw, 18*s+jw, hubPal.innerJoint, false)

	if northWall {
		shadow := hubPal.floorShadow
		vector.DrawFilledRect(screen, sx, sy, ts, ts*0.30, shadow, false)
	}
}

// ── Eau ───────────────────────────────────────────────────────────────────────

func drawHubWater(screen *ebiten.Image, sx, sy, ts float32, tileX, tileY, tick int) {
	v := tileHash(tileX, tileY)
	vf := int(v&5) - 2
	base := color.RGBA{
		R: uint8(clampByte(int(hubPal.waterBase.R) + vf)),
		G: uint8(clampByte(int(hubPal.waterBase.G) + vf)),
		B: uint8(clampByte(int(hubPal.waterBase.B) + vf*2)),
		A: 255,
	}
	vector.DrawFilledRect(screen, sx, sy, ts, ts, base, false)

	// Reflet animé : une ou deux lignes lumineuses qui oscillent
	s := ts / 32.0
	t := float64(tick)*0.04 + float64(v)*0.12
	sheenY := sy + ts*0.4 + float32(math.Sin(t)*float64(ts)*0.12)
	sheenH := fmax(s, 1)
	sheenW := ts * (0.4 + float32(math.Sin(t*0.7+1.2)*0.15))
	sheenX := sx + (ts-sheenW)*0.5
	sheen := color.RGBA{hubPal.waterSheen.R, hubPal.waterSheen.G, hubPal.waterSheen.B, 55}
	vector.DrawFilledRect(screen, sheenX, sheenY, sheenW, sheenH, sheen, false)

	// Second reflet décalé
	sheen2Y := sy + ts*0.65 + float32(math.Cos(t*0.8+0.5)*float64(ts)*0.08)
	sheenW2 := ts * 0.25
	sheen2 := color.RGBA{hubPal.waterSheen.R, hubPal.waterSheen.G, hubPal.waterSheen.B, 30}
	vector.DrawFilledRect(screen, sx+ts*0.3, sheen2Y, sheenW2, sheenH, sheen2, false)
}

// ── Ruines : face visible ─────────────────────────────────────────────────────

func drawHubWallFace(screen *ebiten.Image, sx, sy, ts float32, tileX, tileY int) {
	capH := ts * 7 / 32

	// Dessus de la ruine (légèrement moussu → teinte verte)
	topColor := color.RGBA{
		R: uint8(clampByte(int(hubPal.wallTop.R) - 5)),
		G: uint8(clampByte(int(hubPal.wallTop.G) + 8)),
		B: uint8(clampByte(int(hubPal.wallTop.B) - 8)),
		A: 255,
	}
	vector.DrawFilledRect(screen, sx, sy, ts, capH, topColor, false)

	// Face de pierre
	vector.DrawFilledRect(screen, sx, sy+capH, ts, ts-capH, hubPal.wallFace, false)

	// Ombre gauche
	shadowW := fmax(ts*2/32, 1)
	vector.DrawFilledRect(screen, sx, sy+capH, shadowW, ts-capH, hubPal.wallShadow, false)

	// Ombre bas
	shadowH := fmax(ts*3/32, 1)
	vector.DrawFilledRect(screen, sx, sy+ts-shadowH, ts, shadowH, hubPal.wallShadow, false)

	// Fissure horizontale (ruine un peu endommagée)
	jw := fmax(ts/32, 1)
	crackY := sy + capH + (ts-capH)*0.50
	crackColor := color.RGBA{hubPal.wallShadow.R, hubPal.wallShadow.G, hubPal.wallShadow.B, 140}
	vector.DrawFilledRect(screen, sx+shadowW, crackY, ts-shadowW*2, jw, crackColor, false)

	// Mousse dans les fissures (vert subtil)
	v := tileHash(tileX, tileY)
	if v < 80 {
		mossColor := color.RGBA{45, 72, 25, 90}
		vector.DrawFilledRect(screen, sx+shadowW+ts*0.1, crackY-jw, ts*0.4, jw*2, mossColor, false)
	}
}

// ── Ruines : mur interne ──────────────────────────────────────────────────────

func drawHubWallInner(screen *ebiten.Image, sx, sy, ts float32, tileX, tileY int) {
	vector.DrawFilledRect(screen, sx, sy, ts, ts, hubPal.wallInner, false)

	v := tileHash(tileX, tileY)
	if v < 40 {
		detail := color.RGBA{
			R: uint8(clampByte(int(hubPal.wallInner.R) + 10)),
			G: uint8(clampByte(int(hubPal.wallInner.G) + 10)),
			B: uint8(clampByte(int(hubPal.wallInner.B) + 10)),
			A: 255,
		}
		jw := fmax(ts/32, 1)
		lineY := sy + ts*float32(v%3+1)/4
		vector.DrawFilledRect(screen, sx+ts*0.1, lineY, ts*0.8, jw, detail, false)
	}
}
