package input

import (
	"testing"
)

func TestInputManager_JustPressed(t *testing.T) {
	im := NewInputManager()

	// Test: Aucune touche ne devrait être "just pressed" au départ
	// Note: Ce test est limité car on ne peut pas simuler ebiten.IsKeyPressed
	// sans mocker la bibliothèque. Dans un vrai test, on devrait injecter
	// les dépendances.

	// Vérifier que les maps sont vides au départ
	if len(im.justPressedKeys) != 0 {
		t.Error("justPressedKeys should be empty initially")
	}
}

func TestInputManager_Initialization(t *testing.T) {
	im := NewInputManager()

	if im == nil {
		t.Fatal("InputManager should not be nil")
	}

	if im.pressedKeys == nil {
		t.Error("pressedKeys map should be initialized")
	}

	if im.justPressedKeys == nil {
		t.Error("justPressedKeys map should be initialized")
	}

	if im.prevPressedKeys == nil {
		t.Error("prevPressedKeys map should be initialized")
	}
}

// Note: Pour tester complètement IsKeyJustPressed, il faudrait
// utiliser un système de mock pour ebiten.IsKeyPressed, ce qui
// nécessiterait une refactorisation du code pour l'injecter
// comme dépendance.
