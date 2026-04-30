package entities

import (
	"e-soulslike/world"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
)

// EnemyType définit les différents types d'ennemis
type EnemyType int

const (
	EnemyTypeBasic EnemyType = iota
	EnemyTypeFast
	EnemyTypeTank
)

// EnemyState définit l'état de l'ennemi
type EnemyState int

const (
	StateIdle EnemyState = iota
	StateChasing
	StateAttacking
	StateStunned
)

// Enemy représente un ennemi
type Enemy struct {
	BaseEntity
	Type             EnemyType
	State            EnemyState
	TargetX, TargetY float64
	AttackRange      float64
	DetectionRange   float64
	AttackCooldown   int
	AttackTimer      int
	StunTimer        int
	WanderTimer      int
	WanderTargetX    float64
	WanderTargetY    float64

	// Sprites et animation
	SpritesIdle   []*ebiten.Image
	SpritesMove   []*ebiten.Image
	SpriteAttack  *ebiten.Image
	AnimationFrame int
	AnimationSpeed int
	AnimationTimer int
}

// NewEnemy crée un nouvel ennemi
func NewEnemy(x, y float64, enemyType EnemyType) *Enemy {
	enemy := &Enemy{
		BaseEntity: BaseEntity{
			X:      x,
			Y:      y,
			Width:  28,
			Height: 28,
			Speed:  1.5,
		},
		Type:           enemyType,
		State:          StateIdle,
		AttackRange:    40,
		DetectionRange: 200,
		AttackCooldown: 60,
		WanderTimer:    rand.Intn(120) + 60,
	}

	// Configurer selon le type
	switch enemyType {
	case EnemyTypeBasic:
		enemy.Health = 50
		enemy.MaxHealth = 50
		enemy.Speed = 1.5
	case EnemyTypeFast:
		enemy.Health = 30
		enemy.MaxHealth = 30
		enemy.Speed = 2.5
		enemy.DetectionRange = 250
	case EnemyTypeTank:
		enemy.Health = 150
		enemy.MaxHealth = 150
		enemy.Speed = 1.0
		enemy.Width = 40
		enemy.Height = 40
	}

	enemy.WanderTargetX = x
	enemy.WanderTargetY = y

	// Charger les sprites
	enemy.loadSprites()
	enemy.AnimationSpeed = 10
	enemy.AnimationFrame = 0
	enemy.AnimationTimer = 0

	return enemy
}

// loadSprites charge les sprites de l'ennemi
func (e *Enemy) loadSprites() {
	// Pour l'instant, tous les ennemis utilisent les sprites "spirit"
	// On pourra ajouter d'autres types plus tard
	e.SpritesIdle = loadEnemySpriteFrames("assets/monsters/spirit/idle/%d.png", 4)
	e.SpritesMove = loadEnemySpriteFrames("assets/monsters/spirit/move/%d.png", 4)

	// Sprite d'attaque (une seule frame)
	attackSprite, _, err := ebitenutil.NewImageFromFile("assets/monsters/spirit/attack/0.png")
	if err == nil {
		e.SpriteAttack = attackSprite
	}
}

// loadEnemySpriteFrames charge plusieurs frames d'animation
func loadEnemySpriteFrames(pattern string, count int) []*ebiten.Image {
	sprites := make([]*ebiten.Image, 0, count)
	for i := 0; i < count; i++ {
		path := fmt.Sprintf(pattern, i)
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err == nil {
			sprites = append(sprites, img)
		}
	}
	return sprites
}

// Update met à jour l'ennemi
func (e *Enemy) Update(playerX, playerY float64) {
	e.UpdateInvincibility()

	// Mettre à jour l'animation
	e.AnimationTimer++
	if e.AnimationTimer >= e.AnimationSpeed {
		e.AnimationTimer = 0
		e.AnimationFrame++
		// Boucler l'animation selon l'état
		var maxFrames int
		if e.State == StateChasing {
			maxFrames = len(e.SpritesMove)
		} else {
			maxFrames = len(e.SpritesIdle)
		}
		if maxFrames > 0 && e.AnimationFrame >= maxFrames {
			e.AnimationFrame = 0
		}
	}

	// Gérer le stun
	if e.State == StateStunned {
		e.StunTimer--
		if e.StunTimer <= 0 {
			e.State = StateIdle
		}
		return
	}

	// Calculer la distance au joueur
	dx := playerX - e.X
	dy := playerY - e.Y
	distanceToPlayer := math.Sqrt(dx*dx + dy*dy)

	// Machine à états
	switch e.State {
	case StateIdle:
		// Comportement de vagabondage
		e.WanderTimer--
		if e.WanderTimer <= 0 {
			e.WanderTargetX = e.X + (rand.Float64()-0.5)*200
			e.WanderTargetY = e.Y + (rand.Float64()-0.5)*200
			e.WanderTimer = rand.Intn(120) + 60
		}

		// Se déplacer vers la cible de vagabondage
		wx := e.WanderTargetX - e.X
		wy := e.WanderTargetY - e.Y
		wanderDist := math.Sqrt(wx*wx + wy*wy)
		if wanderDist > 5 {
			e.X += (wx / wanderDist) * e.Speed * 0.5
			e.Y += (wy / wanderDist) * e.Speed * 0.5
		}

		// Détecter le joueur
		if distanceToPlayer < e.DetectionRange {
			e.State = StateChasing
		}

	case StateChasing:
		// Poursuivre le joueur
		if distanceToPlayer > e.AttackRange {
			e.X += (dx / distanceToPlayer) * e.Speed
			e.Y += (dy / distanceToPlayer) * e.Speed
		} else {
			e.State = StateAttacking
		}

		// Perdre le joueur de vue
		if distanceToPlayer > e.DetectionRange*1.5 {
			e.State = StateIdle
		}

	case StateAttacking:
		// Attaquer
		if e.AttackTimer > 0 {
			e.AttackTimer--
		}

		if distanceToPlayer > e.AttackRange {
			e.State = StateChasing
		}
	}
}

