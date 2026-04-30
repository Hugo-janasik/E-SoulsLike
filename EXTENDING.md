# Guide d'extension du jeu

Ce document explique comment ajouter de nouvelles fonctionnalités au jeu E-SoulsLike.

## Exemple 1: Ajouter un système d'items

### Étape 1: Créer le package items

Créer `items/item.go`:

```go
package items

import (
	"e-soulslike/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// ItemType définit les types d'objets
type ItemType int

const (
	ItemHealthPotion ItemType = iota
	ItemStaminaPotion
	ItemKey
	ItemCoin
)

// Item représente un objet au sol
type Item struct {
	X, Y     float64
	Type     ItemType
	IsActive bool
	Radius   float64
}

// NewItem crée un nouvel objet
func NewItem(x, y float64, itemType ItemType) *Item {
	return &Item{
		X:        x,
		Y:        y,
		Type:     itemType,
		IsActive: true,
		Radius:   16,
	}
}

// Draw dessine l'objet
func (i *Item) Draw(screen *ebiten.Image, camera *world.Camera) {
	if !i.IsActive {
		return
	}

	screenX, screenY := camera.WorldToScreen(i.X, i.Y)

	var itemColor color.RGBA
	switch i.Type {
	case ItemHealthPotion:
		itemColor = color.RGBA{255, 50, 50, 255} // Rouge
	case ItemStaminaPotion:
		itemColor = color.RGBA{50, 255, 50, 255} // Vert
	case ItemKey:
		itemColor = color.RGBA{255, 215, 0, 255} // Or
	case ItemCoin:
		itemColor = color.RGBA{255, 255, 0, 255} // Jaune
	}

	vector.DrawFilledCircle(screen, float32(screenX), float32(screenY),
		float32(i.Radius), itemColor, true)
}

// CanPickup vérifie si un joueur peut ramasser l'objet
func (i *Item) CanPickup(playerX, playerY, pickupRadius float64) bool {
	if !i.IsActive {
		return false
	}

	dx := playerX - i.X
	dy := playerY - i.Y
	distance := dx*dx + dy*dy
	maxDistance := (pickupRadius + i.Radius) * (pickupRadius + i.Radius)

	return distance < maxDistance
}

// Pickup ramasse l'objet
func (i *Item) Pickup() {
	i.IsActive = false
}
```

### Étape 2: Créer le système d'inventaire

Créer `items/inventory.go`:

```go
package items

// Inventory représente l'inventaire du joueur
type Inventory struct {
	Items     map[ItemType]int
	MaxSlots  int
	Gold      int
}

// NewInventory crée un nouvel inventaire
func NewInventory(maxSlots int) *Inventory {
	return &Inventory{
		Items:    make(map[ItemType]int),
		MaxSlots: maxSlots,
		Gold:     0,
	}
}

// AddItem ajoute un objet à l'inventaire
func (inv *Inventory) AddItem(itemType ItemType, quantity int) bool {
	inv.Items[itemType] += quantity
	return true
}

// UseItem utilise un objet de l'inventaire
func (inv *Inventory) UseItem(itemType ItemType) bool {
	if inv.Items[itemType] > 0 {
		inv.Items[itemType]--
		return true
	}
	return false
}

// GetItemCount retourne le nombre d'objets d'un certain type
func (inv *Inventory) GetItemCount(itemType ItemType) int {
	return inv.Items[itemType]
}
```

### Étape 3: Intégrer dans le jeu

Modifier `game/game.go`:

