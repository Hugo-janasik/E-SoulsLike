#!/bin/bash

# Script de vérification de la configuration système pour E-SoulsLike

echo "=== Vérification de la configuration système ==="
echo ""

# Vérifier Go
echo "1. Version de Go:"
if command -v go &> /dev/null; then
    go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    REQUIRED_VERSION="1.24.0"

    if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" = "$REQUIRED_VERSION" ]; then
        echo "   ✅ Go version OK (>= 1.24.0)"
    else
        echo "   ⚠️  Go version trop ancienne. Requis: >= 1.24.0"
    fi
else
    echo "   ❌ Go n'est pas installé"
    exit 1
fi

echo ""

# Vérifier le système d'exploitation
echo "2. Système d'exploitation:"
OS=$(uname -s)
echo "   OS: $OS"

case "$OS" in
    Darwin)
        echo "   ✅ macOS détecté"
        echo "   Version: $(sw_vers -productVersion)"

        # Vérifier Xcode Command Line Tools
        if xcode-select -p &> /dev/null; then
            echo "   ✅ Xcode Command Line Tools installés"
        else
            echo "   ⚠️  Xcode Command Line Tools non trouvés"
            echo "      Installer avec: xcode-select --install"
        fi
        ;;
    Linux)
        echo "   ✅ Linux détecté"
        if [ -f /etc/os-release ]; then
            . /etc/os-release
            echo "   Distribution: $NAME $VERSION"
        fi
        ;;
    *)
        echo "   ℹ️  OS: $OS"
        ;;
esac

echo ""

# Vérifier les dépendances Go
echo "3. Dépendances Go:"
if [ -f "go.mod" ]; then
    echo "   ✅ go.mod trouvé"

    if grep -q "github.com/hajimehoshi/ebiten/v2" go.mod; then
        echo "   ✅ Ebiten présent dans go.mod"
        EBITEN_VERSION=$(grep "github.com/hajimehoshi/ebiten/v2" go.mod | awk '{print $2}')
        echo "      Version: $EBITEN_VERSION"
    else
        echo "   ⚠️  Ebiten non trouvé dans go.mod"
    fi
else
    echo "   ❌ go.mod non trouvé"
fi

echo ""

# Vérifier si le binaire existe
echo "4. Binaire du jeu:"
if [ -f "e-soulslike" ]; then
    echo "   ✅ Binaire compilé trouvé"
    ls -lh e-soulslike | awk '{print "      Taille:", $5}'

    if [ -x "e-soulslike" ]; then
        echo "   ✅ Permissions d'exécution OK"
    else
        echo "   ⚠️  Pas de permissions d'exécution"
        echo "      Corriger avec: chmod +x e-soulslike"
    fi
else
    echo "   ℹ️  Binaire non trouvé (normal si pas encore compilé)"
    echo "      Compiler avec: make build"
fi

echo ""

# Recommandations
echo "=== Recommandations ==="

if [ "$OS" = "Darwin" ]; then
    echo "Sur macOS, le jeu utilise OpenGL par défaut pour éviter les bugs Metal."
    echo "Si vous rencontrez des problèmes, consultez TROUBLESHOOTING.md"
fi

echo ""
echo "Pour compiler et lancer le jeu:"
echo "  make build"
echo "  make run"
echo ""
echo "Pour voir toutes les commandes disponibles:"
echo "  make help"
echo ""

# Tester la compilation
echo "=== Test de compilation ==="
read -p "Voulez-vous tester la compilation ? (o/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[OoYy]$ ]]; then
    echo "Compilation en cours..."
    if go build -o e-soulslike-test .; then
        echo "✅ Compilation réussie !"
        rm -f e-soulslike-test
    else
        echo "❌ Échec de la compilation"
        echo "Consultez les erreurs ci-dessus"
    fi
fi

echo ""
echo "=== Fin de la vérification ==="
