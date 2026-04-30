.PHONY: build run clean test

# Nom du binaire
BINARY_NAME=e-soulslike

# Compiler le jeu
build:
	@echo "Compilation du jeu..."
	go build -o $(BINARY_NAME) .
	@echo "Compilation terminée: $(BINARY_NAME)"

# Lancer le jeu (avec OpenGL, recommandé pour macOS)
run: build
	@echo "Lancement du jeu..."
	./$(BINARY_NAME)

# Lancer avec Metal (peut causer des erreurs d'allocation sur certains Mac)
run-metal: build
	@echo "Lancement du jeu avec Metal..."
	EBITEN_GRAPHICS_LIBRARY=metal ./$(BINARY_NAME)

# Lancer explicitement avec OpenGL
run-opengl: build
	@echo "Lancement du jeu avec OpenGL..."
	EBITEN_GRAPHICS_LIBRARY=opengl ./$(BINARY_NAME)

# Nettoyer les fichiers générés
clean:
	@echo "Nettoyage..."
	rm -f $(BINARY_NAME)
	go clean
	@echo "Nettoyage terminé"

# Tests
test:
	@echo "Exécution des tests..."
	go test ./...

# Vérifier le code
lint:
	@echo "Vérification du code..."
	go vet ./...
	go fmt ./...

# Télécharger les dépendances
deps:
	@echo "Téléchargement des dépendances..."
	go mod download
	go mod tidy

# Build pour différentes plateformes
build-all:
	@echo "Compilation multi-plateforme..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .
	@echo "Compilation multi-plateforme terminée"

# Afficher l'aide
help:
	@echo "Commandes disponibles:"
	@echo "  make build       - Compiler le jeu"
	@echo "  make run         - Compiler et lancer le jeu (OpenGL par défaut)"
	@echo "  make run-opengl  - Lancer avec OpenGL explicite"
	@echo "  make run-metal   - Lancer avec Metal (peut avoir des bugs sur macOS)"
	@echo "  make clean       - Nettoyer les fichiers générés"
	@echo "  make test        - Exécuter les tests"
	@echo "  make lint        - Vérifier et formater le code"
	@echo "  make deps        - Télécharger les dépendances"
	@echo "  make build-all   - Compiler pour toutes les plateformes"
	@echo "  make help        - Afficher cette aide"
