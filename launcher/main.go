package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	githubOwner = "YOUR_GITHUB_USERNAME" // À remplacer
	githubRepo  = "e-soulslike"        // À remplacer si le repo a un autre nom
	versionFile = "version.txt"
	apiTimeout  = 5 * time.Second
)

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

func main() {
	fmt.Println("=== E-SoulsLike Launcher ===")

	// Répertoire où se trouve le launcher (et le jeu)
	gameDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	gamePath := filepath.Join(gameDir, gameBinaryName())
	versionPath := filepath.Join(gameDir, versionFile)

	currentVersion := readVersion(versionPath)
	fmt.Printf("Version locale : %s\n", currentVersion)

	fmt.Print("Vérification des mises à jour... ")
	release, err := fetchLatestRelease()
	if err != nil {
		fmt.Printf("impossible (%v)\n", err)
		launchGame(gamePath)
		return
	}

	if release.TagName == currentVersion {
		fmt.Println("à jour !")
		launchGame(gamePath)
		return
	}

	fmt.Printf("\nNouvelle version disponible : %s\n", release.TagName)

	asset := findAsset(release.Assets)
	if asset == nil {
		fmt.Println("Aucun binaire compatible trouvé, lancement de la version actuelle...")
		launchGame(gamePath)
		return
	}

	fmt.Printf("Téléchargement de %s (%s)...\n", asset.Name, formatSize(asset.Size))
	if err := downloadAndReplace(asset, gamePath); err != nil {
		fmt.Printf("Erreur de téléchargement : %v\nLancement de la version actuelle...\n", err)
		launchGame(gamePath)
		return
	}

	writeVersion(versionPath, release.TagName)
	fmt.Printf("Mise à jour vers %s réussie !\n\n", release.TagName)
	launchGame(gamePath)
}

// fetchLatestRelease appelle l'API GitHub pour récupérer la dernière release.
func fetchLatestRelease() (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubOwner, githubRepo)

	client := &http.Client{Timeout: apiTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API GitHub a retourné HTTP %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

// findAsset cherche l'asset correspondant à l'OS et l'architecture courants.
// Convention de nommage attendue : e-soulslike_darwin_arm64, e-soulslike_windows_amd64.exe, etc.
func findAsset(assets []Asset) *Asset {
	goos := runtime.GOOS     // "darwin", "linux", "windows"
	goarch := runtime.GOARCH // "amd64", "arm64"

	// Cherche une correspondance exacte OS + arch
	for i, a := range assets {
		name := strings.ToLower(a.Name)
		if strings.Contains(name, goos) && strings.Contains(name, goarch) {
			return &assets[i]
		}
	}
	// Fallback : juste l'OS
	for i, a := range assets {
		name := strings.ToLower(a.Name)
		if strings.Contains(name, goos) {
			return &assets[i]
		}
	}
	return nil
}

// downloadAndReplace télécharge le binaire dans un .tmp puis remplace l'ancien atomiquement.
// En cas d'échec, l'ancien binaire est restauré.
func downloadAndReplace(asset *Asset, destPath string) error {
	resp, err := http.Get(asset.BrowserDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpPath := destPath + ".tmp"
	f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier temporaire : %w", err)
	}

	written, err := io.Copy(f, resp.Body)
	f.Close()
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("erreur pendant le téléchargement : %w", err)
	}
	if asset.Size > 0 && written != asset.Size {
		os.Remove(tmpPath)
		return fmt.Errorf("taille incorrecte : reçu %d octets, attendu %d", written, asset.Size)
	}

	// Backup de l'ancien binaire
	backupPath := destPath + ".old"
	_ = os.Rename(destPath, backupPath)

	// Remplacement atomique
	if err := os.Rename(tmpPath, destPath); err != nil {
		_ = os.Rename(backupPath, destPath) // rollback
		return fmt.Errorf("impossible de remplacer le binaire : %w", err)
	}

	os.Remove(backupPath)
	return nil
}

// launchGame lance le binaire du jeu et attend qu'il se termine.
func launchGame(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("ERREUR : binaire introuvable à %s\n", path)
		fmt.Println("Appuyez sur Entrée pour quitter...")
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Println("Lancement du jeu...")
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("Le jeu s'est terminé avec une erreur : %v\n", err)
	}
}

func gameBinaryName() string {
	if runtime.GOOS == "windows" {
		return "e-soulslike.exe"
	}
	return "e-soulslike"
}

func readVersion(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "v0.0.0"
	}
	return strings.TrimSpace(string(data))
}

func writeVersion(path, version string) {
	_ = os.WriteFile(path, []byte(version+"\n"), 0644)
}

func formatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}
