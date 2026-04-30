package ui

import (
	"e-soulslike/entities"
	"e-soulslike/input"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// InventoryMenu représente le menu d'inventaire
type InventoryMenu struct {
	screenWidth, screenHeight int
	player                    *entities.Player

	// Navigation
	selectedSlotIndex  int  // 0-11 pour inventaire, 0-6 pour équipement
	isInEquipmentPanel bool // true si on navigue dans le panneau d'équipement

	// Layout
	menuX, menuY          float32
	menuWidth, menuHeight float32
	slotSize              float32

	OnClose func()
}

// NewInventoryMenu crée un nouveau menu d'inventaire
func NewInventoryMenu(width, height int, player *entities.Player, onClose func()) *InventoryMenu {
	centerX := float32(width) / 2
	centerY := float32(height) / 2
	menuWidth := float32(600)
	menuHeight := float32(520)

	return &InventoryMenu{
		screenWidth:        width,
		screenHeight:       height,
		player:             player,
		selectedSlotIndex:  0,
		isInEquipmentPanel: false,
		menuX:              centerX - menuWidth/2,
		menuY:              centerY - menuHeight/2,
		menuWidth:          menuWidth,
		menuHeight:         menuHeight,
		slotSize:           50,
		OnClose:            onClose,
	}
}

// Update met à jour le menu d'inventaire
func (im *InventoryMenu) Update(inputMgr *input.InputManager) {
	if im.player == nil {
		return
	}

	// Flèches directionnelles pour naviguer
	if inputMgr.IsKeyJustPressed(ebiten.KeyArrowUp) {
		im.moveSelection(-4) // Une ligne vers le haut
	}
	if inputMgr.IsKeyJustPressed(ebiten.KeyArrowDown) {
		im.moveSelection(4) // Une ligne vers le bas
	}
	if inputMgr.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		im.moveSelection(-1) // Gauche
	}
	if inputMgr.IsKeyJustPressed(ebiten.KeyArrowRight) {
		im.moveSelection(1) // Droite
	}

	// Entrée pour équiper/déséquiper
	if inputMgr.IsKeyJustPressed(ebiten.KeyEnter) {
		im.handleEquipAction()
	}

	// Seulement Échap pour fermer (I et B sont gérés dans game.go pour éviter le toggle immédiat)
	if inputMgr.IsKeyJustPressed(ebiten.KeyEscape) {
		if im.OnClose != nil {
			im.OnClose()
		}
	}
}

// moveSelection déplace la sélection
func (im *InventoryMenu) moveSelection(delta int) {
	if im.isInEquipmentPanel {
		// Navigation dans les 7 slots d'équipement
		im.selectedSlotIndex += delta
		if im.selectedSlotIndex < 0 {
			im.selectedSlotIndex = 6
		} else if im.selectedSlotIndex > 6 {
			im.selectedSlotIndex = 0
		}
	} else {
		// Navigation dans la grille 4x3 (12 slots)
		newIndex := im.selectedSlotIndex + delta

		// Gérer les bords horizontaux
		if delta == -1 && im.selectedSlotIndex%4 == 0 {
			newIndex = im.selectedSlotIndex + 3
		} else if delta == 1 && im.selectedSlotIndex%4 == 3 {
			newIndex = im.selectedSlotIndex - 3
		}

		// Gérer les bords verticaux
		if newIndex < 0 {
			newIndex = im.selectedSlotIndex + 8
		} else if newIndex >= 12 {
			newIndex = im.selectedSlotIndex - 8
		}

		im.selectedSlotIndex = newIndex
	}
}

// handleEquipAction gère l'action d'équiper/déséquiper
func (im *InventoryMenu) handleEquipAction() {
	if im.isInEquipmentPanel {
		// Déséquiper l'item sélectionné
		slot := im.getEquipmentSlotFromIndex(im.selectedSlotIndex)
		item := im.player.Equipment.Unequip(slot)
		if item != nil {
			// Ajouter à l'inventaire si possible
			if !im.player.Inventory.AddItem(item, 1) {
				// Si l'inventaire est plein, ré-équiper
				im.player.Equipment.Equip(item)
			} else {
				// Recalculer les stats après déséquipement
				im.player.UpdateDerivedStats()
			}
		}
	} else {
		// Équiper l'item sélectionné de l'inventaire
		invSlot := im.player.Inventory.Slots[im.selectedSlotIndex]
		if invSlot != nil && invSlot.Item != nil {
			previousItem := im.player.Equipment.Equip(invSlot.Item)
			if previousItem != nil {
				// Remettre l'ancien item dans l'inventaire
				im.player.Inventory.AddItem(previousItem, 1)
			}
			// Retirer de l'inventaire
			im.player.Inventory.RemoveItem(im.selectedSlotIndex)

			// Recalculer les stats après équipement
			im.player.UpdateDerivedStats()
		}
	}
}