```go
import (
	// ... autres imports
	"e-soulslike/items"
)

type Game struct {
	// ... champs existants
	items     []*items.Item
	inventory *items.Inventory
}

func NewGame(width, height int) *Game {
	g := &Game{
		// ... initialisation existante
		inventory: items.NewInventory(20),
	}

	// Ajouter quelques objets dans le monde
	g.items = []*items.Item{
		items.NewItem(500, 300, items.ItemHealthPotion),
		items.NewItem(700, 400, items.ItemStaminaPotion),
		items.NewItem(900, 500, items.ItemCoin),
	}

	return g
}

func (g *Game) Update() error {
	// ... code existant

	// Vérifier le ramassage des objets
	for _, item := range g.items {
		if item.CanPickup(g.player.X, g.player.Y, 40) {
			// Ramasser l'objet
			item.Pickup()
			g.inventory.AddItem(item.Type, 1)

			// Appliquer l'effet immédiatement (optionnel)
			g.applyItemEffect(item.Type)
		}
	}

	return nil
}

func (g *Game) applyItemEffect(itemType items.ItemType) {
	switch itemType {
	case items.ItemHealthPotion:
		g.player.Health += 30
		if g.player.Health > g.player.MaxHealth {
			g.player.Health = g.player.MaxHealth
		}
	case items.ItemStaminaPotion:
		g.player.Stamina += 50
		if g.player.Stamina > g.player.MaxStamina {
			g.player.Stamina = g.player.MaxStamina
		}
	case items.ItemCoin:
		g.inventory.Gold += 10
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ... code de rendu existant

	// Dessiner les objets
	for _, item := range g.items {
		item.Draw(screen, g.camera)
	}

	// Afficher l'inventaire (simple)
	debugInfo := fmt.Sprintf(
		"HP Potions: %d | Stamina Potions: %d | Gold: %d",
		g.inventory.GetItemCount(items.ItemHealthPotion),
		g.inventory.GetItemCount(items.ItemStaminaPotion),
		g.inventory.Gold,
	)
	ebitenutil.DebugPrintAt(screen, debugInfo, 10, 80)
}
```

## Exemple 2: Ajouter des projectiles

### Créer `combat/projectile.go`:

```go
package combat

import (
	"e-soulslike/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

type Projectile struct {
	X, Y       float64
	VelX, VelY float64
	Damage     int
	Radius     float64
	Lifetime   int
	IsActive   bool
	FromPlayer bool // true si tiré par le joueur
}

func NewProjectile(x, y, angle, speed float64, damage int, fromPlayer bool) *Projectile {
	return &Projectile{
		X:          x,
		Y:          y,
		VelX:       math.Cos(angle) * speed,
		VelY:       math.Sin(angle) * speed,
		Damage:     damage,
		Radius:     8,
		Lifetime:   180, // 3 secondes à 60 FPS
		IsActive:   true,
		FromPlayer: fromPlayer,
	}
}

func (p *Projectile) Update() {
	p.X += p.VelX
	p.Y += p.VelY
	p.Lifetime--

	if p.Lifetime <= 0 {
		p.IsActive = false
	}
}

func (p *Projectile) Draw(screen *ebiten.Image, camera *world.Camera) {
	if !p.IsActive {
		return
	}

	screenX, screenY := camera.WorldToScreen(p.X, p.Y)

	projectileColor := color.RGBA{255, 200, 50, 255}
	if !p.FromPlayer {
		projectileColor = color.RGBA{200, 50, 255, 255}
	}

	vector.DrawFilledCircle(screen, float32(screenX), float32(screenY),
		float32(p.Radius), projectileColor, true)
}
```

## Exemple 3: Ajouter un système de particules

### Créer `effects/particles.go`:

```go
package effects

import (
	"e-soulslike/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
)

type Particle struct {
	X, Y       float64
	VelX, VelY float64
	Life       int
	MaxLife    int
	Size       float64
	Color      color.RGBA
}

type ParticleSystem struct {
	Particles []*Particle
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		Particles: make([]*Particle, 0),
	}
}

// SpawnBloodEffect crée un effet de sang
func (ps *ParticleSystem) SpawnBloodEffect(x, y float64) {
	for i := 0; i < 10; i++ {
		particle := &Particle{
			X:       x,
			Y:       y,
			VelX:    (rand.Float64() - 0.5) * 6,
			VelY:    (rand.Float64() - 0.5) * 6,
			Life:    30 + rand.Intn(20),
			MaxLife: 30 + rand.Intn(20),
			Size:    2 + rand.Float64()*3,
			Color:   color.RGBA{200 + uint8(rand.Intn(55)), 0, 0, 255},
		}
		ps.Particles = append(ps.Particles, particle)
	}
}

func (ps *ParticleSystem) Update() {
	for i := len(ps.Particles) - 1; i >= 0; i-- {
		p := ps.Particles[i]
		p.X += p.VelX
		p.Y += p.VelY
		p.VelY += 0.2 // Gravité
		p.Life--

		if p.Life <= 0 {
			// Retirer la particule
			ps.Particles = append(ps.Particles[:i], ps.Particles[i+1:]...)
		}
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image, camera *world.Camera) {
	for _, p := range ps.Particles {
		screenX, screenY := camera.WorldToScreen(p.X, p.Y)

		// Fade out
		alpha := float64(p.Life) / float64(p.MaxLife)
		col := p.Color
		col.A = uint8(float64(col.A) * alpha)

		vector.DrawFilledCircle(screen, float32(screenX), float32(screenY),
			float32(p.Size), col, true)
	}
}
```

