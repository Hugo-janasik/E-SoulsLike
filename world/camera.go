package world

import "math"

// Camera gère la vue de la caméra 2D
type Camera struct {
	X, Y          float64
	Width, Height float64
	Zoom          float64
	FollowSpeed   float64
}

// NewCamera crée une nouvelle caméra
func NewCamera(width, height float64) *Camera {
	return &Camera{
		X:           0,
		Y:           0,
		Width:       width,
		Height:      height,
		Zoom:        1.0,
		FollowSpeed: 0.1,
	}
}

// Follow fait suivre la caméra à une position (smooth follow)
func (c *Camera) Follow(targetX, targetY float64) {
	// Position cible de la caméra (centrée sur la cible)
	targetCamX := targetX - c.Width/2
	targetCamY := targetY - c.Height/2

	// Interpolation linéaire pour un mouvement fluide
	c.X += (targetCamX - c.X) * c.FollowSpeed
	c.Y += (targetCamY - c.Y) * c.FollowSpeed
}

// WorldToScreen convertit des coordonnées monde en coordonnées écran
func (c *Camera) WorldToScreen(worldX, worldY float64) (float64, float64) {
	screenX := (worldX - c.X) * c.Zoom
	screenY := (worldY - c.Y) * c.Zoom
	return screenX, screenY
}

// ScreenToWorld convertit des coordonnées écran en coordonnées monde
func (c *Camera) ScreenToWorld(screenX, screenY float64) (float64, float64) {
	worldX := screenX/c.Zoom + c.X
	worldY := screenY/c.Zoom + c.Y
	return worldX, worldY
}

// IsVisible vérifie si un rectangle est visible par la caméra
func (c *Camera) IsVisible(x, y, width, height float64) bool {
	screenX, screenY := c.WorldToScreen(x, y)

	return screenX+width*c.Zoom >= 0 &&
		screenX <= c.Width &&
		screenY+height*c.Zoom >= 0 &&
		screenY <= c.Height
}

// Shake fait trembler la caméra (utile pour les impacts)
func (c *Camera) Shake(intensity float64) {
	// TODO: Implémenter le shake de caméra
}

// SetZoom définit le niveau de zoom
func (c *Camera) SetZoom(zoom float64) {
	c.Zoom = math.Max(0.5, math.Min(3.0, zoom))
}
