package p2t

import (
	"fmt"
	"io"
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
		return "", fmt.Errorf("impossible de trouver l'exécutable PdfTextExtractor.exe: %w\n\nPour résoudre ce problème:\n1. Utilisez InstallExtractor() pour installer l'exécutable\n2. Ou définissez la variable d'environnement MYPDF2TXT_EXTRACTOR_PATH\n3. Ou consultez le README: https://github.com/xavier268/mypdf2txt#installation", err)
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

// GetUserSpaceExtractorPath retourne le chemin vers l'emplacement fixe de l'exécutable dans l'espace utilisateur
func GetUserSpaceExtractorPath() (string, error) {
	// Sur Windows, utiliser %LOCALAPPDATA%\mypdf2txt\bin\PdfTextExtractor.exe
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return "", fmt.Errorf("variable d'environnement LOCALAPPDATA non définie")
	}

	extractorDir := filepath.Join(localAppData, "mypdf2txt", "bin")
	extractorPath := filepath.Join(extractorDir, "PdfTextExtractor.exe")

	return extractorPath, nil
}

// findExtractorExecutable trouve le chemin vers l'exécutable PdfTextExtractor.exe
// Ordre de recherche:
// 1. Variable d'environnement MYPDF2TXT_EXTRACTOR_PATH
// 2. Emplacement fixe dans l'espace utilisateur (%LOCALAPPDATA%\mypdf2txt\bin\)
// 3. Chemins relatifs (pour le développement)
func findExtractorExecutable() (string, error) {
	// 1. Vérifier la variable d'environnement MYPDF2TXT_EXTRACTOR_PATH
	if envPath := os.Getenv("MYPDF2TXT_EXTRACTOR_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
	}

	// 2. Vérifier l'emplacement fixe dans l'espace utilisateur
	userSpacePath, err := GetUserSpaceExtractorPath()
	if err == nil {
		if _, err := os.Stat(userSpacePath); err == nil {
			return userSpacePath, nil
		}
	}

	// 3. Essayer plusieurs emplacements relatifs (pour le développement)
	possiblePaths := []string{
		"bin/PdfTextExtractor.exe",
		"../bin/PdfTextExtractor.exe",
		"../../bin/PdfTextExtractor.exe",
		"tools/PdfTextExtractor/bin/Release/net8.0-windows10.0.19041.0/win-x64/publish/PdfTextExtractor.exe",
		"tools/PdfTextExtractor/bin/Release/net8.0-windows10.0.19041.0/win-x64/PdfTextExtractor.exe",
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

	// Fournir un message d'erreur utile
	userSpacePathStr := "N/A"
	if userSpacePath, err := GetUserSpaceExtractorPath(); err == nil {
		userSpacePathStr = userSpacePath
	}

	return "", fmt.Errorf("PdfTextExtractor.exe introuvable. Emplacements vérifiés:\n"+
		"  1. Variable d'environnement MYPDF2TXT_EXTRACTOR_PATH: %s\n"+
		"  2. Espace utilisateur: %s\n"+
		"  3. Chemins relatifs: %v\n"+
		"Utilisez InstallExtractor() pour installer l'exécutable ou définissez MYPDF2TXT_EXTRACTOR_PATH",
		os.Getenv("MYPDF2TXT_EXTRACTOR_PATH"), userSpacePathStr, possiblePaths)
}

// InstallExtractor copie l'exécutable PdfTextExtractor.exe et ses dépendances
// depuis l'emplacement de build vers l'emplacement fixe dans l'espace utilisateur.
//
// Paramètres:
//   - sourcePath: chemin vers l'exécutable source (optionnel, auto-détecté si vide)
//
// Retourne:
//   - installedPath: le chemin où l'exécutable a été installé
//   - err: une erreur si l'installation échoue
func InstallExtractor(sourcePath string) (installedPath string, err error) {
	// Vérifier que nous sommes sur Windows
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("cette fonction n'est disponible que sur Windows")
	}

	// Si aucun chemin source n'est fourni, essayer de le trouver
	if sourcePath == "" {
		possibleSources := []string{
			"tools/PdfTextExtractor/bin/Release/net8.0-windows10.0.19041.0/win-x64/publish/PdfTextExtractor.exe",
			"tools/PdfTextExtractor/bin/Release/net8.0-windows10.0.19041.0/win-x64/PdfTextExtractor.exe",
			"bin/PdfTextExtractor.exe",
			"../bin/PdfTextExtractor.exe",
		}

		for _, relPath := range possibleSources {
			absPath, err := filepath.Abs(relPath)
			if err != nil {
				continue
			}
			if _, err := os.Stat(absPath); err == nil {
				sourcePath = absPath
				break
			}
		}

		if sourcePath == "" {
			return "", fmt.Errorf("impossible de trouver l'exécutable source. Spécifiez explicitement le chemin ou compilez d'abord avec 'powershell.exe -ExecutionPolicy Bypass -File tools/build.ps1'")
		}
	}

	// Vérifier que le fichier source existe
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return "", fmt.Errorf("le fichier source n'existe pas: %s", sourcePath)
	}

	// Obtenir le chemin de destination
	destPath, err := GetUserSpaceExtractorPath()
	if err != nil {
		return "", fmt.Errorf("impossible d'obtenir le chemin de destination: %w", err)
	}

	// Créer le répertoire de destination s'il n'existe pas
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("impossible de créer le répertoire de destination: %w", err)
	}

	// Copier l'exécutable
	if err := copyFile(sourcePath, destPath); err != nil {
		return "", fmt.Errorf("impossible de copier l'exécutable: %w", err)
	}

	// Copier les DLL nécessaires (si elles existent dans le même répertoire que l'exécutable source)
	sourceDir := filepath.Dir(sourcePath)
	dllNames := []string{
		"PdfTextExtractor.dll",
		"WinRT.Runtime.dll",
		"Microsoft.Windows.SDK.NET.dll",
	}

	for _, dllName := range dllNames {
		srcDll := filepath.Join(sourceDir, dllName)
		if _, err := os.Stat(srcDll); err == nil {
			destDll := filepath.Join(destDir, dllName)
			if err := copyFile(srcDll, destDll); err != nil {
				// Ne pas échouer si une DLL ne peut pas être copiée (elle n'est peut-être pas nécessaire)
				fmt.Fprintf(os.Stderr, "Avertissement: impossible de copier %s: %v\n", dllName, err)
			}
		}
	}

	return destPath, nil
}

// copyFile copie un fichier de src vers dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Copier les permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, sourceInfo.Mode())
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