## Exemple 4: Ajouter un système de quêtes

### Créer `quests/quest.go`:

```go
package quests

type QuestType int

const (
	QuestKillEnemies QuestType = iota
	QuestCollectItems
	QuestReachLocation
)

type Quest struct {
	ID          int
	Name        string
	Description string
	Type        QuestType
	Target      int
	Current     int
	Completed   bool
	Reward      int
}

type QuestManager struct {
	ActiveQuests    []*Quest
	CompletedQuests []*Quest
}

func NewQuestManager() *QuestManager {
	return &QuestManager{
		ActiveQuests:    make([]*Quest, 0),
		CompletedQuests: make([]*Quest, 0),
	}
}

func (qm *QuestManager) AddQuest(quest *Quest) {
	qm.ActiveQuests = append(qm.ActiveQuests, quest)
}

func (qm *QuestManager) UpdateProgress(questID int, amount int) {
	for i, quest := range qm.ActiveQuests {
		if quest.ID == questID {
			quest.Current += amount
			if quest.Current >= quest.Target {
				quest.Completed = true
				// Déplacer vers les quêtes complétées
				qm.CompletedQuests = append(qm.CompletedQuests, quest)
				qm.ActiveQuests = append(qm.ActiveQuests[:i], qm.ActiveQuests[i+1:]...)
			}
			break
		}
	}
}
```

## Bonnes pratiques

### 1. Séparation des responsabilités

Chaque système doit avoir sa propre responsabilité:
- `entities/` : Tout ce qui est vivant (joueur, ennemis, PNJ)
- `items/` : Objets, inventaire, équipement
- `combat/` : Système de combat, projectiles, dégâts
- `world/` : Monde, carte, caméra
- `effects/` : Effets visuels, particules
- `audio/` : Musique et sons

### 2. Interface pour l'extensibilité

Utilisez des interfaces pour permettre différentes implémentations:

```go
type Damageable interface {
	TakeDamage(amount int)
	GetHealth() int
	IsAlive() bool
}

type Drawable interface {
	Draw(screen *ebiten.Image, camera *world.Camera)
}

type Updatable interface {
	Update()
}
```

### 3. Event System (optionnel)

Pour découpler les systèmes:

```go
type EventType int

const (
	EventEnemyKilled EventType = iota
	EventItemPickup
	EventPlayerDamaged
)

type Event struct {
	Type EventType
	Data interface{}
}

type EventBus struct {
	listeners map[EventType][]func(Event)
}

func (eb *EventBus) Subscribe(eventType EventType, callback func(Event)) {
	eb.listeners[eventType] = append(eb.listeners[eventType], callback)
}

func (eb *EventBus) Publish(event Event) {
	for _, callback := range eb.listeners[event.Type] {
		callback(event)
	}
}
```

### 4. Configuration externe

Créer un fichier `config.go`:

```go
package config

type GameConfig struct {
	ScreenWidth     int
	ScreenHeight    int
	PlayerSpeed     float64
	PlayerMaxHealth int
	EnemySpawnRate  int
}

func LoadConfig() *GameConfig {
	return &GameConfig{
		ScreenWidth:     1280,
		ScreenHeight:    720,
		PlayerSpeed:     3.0,
		PlayerMaxHealth: 100,
		EnemySpawnRate:  300,
	}
}
```

### 5. Sauvegarde/Chargement

Créer `save/save.go`:

```go
package save

import (
	"encoding/json"
	"os"
)

type SaveData struct {
	PlayerX       float64
	PlayerY       float64
	PlayerHealth  int
	Inventory     map[string]int
	CompletedQuests []int
}

func SaveGame(data *SaveData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func LoadGame(filename string) (*SaveData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data SaveData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return &data, err
}
```

## Checklist pour ajouter une fonctionnalité

- [ ] Créer le package approprié
- [ ] Définir les structures de données
- [ ] Implémenter les méthodes Update() et Draw()
- [ ] Intégrer dans game.go
- [ ] Ajouter des tests
- [ ] Documenter dans ARCHITECTURE.md
- [ ] Mettre à jour README.md
- [ ] Tester en jeu

## Ressources utiles

- [Ebiten Examples](https://ebitengine.org/en/examples/)
- [Game Programming Patterns](https://gameprogrammingpatterns.com/)
- [Red Blob Games](https://www.redblobgames.com/) - Algorithmes de jeux