// getEquipmentSlotFromIndex convertit un index (0-6) en EquipmentSlot
func (im *InventoryMenu) getEquipmentSlotFromIndex(index int) entities.EquipmentSlot {
	slots := []entities.EquipmentSlot{
		entities.SlotWeapon,
		entities.SlotChest,
		entities.SlotHelmet,
		entities.SlotRing,
		entities.SlotAmulet,
		entities.SlotBelt,
		entities.SlotBoots,
	}
	if index >= 0 && index < len(slots) {
		return slots[index]
	}
	return entities.SlotWeapon
}

// Draw dessine le menu d'inventaire
func (im *InventoryMenu) Draw(screen *ebiten.Image) {
	if im.player == nil {
		return
	}

	// Overlay semi-transparent
	overlay := color.RGBA{0, 0, 0, 180}
	vector.DrawFilledRect(screen, 0, 0, float32(im.screenWidth), float32(im.screenHeight), overlay, false)

	// Fond du menu
	vector.DrawFilledRect(screen, im.menuX, im.menuY, im.menuWidth, im.menuHeight, color.RGBA{40, 40, 50, 255}, false)
	vector.StrokeRect(screen, im.menuX, im.menuY, im.menuWidth, im.menuHeight, 2, color.RGBA{100, 150, 255, 255}, false)

	// Titres des panneaux
	ebitenutil.DebugPrintAt(screen, "=== INVENTAIRE ===", int(im.menuX+20), int(im.menuY+12))
	ebitenutil.DebugPrintAt(screen, "=== EQUIPEMENT ===", int(im.menuX+im.menuWidth-240), int(im.menuY+12))

	// Séparateur vertical
	vector.StrokeLine(screen, im.menuX+im.menuWidth-250, im.menuY+5, im.menuX+im.menuWidth-250, im.menuY+im.menuHeight-95, 1, color.RGBA{80, 80, 100, 255}, false)

	// Grille inventaire (4x3, côté gauche)
	im.drawInventoryGrid(screen, im.menuX+20, im.menuY+40)

	// Panneau équipement (liste style magasin, côté droit)
	im.drawEquipmentPanel(screen, im.menuX+im.menuWidth-245, im.menuY+35)

	// Détails de l'item sélectionné
	im.drawItemDetails(screen, im.menuX+10, im.menuY+im.menuHeight-95)

	// Contrôles
	controlsY := int(im.menuY + im.menuHeight - 20)
	ebitenutil.DebugPrintAt(screen, "Fleches: Nav | Entree: Equiper | I/Echap: Fermer", int(im.menuX+20), controlsY)
}

// drawInventoryGrid dessine la grille d'inventaire (4x3)
func (im *InventoryMenu) drawInventoryGrid(screen *ebiten.Image, startX, startY float32) {
	for row := 0; row < 3; row++ {
		for col := 0; col < 4; col++ {
			slotIndex := row*4 + col
			x := startX + float32(col)*(im.slotSize+5)
			y := startY + float32(row)*(im.slotSize+5)

			slotColor := color.RGBA{60, 60, 70, 255}
			if !im.isInEquipmentPanel && im.selectedSlotIndex == slotIndex {
				slotColor = color.RGBA{100, 150, 255, 255}
			}
			vector.DrawFilledRect(screen, x, y, im.slotSize, im.slotSize, slotColor, false)
			vector.StrokeRect(screen, x, y, im.slotSize, im.slotSize, 1, color.RGBA{100, 100, 110, 255}, false)

			slot := im.player.Inventory.Slots[slotIndex]
			if slot != nil && slot.Item != nil {
				baseItem := getBaseItem(slot.Item)
				if baseItem != nil {
					if baseItem.IconSprite != nil {
						op := &ebiten.DrawImageOptions{}
						bounds := baseItem.IconSprite.Bounds()
						scale := float64(im.slotSize-10) / float64(bounds.Dx())
						if float64(bounds.Dy())*scale > float64(im.slotSize-10) {
							scale = float64(im.slotSize-10) / float64(bounds.Dy())
						}
						op.GeoM.Scale(scale, scale)
						op.GeoM.Translate(float64(x+5), float64(y+5))
						screen.DrawImage(baseItem.IconSprite, op)
					} else {
						itemColor := getItemColor(baseItem.Type)
						vector.DrawFilledRect(screen, x+10, y+10, im.slotSize-20, im.slotSize-20, itemColor, false)
					}
					if baseItem.Stackable && slot.Quantity > 1 {
						ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", slot.Quantity), int(x+35), int(y+35))
					}
				}
			}
		}
	}
}

