# Guide de contribution

Merci de ton intérêt pour contribuer à E-SoulsLike !

## Comment commencer

1. **Familiarise-toi avec le code**
   - Lis [README.md](README.md) pour comprendre le jeu
   - Consulte [ARCHITECTURE.md](ARCHITECTURE.md) pour comprendre la structure
   - Lis [EXTENDING.md](EXTENDING.md) pour voir des exemples

2. **Configure ton environnement**
   ```bash
   # Cloner le projet
   git clone [url-du-repo]
   cd e-soulslike

   # Installer les dépendances
   make deps

   # Compiler le jeu
   make build

   # Lancer le jeu
   make run
   ```

3. **Vérifie que tout fonctionne**
   ```bash
   # Exécuter les tests
   make test

   # Vérifier le code
   make lint
   ```

## Processus de contribution

### 1. Créer une branche

```bash
git checkout -b feature/ma-nouvelle-fonctionnalite
# ou
git checkout -b fix/correction-bug
```

### 2. Développer

- Écris du code propre et commenté
- Suis les conventions de code Go
- Ajoute des tests pour tes nouvelles fonctionnalités
- Documente les fonctions publiques

### 3. Tester

```bash
# Tests unitaires
go test ./...

# Tests avec couverture
go test -cover ./...

# Vérifier le formatage
go fmt ./...

# Vérifier avec vet
go vet ./...
```

### 4. Commit

Utilise des messages de commit clairs:

```
feat: Ajout du système d'inventaire
fix: Correction du bug de collision
docs: Mise à jour du README
refactor: Réorganisation du système de combat
test: Ajout de tests pour les ennemis
```

### 5. Push et Pull Request

```bash
git push origin feature/ma-nouvelle-fonctionnalite
```

Ensuite, crée une Pull Request sur GitHub avec:
- Une description claire de ce que tu as fait
- Pourquoi ces changements sont nécessaires
- Des captures d'écran si pertinent

## Conventions de code

### Nommage

```go
// Exporté (public) - PascalCase
type PlayerHealth struct { }
func CalculateDamage() { }

// Non exporté (privé) - camelCase
type internalState struct { }
func calculateDistance() { }
```

### Commentaires

```go
// Package combat gère le système de combat du jeu
package combat

// CalculateDamage calcule les dégâts infligés en tenant compte
// des modificateurs et des résistances
func CalculateDamage(base int, modifier float64) int {
    // Code...
}
```

### Organisation des imports

```go
import (
    // Standard library en premier
    "fmt"
    "math"

    // Packages tiers ensuite
    "github.com/hajimehoshi/ebiten/v2"

    // Packages locaux en dernier
    "e-soulslike/entities"
    "e-soulslike/combat"
)
```

### Structure des fichiers

```go
package example

// 1. Imports

// 2. Constantes
const (
    MaxHealth = 100
)

// 3. Types
type Player struct { }

// 4. Constructeurs
func NewPlayer() *Player { }

// 5. Méthodes
func (p *Player) Update() { }

// 6. Fonctions utilitaires
func calculateSomething() { }
```

## Types de contributions

### Nouvelles fonctionnalités

Avant d'implémenter une grosse fonctionnalité:
1. Ouvre une issue pour en discuter
2. Attends un retour avant de commencer
3. Découpe en petites Pull Requests si possible

Exemples:
- Système d'inventaire
- Nouveaux types d'ennemis
- Système de quêtes
- Effets visuels

### Corrections de bugs

1. Décris le bug dans une issue
2. Explique comment le reproduire
3. Propose une solution dans ta PR

### Documentation

- Amélioration du README
- Correction de typos
- Ajout d'exemples
- Traductions

### Tests

- Ajout de tests unitaires
- Tests d'intégration
- Tests de performance

## Zones à améliorer

### Court terme
- [ ] Sprites et animations (partiellement : assets présents, intégration à finaliser)
- [ ] Système de sons et musiques
- [ ] Plus de types d'ennemis
- [x] Système d'inventaire
- [x] UI améliorée (menus complets)

### Moyen terme
- [x] Système de sauvegarde
- [x] Génération procédurale (zones et donjons)
- [ ] Boss avec patterns d'attaque
- [ ] Quêtes
- [ ] Dialogues et PNJ

### Long terme
- [ ] Multijoueur
- [ ] Éditeur de niveaux
- [ ] Modding support

## Standards de qualité

### Code

- [ ] Le code compile sans erreur
- [ ] `go vet` ne retourne aucune erreur
- [ ] Le code est formaté avec `go fmt`
- [ ] Les tests passent (`go test ./...`)
- [ ] Pas de code mort ou commenté
- [ ] Variables et fonctions bien nommées

### Documentation

- [ ] Fonctions publiques documentées
- [ ] README à jour si nécessaire
- [ ] ARCHITECTURE.md à jour si changement structurel
- [ ] Commentaires pour la logique complexe

### Tests

- [ ] Tests pour les nouvelles fonctionnalités
- [ ] Tests pour les corrections de bugs
- [ ] Couverture de code > 70% pour le nouveau code

## Exemples de bonnes contributions

### Exemple 1: Ajouter un nouveau type d'ennemi

```go
// entities/enemy.go

const (
    EnemyTypeBasic EnemyType = iota
    EnemyTypeFast
    EnemyTypeTank
    EnemyTypeRanged  // NOUVEAU
)

func NewEnemy(x, y float64, enemyType EnemyType) *Enemy {
    enemy := &Enemy{
        // ... code existant
    }

    switch enemyType {
    // ... cas existants
    case EnemyTypeRanged:
        enemy.Health = 40
        enemy.MaxHealth = 40
        enemy.Speed = 1.2
        enemy.AttackRange = 150  // Portée plus grande
        enemy.DetectionRange = 300
    }

    return enemy
}
```

### Exemple 2: Ajouter un test

```go
// combat/combat_test.go

func TestDamageCalculation(t *testing.T) {
    tests := []struct {
        name     string
        base     int
        modifier float64
        expected int
    }{
        {"No modifier", 10, 1.0, 10},
        {"Double damage", 10, 2.0, 20},
        {"Half damage", 10, 0.5, 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateDamage(tt.base, tt.modifier)
            if result != tt.expected {
                t.Errorf("got %d, want %d", result, tt.expected)
            }
        })
    }
}
```

## Besoin d'aide ?

- Ouvre une issue avec la question
- Consulte la documentation Ebiten: https://ebitengine.org/
- Lis le code existant pour comprendre les patterns

## Code de conduite

- Sois respectueux et constructif
- Accepte les critiques et feedbacks
- Aide les autres contributeurs
- Privilégie la qualité à la quantité

## Questions fréquentes

**Q: Comment tester mes changements visuels ?**
R: Lance le jeu avec `make run` et teste manuellement. Ajoute des captures d'écran dans ta PR.

**Q: Puis-je utiliser des bibliothèques externes ?**
R: Oui, mais discutes-en d'abord dans une issue. Privilégie les bibliothèques légères et maintenues.

**Q: Comment déboguer le jeu ?**
R: Utilise `ebitenutil.DebugPrint()` pour afficher des infos. Tu peux aussi utiliser un debugger Go.

**Q: Le jeu lag, que faire ?**
R: Profile avec `go tool pprof` pour identifier les goulots d'étranglement.

## Remerciements

Merci à tous les contributeurs qui rendent ce projet meilleur !

---

Bon code ! 🎮
