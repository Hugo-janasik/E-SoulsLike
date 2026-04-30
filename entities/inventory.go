package entities

// InventorySlot représente un slot d'inventaire
type InventorySlot struct {
	Item     interface{} // Peut être *Weapon, *Armor, *Accessory, ou *Item
	Quantity int
}

// Inventory représente l'inventaire du joueur
type Inventory struct {
	Slots    [12]*InventorySlot // 4x3 grille
	MaxSlots int
}

// NewInventory crée un nouvel inventaire
func NewInventory() *Inventory {
	return &Inventory{
		Slots:    [12]*InventorySlot{},
		MaxSlots: 12,
	}
}

// AddItem ajoute un item à l'inventaire
// item peut être *Weapon, *Armor, *Accessory, ou *Item
func (inv *Inventory) AddItem(item interface{}, quantity int) bool {
	if item == nil || quantity <= 0 {
		return false
	}

	// Extraire l'Item de base pour vérifier stackable
	baseItem := inv.getBaseItem(item)
	if baseItem == nil {
		return false
	}

	// Si l'item est stackable, chercher un slot existant
	if baseItem.Stackable {
		for i, slot := range inv.Slots {
			if slot != nil && slot.Item != nil {
				slotBaseItem := inv.getBaseItem(slot.Item)
				if slotBaseItem != nil && slotBaseItem.ID == baseItem.ID {
					// Vérifier si on peut ajouter plus dans ce slot
					if slot.Quantity+quantity <= baseItem.MaxStack {
						inv.Slots[i].Quantity += quantity
						return true
					} else {
						// Remplir ce slot et continuer avec le reste
						remaining := (slot.Quantity + quantity) - baseItem.MaxStack
						inv.Slots[i].Quantity = baseItem.MaxStack
						quantity = remaining
						// Continuer pour trouver/créer un autre slot
					}
				}
			}
		}
	}

	// Chercher un slot vide
	for i, slot := range inv.Slots {
		if slot == nil {
			inv.Slots[i] = &InventorySlot{
				Item:     item,
				Quantity: quantity,
			}
			return true
		}
	}

	// Inventaire plein
	return false
}

// getBaseItem extrait l'Item de base d'un item (gère *Weapon, *Armor, *Accessory, *Item)
func (inv *Inventory) getBaseItem(item interface{}) *Item {
	switch i := item.(type) {
	case *Weapon:
		return &i.Item
	case *Armor:
		return &i.Item
	case *Accessory:
		return &i.Item
	case *Item:
		return i
	default:
		return nil
	}
}

// RemoveItem retire un item du slot donné
// Retourne l'item (peut être *Weapon, *Armor, *Accessory, ou *Item)
func (inv *Inventory) RemoveItem(slotIndex int) interface{} {
	if slotIndex < 0 || slotIndex >= inv.MaxSlots {
		return nil
	}

	slot := inv.Slots[slotIndex]
	if slot == nil {
		return nil
	}

	item := slot.Item
	slot.Quantity--

	// Si la quantité tombe à 0, vider le slot
	if slot.Quantity <= 0 {
		inv.Slots[slotIndex] = nil
	}

	return item
}

// GetItem retourne l'item et la quantité dans un slot
// L'item peut être *Weapon, *Armor, *Accessory, ou *Item
func (inv *Inventory) GetItem(slotIndex int) (interface{}, int) {
	if slotIndex < 0 || slotIndex >= inv.MaxSlots {
		return nil, 0
	}

	slot := inv.Slots[slotIndex]
	if slot == nil {
		return nil, 0
	}

	return slot.Item, slot.Quantity
}

// HasItemByID vérifie si un item avec l'ID donné est présent dans l'inventaire
func (inv *Inventory) HasItemByID(id string) bool {
	for _, slot := range inv.Slots {
		if slot != nil && slot.Item != nil {
			if base := inv.getBaseItem(slot.Item); base != nil && base.ID == id {
				return true
			}
		}
	}
	return false
}

// HasSpace vérifie si l'inventaire a au moins un slot libre
func (inv *Inventory) HasSpace() bool {
	for _, slot := range inv.Slots {
		if slot == nil {
			return true
		}
	}
	return false
}
