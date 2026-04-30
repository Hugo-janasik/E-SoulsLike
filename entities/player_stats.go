package entities

import "math"

// PlayerStats représente les statistiques du joueur (leveling)
type PlayerStats struct {
	// Niveau et progression
	Level int
	Souls int // Monnaie (comme les âmes de Dark Souls)

	// Caractéristiques de base
	Force   int // Augmente les dégâts
	Stamina int // Augmente l'endurance max
	Vie     int // Augmente la santé max

	// Stats dérivées (calculées)
	MaxHealth  float64
	MaxStamina float64
	Damage     float64
}

// NewPlayerStats crée les stats initiales d'un nouveau joueur
func NewPlayerStats() *PlayerStats {
	ps := &PlayerStats{
		Level:   1,
		Souls:   0,
		Force:   10,
		Stamina: 10,
		Vie:     10,
	}
	ps.RecalculateStats()
	return ps
}

// RecalculateStats recalcule les stats dérivées basées sur les caractéristiques
func (ps *PlayerStats) RecalculateStats() {
	// Vie -> MaxHealth (chaque point = 10 HP)
	ps.MaxHealth = float64(ps.Vie * 10)

	// Stamina -> MaxStamina (chaque point = 5 stamina)
	ps.MaxStamina = float64(ps.Stamina * 5)

	// Force -> Damage (chaque point = 2 dégâts)
	ps.Damage = float64(ps.Force * 2)
}

// GetSoulsForNextLevel retourne le coût en âmes pour monter au niveau suivant
func (ps *PlayerStats) GetSoulsForNextLevel() int {
	// Formule Dark Souls-like: coût augmente avec le niveau
	// Niveau 2: ~100 souls, Niveau 10: ~500 souls, etc.
	return int(100 * math.Pow(1.15, float64(ps.Level-1)))
}

// GetSoulsForStatUpgrade retourne le coût pour augmenter une stat
func (ps *PlayerStats) GetSoulsForStatUpgrade() int {
	// Coût basé sur le niveau actuel
	return ps.Level * 50
}

// CanLevelUp vérifie si le joueur peut monter de niveau
func (ps *PlayerStats) CanLevelUp() bool {
	return ps.Souls >= ps.GetSoulsForNextLevel()
}

// LevelUp monte le joueur d'un niveau
func (ps *PlayerStats) LevelUp() bool {
	cost := ps.GetSoulsForNextLevel()
	if ps.Souls >= cost {
		ps.Souls -= cost
		ps.Level++
		// Bonus automatique à chaque niveau
		ps.Force++
		ps.Stamina++
		ps.Vie++
		ps.RecalculateStats()
		return true
	}
	return false
}

// CanUpgradeStat vérifie si le joueur peut améliorer une stat
func (ps *PlayerStats) CanUpgradeStat() bool {
	return ps.Souls >= ps.GetSoulsForStatUpgrade()
}

// UpgradeForce améliore la Force
func (ps *PlayerStats) UpgradeForce() bool {
	cost := ps.GetSoulsForStatUpgrade()
	if ps.Souls >= cost {
		ps.Souls -= cost
		ps.Force++
		ps.RecalculateStats()
		return true
	}
	return false
}

// UpgradeStamina améliore la Stamina
func (ps *PlayerStats) UpgradeStamina() bool {
	cost := ps.GetSoulsForStatUpgrade()
	if ps.Souls >= cost {
		ps.Souls -= cost
		ps.Stamina++
		ps.RecalculateStats()
		return true
	}
	return false
}

// UpgradeVie améliore la Vie
func (ps *PlayerStats) UpgradeVie() bool {
	cost := ps.GetSoulsForStatUpgrade()
	if ps.Souls >= cost {
		ps.Souls -= cost
		ps.Vie++
		ps.RecalculateStats()
		return true
	}
	return false
}

// AddSouls ajoute des âmes (quand on tue un ennemi)
func (ps *PlayerStats) AddSouls(amount int) {
	ps.Souls += amount
}

// GetSoulDropForEnemy retourne le nombre d'âmes qu'un ennemi donne
func GetSoulDropForEnemy(enemyType EnemyType) int {
	switch enemyType {
	case EnemyTypeBasic:
		return 50
	default:
		return 50
	}
}
