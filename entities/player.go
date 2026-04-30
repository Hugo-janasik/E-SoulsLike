package entities

import (
	"e-soulslike/input"
	"e-soulslike/world"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PlayerDirection représente la direction du joueur
type PlayerDirection int

const (
	DirectionDown  PlayerDirection = iota
	DirectionUp
	DirectionLeft
	DirectionRight
)

// Constantes du joueur — remplace les nombres magiques dispersés dans le code
const (
	PlayerSpeed = 3.0 // Vitesse de déplacement de base (pixels/frame)
	SprintMultiplier = 1.8 // Multiplicateur de vitesse en sprint

	AnimationSpeedWalk  = 6  // Frames entre chaque sprite en marche (~10 fps à 60 fps)
	AnimationSpeedIdle  = 10 // Frames entre chaque sprite en idle (~6 fps à 60 fps)
	AnimationFrameCount = 4  // Nombre de frames par animation

	AttackDuration  = 30 // Durée d'une attaque en frames
	StaminaCostAttack = 10 // Coût en stamina d'une attaque
	StaminaMinAttack  = 20 // Stamina minimale pour attaquer

	DodgeDuration   = 20 // Durée d'une esquive en frames
	StaminaCostDodge = 15 // Coût en stamina d'une esquive
	StaminaMinDodge  = 15 // Stamina minimale pour esquiver

	StaminaRegenRate = 10 // Frames entre chaque point de stamina régénérée
)

// Player représente le joueur
type Player struct {
	BaseEntity
	Stats             *PlayerStats
	Stamina           int
	MaxStamina        int
	StaminaRegenTimer int
	IsAttacking       bool
	AttackTimer       int
	IsDodging         bool
	DodgeTimer        int
	DodgeDirection    float64
	Facing            float64 // Direction en radians

	// Animation et sprites
	Direction      PlayerDirection
	AnimationTimer int
	AnimationFrame int
	IsMoving       bool

	// Sprites pour chaque direction (AnimationFrameCount frames par animation)
	SpritesDown      []*ebiten.Image
	SpritesUp        []*ebiten.Image
	SpritesLeft      []*ebiten.Image
	SpritesRight     []*ebiten.Image
	SpritesDownIdle  []*ebiten.Image
	SpritesUpIdle    []*ebiten.Image
	SpritesLeftIdle  []*ebiten.Image
	SpritesRightIdle []*ebiten.Image

	// Inventaire et équipement
	Equipment *Equipment
	Inventory *Inventory
}

// NewPlayer crée un nouveau joueur
func NewPlayer(x, y float64) *Player {
	stats := NewPlayerStats()

	p := &Player{
		BaseEntity: BaseEntity{
			X:         x,
			Y:         y,
			Width:     32,
			Height:    32,
			Health:    int(stats.MaxHealth),
			MaxHealth: int(stats.MaxHealth),
			Speed:     PlayerSpeed,
		},
		Stats:      stats,
		Stamina:    int(stats.MaxStamina),
		MaxStamina: int(stats.MaxStamina),
		Facing:     0,
		Direction:  DirectionDown,
	}

	p.loadSprites()

	p.Equipment = NewEquipment()
	p.Inventory = NewInventory()

	startSword := NewStartingSword()
	p.Equipment.Weapon = startSword

	return p
}

// loadSprites charge tous les sprites du joueur
func (p *Player) loadSprites() {
	p.SpritesDown = loadSpriteFrames("assets/player/down/down_%d.png", AnimationFrameCount)
	p.SpritesUp = loadSpriteFrames("assets/player/up/up_%d.png", AnimationFrameCount)
	p.SpritesLeft = loadSpriteFrames("assets/player/left/left_%d.png", AnimationFrameCount)
	p.SpritesRight = loadSpriteFrames("assets/player/right/right_%d.png", AnimationFrameCount)

	p.SpritesDownIdle = loadSingleSprite("assets/player/down_idle/idle_down.png")
	p.SpritesUpIdle = loadSingleSprite("assets/player/up_idle/idle_up.png")
	p.SpritesLeftIdle = loadSingleSprite("assets/player/left_idle/idle_left.png")
	p.SpritesRightIdle = loadSingleSprite("assets/player/right_idle/idle_right.png")
}

// UpdateDerivedStats recalcule les stats dérivées avec les bonus d'équipement
func (p *Player) UpdateDerivedStats() {
	p.Stats.RecalculateStats()

	bonusVie := p.Equipment.GetTotalBonusVie()
	p.MaxHealth = int((p.Stats.Vie + bonusVie) * 10)
	p.MaxStamina = int(p.Stats.MaxStamina)
}

// loadSpriteFrames charge une séquence de sprites numérotés.
// Loggue un avertissement pour chaque sprite manquant.
func loadSpriteFrames(pathPattern string, frameCount int) []*ebiten.Image {
	frames := make([]*ebiten.Image, 0, frameCount)
	for i := 0; i < frameCount; i++ {
		path := fmt.Sprintf(pathPattern, i)
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Printf("[WARN] Sprite introuvable: %s (%v)", path, err)
			continue
		}
		frames = append(frames, img)
	}
	if len(frames) == 0 {
		return nil
	}
	return frames
}

