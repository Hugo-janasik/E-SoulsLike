package save

import (
	"encoding/json"
	"os"
	"time"
)

// SaveData contient toutes les données à sauvegarder
type SaveData struct {
	// Informations de sauvegarde
	SaveName    string    `json:"save_name"`
	SaveDate    time.Time `json:"save_date"`
	Version     string    `json:"version"`

	// Joueur - Position et ressources
	PlayerX       float64 `json:"player_x"`
	PlayerY       float64 `json:"player_y"`
	PlayerHealth  int     `json:"player_health"`
	PlayerStamina int     `json:"player_stamina"`

	// Joueur - Stats et progression
	PlayerLevel   int `json:"player_level"`
	PlayerSouls   int `json:"player_souls"`
	PlayerForce   int `json:"player_force"`
	PlayerStaminaStat int `json:"player_stamina_stat"` // Stat (différent de la stamina actuelle)
	PlayerVie     int `json:"player_vie"`

	// Inventaire et équipement
	EquippedWeaponID  string         `json:"equipped_weapon_id"`
	EquippedChestID   string         `json:"equipped_chest_id"`
	EquippedHelmetID  string         `json:"equipped_helmet_id"`
	EquippedBootsID   string         `json:"equipped_boots_id"`
	EquippedRingID    string         `json:"equipped_ring_id"`
	EquippedAmuletID  string         `json:"equipped_amulet_id"`
	EquippedBeltID    string         `json:"equipped_belt_id"`
	InventoryItems    map[string]int `json:"inventory_items"` // ItemID -> Quantité

	// Statistiques
	EnemiesKilled    int           `json:"enemies_killed"`
	PlayTime         time.Duration `json:"play_time"`
	TotalDamageDealt int           `json:"total_damage_dealt"`
	TotalDamageTaken int           `json:"total_damage_taken"`

	// Monde
	MapSeed int64 `json:"map_seed"`

	// Donjons
	IsInDungeon      bool   `json:"is_in_dungeon"`
	CurrentDungeonID string `json:"current_dungeon_id"`
	CurrentFloor     int    `json:"current_floor"`
}

// NewSaveData crée une nouvelle sauvegarde
func NewSaveData(saveName string) *SaveData {
	return &SaveData{
		SaveName: saveName,
		SaveDate: time.Now(),
		Version:  "1.0.0",
	}
}

// SaveToFile sauvegarde les données dans un fichier
func (sd *SaveData) SaveToFile(filename string) error {
	// Créer le dossier saves s'il n'existe pas
	if err := os.MkdirAll("saves", 0755); err != nil {
		return err
	}

	file, err := os.Create("saves/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(sd)
}

// LoadFromFile charge une sauvegarde depuis un fichier
func LoadFromFile(filename string) (*SaveData, error) {
	file, err := os.Open("saves/" + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var saveData SaveData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&saveData)
	if err != nil {
		return nil, err
	}

	return &saveData, nil
}

// ListSaves retourne la liste des sauvegardes disponibles
func ListSaves() ([]string, error) {
	files, err := os.ReadDir("saves")
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var saves []string
	for _, file := range files {
		if !file.IsDir() && len(file.Name()) > 5 && file.Name()[len(file.Name())-5:] == ".json" {
			saves = append(saves, file.Name())
		}
	}

	return saves, nil
}

// DeleteSave supprime une sauvegarde
func DeleteSave(filename string) error {
	return os.Remove("saves/" + filename)
}

// QuickSave crée une sauvegarde rapide
func QuickSave(data *SaveData) error {
	data.SaveName = "QuickSave"
	data.SaveDate = time.Now()
	return data.SaveToFile("quicksave.json")
}

// LoadQuickSave charge la sauvegarde rapide
func LoadQuickSave() (*SaveData, error) {
	return LoadFromFile("quicksave.json")
}
