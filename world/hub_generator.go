package world

import "math"

// GenerateHubTilemap génère la tilemap fixe de la Place du Feu (hub).
//
// Layout (100×100, centre à (50,50)) :
//   - Plaza centrale circulaire (radius 12) en pierre chaude
//   - Cercle intérieur (radius 5) autour du feu
//   - 3 chambres rectangulaires centrées exactement sur les portails :
//     Ouest (31,50) · Nord (50,31) · Est (69,50)
//   - Corridors larges reliant directement la plaza à chaque chambre
//   - 4 pools d'eau dans les quadrants entre les corridors
//   - Herbe en dehors des zones piétonnes, aucune ruine
func GenerateHubTilemap(width, height int) *TileMap {
	tm := &TileMap{
		Width:         width,
		Height:        height,
		Tiles:         make([][]Tile, height),
		WalkableWater: true,
		IsDungeon:     true,
		Theme:         ThemeHub,
	}
	for y := 0; y < height; y++ {
		tm.Tiles[y] = make([]Tile, width)
	}

	cx, cy := width/2, height/2 // (50, 50)

	// ── 1. Base : tout en herbe ───────────────────────────────────────────────
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			tm.Tiles[y][x].Type = TileGrass
		}
	}

	// ── 2. Pools d'eau dans les 4 quadrants (avant les zones piétonnes) ───────
	// Positionnés entre les 3 corridors pour remplir les espaces verts
	hubFillCircle(tm, cx-20, cy-20, 7, TileWater) // NW
	hubFillCircle(tm, cx+20, cy-20, 7, TileWater) // NE
	hubFillCircle(tm, cx-20, cy+20, 7, TileWater) // SW
	hubFillCircle(tm, cx+20, cy+20, 7, TileWater) // SE

	// ── 3. Chambres des portails (centrées exactement sur le portail) ─────────
	// Les portails sont en world coords à ±600 du centre (1600,1600) → ±18.75 tiles
	// Arrondi : Ouest=(31,50), Nord=(50,31), Est=(69,50)

	// Chambre Ouest  — portail à tile (31, 50)
	// Taille 16×14 → x: 23–38, y: 43–57   centre: (30.5, 50) ≈ (31, 50) ✓
	hubFillRect(tm, 23, 43, 38, 57, TileDirt)

	// Chambre Nord — portail à tile (50, 31)
	// Taille 14×16 → x: 43–57, y: 23–38   centre: (50, 30.5) ≈ (50, 31) ✓
	hubFillRect(tm, 43, 23, 57, 38, TileDirt)

	// Chambre Est — portail à tile (69, 50)
	// Taille 16×14 → x: 62–77, y: 43–57   centre: (69.5, 50) ≈ (69, 50) ✓
	hubFillRect(tm, 62, 43, 77, 57, TileDirt)

	// ── 4. Corridors reliant directement la plaza aux chambres ───────────────
	// Largeur 8 tiles (±4 du centre) pour un accès confortable

	// Corridor Ouest : de la chambre (x=38) jusqu'à la plaza (x=38 chevauchement)
	// On tire un rectangle large de la chambre jusqu'au centre, la plaza l'écrase
	hubFillRect(tm, 23, cy-4, cx, cy+4, TileDirt)

	// Corridor Nord
	hubFillRect(tm, cx-4, 23, cx+4, cy, TileDirt)

	// Corridor Est
	hubFillRect(tm, cx, cy-4, 77, cy+4, TileDirt)

	// ── 5. Plaza centrale (circle radius 12) ─────────────────────────────────
	// Dessinée APRÈS les corridors pour les écraser proprement avec TileDirt
	hubFillCircle(tm, cx, cy, 12, TileDirt)

	// ── 6. Cercle intérieur du feu (radius 5) — pierre spéciale ──────────────
	hubFillCircle(tm, cx, cy, 5, TileStone)

	return tm
}

// ─── Helpers géométriques ─────────────────────────────────────────────────────

func hubFillCircle(tm *TileMap, cx, cy, radius int, t TileType) {
	r2 := radius * radius
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			if x < 0 || x >= tm.Width || y < 0 || y >= tm.Height {
				continue
			}
			dx, dy := x-cx, y-cy
			if dx*dx+dy*dy <= r2 {
				tm.Tiles[y][x].Type = t
			}
		}
	}
}

func hubFillRect(tm *TileMap, x1, y1, x2, y2 int, t TileType) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			if x < 0 || x >= tm.Width || y < 0 || y >= tm.Height {
				continue
			}
			tm.Tiles[y][x].Type = t
		}
	}
}

// abs retourne la valeur absolue d'un entier
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// HubNPCPositions retourne les positions monde des PNJs décoratifs.
// Chaque PNJ est placé dans ou près de la chambre de son portail.
func HubNPCPositions() []struct {
	X, Y float64
	Name string
} {
	return []struct {
		X, Y float64
		Name string
	}{
		// Marchand — chambre Ouest, côté sud de la pièce
		{X: 1000, Y: 1720, Name: "Marchand"},
		// Guerrier blessé — chambre Nord, côté est
		{X: 1680, Y: 1040, Name: "Guerrier blessé"},
		// Voyageur — chambre Est, côté nord
		{X: 2180, Y: 1490, Name: "Voyageur mystérieux"},
	}
}

// HubFirePosition retourne la position monde du feu central.
func HubFirePosition() (float64, float64) {
	return 1600, 1600
}

// HubDistFromCenter retourne la distance euclidienne d'une position monde au centre du hub.
func HubDistFromCenter(worldX, worldY float64) float64 {
	dx := worldX - 1600
	dy := worldY - 1600
	return math.Sqrt(dx*dx + dy*dy)
}
