package entities

import "github.com/hajimehoshi/ebiten/v2/ebitenutil"

// NewStartingSword crée l'épée de départ du joueur
func NewStartingSword() *Weapon {
	iconSprite, _, _ := ebitenutil.NewImageFromFile("assets/weapons/sword/full.png")

	return &Weapon{
		Item: Item{
			ID:          "sword_basic",
			Name:        "Épée de base",
			Description: "Une simple épée en fer",
			Type:        ItemTypeWeapon,
			IconSprite:  iconSprite,
			Stackable:   false,
			MaxStack:    1,
		},
		Damage:      5,
		AttackRange: 64.0,
		AttackSpeed: 1.0,
		Slot:        SlotWeapon,
		SpriteDown:  loadWeaponSprite("assets/weapons/sword/down.png"),
		SpriteUp:    loadWeaponSprite("assets/weapons/sword/up.png"),
		SpriteLeft:  loadWeaponSprite("assets/weapons/sword/left.png"),
		SpriteRight: loadWeaponSprite("assets/weapons/sword/right.png"),
	}
}

// ========== SHOP ITEMS ==========

// NewBeginnerHelmet crée le casque du débutant
func NewBeginnerHelmet() *Armor {
	iconSprite, _, _ := ebitenutil.NewImageFromFile("assets/object/casque.png")

	return &Armor{
		Item: Item{
			ID:          "helmet_beginner",
			Name:        "Casque du débutant",
			Description: "Un casque simple mais solide",
			Type:        ItemTypeArmor,
			IconSprite:  iconSprite,
			Stackable:   false,
			MaxStack:    1,
		},
		Defense:   2,
		BonusVie:  10,
		BonusForce: 0,
		Slot:      SlotHelmet,
	}
}

// NewBeginnerChestplate crée le plastron du débutant
func NewBeginnerChestplate() *Armor {
	iconSprite, _, _ := ebitenutil.NewImageFromFile("assets/object/plastron.png")

	return &Armor{
		Item: Item{
			ID:          "chest_beginner",
			Name:        "Plastron du débutant",
			Description: "Une protection robuste pour le torse",
			Type:        ItemTypeArmor,
			IconSprite:  iconSprite,
			Stackable:   false,
			MaxStack:    1,
		},
		Defense:    5,
		BonusVie:   20,
		BonusForce: 0,
		Slot:       SlotChest,
	}
}

// NewBeginnerBoots crée les bottes du débutant
func NewBeginnerBoots() *Armor {
	iconSprite, _, _ := ebitenutil.NewImageFromFile("assets/object/botte.png")

	return &Armor{
		Item: Item{
			ID:          "boots_beginner",
			Name:        "Bottes du débutant",
			Description: "Des bottes confortables et résistantes",
			Type:        ItemTypeArmor,
			IconSprite:  iconSprite,
			Stackable:   false,
			MaxStack:    1,
		},
		Defense:   1,
		BonusVie:  5,
		BonusForce: 0,
		Slot:      SlotBoots,
	}
}

// NewBeginnerAmulet crée l'amulette du débutant
func NewBeginnerAmulet() *Accessory {
	iconSprite, _, _ := ebitenutil.NewImageFromFile("assets/weapons/sword/full.png") // TODO: icône dédiée

	return &Accessory{
		Item: Item{
			ID:          "amulet_beginner",
			Name:        "Amulette du débutant",
			Description: "Une amulette qui renforce la force",
			Type:        ItemTypeAccessory,
			IconSprite:  iconSprite,
			Stackable:   false,
			MaxStack:    1,
		},
		BonusVie:   0,
		BonusForce: 3,
		Slot:       SlotAmulet,
	}
}

// NewBeginnerRing crée l'anneau du débutant
func NewBeginnerRing() *Accessory {
	iconSprite, _, _ := ebitenutil.NewImageFromFile("assets/object/anneau.png")

	return &Accessory{
		Item: Item{
			ID:          "ring_beginner",
			Name:        "Anneau du débutant",
			Description: "Un anneau qui augmente la force",
			Type:        ItemTypeAccessory,
			IconSprite:  iconSprite,
			Stackable:   false,
			MaxStack:    1,
		},
		BonusVie:   0,
		BonusForce: 2,
		Slot:       SlotRing,
	}
}

// NewBeginnerLeggings crée les jambières du débutant (via slot Belt)
func NewBeginnerLeggings() *Accessory {
	iconSprite, _, _ := ebitenutil.NewImageFromFile("assets/object/jambière.png")

	return &Accessory{
		Item: Item{
			ID:          "leggings_beginner",
			Name:        "Jambe du débutant",
			Description: "Protection pour les jambes",
			Type:        ItemTypeAccessory,
			IconSprite:  iconSprite,
			Stackable:   false,
			MaxStack:    1,
		},
		BonusVie:   15,
		BonusForce: 0,
		Slot:       SlotBelt,
	}
}
