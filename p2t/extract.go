package p2t

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ExtractTextOptions définit les options pour l'extraction de texte
type ExtractTextOptions struct {
	// Language spécifie la langue pour l'OCR (ex: "fr-FR", "en-US")
	Language string
	// DPI spécifie la résolution pour le rendu des pages PDF en images (défaut: 300)
	DPI int
}

// DefaultExtractTextOptions retourne les options par défaut
func DefaultExtractTextOptions() ExtractTextOptions {
	return ExtractTextOptions{
		Language: "fr-FR",
		DPI:      300,
	}
}

// ExtractText extrait le texte d'un fichier PDF (y compris les PDF scannés)
// en utilisant l'exécutable C# natif Windows 11.
// Cette fonction utilise les API Windows natives (Windows.Data.Pdf + Windows.Media.Ocr)
// et ne nécessite aucune dépendance externe.
//
// Paramètres:
//   - filename: chemin vers le fichier PDF
//   - options: options d'extraction (peut être nil pour utiliser les valeurs par défaut)
//
// Retourne:
//   - text: le texte extrait de toutes les pages
//   - err: une erreur si l'extraction échoue
func ExtractText(filename string, options *ExtractTextOptions) (text string, err error) {
	// Vérifier que nous sommes sur Windows
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("cette fonction n'est disponible que sur Windows")
	}

	// Vérifier que le fichier existe
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", fmt.Errorf("le fichier n'existe pas: %s", filename)
	}

	// Utiliser les options par défaut si non spécifiées
	if options == nil {
		defaultOpts := DefaultExtractTextOptions()
		options = &defaultOpts
	}

	// Obtenir le chemin absolu du fichier PDF
	absFilename, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("impossible d'obtenir le chemin absolu: %w", err)
	}

	// Trouver l'exécutable PdfTextExtractor
	extractorPath, err := findExtractorExecutable()
	if err != nil {
		return "", fmt.Errorf("impossible de trouver l'exécutable PdfTextExtractor.exe: %w\nVeuillez exécuter 'powershell.exe -ExecutionPolicy Bypass -File tools/build.ps1' pour le compiler", err)
	}

	// Préparer la commande
	args := []string{
		absFilename,
		options.Language,
		fmt.Sprintf("%d", options.DPI),
	}

	// Exécuter l'extracteur
	cmd := exec.Command(extractorPath, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("erreur lors de l'exécution de PdfTextExtractor: %w\nSortie: %s", err, string(output))
	}

	// Extraire le texte de la sortie
	// Le texte est entre les marqueurs "========== TEXTE EXTRAIT =========="
	outputStr := string(output)
	return extractTextFromOutput(outputStr), nil
}

// findExtractorExecutable trouve le chemin vers l'exécutable PdfTextExtractor.exe
func findExtractorExecutable() (string, error) {
	// Essayer plusieurs emplacements possibles
	possiblePaths := []string{
		"bin/PdfTextExtractor.exe",
		"../bin/PdfTextExtractor.exe",
		"../../bin/PdfTextExtractor.exe",
		"tools/PdfTextExtractor/bin/Release/net8.0-windows10.0.19041.0/win-x64/publish/PdfTextExtractor.exe",
	}

	for _, relPath := range possiblePaths {
		absPath, err := filepath.Abs(relPath)
		if err != nil {
			continue
		}
		if _, err := os.Stat(absPath); err == nil {
			return absPath, nil
		}
	}

	return "", fmt.Errorf("PdfTextExtractor.exe introuvable dans les emplacements: %v", possiblePaths)
}

// extractTextFromOutput extrait le texte de la sortie du script PowerShell
func extractTextFromOutput(output string) string {
	// Chercher le début du texte extrait
	startMarker := "========== TEXTE EXTRAIT =========="
	endMarker := "==================================="

	startIdx := strings.Index(output, startMarker)
	if startIdx == -1 {
		// Si le marqueur n'est pas trouvé, retourner toute la sortie
		return strings.TrimSpace(output)
	}

	// Avancer après le marqueur de début
	startIdx += len(startMarker)

	// Chercher le marqueur de fin
	endIdx := strings.Index(output[startIdx:], endMarker)
	if endIdx == -1 {
		// Si le marqueur de fin n'est pas trouvé, prendre jusqu'à la fin
		return strings.TrimSpace(output[startIdx:])
	}

	// Extraire et nettoyer le texte
	text := output[startIdx : startIdx+endIdx]
	return strings.TrimSpace(text)
}
