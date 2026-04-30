# Guide de dépannage

Ce document liste les problèmes courants et leurs solutions.

## Problèmes macOS

### ❌ Erreur: `[CAMetalLayer nextDrawable] returning nil because allocation failed`

**Symptômes:**
```
[CAMetalLayer nextDrawable] returning nil because allocation failed.
```

**Cause:**
Bug connu avec le backend Metal d'Ebiten sur certains Mac, particulièrement avec des résolutions élevées ou certaines versions de macOS.

**Solutions:**

#### Solution 1: Utiliser OpenGL (Recommandé) ✅
Le jeu utilise maintenant OpenGL par défaut. Relance simplement le jeu:

```bash
make run
# ou
./e-soulslike
```

#### Solution 2: Variable d'environnement manuelle
Si tu veux forcer OpenGL manuellement:

```bash
EBITEN_GRAPHICS_LIBRARY=opengl ./e-soulslike
```

#### Solution 3: Réduire la résolution
Modifie [main.go](main.go:11-12) et réduis la résolution:

```go
const (
    screenWidth  = 1024  // Au lieu de 1280
    screenHeight = 576   // Au lieu de 720
    title        = "E-SoulsLike - Dark Zelda"
)
```

### ❌ Le jeu lag ou les FPS sont bas

**Solutions:**

1. **Vérifier les processus en arrière-plan**
   ```bash
   # Voir l'utilisation CPU
   top -o cpu
   ```

2. **Réduire le nombre d'ennemis**
   Modifie [game/game.go](game/game.go:42-46):
   ```go
   g.enemies = []*entities.Enemy{
       entities.NewEnemy(400, 300, entities.EnemyTypeBasic),
       // Commenter les autres ennemis pour tester
   }
   ```

3. **Profiler le jeu**
   ```bash
   go build -o e-soulslike
   # Lancer avec profiling
   EBITEN_INTERNAL_PROFILE=1 ./e-soulslike
   ```

### ❌ Erreur: `permission denied`

**Cause:**
Le binaire n'a pas les permissions d'exécution.

**Solution:**
```bash
chmod +x e-soulslike
./e-soulslike
```

## Problèmes de compilation

### ❌ Erreur: `package github.com/hajimehoshi/ebiten/v2: cannot find package`

**Cause:**
Ebiten n'est pas installé.

**Solution:**
```bash
go mod download
go mod tidy
# ou
make deps
```

### ❌ Erreur: `go version too old`

**Cause:**
Ebiten v2.9+ requiert Go 1.24+.

**Solution:**
```bash
# Vérifier ta version
go version

# Mettre à jour Go si nécessaire
# Sur macOS avec Homebrew:
brew upgrade go

# Ou télécharger depuis https://go.dev/dl/
```

### ❌ Erreur de compilation C/CGO

**Cause:**
Problème avec les dépendances système d'Ebiten.

**Solution sur macOS:**
```bash
# Installer Xcode Command Line Tools
xcode-select --install

# Si tu utilises Homebrew:
brew install pkg-config
```

## Problèmes de jeu

### ❌ Le joueur ne bouge pas

**Vérifications:**

1. **Vérifier les touches du clavier**
   - ZQSD (AZERTY) ou WASD (QWERTY)
   - Flèches directionnelles

2. **La fenêtre a le focus**
   - Clique sur la fenêtre du jeu

3. **Vérifier les logs**
   - Cherche des erreurs dans le terminal

### ❌ Les attaques ne fonctionnent pas

**Vérifications:**

1. **Vérifier la stamina**
   - L'attaque coûte 20 stamina (barre verte)
   - Attends que la stamina se régénère

2. **Touche correcte**
   - Utilise la touche Espace
   - Pas Entrée ou autre

### ❌ Les ennemis ne réagissent pas

**Cause possible:**
Les ennemis sont trop loin.

**Solution:**
Approche-toi d'eux. Le rayon de détection est de 200 pixels.

## Problèmes de performance

### 🐌 FPS bas (< 60)

**Diagnostics:**

