package input

import "github.com/hajimehoshi/ebiten/v2"

// keysToCheck liste les touches surveillées (déclarée une seule fois pour éviter les allocations)
// keysToCheck utilise les positions physiques QWERTY.
// Sur un clavier AZERTY : Z (affiché) = position W, Q (affiché) = position A.
var keysToCheck = []ebiten.Key{
	ebiten.KeyW, ebiten.KeyA, ebiten.KeyS, ebiten.KeyD, // ZQSD sur AZERTY
	ebiten.KeySpace, ebiten.KeyShift, ebiten.KeyControl,
	ebiten.KeyI, ebiten.KeyB,
	ebiten.KeyUp, ebiten.KeyDown, ebiten.KeyLeft, ebiten.KeyRight,
	ebiten.KeyEscape, ebiten.KeyEnter,
	ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4,
	ebiten.KeyE, ebiten.KeyR, ebiten.KeyF,
}

// InputManager gère les entrées clavier et souris
type InputManager struct {
	pressedKeys      map[ebiten.Key]bool
	justPressedKeys  map[ebiten.Key]bool
	prevPressedKeys  map[ebiten.Key]bool
	mouseX           int
	mouseY           int
	mousePressed     bool
	mouseJustPressed bool
	prevMousePressed bool
}

// NewInputManager crée un nouveau gestionnaire d'entrées
func NewInputManager() *InputManager {
	return &InputManager{
		pressedKeys:     make(map[ebiten.Key]bool),
		justPressedKeys: make(map[ebiten.Key]bool),
		prevPressedKeys: make(map[ebiten.Key]bool),
	}
}

// Update met à jour l'état des entrées (appelé chaque frame)
// Optimisation : les maps sont réutilisées pour éviter les allocations à 60 FPS
func (im *InputManager) Update() {
	// Copier pressedKeys -> prevPressedKeys sans ré-allouer
	for k := range im.prevPressedKeys {
		delete(im.prevPressedKeys, k)
	}
	for k, v := range im.pressedKeys {
		im.prevPressedKeys[k] = v
	}

	// Vider les maps courantes sans ré-allouer
	for k := range im.pressedKeys {
		delete(im.pressedKeys, k)
	}
	for k := range im.justPressedKeys {
		delete(im.justPressedKeys, k)
	}

	// Remplir l'état courant
	for _, key := range keysToCheck {
		if ebiten.IsKeyPressed(key) {
			im.pressedKeys[key] = true

			// Nouvelle pression si la touche n'était pas enfoncée la frame précédente
			if !im.prevPressedKeys[key] {
				im.justPressedKeys[key] = true
			}
		}
	}

	// Gérer la souris
	im.prevMousePressed = im.mousePressed
	im.mouseX, im.mouseY = ebiten.CursorPosition()
	im.mousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	im.mouseJustPressed = im.mousePressed && !im.prevMousePressed
}

// IsKeyPressed vérifie si une touche est enfoncée
func (im *InputManager) IsKeyPressed(key ebiten.Key) bool {
	return im.pressedKeys[key]
}

// IsKeyJustPressed vérifie si une touche vient d'être enfoncée (une seule frame)
func (im *InputManager) IsKeyJustPressed(key ebiten.Key) bool {
	return im.justPressedKeys[key]
}

// GetMousePosition retourne la position de la souris
func (im *InputManager) GetMousePosition() (int, int) {
	return im.mouseX, im.mouseY
}

// IsMousePressed vérifie si le bouton gauche de la souris est enfoncé
func (im *InputManager) IsMousePressed() bool {
	return im.mousePressed
}

// IsMouseJustPressed vérifie si le bouton gauche vient d'être enfoncé
func (im *InputManager) IsMouseJustPressed() bool {
	return im.mouseJustPressed
}

// GetMovementVector retourne un vecteur de mouvement normalisé basé sur WASD/ZQSD
func (im *InputManager) GetMovementVector() (float64, float64) {
	dx, dy := 0.0, 0.0

	if im.IsKeyPressed(ebiten.KeyW) || im.IsKeyPressed(ebiten.KeyUp) { // Z sur AZERTY
		dy = -1
	}
	if im.IsKeyPressed(ebiten.KeyS) || im.IsKeyPressed(ebiten.KeyDown) {
		dy = 1
	}
	if im.IsKeyPressed(ebiten.KeyA) || im.IsKeyPressed(ebiten.KeyLeft) { // Q sur AZERTY
		dx = -1
	}
	if im.IsKeyPressed(ebiten.KeyD) || im.IsKeyPressed(ebiten.KeyRight) {
		dx = 1
	}

	return dx, dy
}
