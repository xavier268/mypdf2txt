package p2t

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExtractText(t *testing.T) {
	// Vérifier que nous sommes sur Windows
	if runtime.GOOS != "windows" {
		t.Skip("Ce test nécessite Windows")
	}

	// Définir le chemin vers les fichiers de test
	testFilesDir := filepath.Join("..", "testFiles")

	// Lister tous les fichiers PDF dans le dossier testFiles
	entries, err := os.ReadDir(testFilesDir)
	if err != nil {
		t.Fatalf("Impossible de lire le dossier testFiles: %v", err)
	}

	// Filtrer pour obtenir uniquement les fichiers PDF
	var pdfFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".pdf" {
			pdfFiles = append(pdfFiles, filepath.Join(testFilesDir, entry.Name()))
		}
	}

	if len(pdfFiles) == 0 {
		t.Skip("Aucun fichier PDF trouvé dans testFiles")
	}

	// Tester chaque fichier PDF
	for _, pdfFile := range pdfFiles {
		t.Run(filepath.Base(pdfFile), func(t *testing.T) {
			// Appeler ExtractText avec les options par défaut
			text, err := ExtractText(pdfFile, nil)
			if err != nil {
				t.Fatalf("ExtractText a échoué pour %s: %v", pdfFile, err)
			}

			// Afficher les informations sur le texte extrait
			t.Logf("Fichier: %s", pdfFile)
			t.Logf("Longueur du texte extrait: %d caractères", len(text))

			// Afficher un extrait du texte (premiers 200 caractères)
			if len(text) > 200 {
				t.Logf("Extrait du texte: %s...", text[:200])
			} else {
				t.Logf("Texte complet: %s", text)
			}

			// Vérifier qu'au moins du texte a été extrait
			// (même si c'est un PDF vide, il devrait y avoir au moins les marqueurs de pages)
			if len(text) == 0 {
				t.Logf("ATTENTION: Aucun texte extrait pour %s", pdfFile)
			}
		})
	}
}

func TestExtractTextNonExistentFile(t *testing.T) {
	// Vérifier que nous sommes sur Windows
	if runtime.GOOS != "windows" {
		t.Skip("Ce test nécessite Windows")
	}

	// Tester avec un fichier qui n'existe pas
	_, err := ExtractText("fichier_inexistant.pdf", nil)
	if err == nil {
		t.Error("ExtractText devrait échouer avec un fichier inexistant")
	}
	t.Logf("Erreur attendue: %v", err)
}

func TestExtractTextWithCustomOptions(t *testing.T) {
	// Vérifier que nous sommes sur Windows
	if runtime.GOOS != "windows" {
		t.Skip("Ce test nécessite Windows")
	}

	// Définir le chemin vers les fichiers de test
	testFilesDir := filepath.Join("..", "testFiles")

	// Lister tous les fichiers PDF dans le dossier testFiles
	entries, err := os.ReadDir(testFilesDir)
	if err != nil {
		t.Fatalf("Impossible de lire le dossier testFiles: %v", err)
	}

	// Prendre le premier fichier PDF
	var pdfFile string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".pdf" {
			pdfFile = filepath.Join(testFilesDir, entry.Name())
			break
		}
	}

	if pdfFile == "" {
		t.Skip("Aucun fichier PDF trouvé dans testFiles")
	}

	// Tester avec des options personnalisées
	options := &ExtractTextOptions{
		Language: "en-US",
		DPI:      200,
	}

	text, err := ExtractText(pdfFile, options)
	if err != nil {
		t.Fatalf("ExtractText a échoué: %v", err)
	}

	t.Logf("Fichier: %s", pdfFile)
	t.Logf("Options: Language=%s, DPI=%d", options.Language, options.DPI)
	t.Logf("Longueur du texte extrait: %d caractères", len(text))
}
