# Architecture du jeu E-SoulsLike

Ce document décrit l'architecture et la structure du jeu E-SoulsLike.

## Vue d'ensemble

E-SoulsLike est un jeu d'action-aventure 2D développé en Go avec la bibliothèque Ebiten. Le jeu suit une architecture modulaire avec séparation des responsabilités.

## Architecture globale

```
┌─────────────────────────────────────────────┐
│              Main (main.go)                 │
│  - Point d'entrée                           │
│  - Configuration de la fenêtre              │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│            Game Loop (game.go)              │
│  - Update() : Logique du jeu (60 FPS)      │
│  - Draw() : Rendu graphique                │
│  - Layout() : Gestion de la taille         │
└──────────────────┬──────────────────────────┘
                   │
     ┌─────────────┼──────────────┐
     ▼             ▼              ▼
 ┌────────┐  ┌──────────┐  ┌──────────┐
 │entities│  │  world   │  │    ui    │
 └────────┘  └──────────┘  └──────────┘
     │             │              │
     ▼             ▼              ▼
 ┌────────┐  ┌──────────┐  ┌──────────┐
 │ combat │  │  input   │  │  stats   │
 └────────┘  └──────────┘  └──────────┘
                                │
                    ┌───────────┴───────────┐
                    ▼                       ▼
              ┌──────────┐           ┌──────────┐
              │ settings │           │   save   │
              └──────────┘           └──────────┘
```

## Modules principaux

### 1. Package `main`

**Fichier**: [main.go](main.go)

Responsabilités:
- Initialisation de la fenêtre Ebiten
- Configuration des paramètres globaux (résolution, titre)
- Lancement de la boucle de jeu

### 2. Package `game`

**Fichier**: [game/game.go](game/game.go)

Responsabilités:
- Gestion de la boucle de jeu principale
- Coordination entre tous les systèmes
- Méthode `Update()`: Met à jour la logique (60 fois/seconde)
- Méthode `Draw()`: Rendu graphique
- Méthode `Layout()`: Gestion de la taille de l'écran

Structure:
```go
type Game struct {
    width          int
    height         int
    currentState   GameState          // MainMenu, Playing, Dead
    player         *entities.Player
    camera         *world.Camera
    inputMgr       *input.InputManager
    combatSystem   *combat.CombatSystem
    zoneManager    *world.ZoneManager
    dungeonManager *world.DungeonManager
    mainMenu       *ui.MainMenu
    pauseMenu      *ui.PauseMenu
    campfireMenu   *ui.CampfireMenu
    inventoryMenu  *ui.InventoryMenu
    hud            *ui.HUD
    deathScreen    *ui.DeathScreen
    gameStats      *stats.GameStats
    settings       *settings.Settings
    isPaused       bool
    isAtCampfire   bool
    isInInventory  bool
    saveMessage    string
    saveMessageTimer int
}
```

### 3. Package `entities`

**Fichiers**:
- [entities/entity.go](entities/entity.go) - Interface et structure de base
- [entities/player.go](entities/player.go) - Logique du joueur
- [entities/enemy.go](entities/enemy.go) - Logique des ennemis

Responsabilités:
- Définition de toutes les entités du jeu
- Gestion des comportements (IA, mouvements)
- États et statistiques

#### Structure BaseEntity

```go
type BaseEntity struct {
    X, Y              float64  // Position
    Width, Height     float64  // Taille
    Health, MaxHealth int      // Santé
    Speed             float64  // Vitesse
    IsInvincible      bool     // État invincible
    InvincibleTimer   int      // Timer d'invincibilité
}
```

#### Player

États du joueur:
- Normal: Peut se déplacer et agir
- Attacking: Animation d'attaque en cours
- Dodging: Esquive avec invincibilité

Ressources:
- Health: Points de vie
- Stamina: Pour sprint, attaque, esquive

#### Enemy

Machine à états:
```
Idle ─────► Chasing ─────► Attacking
  ▲            │               │
  │            ▼               │
  └────────── Lost Player ◄────┘
               │
               ▼
            Stunned (temporaire)
```