// loadSingleSprite charge un sprite unique et le retourne dans un tableau.
// Loggue un avertissement si le fichier est introuvable.
func loadSingleSprite(path string) []*ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Printf("[WARN] Sprite introuvable: %s (%v)", path, err)
		return nil
	}
	return []*ebiten.Image{img}
}

// ApplyStatsToPlayer met à jour les stats du joueur après un level up
func (p *Player) ApplyStatsToPlayer() {
	oldMaxHealth := p.MaxHealth
	oldMaxStamina := p.MaxStamina

	p.MaxHealth = int(p.Stats.MaxHealth)
	p.MaxStamina = int(p.Stats.MaxStamina)

	if p.MaxHealth > oldMaxHealth {
		p.Health += (p.MaxHealth - oldMaxHealth)
		if p.Health > p.MaxHealth {
			p.Health = p.MaxHealth
		}
	}

	if p.MaxStamina > oldMaxStamina {
		p.Stamina += (p.MaxStamina - oldMaxStamina)
		if p.Stamina > p.MaxStamina {
			p.Stamina = p.MaxStamina
		}
	}
}

// GetDamage retourne les dégâts du joueur basés sur la Force
func (p *Player) GetDamage() int {
	return int(p.Stats.Damage)
}

// GetDamageWithEquipment retourne les dégâts du joueur incluant les bonus d'équipement
func (p *Player) GetDamageWithEquipment() int {
	baseDamage := int(p.Stats.Damage)
	bonusForce := p.Equipment.GetTotalBonusForce()
	weaponDamage := p.Equipment.GetTotalDamage()
	return baseDamage + (bonusForce * 2) + weaponDamage
}

// canMoveTo vérifie si le joueur peut se déplacer vers la position donnée
func (p *Player) canMoveTo(newX, newY float64, tilemap *world.TileMap) bool {
	if tilemap == nil {
		return true
	}

	halfWidth := p.Width / 2
	halfHeight := p.Height / 2

	// Vérifier les 4 coins du joueur (IsWalkable gère les bornes via GetTile)
	corners := [4][2]float64{
		{newX - halfWidth, newY - halfHeight}, // haut gauche
		{newX + halfWidth, newY - halfHeight}, // haut droit
		{newX - halfWidth, newY + halfHeight}, // bas gauche
		{newX + halfWidth, newY + halfHeight}, // bas droit
	}
	for _, c := range corners {
		tileX := int(c[0]) / world.TileSize
		tileY := int(c[1]) / world.TileSize
		if !tilemap.IsWalkable(tileX, tileY) {
			return false
		}
	}

	return true
}

// updateAnimation avance l'animation selon que le joueur est en mouvement ou idle
func (p *Player) updateAnimation() {
	speed := AnimationSpeedIdle
	if p.IsMoving {
		speed = AnimationSpeedWalk
	}
	p.AnimationTimer++
	if p.AnimationTimer >= speed {
		p.AnimationTimer = 0
		p.AnimationFrame = (p.AnimationFrame + 1) % AnimationFrameCount
	}
}

