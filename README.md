# E-SoulsLike - Jeu 2D style Zelda/Dark Souls

Un jeu d'action-aventure 2D inspiré de Zelda et Dark Souls, développé en Go avec la bibliothèque Ebiten.

## Fonctionnalités actuelles

- **Système de joueur** avec déplacements fluides
- **Système de combat** avec attaques et esquives
- **Gestion de la stamina** pour les actions spéciales
- **Ennemis avec IA** (poursuite, attaque, vagabondage)
- **Caméra qui suit le joueur** avec interpolation smooth
- **Monde procédural** avec différents types de tuiles
- **Système de santé** avec invincibilité temporaire après dégâts
- **Collisions et knockback**
- **Menu principal** (nouvelle partie, charger, quitter)
- **Menu pause** avec statistiques, options, tutoriel et sauvegarde rapide
- **Menu campfire** avec niveau up et régénération
- **Inventaire** (objets, équipements)
- **Système de niveau et d'expérience** via le campfire
- **Système de sauvegarde** rapide et automatique
- **Plusieurs zones** et donjons avec transitions
- **Écran de mort** (YOU DIED) avec respawn au dernier feu de camp
- **HUD** (vie, stamina, notifications)

## Contrôles

### Déplacements
- **ZQSD** ou **WASD** : Se déplacer
- **Flèches directionnelles** : Se déplacer (alternative)

### Actions
- **Espace** : Attaquer (coûte 20 stamina)
- **Shift** : Sprint (consomme stamina)
- **Shift + Direction** : Esquive/Roll (coûte 30 stamina, donne invincibilité temporaire)
- **E** : Interagir (feux de camp, portails)
- **I** ou **B** : Ouvrir/Fermer l'inventaire
- **Échap** : Pause / Reprendre le jeu

## Comment jouer

### 1. Compiler le jeu

```bash
go build -o e-soulslike
```

### 2. Lancer le jeu

```bash
./e-soulslike
# ou
make run
```

### Problèmes courants

**Sur macOS:** Si tu vois des erreurs `[CAMetalLayer nextDrawable]`, c'est normal, le jeu utilise maintenant OpenGL par défaut pour éviter ce problème. Si le jeu ne se lance toujours pas, consulte [TROUBLESHOOTING.md](TROUBLESHOOTING.md).

```bash
# Forcer OpenGL si besoin
make run-opengl
```

Voir [TROUBLESHOOTING.md](TROUBLESHOOTING.md) pour plus de solutions.

## Structure du projet

```
e-soulslike/
├── main.go                 # Point d'entrée du jeu
├── game/                   # Package principal du jeu
│   └── game.go            # Boucle principale (Update/Draw)
├── entities/              # Entités du jeu
│   ├── entity.go          # Interface de base
│   ├── player.go          # Joueur
│   ├── player_stats.go    # Stats et niveau du joueur
│   ├── enemy.go           # Ennemis
│   ├── fireguardian.go    # Boss / Gardien de feu
│   ├── campfire.go        # Feux de camp
│   ├── inventory.go       # Inventaire
│   ├── item.go            # Items
│   ├── items_registry.go  # Registre des items
│   └── equipment.go       # Équipements
├── combat/                # Système de combat
│   └── combat.go          # Gestion des collisions et dégâts
├── world/                 # Gestion du monde
│   ├── camera.go          # Caméra 2D
│   ├── tilemap.go         # Carte du monde
│   ├── zone.go            # Zones
│   ├── zone_manager.go    # Gestionnaire de zones
│   ├── dungeon.go         # Donjons
│   ├── dungeon_generator.go
│   ├── dungeon_manager.go
│   └── portal.go          # Portails entre zones
├── ui/                    # Interface utilisateur
│   ├── button.go
│   ├── main_menu.go
│   ├── pause_menu.go
│   ├── campfire_menu.go
│   ├── inventory_menu.go
│   ├── levelup_menu.go
│   ├── shop_menu.go
│   ├── death_screen.go
│   ├── hud.go
│   ├── stats_display.go
│   ├── tutorial.go
│   └── options_menu.go
├── stats/                 # Statistiques de partie
│   └── stats.go
├── settings/              # Paramètres du jeu
│   └── settings.go
├── save/                  # Système de sauvegarde
│   └── save.go
├── input/                 # Gestion des entrées
│   └── input.go
├── launcher/              # Launcher avec auto-update
│   └── main.go
└── assets/                # Ressources (sprites, etc.)
    ├── monsters/
    ├── npcs/
    ├── object/
    ├── player/
    ├── tiled/
    └── weapons/
```

## Systèmes de jeu

### Système de combat
- Attaques avec zone d'effet devant le joueur
- Dégâts et knockback sur les ennemis
- Système de stun pour les ennemis touchés
- Invincibilité temporaire après avoir reçu des dégâts

### Système de stamina
- Sprint : consomme 2 stamina par frame
- Attaque : coûte 20 stamina
- Esquive : coûte 30 stamina
- Régénération automatique de 1 stamina par frame

### IA des ennemis
- **État Idle** : Vagabondage aléatoire
- **État Chasing** : Poursuite du joueur dans le rayon de détection
- **État Attacking** : Attaque quand à portée
- **État Stunned** : Étourdi après avoir reçu des dégâts

### Types d'ennemis
- **Basic** : Ennemi standard (50 HP, vitesse normale)
- **Fast** : Ennemi rapide (30 HP, vitesse élevée)
- **Tank** : Ennemi résistant (150 HP, vitesse lente)

## Prochaines fonctionnalités à ajouter

### Court terme
- [x] Système d'inventaire
- [x] Objets à ramasser (potions, équipements)
- [x] Système de niveau et d'expérience
- [ ] Plus de types d'ennemis
- [ ] Boss avec patterns d'attaque
- [ ] Sons et musiques
- [ ] Sprites et animations (partiellement)

### Moyen terme
- [x] Système de sauvegarde
- [x] Plusieurs zones/niveaux
- [ ] Dialogues et PNJ
- [ ] Quêtes
- [ ] Système de craft
- [ ] Talents et compétences

### Long terme
- [ ] Multijoueur coopératif
- [ ] Système de météo et jour/nuit
- [ ] Histoire et cinématiques

## Dépendances

- [Ebiten v2](https://github.com/hajimehoshi/ebiten) - Moteur de jeu 2D pour Go
- Go 1.24+

## Configuration recommandée

- Résolution : 1280x720 (redimensionnable)
- FPS : 60 (géré automatiquement par Ebiten)

## Développement

### Ajouter un nouvel ennemi

1. Ouvrir `entities/enemy.go`
2. Ajouter un nouveau type dans `EnemyType`
3. Configurer les stats dans `NewEnemy()`
4. Créer l'instance dans `game/game.go`

### Modifier les contrôles

1. Ouvrir `input/input.go`
2. Ajouter les touches dans `keysToCheck`
3. Utiliser dans `entities/player.go`

### Ajuster la difficulté

Fichiers à modifier :
- `entities/player.go` : Santé, stamina, vitesse du joueur
- `entities/enemy.go` : Stats des ennemis
- `combat/combat.go` : Dégâts, knockback

## Licence

Projet personnel - Libre d'utilisation et de modification

## Auteur

Hugo Janasik - Décembre 2024