Types d'ennemis:
- **Basic**: Équilibré (50 HP, vitesse 1.5)
- **Fast**: Rapide mais fragile (30 HP, vitesse 2.5)
- **Tank**: Lent mais résistant (150 HP, vitesse 1.0)

### 4. Package `combat`

**Fichier**: [combat/combat.go](combat/combat.go)

Responsabilités:
- Détection des collisions
- Calcul des dégâts
- Gestion des knockbacks
- Système de coups critiques

Fonctions utilitaires:
```go
CheckCollision(x1, y1, r1, x2, y2, r2 float64) bool
GetDistance(x1, y1, x2, y2 float64) float64
GetDirection(x1, y1, x2, y2 float64) (float64, float64)
```

### 5. Package `world`

**Fichiers**:
- [world/camera.go](world/camera.go) - Système de caméra 2D
- [world/tilemap.go](world/tilemap.go) - Génération et rendu de la carte
- [world/zone.go](world/zone.go) - Zones du monde
- [world/zone_manager.go](world/zone_manager.go) - Transitions entre zones
- [world/dungeon.go](world/dungeon.go) - Donjons et étages
- [world/dungeon_generator.go](world/dungeon_generator.go) - Génération procédurale
- [world/dungeon_manager.go](world/dungeon_manager.go) - Gestion des donjons
- [world/portal.go](world/portal.go) - Portails de transition

#### Camera