// Update met à jour l'état du joueur
func (p *Player) Update(inputMgr *input.InputManager, tilemap *world.TileMap) {
	p.UpdateInvincibility()

	// Gérer l'esquive
	if p.IsDodging {
		p.DodgeTimer--
		if p.DodgeTimer <= 0 {
			p.IsDodging = false
		} else {
			newX := p.X + math.Cos(p.DodgeDirection)*p.Speed*3
			newY := p.Y + math.Sin(p.DodgeDirection)*p.Speed*3
			if p.canMoveTo(newX, newY, tilemap) {
				p.X = newX
				p.Y = newY
			}
			return
		}
	}

	// Gérer l'attaque
	if p.IsAttacking {
		p.AttackTimer--
		if p.AttackTimer <= 0 {
			p.IsAttacking = false
		}
		return
	}

	// Régénération de la stamina
	p.StaminaRegenTimer++
	if p.StaminaRegenTimer >= StaminaRegenRate {
		p.StaminaRegenTimer = 0
		if p.Stamina < p.MaxStamina {
			p.Stamina++
		}
	}

	// Mouvement
	dx, dy := 0.0, 0.0

	if inputMgr.IsKeyPressed(ebiten.KeyW) { // Z sur AZERTY
		dy = -1
	}
	if inputMgr.IsKeyPressed(ebiten.KeyS) {
		dy = 1
	}
	if inputMgr.IsKeyPressed(ebiten.KeyA) { // Q sur AZERTY
		dx = -1
	}
	if inputMgr.IsKeyPressed(ebiten.KeyD) {
		dx = 1
	}

	if dx != 0 || dy != 0 {
		// Normaliser le vecteur de déplacement (déplacement diagonal identique)
		length := math.Sqrt(dx*dx + dy*dy)
		dx /= length
		dy /= length

		p.Facing = math.Atan2(dy, dx)

		// Direction du sprite selon l'axe dominant
		if math.Abs(dx) > math.Abs(dy) {
			if dx > 0 {
				p.Direction = DirectionRight
			} else {
				p.Direction = DirectionLeft
			}
		} else {
			if dy > 0 {
				p.Direction = DirectionDown
			} else {
				p.Direction = DirectionUp
			}
		}

		p.IsMoving = true

		// Sprint avec Shift
		speed := p.Speed
		if inputMgr.IsKeyPressed(ebiten.KeyShift) && p.Stamina > 0 {
			speed *= SprintMultiplier
			p.Stamina--
		}

		newX := p.X + dx*speed
		newY := p.Y + dy*speed
		if p.canMoveTo(newX, newY, tilemap) {
			p.X = newX
			p.Y = newY
		}
	} else {
		p.IsMoving = false
	}

	// Mettre à jour l'animation (helper mutualisé)
	p.updateAnimation()

	// Attaque avec Espace
	if inputMgr.IsKeyJustPressed(ebiten.KeySpace) && !p.IsAttacking && p.Stamina >= StaminaMinAttack {
		p.IsAttacking = true
		p.AttackTimer = AttackDuration
		p.Stamina -= StaminaCostAttack
	}

	// Esquive avec Shift + direction
	if inputMgr.IsKeyJustPressed(ebiten.KeyShift) && (dx != 0 || dy != 0) && p.Stamina >= StaminaMinDodge && !p.IsDodging {
		p.IsDodging = true
		p.DodgeTimer = DodgeDuration
		p.DodgeDirection = math.Atan2(dy, dx)
		p.Stamina -= StaminaCostDodge
		p.IsInvincible = true
		p.InvincibleTimer = DodgeDuration
	}
}

// Draw dessine le joueur
func (p *Player) Draw(screen *ebiten.Image, camera *world.Camera) {
	screenX, screenY := camera.WorldToScreen(p.X, p.Y)

	// Sélectionner les sprites appropriés selon direction et état
	var sprites []*ebiten.Image
	if p.IsMoving {
		switch p.Direction {
		case DirectionDown:
			sprites = p.SpritesDown
		case DirectionUp:
			sprites = p.SpritesUp
		case DirectionLeft:
			sprites = p.SpritesLeft
		case DirectionRight:
			sprites = p.SpritesRight
		}
	} else {
		switch p.Direction {
		case DirectionDown:
			sprites = p.SpritesDownIdle
		case DirectionUp:
			sprites = p.SpritesUpIdle
		case DirectionLeft:
			sprites = p.SpritesLeftIdle
		case DirectionRight:
			sprites = p.SpritesRightIdle
		}
	}

	// Épée dans le dos (toujours derrière le joueur, sauf pendant l'attaque)
	if p.Equipment.Weapon != nil && !p.IsAttacking {
		p.drawWeaponOnBack(screen, screenX, screenY)
	}

	if sprites != nil && len(sprites) > 0 {
		frame := p.AnimationFrame % len(sprites)
		sprite := sprites[frame]

		op := &ebiten.DrawImageOptions{}
		bounds := sprite.Bounds()
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		op.GeoM.Translate(screenX, screenY)

		if p.IsInvincible && p.InvincibleTimer%10 < 5 {
			op.ColorScale.Scale(1, 1, 1, 0.5)
		}
		if p.IsDodging {
			op.ColorScale.Scale(0.6, 0.6, 1.0, 0.8)
		}

		screen.DrawImage(sprite, op)
	} else {
		// Fallback : cercle coloré si les sprites ne sont pas chargés
		playerColor := color.RGBA{100, 200, 255, 255}
		if p.IsInvincible && p.InvincibleTimer%10 < 5 {
			playerColor = color.RGBA{255, 255, 255, 128}
		}
		if p.IsDodging {
			playerColor = color.RGBA{150, 150, 255, 200}
		}
		vector.DrawFilledCircle(screen, float32(screenX), float32(screenY), float32(p.Width/2), playerColor, true)

		endX := screenX + math.Cos(p.Facing)*p.Width*0.7
		endY := screenY + math.Sin(p.Facing)*p.Height*0.7
		vector.StrokeLine(screen, float32(screenX), float32(screenY), float32(endX), float32(endY), 3, color.White, true)
	}

	// Pendant l'attaque, dessiner l'arme en main (devant le joueur)
	if p.Equipment.Weapon != nil && p.IsAttacking {
		p.drawWeapon(screen, screenX, screenY)
	}

	// Arc d'attaque visuel
	if p.IsAttacking {
		attackColor := color.RGBA{255, 100, 100, 180}
		attackRange := p.Width * 2
		for i := -30.0; i <= 30.0; i += 5 {
			angle := p.Facing + i*math.Pi/180
			attackX := screenX + math.Cos(angle)*attackRange
			attackY := screenY + math.Sin(angle)*attackRange
			vector.DrawFilledCircle(screen, float32(attackX), float32(attackY), 4, attackColor, true)
		}
	}
}