// Draw dessine l'ennemi
func (e *Enemy) Draw(screen *ebiten.Image, camera *world.Camera) {
	if !e.IsAlive() {
		return
	}

	screenX, screenY := camera.WorldToScreen(e.X, e.Y)

	// Sélectionner le sprite selon l'état
	var sprite *ebiten.Image
	if e.State == StateAttacking && e.SpriteAttack != nil {
		sprite = e.SpriteAttack
	} else if e.State == StateChasing && len(e.SpritesMove) > 0 {
		frame := e.AnimationFrame % len(e.SpritesMove)
		sprite = e.SpritesMove[frame]
	} else if len(e.SpritesIdle) > 0 {
		frame := e.AnimationFrame % len(e.SpritesIdle)
		sprite = e.SpritesIdle[frame]
	}

	// Dessiner le sprite si disponible
	if sprite != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := sprite.Bounds()
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		op.GeoM.Translate(screenX, screenY)

		// Effets visuels
		if e.IsInvincible && e.InvincibleTimer%10 < 5 {
			op.ColorScale.Scale(1, 1, 1, 0.5)
		}
		if e.State == StateStunned {
			op.ColorScale.Scale(0.6, 0.6, 1.0, 0.8)
		}

		screen.DrawImage(sprite, op)
	} else {
		// Fallback: dessiner un cercle si les sprites ne sont pas chargés
		var enemyColor color.RGBA
		switch e.Type {
		case EnemyTypeBasic:
			enemyColor = color.RGBA{255, 100, 100, 255}
		case EnemyTypeFast:
			enemyColor = color.RGBA{255, 200, 100, 255}
		case EnemyTypeTank:
			enemyColor = color.RGBA{150, 50, 50, 255}
		}
		if e.IsInvincible && e.InvincibleTimer%10 < 5 {
			enemyColor = color.RGBA{255, 255, 255, 200}
		}
		if e.State == StateStunned {
			enemyColor = color.RGBA{200, 200, 255, 200}
		}
		vector.DrawFilledCircle(screen, float32(screenX), float32(screenY), float32(e.Width/2), enemyColor, true)
	}

	// Indicateur d'état
	if e.State == StateChasing || e.State == StateAttacking {
		// Point d'exclamation au-dessus de la tête
		vector.DrawFilledCircle(screen, float32(screenX), float32(screenY-e.Height), 3, color.RGBA{255, 255, 0, 255}, true)
	}

	// Barre de vie
	if e.Health < e.MaxHealth {
		barWidth := e.Width * 1.2
		barHeight := 3.0
		barX := screenX - barWidth/2
		barY := screenY - e.Height - 8

		// Fond
		vector.DrawFilledRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{50, 50, 50, 255}, true)
		// Vie
		healthWidth := barWidth * float64(e.Health) / float64(e.MaxHealth)
		vector.DrawFilledRect(screen, float32(barX), float32(barY), float32(healthWidth), float32(barHeight), color.RGBA{255, 50, 50, 255}, true)
	}
}

// CanAttack vérifie si l'ennemi peut attaquer
func (e *Enemy) CanAttack() bool {
	return e.State == StateAttacking && e.AttackTimer <= 0
}

// PerformAttack lance une attaque
func (e *Enemy) PerformAttack() {
	e.AttackTimer = e.AttackCooldown
}

// GetAttackDamage retourne les dégâts de l'attaque
func (e *Enemy) GetAttackDamage() int {
	switch e.Type {
	case EnemyTypeBasic:
		return 10
	case EnemyTypeFast:
		return 7
	case EnemyTypeTank:
		return 20
	default:
		return 10
	}
}

// Stun étourdit l'ennemi
func (e *Enemy) Stun(duration int) {
	e.State = StateStunned
	e.StunTimer = duration
}

// Respawn réinitialise l'ennemi à son état initial (utilisé au respawn du joueur)
func (e *Enemy) Respawn() {
	e.Health = e.MaxHealth
	e.State = StateIdle
	e.AttackTimer = 0
	e.StunTimer = 0
	e.IsInvincible = false
	e.InvincibleTimer = 0
	e.WanderTimer = rand.Intn(120) + 60
}