Fonctionnalités:
- Suivi fluide du joueur (smooth follow)
- Conversion coordonnées monde ↔ écran
- Culling de frustum (n'afficher que le visible)
- Support du zoom

```go
type Camera struct {
    X, Y          float64  // Position de la caméra
    Width, Height float64  // Taille de la vue
    Zoom          float64  // Niveau de zoom
    FollowSpeed   float64  // Vitesse de suivi (0-1)
}
```

#### TileMap

Génération procédurale:
- Différents types de tuiles (herbe, terre, pierre, eau)
- Variation de couleur pour plus de réalisme
- Rendu optimisé (seulement les tuiles visibles)

```go
const TileSize = 32  // Taille d'une tuile en pixels

type Tile struct {
    Type  TileType
    Color color.RGBA
}
```

### 6. Package `input`

**Fichier**: [input/input.go](input/input.go)

Responsabilités:
- Gestion centralisée des entrées clavier/souris
- Détection des pressions (continues et instantanées)
- Support QWERTY et AZERTY

```go
type InputManager struct {
    pressedKeys      map[ebiten.Key]bool  // Touches enfoncées
    justPressedKeys  map[ebiten.Key]bool  // Nouvelles pressions
    prevPressedKeys  map[ebiten.Key]bool  // État frame précédente
    mouseX, mouseY   int                   // Position souris
    // ...
}
```

## Flux de données

### Boucle Update (Logique)

```
1. InputManager.Update()
   └─> Capture les entrées clavier/souris

2. Gestion de l'état (GameState)
   ├─> MainMenu → MainMenu.Update()
   ├─> Dead → DeathScreen.Update()
   └─> Playing → suite...

3. Gestion pause / inventaire / campfire
   ├─> Échap → PauseMenu.Update()
   ├─> E → Interaction (campfire, portail)
   └─> I/B → InventoryMenu.Update()

4. Player.Update(inputMgr)
   ├─> Traite les mouvements
   ├─> Traite les actions (attaque, esquive)
   └─> Met à jour la stamina

5. Camera.Follow(player.X, player.Y)
   └─> Suit le joueur avec interpolation

6. Enemy.Update(playerX, playerY) [pour chaque ennemi]
   ├─> Machine à états (Idle/Chase/Attack/Stun)
   └─> IA de comportement

7. CombatSystem.Update()
   ├─> Vérifie attaques joueur → ennemis
   ├─> Vérifie attaques ennemis → joueur
   └─> Enregistre les stats (gameStats)
```

### Boucle Draw (Rendu)

```
1. Remplir l'écran avec couleur de fond

2. TileMap.Draw(screen, camera)
   └─> Dessine uniquement les tuiles visibles

3. Enemy.Draw(screen, camera) [pour chaque ennemi]
   ├─> Convertit coordonnées monde → écran
   ├─> Dessine le sprite/cercle
   └─> Dessine les barres de vie

4. Player.Draw(screen, camera)
   ├─> Dessine le joueur
   ├─> Affiche direction du regard
   ├─> Affiche zone d'attaque si en attaque
   └─> Dessine barres vie/stamina

5. Debug.Draw(screen)
   └─> Affiche FPS et informations de debug
```

## Systèmes de jeu

### Système de Combat

1. **Attaque du joueur**
   - Hitbox en arc devant le joueur
   - Rayon: largeur du joueur × 2
   - Dégâts: 15 + bonus critique
   - Effet: Stun + knockback sur ennemi

2. **Attaque des ennemis**
   - Cooldown entre attaques
   - Portée définie par type d'ennemi
   - Dégâts variables selon type
   - Esquive du joueur = ignore les dégâts

### Système de Stamina

- Régénération: 1 point/frame (auto)
- Consommation:
  - Sprint: 2 points/frame
  - Attaque: 20 points (instantané)
  - Esquive: 30 points (instantané)

### Système d'Invincibilité

Après dégâts:
- Durée: 60 frames (1 seconde)
- Effet visuel: Clignotement
- Empêche nouveaux dégâts

Pendant esquive:
- Durée: 20 frames
- Mouvement rapide
- Invincible aux attaques

## Performance

### Optimisations

1. **Rendu de la TileMap**
   - Culling: Seulement les tuiles visibles
   - Calcul: ~40-50 tuiles au lieu de 10000

2. **Détection de collisions**
   - Collisions circulaires (rapide)
   - Vérification de distance avant calcul précis

3. **Machine à états des ennemis**
   - Pas de pathfinding complexe
   - IA simple mais efficace

### Profilage

Pour profiler le jeu:
```bash
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

## Extension du jeu

### Ajouter une nouvelle entité

1. Créer une structure qui contient `BaseEntity`
2. Implémenter l'interface `Entity`
3. Ajouter au slice dans `Game`
4. Appeler `Update()` et `Draw()` dans la boucle

### Ajouter un nouveau système

1. Créer un nouveau package (ex: `inventory/`)
2. Créer la structure du système
3. Initialiser dans `game.NewGame()`
4. Appeler dans `game.Update()` ou `game.Draw()`

### Ajouter des sprites

1. Placer les images dans `assets/sprites/`
2. Charger avec `ebitenutil.NewImageFromFile()`
3. Stocker dans la structure de l'entité
4. Dessiner avec `screen.DrawImage()` au lieu de `vector.DrawFilledCircle()`

## Tests

Les tests sont organisés par package:
```bash
go test ./...                    # Tous les tests
go test ./combat/                # Tests du système de combat
go test -v ./...                 # Mode verbose
go test -cover ./...             # Couverture de code
```

## Conventions de code

- **Nommage**: PascalCase pour exports, camelCase pour privé
- **Commentaires**: Toutes les fonctions publiques documentées
- **Organisation**: Un système = un package
- **Types**: Utiliser des types personnalisés (EnemyType, EnemyState)
- **Constantes**: Définir les valeurs magiques en constantes

## Diagramme de dépendances

```
main
 └─> game
      ├─> entities
      ├─> combat
      │    └─> entities, stats
      ├─> world
      │    ├─> camera, tilemap
      │    ├─> zone, zone_manager
      │    └─> dungeon, dungeon_manager, portal
      ├─> ui
      │    ├─> main_menu, pause_menu
      │    ├─> campfire_menu, inventory_menu
      │    ├─> levelup_menu, shop_menu
      │    ├─> death_screen, hud, tutorial
      │    └─> options_menu, stats_display
      ├─> input
      ├─> stats
      ├─> settings
      └─> save
```

## Ressources

- [Ebiten Documentation](https://ebitengine.org/en/documents/)
- [Go by Example](https://gobyexample.com/)
- [Game Programming Patterns](https://gameprogrammingpatterns.com/)
