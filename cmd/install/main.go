package main

import (
	"fmt"
	"os"

	"github.com/xavier268/mypdf2txt/p2t"
)

func main() {
	fmt.Println("=== Installation de PdfTextExtractor ===")
	fmt.Println()

	// Installer l'extracteur
	installedPath, err := p2t.InstallExtractor("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de l'installation: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Installation réussie!\n")
	fmt.Printf("✓ Exécutable installé dans: %s\n", installedPath)
	fmt.Println()
	fmt.Println("Vous pouvez maintenant utiliser le package p2t depuis n'importe quelle application Go.")
}
