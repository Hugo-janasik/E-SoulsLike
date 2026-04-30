package entities

import (
	"e-soulslike/world"
	"github.com/hajimehoshi/ebiten/v2"
)

// Entity est l'interface de base pour toutes les entités du jeu
type Entity interface {
	Update()
	Draw(screen *ebiten.Image, camera *world.Camera)
	GetPosition() (float64, float64)
	GetHealth() int
	TakeDamage(damage int)
	IsAlive() bool
}

// BaseEntity contient les propriétés communes à toutes les entités
type BaseEntity struct {
	X, Y            float64
	Width, Height   float64
	Health          int
	MaxHealth       int
	Speed           float64
	IsInvincible    bool
	InvincibleTimer int
}

// GetPosition retourne la position de l'entité
func (e *BaseEntity) GetPosition() (float64, float64) {
	return e.X, e.Y
}

// GetHealth retourne la santé actuelle
func (e *BaseEntity) GetHealth() int {
	return e.Health
}

// TakeDamage inflige des dégâts à l'entité
func (e *BaseEntity) TakeDamage(damage int) {
	if e.IsInvincible {
		return
	}
	e.Health -= damage
	if e.Health < 0 {
		e.Health = 0
	}
	// Activer l'invincibilité temporaire
	e.IsInvincible = true
	e.InvincibleTimer = 60 // 1 seconde à 60 FPS
}

// IsAlive retourne true si l'entité est en vie
func (e *BaseEntity) IsAlive() bool {
	return e.Health > 0
}

// UpdateInvincibility met à jour le timer d'invincibilité
func (e *BaseEntity) UpdateInvincibility() {
	if e.IsInvincible {
		e.InvincibleTimer--
		if e.InvincibleTimer <= 0 {
			e.IsInvincible = false
		}
	}
}
