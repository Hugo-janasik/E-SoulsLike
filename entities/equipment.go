package entities

// Equipment représente l'équipement équipé du joueur
type Equipment struct {
	Weapon *Weapon
	Chest  *Armor
	Ring   *Accessory
	Amulet *Accessory
	Belt   *Accessory
	Helmet *Armor
	Boots  *Armor
}

// NewEquipment crée un nouvel équipement vide
func NewEquipment() *Equipment {
	return &Equipment{}
}

// Equip équipe un item dans le bon slot
// item peut être *Weapon, *Armor, ou *Accessory
// Retourne l'item précédemment équipé (nil si aucun)
func (e *Equipment) Equip(item interface{}) interface{} {
	if item == nil {
		return nil
	}

	// Type assertion directe sur l'interface{}
	switch i := item.(type) {
	case *Weapon:
		oldWeapon := e.Weapon
		e.Weapon = i
		// Éviter de retourner un "typed nil" qui n'est pas égal à nil interface
		if oldWeapon == nil {
			return nil
		}
		return oldWeapon

	case *Armor:
		switch i.Slot {
		case SlotChest:
			oldArmor := e.Chest
			e.Chest = i
			if oldArmor == nil {
				return nil
			}
			return oldArmor
		case SlotHelmet:
			oldArmor := e.Helmet
			e.Helmet = i
			if oldArmor == nil {
				return nil
			}
			return oldArmor
		case SlotBoots:
			oldArmor := e.Boots
			e.Boots = i
			if oldArmor == nil {
				return nil
			}
			return oldArmor
		}

	case *Accessory:
		switch i.Slot {
		case SlotRing:
			oldAccessory := e.Ring
			e.Ring = i
			if oldAccessory == nil {
				return nil
			}
			return oldAccessory
		case SlotAmulet:
			oldAccessory := e.Amulet
			e.Amulet = i
			if oldAccessory == nil {
				return nil
			}
			return oldAccessory
		case SlotBelt:
			oldAccessory := e.Belt
			e.Belt = i
			if oldAccessory == nil {
				return nil
			}
			return oldAccessory
		}
	}

	return nil
}

// Unequip déséquipe l'item du slot donné
// Retourne l'item déséquipé (*Weapon, *Armor, ou *Accessory) (nil si aucun)
func (e *Equipment) Unequip(slot EquipmentSlot) interface{} {
	switch slot {
	case SlotWeapon:
		if e.Weapon != nil {
			item := e.Weapon
			e.Weapon = nil
			return item
		}
	case SlotChest:
		if e.Chest != nil {
			item := e.Chest
			e.Chest = nil
			return item
		}
	case SlotHelmet:
		if e.Helmet != nil {
			item := e.Helmet
			e.Helmet = nil
			return item
		}
	case SlotBoots:
		if e.Boots != nil {
			item := e.Boots
			e.Boots = nil
			return item
		}
	case SlotRing:
		if e.Ring != nil {
			item := e.Ring
			e.Ring = nil
			return item
		}
	case SlotAmulet:
		if e.Amulet != nil {
			item := e.Amulet
			e.Amulet = nil
			return item
		}
	case SlotBelt:
		if e.Belt != nil {
			item := e.Belt
			e.Belt = nil
			return item
		}
	}
	return nil
}

// HasItemByID vérifie si un item avec l'ID donné est actuellement équipé
func (e *Equipment) HasItemByID(id string) bool {
	if e.Weapon != nil && e.Weapon.Item.ID == id {
		return true
	}
	if e.Chest != nil && e.Chest.Item.ID == id {
		return true
	}
	if e.Helmet != nil && e.Helmet.Item.ID == id {
		return true
	}
	if e.Boots != nil && e.Boots.Item.ID == id {
		return true
	}
	if e.Ring != nil && e.Ring.Item.ID == id {
		return true
	}
	if e.Amulet != nil && e.Amulet.Item.ID == id {
		return true
	}
	if e.Belt != nil && e.Belt.Item.ID == id {
		return true
	}
	return false
}

// GetTotalDamage retourne les dégâts bonus de l'équipement
func (e *Equipment) GetTotalDamage() int {
	damage := 0
	if e.Weapon != nil {
		damage += e.Weapon.Damage
	}
	return damage
}

// GetTotalDefense retourne la défense totale de toutes les armures
func (e *Equipment) GetTotalDefense() int {
	defense := 0
	if e.Chest != nil {
		defense += e.Chest.Defense
	}
	if e.Helmet != nil {
		defense += e.Helmet.Defense
	}
	if e.Boots != nil {
		defense += e.Boots.Defense
	}
	return defense
}

// GetTotalBonusVie retourne le bonus total de Vie de l'équipement
func (e *Equipment) GetTotalBonusVie() int {
	bonus := 0
	if e.Chest != nil {
		bonus += e.Chest.BonusVie
	}
	if e.Helmet != nil {
		bonus += e.Helmet.BonusVie
	}
	if e.Boots != nil {
		bonus += e.Boots.BonusVie
	}
	if e.Ring != nil {
		bonus += e.Ring.BonusVie
	}
	if e.Amulet != nil {
		bonus += e.Amulet.BonusVie
	}
	if e.Belt != nil {
		bonus += e.Belt.BonusVie
	}
	return bonus
}

// GetTotalBonusForce retourne le bonus total de Force de l'équipement
func (e *Equipment) GetTotalBonusForce() int {
	bonus := 0
	if e.Chest != nil {
		bonus += e.Chest.BonusForce
	}
	if e.Helmet != nil {
		bonus += e.Helmet.BonusForce
	}
	if e.Boots != nil {
		bonus += e.Boots.BonusForce
	}
	if e.Ring != nil {
		bonus += e.Ring.BonusForce
	}
	if e.Amulet != nil {
		bonus += e.Amulet.BonusForce
	}
	if e.Belt != nil {
		bonus += e.Belt.BonusForce
	}
	return bonus
}
