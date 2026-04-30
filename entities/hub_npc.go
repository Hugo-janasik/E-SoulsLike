package entities

import (
	"e-soulslike/world"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// HubNPC représente un PNJ décoratif statique dans le hub.
// Non interactif pour l'instant — présence visuelle uniquement.
type HubNPC struct {
	X, Y   float64
	Name   string
	sprite *ebiten.Image
}

// NewHubNPC crée un PNJ décoratif à la position donnée.
// spritePath : chemin vers le PNG du sprite (relatif au dossier du jeu).
func NewHubNPC(x, y float64, name, spritePath string) *HubNPC {
	npc := &HubNPC{X: x, Y: y, Name: name}
	if img, _, err := ebitenutil.NewImageFromFile(spritePath); err == nil {
		npc.sprite = img
	}
	return npc
}

// Draw dessine le PNJ à l'écran avec la caméra donnée.
func (n *HubNPC) Draw(screen *ebiten.Image, camera *world.Camera) {
	sx, sy := camera.WorldToScreen(n.X, n.Y)

	if n.sprite != nil {
		w, h := n.sprite.Size()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(w)/2, -float64(h))
		op.GeoM.Translate(sx, sy)
		screen.DrawImage(n.sprite, op)
	} else {
		// Placeholder : petit cercle coloré si le sprite n'a pas chargé
		vector.DrawFilledCircle(screen, float32(sx), float32(sy)-16, 10, color.RGBA{180, 150, 100, 200}, false)
		vector.DrawFilledRect(screen, float32(sx)-6, float32(sy)-16, 12, 20, color.RGBA{140, 110, 70, 200}, false)
	}
}

// Update est un no-op pour l'instant (les HubNPC sont statiques).
func (n *HubNPC) Update(_, _ float64) {}