// drawWeaponOnBack dessine l'arme accrochée dans le dos du joueur (derrière lui, en diagonale)
func (p *Player) drawWeaponOnBack(screen *ebiten.Image, screenX, screenY float64) {
	if p.Equipment.Weapon == nil {
		return
	}

	weaponSprite := p.Equipment.Weapon.IconSprite
	if weaponSprite == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	bounds := weaponSprite.Bounds()
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	// Réduire la taille à ~22px pour qu'elle reste proportionnelle au personnage
	const backWeaponSize = 22.0
	scale := backWeaponSize / math.Max(w, h)

	// Centrer le sprite sur son propre axe, puis faire pivoter de -45° (épée en diagonale)
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Rotate(-math.Pi / 4)
	// Décaler légèrement vers le haut-gauche pour simuler le dos
	op.GeoM.Translate(screenX-5, screenY-8)

	if p.IsInvincible && p.InvincibleTimer%10 < 5 {
		op.ColorScale.Scale(1, 1, 1, 0.5)
	}
	if p.IsDodging {
		op.ColorScale.Scale(0.6, 0.6, 1.0, 0.8)
	}

	screen.DrawImage(weaponSprite, op)
}

// drawWeapon dessine l'arme équipée du joueur
func (p *Player) drawWeapon(screen *ebiten.Image, screenX, screenY float64) {
	if p.Equipment.Weapon == nil {
		return
	}

	var weaponSprite *ebiten.Image
	var offsetX, offsetY float64

	switch p.Direction {
	case DirectionDown:
		weaponSprite = p.Equipment.Weapon.SpriteDown
		offsetX = 12
		offsetY = 12
	case DirectionUp:
		weaponSprite = p.Equipment.Weapon.SpriteUp
		offsetX = 12
		offsetY = -12
	case DirectionLeft:
		weaponSprite = p.Equipment.Weapon.SpriteLeft
		offsetX = -12
		offsetY = 2
	case DirectionRight:
		weaponSprite = p.Equipment.Weapon.SpriteRight
		offsetX = 12
		offsetY = 2
	}

	if weaponSprite != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := weaponSprite.Bounds()
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		op.GeoM.Translate(screenX+offsetX, screenY+offsetY)

		if p.IsInvincible && p.InvincibleTimer%10 < 5 {
			op.ColorScale.Scale(1, 1, 1, 0.5)
		}
		if p.IsDodging {
			op.ColorScale.Scale(0.6, 0.6, 1.0, 0.8)
		}

		screen.DrawImage(weaponSprite, op)
	}
}

// GetAttackHitbox retourne la hitbox de l'attaque si le joueur attaque
func (p *Player) GetAttackHitbox() (float64, float64, float64, bool) {
	if !p.IsAttacking {
		return 0, 0, 0, false
	}

	attackRange := p.Width * 2
	attackX := p.X + math.Cos(p.Facing)*attackRange
	attackY := p.Y + math.Sin(p.Facing)*attackRange

	return attackX, attackY, attackRange, true
}
