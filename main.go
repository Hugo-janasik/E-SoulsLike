package main

import (
	"log"
	"os"

	"e-soulslike/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	title        = "E-SoulsLike - Dark Zelda"
)

func main() {
	// Fix pour macOS: Forcer OpenGL au lieu de Metal si pas déjà défini
	if os.Getenv("EBITEN_GRAPHICS_LIBRARY") == "" {
		os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame(screenWidth, screenHeight)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
