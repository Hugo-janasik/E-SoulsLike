package ui

import (
	"e-soulslike/entities"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const shopIconSize = float32(40) // Taille de l'icône dans le shop

// ShopItem représente un item à vendre dans le shop
type ShopItem struct {
	CreateItem func() interface{} // Fonction pour créer l'item
	ItemName   string
	ItemDesc   string
	Price      int
	Icon       *ebiten.Image // Icône affichée dans le shop
}

// ShopMenu représente le menu du magasin
type ShopMenu struct {
	x, y, width, height float32
	player              *entities.Player
	items               []ShopItem
	selectedIndex       int
	OnClose             func()
	feedbackMessage     string
	feedbackTimer       int
}

// NewShopMenu crée un nouveau menu de shop
func NewShopMenu(x, y, width, height float32, player *entities.Player, onClose func()) *ShopMenu {
	sm := &ShopMenu{
		x:             x,
		y:             y,
		width:         width,
		height:        height,
		player:        player,
		selectedIndex: 0,
		OnClose:       onClose,
	}

	// Définir les items du shop
	rawItems := []ShopItem{
		{
			CreateItem: func() interface{} { return entities.NewBeginnerHelmet() },
			ItemName:   "Casque du débutant",
			ItemDesc:   "+10 Vie",
			Price:      50,
		},
		{
			CreateItem: func() interface{} { return entities.NewBeginnerChestplate() },
			ItemName:   "Plastron du débutant",
			ItemDesc:   "+20 Vie",
			Price:      50,
		},
		{
			CreateItem: func() interface{} { return entities.NewBeginnerLeggings() },
			ItemName:   "Jambe du débutant",
			ItemDesc:   "+15 Vie",
			Price:      50,
		},
		{
			CreateItem: func() interface{} { return entities.NewBeginnerBoots() },
			ItemName:   "Bottes du débutant",
			ItemDesc:   "+5 Vie",
			Price:      50,
		},
		{
			CreateItem: func() interface{} { return entities.NewBeginnerAmulet() },
			ItemName:   "Amulette du débutant",
			ItemDesc:   "+3 Force",
			Price:      50,
		},
		{
			CreateItem: func() interface{} { return entities.NewBeginnerRing() },
			ItemName:   "Anneau du débutant",
			ItemDesc:   "+2 Force",
			Price:      50,
		},
	}

	// Charger les icônes une seule fois à l'initialisation du shop
	// (on crée une instance temporaire juste pour extraire l'IconSprite)
	for i, raw := range rawItems {
		instance := raw.CreateItem()
		if baseItem := extractBaseItem(instance); baseItem != nil {
			rawItems[i].Icon = baseItem.IconSprite
		}
	}

	sm.items = rawItems
	return sm
}

// extractBaseItem extrait l'*entities.Item d'un item quelconque
func extractBaseItem(item interface{}) *entities.Item {
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
	}
	return nil
}

// Update met à jour le menu du shop
func (sm *ShopMenu) Update(mouseX, mouseY int, mousePressed, mouseJustPressed bool) {
	if sm.feedbackTimer > 0 {
		sm.feedbackTimer--
		if sm.feedbackTimer == 0 {
			sm.feedbackMessage = ""
		}
	}

	itemHeight := float32(60)
	itemY := sm.y + 100

	for i := range sm.items {
		itemRect := struct{ x, y, width, height float32 }{
			x:      sm.x + 20,
			y:      itemY + float32(i)*itemHeight,
			width:  sm.width - 40,
			height: itemHeight - 5,
		}

		if float32(mouseX) >= itemRect.x && float32(mouseX) <= itemRect.x+itemRect.width &&
			float32(mouseY) >= itemRect.y && float32(mouseY) <= itemRect.y+itemRect.height {
			sm.selectedIndex = i
			if mouseJustPressed {
				sm.buyItem(i)
			}
		}
	}
}

// playerHasItem vérifie si le joueur possède déjà l'item (inventaire ou équipé)
func (sm *ShopMenu) playerHasItem(index int) bool {
	if index < 0 || index >= len(sm.items) {
		return false
	}
	temp := sm.items[index].CreateItem()
	base := extractBaseItem(temp)
	if base == nil {
		return false
	}
	return sm.player.Inventory.HasItemByID(base.ID) || sm.player.Equipment.HasItemByID(base.ID)
}

