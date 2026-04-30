package stats

import (
	"time"
)

// GameStats contient les statistiques de la partie
type GameStats struct {
	// Combat
	EnemiesKilled    int
	TotalDamageDealt int
	TotalDamageTaken int
	AttacksLanded    int
	AttacksMissed    int
	DodgesPerformed  int

	// Temps
	StartTime    time.Time
	PlayTime     time.Duration
	TotalPaused  time.Duration

	// Mouvement
	DistanceTraveled float64

	// Stamina
	StaminaUsed      int
	TimesOutOfStamina int
}

// NewGameStats crée de nouvelles statistiques
func NewGameStats() *GameStats {
	return &GameStats{
		StartTime: time.Now(),
	}
}

// RecordEnemyKill enregistre la mort d'un ennemi
func (gs *GameStats) RecordEnemyKill() {
	gs.EnemiesKilled++
}

// RecordDamageDealt enregistre les dégâts infligés
func (gs *GameStats) RecordDamageDealt(damage int) {
	gs.TotalDamageDealt += damage
	gs.AttacksLanded++
}

// RecordDamageTaken enregistre les dégâts reçus
func (gs *GameStats) RecordDamageTaken(damage int) {
	gs.TotalDamageTaken += damage
}

// RecordAttackMiss enregistre une attaque ratée
func (gs *GameStats) RecordAttackMiss() {
	gs.AttacksMissed++
}

// RecordDodge enregistre une esquive
func (gs *GameStats) RecordDodge() {
	gs.DodgesPerformed++
}

// RecordDistance enregistre la distance parcourue
func (gs *GameStats) RecordDistance(distance float64) {
	gs.DistanceTraveled += distance
}

// RecordStaminaUsed enregistre l'utilisation de stamina
func (gs *GameStats) RecordStaminaUsed(amount int) {
	gs.StaminaUsed += amount
}

// RecordOutOfStamina enregistre qu'on est à court de stamina
func (gs *GameStats) RecordOutOfStamina() {
	gs.TimesOutOfStamina++
}

// UpdatePlayTime met à jour le temps de jeu
func (gs *GameStats) UpdatePlayTime() {
	gs.PlayTime = time.Since(gs.StartTime) - gs.TotalPaused
}

// GetAccuracy retourne le taux de précision des attaques
func (gs *GameStats) GetAccuracy() float64 {
	total := gs.AttacksLanded + gs.AttacksMissed
	if total == 0 {
		return 0
	}
	return float64(gs.AttacksLanded) / float64(total) * 100
}

// GetAverageDamagePerHit retourne les dégâts moyens par coup
func (gs *GameStats) GetAverageDamagePerHit() float64 {
	if gs.AttacksLanded == 0 {
		return 0
	}
	return float64(gs.TotalDamageDealt) / float64(gs.AttacksLanded)
}
