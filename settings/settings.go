package settings

import (
	"encoding/json"
	"os"
)

// Settings contient les paramètres du jeu
type Settings struct {
	// Audio
	MasterVolume float64 `json:"master_volume"`
	MusicVolume  float64 `json:"music_volume"`
	SFXVolume    float64 `json:"sfx_volume"`
	AudioEnabled bool    `json:"audio_enabled"`

	// Gameplay
	ShowFPS      bool `json:"show_fps"`
	ShowDamageNumbers bool `json:"show_damage_numbers"`
	CameraSmoothing float64 `json:"camera_smoothing"`

	// Controls
	MouseSensitivity float64 `json:"mouse_sensitivity"`
}

// NewDefaultSettings crée les paramètres par défaut
func NewDefaultSettings() *Settings {
	return &Settings{
		MasterVolume:     0.8,
		MusicVolume:      0.7,
		SFXVolume:        0.9,
		AudioEnabled:     true,
		ShowFPS:          false,
		ShowDamageNumbers: true,
		CameraSmoothing:  0.1,
		MouseSensitivity: 1.0,
	}
}

// SaveToFile sauvegarde les paramètres dans un fichier
func (s *Settings) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s)
}

// LoadFromFile charge les paramètres depuis un fichier
func LoadSettingsFromFile(filename string) (*Settings, error) {
	file, err := os.Open(filename)
	if err != nil {
		// Si le fichier n'existe pas, retourner les paramètres par défaut
		if os.IsNotExist(err) {
			return NewDefaultSettings(), nil
		}
		return nil, err
	}
	defer file.Close()

	var settings Settings
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	if err != nil {
		return NewDefaultSettings(), err
	}

	return &settings, nil
}