1. **Vérifier le nombre d'entités**
   ```go
   // Dans game.go, afficher le nombre d'ennemis
   fmt.Printf("Ennemis actifs: %d\n", len(g.enemies))
   ```

2. **Profiler avec pprof**
   ```bash
   go test -cpuprofile=cpu.prof -bench=.
   go tool pprof cpu.prof
   ```

3. **Optimisations possibles:**
   - Réduire la taille de la carte (TileMap)
   - Limiter le nombre d'ennemis
   - Désactiver le redimensionnement de fenêtre

### 💾 Utilisation mémoire élevée

**Solutions:**

1. **Vérifier les fuites mémoire**
   ```bash
   go test -memprofile=mem.prof -bench=.
   go tool pprof mem.prof
   ```

2. **Limiter la taille de la carte**
   Modifie [game/game.go](game/game.go:37):
   ```go
   g.worldMap = world.NewTileMap(50, 50)  // Au lieu de 100, 100
   ```

## Problèmes Linux

### ❌ Erreur: `X11` ou `libGL` manquant

**Solution sur Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install libgl1-mesa-dev xorg-dev
```

**Solution sur Fedora:**
```bash
sudo dnf install mesa-libGL-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel
```

## Problèmes Windows

### ❌ Fenêtre noire ou crash au démarrage

**Solutions:**

1. **Mettre à jour les drivers graphiques**
   - NVIDIA: https://www.nvidia.com/drivers
   - AMD: https://www.amd.com/support
   - Intel: https://www.intel.com/content/www/us/en/support/detect.html

2. **Essayer OpenGL**
   ```cmd
   set EBITEN_GRAPHICS_LIBRARY=opengl
   e-soulslike.exe
   ```

### ❌ Erreur: `VCRUNTIME140.dll` manquant

**Solution:**
Installer Visual C++ Redistributable:
https://aka.ms/vs/17/release/vc_redist.x64.exe

## Debug avancé

### Activer les logs de debug d'Ebiten

```bash
EBITEN_LOG_LEVEL=debug ./e-soulslike
```

### Afficher les FPS en jeu

Ajoute dans [game/game.go](game/game.go:94-97):

```go
import "github.com/hajimehoshi/ebiten/v2"

// Dans Draw()
fps := ebiten.ActualFPS()
ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", fps))
```

### Profiling détaillé

```bash
# CPU profiling
go build -o e-soulslike
CPUPROFILE=cpu.prof ./e-soulslike
go tool pprof cpu.prof

# Memory profiling
MEMPROFILE=mem.prof ./e-soulslike
go tool pprof mem.prof
```

## Obtenir de l'aide

Si ton problème n'est pas listé ici:

1. **Vérifier les issues Ebiten**
   - https://github.com/hajimehoshi/ebiten/issues

2. **Créer une issue avec:**
   - Ta version de Go (`go version`)
   - Ton OS et version
   - Le message d'erreur complet
   - Les étapes pour reproduire

3. **Logs utiles à inclure:**
   ```bash
   go version
   go env
   uname -a  # Linux/macOS
   systeminfo  # Windows
   ```

## Checklist de dépannage générale

- [ ] Go est à jour (>= 1.24)
- [ ] Les dépendances sont installées (`make deps`)
- [ ] Le code compile (`make build`)
- [ ] Les tests passent (`make test`)
- [ ] La fenêtre du jeu a le focus
- [ ] Pas d'autre programme intensif en CPU/GPU
- [ ] Drivers graphiques à jour

## Problèmes connus

### macOS Ventura/Sonoma + Metal
- **Problème**: Allocation failures avec Metal
- **Status**: Connu, utilise OpenGL
- **Workaround**: Activé par défaut dans le code

### Résolutions très élevées (> 4K)
- **Problème**: Peut causer des lags
- **Solution**: Réduire la résolution de la fenêtre

### Multiples moniteurs
- **Problème**: Fenêtre peut apparaître sur le mauvais écran
- **Solution**: Déplacer manuellement ou modifier les paramètres système

## Ressources

- [Documentation Ebiten](https://ebitengine.org/)
- [FAQ Ebiten](https://ebitengine.org/en/documents/faq.html)
- [Go Documentation](https://go.dev/doc/)