// buyItem achète un item
func (sm *ShopMenu) buyItem(index int) {
	if index < 0 || index >= len(sm.items) {
		return
	}

	shopItem := sm.items[index]

	if sm.playerHasItem(index) {
		sm.feedbackMessage = "Deja possede!"
		sm.feedbackTimer = 120
		return
	}

	if sm.player.Stats.Souls < shopItem.Price {
		sm.feedbackMessage = "Pas assez d'ames!"
		sm.feedbackTimer = 120
		return
	}

	itemInterface := shopItem.CreateItem()
	success := sm.player.Inventory.AddItem(itemInterface, 1)

	if success {
		sm.player.Stats.Souls -= shopItem.Price
		sm.feedbackMessage = "Achat reussi!"
		sm.feedbackTimer = 120
	} else {
		sm.feedbackMessage = "Inventaire plein!"
		sm.feedbackTimer = 120
	}
}

// Draw dessine le menu du shop
func (sm *ShopMenu) Draw(screen *ebiten.Image) {
	// Fond du menu
	vector.DrawFilledRect(screen, sm.x, sm.y, sm.width, sm.height, color.RGBA{40, 40, 50, 255}, false)
	vector.StrokeRect(screen, sm.x, sm.y, sm.width, sm.height, 3, color.RGBA{200, 200, 50, 255}, false)

	// Titre
	titleText := "=== MAGASIN ==="
	titleX := int(sm.x+sm.width/2) - len(titleText)*6/2
	ebitenutil.DebugPrintAt(screen, titleText, titleX, int(sm.y)+30)

	// Âmes du joueur
	soulsText := fmt.Sprintf("Ames: %d", sm.player.Stats.Souls)
	ebitenutil.DebugPrintAt(screen, soulsText, int(sm.x+20), int(sm.y+60))

	// Items
	itemHeight := float32(60)
	itemY := sm.y + 100

	for i, shopItem := range sm.items {
		y := itemY + float32(i)*itemHeight

		// Couleur de fond selon sélection et possession
		alreadyOwned := sm.playerHasItem(i)
		bgColor := color.RGBA{60, 60, 70, 255}
		if alreadyOwned {
			bgColor = color.RGBA{40, 40, 40, 255}
		} else if i == sm.selectedIndex {
			bgColor = color.RGBA{80, 100, 255, 255}
		}
		vector.DrawFilledRect(screen, sm.x+20, y, sm.width-40, itemHeight-5, bgColor, false)
		borderColor := color.RGBA{100, 100, 110, 255}
		if alreadyOwned {
			borderColor = color.RGBA{80, 80, 80, 255}
		}
		vector.StrokeRect(screen, sm.x+20, y, sm.width-40, itemHeight-5, 1, borderColor, false)

		// --- Icône ---
		iconX := sm.x + 28
		iconY := y + (itemHeight-5-shopIconSize)/2
		if shopItem.Icon != nil {
			// Scaler l'icône pour tenir dans shopIconSize × shopIconSize
			bounds := shopItem.Icon.Bounds()
			scaleX := float64(shopIconSize) / float64(bounds.Dx())
			scaleY := float64(shopIconSize) / float64(bounds.Dy())
			scale := scaleX
			if scaleY < scale {
				scale = scaleY
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(iconX), float64(iconY))
			screen.DrawImage(shopItem.Icon, op)
		} else {
			// Fallback : carré de couleur si pas d'icône
			vector.DrawFilledRect(screen, iconX, iconY, shopIconSize, shopIconSize, color.RGBA{100, 120, 180, 180}, false)
		}
		// Bordure de l'icône
		vector.StrokeRect(screen, iconX, iconY, shopIconSize, shopIconSize, 1, color.RGBA{150, 150, 160, 200}, false)

		// --- Texte (décalé après l'icône) ---
		textX := int(iconX + shopIconSize + 8)
		ebitenutil.DebugPrintAt(screen, shopItem.ItemName, textX, int(y+10))
		ebitenutil.DebugPrintAt(screen, shopItem.ItemDesc, textX, int(y+28))

		// Prix (aligné à droite) ou indicateur "Possede"
		if alreadyOwned {
			ebitenutil.DebugPrintAt(screen, "[Possede]", int(sm.x+sm.width-110), int(y+25))
		} else {
			priceText := fmt.Sprintf("Prix: %d ames", shopItem.Price)
			ebitenutil.DebugPrintAt(screen, priceText, int(sm.x+sm.width-150), int(y+25))
		}
	}

	// Message de feedback
	if sm.feedbackMessage != "" {
		feedbackX := int(sm.x+sm.width/2) - len(sm.feedbackMessage)*6/2
		ebitenutil.DebugPrintAt(screen, sm.feedbackMessage, feedbackX, int(sm.y+sm.height-60))
	}

	// Instructions
	instructionsText := "Cliquez sur un item pour l'acheter | ECHAP pour revenir"
	instructionsX := int(sm.x+sm.width/2) - len(instructionsText)*6/2
	ebitenutil.DebugPrintAt(screen, instructionsText, instructionsX, int(sm.y+sm.height-30))
}
