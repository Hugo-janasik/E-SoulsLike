package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ItemType définit les types d'items
type ItemType int

const (
	ItemTypeWeapon ItemType = iota
	ItemTypeArmor
	ItemTypeAccessory
	ItemTypeConsumable
)

// EquipmentSlot définit les slots d'équipement
type EquipmentSlot int

const (
	SlotWeapon EquipmentSlot = iota
	SlotChest  // Plastron
	SlotRing   // Anneau
	SlotAmulet // Amulette
	SlotBelt   // Ceinture
	SlotHelmet // Chapeau
	SlotBoots  // Bottes
)

// Item représente un item générique
type Item struct {
	ID          string
	Name        string
	Description string
	Type        ItemType
	IconSprite  *ebiten.Image // Sprite pour l'inventaire
	Stackable   bool
	MaxStack    int
}

// Weapon représente une arme
type Weapon struct {
	Item
	Damage      int
	AttackRange float64
	AttackSpeed float64
	Slot        EquipmentSlot // Toujours SlotWeapon

	// Sprites pour chaque direction
	SpriteDown  *ebiten.Image
	SpriteUp    *ebiten.Image
	SpriteLeft  *ebiten.Image
	SpriteRight *ebiten.Image
}

// Armor représente une pièce d'armure
type Armor struct {
	Item
	Defense   int
	BonusVie  int // Bonus de points de Vie (chaque point = +10 HP max)
	BonusForce int // Bonus de Force (augmente les dégâts)
	Slot      EquipmentSlot // Plastron, Chapeau, Bottes
}

// Accessory représente un accessoire
type Accessory struct {
	Item
	BonusVie   int // Bonus de points de Vie
	BonusForce int // Bonus de Force
	Slot       EquipmentSlot // Anneau, Amulette, Ceinture
}

// loadWeaponSprite charge un sprite d'arme
func loadWeaponSprite(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil
	}
	return img
}