// drawEquipmentPanel dessine le panneau d'équipement en liste (style magasin)
func (im *InventoryMenu) drawEquipmentPanel(screen *ebiten.Image, startX, startY float32) {
	const rowH = float32(50)
	const iconSz = float32(40)
	const panelW = float32(220)

	equipmentSlots := []struct {
		name string
		slot entities.EquipmentSlot
	}{
		{"Arme", entities.SlotWeapon},
		{"Plastron", entities.SlotChest},
		{"Casque", entities.SlotHelmet},
		{"Anneau", entities.SlotRing},
		{"Amulette", entities.SlotAmulet},
		{"Ceinture", entities.SlotBelt},
		{"Bottes", entities.SlotBoots},
	}

	for i, equip := range equipmentSlots {
		y := startY + float32(i)*rowH

		bgColor := color.RGBA{60, 60, 70, 255}
		if im.isInEquipmentPanel && im.selectedSlotIndex == i {
			bgColor = color.RGBA{80, 100, 255, 255}
		}
		vector.DrawFilledRect(screen, startX, y, panelW, rowH-4, bgColor, false)
		vector.StrokeRect(screen, startX, y, panelW, rowH-4, 1, color.RGBA{100, 100, 110, 255}, false)

		var itemName string
		var itemIcon *ebiten.Image
		switch equip.slot {
		case entities.SlotWeapon:
			if im.player.Equipment.Weapon != nil {
				itemName = im.player.Equipment.Weapon.Name
				itemIcon = im.player.Equipment.Weapon.IconSprite
			}
		case entities.SlotChest:
			if im.player.Equipment.Chest != nil {
				itemName = im.player.Equipment.Chest.Name
				itemIcon = im.player.Equipment.Chest.IconSprite
			}
		case entities.SlotHelmet:
			if im.player.Equipment.Helmet != nil {
				itemName = im.player.Equipment.Helmet.Name
				itemIcon = im.player.Equipment.Helmet.IconSprite
			}
		case entities.SlotBoots:
			if im.player.Equipment.Boots != nil {
				itemName = im.player.Equipment.Boots.Name
				itemIcon = im.player.Equipment.Boots.IconSprite
			}
		case entities.SlotRing:
			if im.player.Equipment.Ring != nil {
				itemName = im.player.Equipment.Ring.Name
				itemIcon = im.player.Equipment.Ring.IconSprite
			}
		case entities.SlotAmulet:
			if im.player.Equipment.Amulet != nil {
				itemName = im.player.Equipment.Amulet.Name
				itemIcon = im.player.Equipment.Amulet.IconSprite
			}
		case entities.SlotBelt:
			if im.player.Equipment.Belt != nil {
				itemName = im.player.Equipment.Belt.Name
				itemIcon = im.player.Equipment.Belt.IconSprite
			}
		}

		// Icône
		iconX := startX + 5
		iconY := y + (rowH-4-iconSz)/2
		if itemIcon != nil {
			b := itemIcon.Bounds()
			scale := float64(iconSz) / float64(b.Dx())
			if s2 := float64(iconSz) / float64(b.Dy()); s2 < scale {
				scale = s2
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(iconX), float64(iconY))
			screen.DrawImage(itemIcon, op)
		} else {
			// Carré vide stylisé
			vector.DrawFilledRect(screen, iconX, iconY, iconSz, iconSz, color.RGBA{45, 45, 55, 255}, false)
			vector.StrokeRect(screen, iconX, iconY, iconSz, iconSz, 1, color.RGBA{80, 80, 95, 255}, false)
		}

		// Texte : label du slot + nom de l'item
		textX := int(startX + iconSz + 10)
		ebitenutil.DebugPrintAt(screen, equip.name, textX, int(y+6))
		if itemName != "" {
			ebitenutil.DebugPrintAt(screen, itemName, textX, int(y+22))
		} else {
			ebitenutil.DebugPrintAt(screen, "[ vide ]", textX, int(y+22))
		}
	}
}

