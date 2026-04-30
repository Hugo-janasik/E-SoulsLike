package combat

import (
	"e-soulslike/entities"
	"e-soulslike/stats"
	"math"
	"math/rand"
)

// Constantes de combat
const (
	KnockbackForceEnemy  = 15.0 // Force de knockback infligée aux ennemis
	KnockbackForcePlayer = 8.0  // Force de knockback infligée au joueur
	StunDuration         = 20   // Durée de l'étourdissement en frames
	CriticalHitThreshold = 0.8  // Seuil de coup critique (20% de chance)
	CriticalHitDamage    = 10   // Dégâts bonus d'un coup critique
)

// CombatSystem gère les interactions de combat
type CombatSystem struct {
	Player    *entities.Player
	Enemies   []*entities.Enemy
	GameStats *stats.GameStats // pour enregistrer les statistiques en temps réel
}

// NewCombatSystem crée un nouveau système de combat
func NewCombatSystem(player *entities.Player, enemies []*entities.Enemy, gameStats *stats.GameStats) *CombatSystem {
	return &CombatSystem{
		Player:    player,
		Enemies:   enemies,
		GameStats: gameStats,
	}
}

// Update met à jour le système de combat
func (cs *CombatSystem) Update() {
	cs.checkPlayerAttacks()
	cs.checkEnemyAttacks()
}

// checkPlayerAttacks vérifie si le joueur touche des ennemis
func (cs *CombatSystem) checkPlayerAttacks() {
	attackX, attackY, attackRange, isAttacking := cs.Player.GetAttackHitbox()

	if !isAttacking {
		return
	}

	for _, enemy := range cs.Enemies {
		if !enemy.IsAlive() {
			continue
		}

		dx := enemy.X - attackX
		dy := enemy.Y - attackY
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance < attackRange+enemy.Width/2 {
			wasAlive := enemy.IsAlive()

			// Dégâts = équipement + éventuel coup critique
			damage := cs.Player.GetDamageWithEquipment() + calculateCriticalHit()
			enemy.TakeDamage(damage)

			// Enregistrer les dégâts infligés
			if cs.GameStats != nil {
				cs.GameStats.RecordDamageDealt(damage)
			}

			// Donner des âmes et enregistrer le kill si l'ennemi vient de mourir
			if wasAlive && !enemy.IsAlive() {
				souls := entities.GetSoulDropForEnemy(enemy.Type)
				cs.Player.Stats.AddSouls(souls)
				if cs.GameStats != nil {
					cs.GameStats.RecordEnemyKill()
				}
			}

			// Étourdir l'ennemi s'il est encore en vie
			if enemy.IsAlive() {
				enemy.Stun(StunDuration)
			}

			// Knockback de l'ennemi
			if distance > 0 {
				enemy.X += (dx / distance) * KnockbackForceEnemy
				enemy.Y += (dy / distance) * KnockbackForceEnemy
			}
		}
	}
}

// checkEnemyAttacks vérifie si les ennemis touchent le joueur
func (cs *CombatSystem) checkEnemyAttacks() {
	for _, enemy := range cs.Enemies {
		if !enemy.IsAlive() {
			continue
		}

		if !enemy.CanAttack() {
			continue
		}

		dx := cs.Player.X - enemy.X
		dy := cs.Player.Y - enemy.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance < enemy.AttackRange {
			enemy.PerformAttack()

			// Infliger des dégâts uniquement si le joueur n'esquive pas
			if !cs.Player.IsDodging {
				damage := enemy.GetAttackDamage()
				cs.Player.TakeDamage(damage)

				// Enregistrer les dégâts reçus
				if cs.GameStats != nil {
					cs.GameStats.RecordDamageTaken(damage)
				}

				// Petit knockback pour le joueur
				if distance > 0 {
					cs.Player.X += (dx / distance) * KnockbackForcePlayer
					cs.Player.Y += (dy / distance) * KnockbackForcePlayer
				}
			}
		}
	}
}

// calculateCriticalHit retourne des dégâts bonus en cas de coup critique (20% de chance)
func calculateCriticalHit() int {
	if rand.Float64() > CriticalHitThreshold {
		return CriticalHitDamage
	}
	return 0
}

// CheckCollision vérifie une collision circulaire entre deux entités
func CheckCollision(x1, y1, r1, x2, y2, r2 float64) bool {
	dx := x2 - x1
	dy := y2 - y1
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < (r1 + r2)
}

// GetDistance retourne la distance entre deux points
func GetDistance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// GetDirection retourne la direction normalisée de (x1,y1) vers (x2,y2)
func GetDirection(x1, y1, x2, y2 float64) (float64, float64) {
	dx := x2 - x1
	dy := y2 - y1
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance == 0 {
		return 0, 0
	}

	return dx / distance, dy / distance
}