// drawItemDetails dessine les détails de l'item sélectionné
func (im *InventoryMenu) drawItemDetails(screen *ebiten.Image, startX, startY float32) {
	var item *entities.Item
	var itemName string
	var itemDesc string
	var stats string

	if im.isInEquipmentPanel {
		// Afficher l'item équipé sélectionné
		slot := im.getEquipmentSlotFromIndex(im.selectedSlotIndex)
		switch slot {
		case entities.SlotWeapon:
			if im.player.Equipment.Weapon != nil {
				item = &im.player.Equipment.Weapon.Item
				itemName = im.player.Equipment.Weapon.Name
				itemDesc = im.player.Equipment.Weapon.Description
				stats = fmt.Sprintf("Degats: +%d | Portee: %.0f", im.player.Equipment.Weapon.Damage, im.player.Equipment.Weapon.AttackRange)
			}
		case entities.SlotChest:
			if im.player.Equipment.Chest != nil {
				item = &im.player.Equipment.Chest.Item
				itemName = im.player.Equipment.Chest.Name
				itemDesc = im.player.Equipment.Chest.Description
				stats = fmt.Sprintf("Defense: +%d | Vie: +%d | Force: +%d", im.player.Equipment.Chest.Defense, im.player.Equipment.Chest.BonusVie, im.player.Equipment.Chest.BonusForce)
			}
		case entities.SlotHelmet:
			if im.player.Equipment.Helmet != nil {
				item = &im.player.Equipment.Helmet.Item
				itemName = im.player.Equipment.Helmet.Name
				itemDesc = im.player.Equipment.Helmet.Description
				stats = fmt.Sprintf("Defense: +%d | Vie: +%d | Force: +%d", im.player.Equipment.Helmet.Defense, im.player.Equipment.Helmet.BonusVie, im.player.Equipment.Helmet.BonusForce)
			}
		case entities.SlotBoots:
			if im.player.Equipment.Boots != nil {
				item = &im.player.Equipment.Boots.Item
				itemName = im.player.Equipment.Boots.Name
				itemDesc = im.player.Equipment.Boots.Description
				stats = fmt.Sprintf("Defense: +%d | Vie: +%d | Force: +%d", im.player.Equipment.Boots.Defense, im.player.Equipment.Boots.BonusVie, im.player.Equipment.Boots.BonusForce)
			}
		case entities.SlotRing, entities.SlotAmulet, entities.SlotBelt:
			// Les accessoires
			var accessory *entities.Accessory
			switch slot {
			case entities.SlotRing:
				accessory = im.player.Equipment.Ring
			case entities.SlotAmulet:
				accessory = im.player.Equipment.Amulet
			case entities.SlotBelt:
				accessory = im.player.Equipment.Belt
			}
			if accessory != nil {
				item = &accessory.Item
				itemName = accessory.Name
				itemDesc = accessory.Description
				stats = fmt.Sprintf("Vie: +%d | Force: +%d", accessory.BonusVie, accessory.BonusForce)
			}
		}
	} else {
		// Afficher l'item de l'inventaire sélectionné
		invSlot := im.player.Inventory.Slots[im.selectedSlotIndex]
		if invSlot != nil && invSlot.Item != nil {
			item = getBaseItem(invSlot.Item)
			if item != nil {
				itemName = item.Name
				itemDesc = item.Description
			}
		}
	}

	if item != nil {
		// Fond des détails
		detailsHeight := float32(70)
		vector.DrawFilledRect(screen, startX, startY, im.menuWidth-20, detailsHeight, color.RGBA{30, 30, 40, 255}, false)
		vector.StrokeRect(screen, startX, startY, im.menuWidth-20, detailsHeight, 1, color.RGBA{80, 80, 90, 255}, false)

		// Nom de l'item
		ebitenutil.DebugPrintAt(screen, itemName, int(startX+5), int(startY+5))

		// Description
		ebitenutil.DebugPrintAt(screen, itemDesc, int(startX+5), int(startY+25))

		// Stats (si disponibles)
		if stats != "" {
			ebitenutil.DebugPrintAt(screen, stats, int(startX+5), int(startY+45))
		}
	}
}

// getBaseItem extrait l'Item de base d'un item (peut être *Weapon, *Armor, *Accessory, ou *Item)
func getBaseItem(item interface{}) *entities.Item {
	if item == nil {
		return nil
	}
	switch i := item.(type) {
	case *entities.Weapon:
		return &i.Item
	case *entities.Armor:
		return &i.Item
	case *entities.Accessory:
		return &i.Item
	case *entities.Item:
		return i
	default:
		return nil
	}
}

// getItemColor retourne une couleur par défaut pour un type d'item
func getItemColor(itemType entities.ItemType) color.RGBA {
	switch itemType {
	case entities.ItemTypeWeapon:
		return color.RGBA{200, 100, 100, 255} // Rouge pour les armes
	case entities.ItemTypeArmor:
		return color.RGBA{100, 150, 200, 255} // Bleu pour les armures
	case entities.ItemTypeAccessory:
		return color.RGBA{150, 100, 200, 255} // Violet pour les accessoires
	case entities.ItemTypeConsumable:
		return color.RGBA{100, 200, 100, 255} // Vert pour les consommables
	default:
		return color.RGBA{150, 150, 150, 255} // Gris par défaut
	}
}
