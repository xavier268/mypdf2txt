package main

import (
	"fmt"
	"os"

	"github.com/xavier268/mypdf2txt/p2t"
)

func main() {
	// Vérifier les arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: example <fichier.pdf>")
		fmt.Println("Exemple: example document.pdf")
		os.Exit(1)
	}

	pdfFile := os.Args[1]

	// Vérifier que le fichier existe
	if _, err := os.Stat(pdfFile); os.IsNotExist(err) {
		fmt.Printf("Erreur: le fichier %s n'existe pas\n", pdfFile)
		os.Exit(1)
	}

	fmt.Printf("Extraction du texte de: %s\n\n", pdfFile)

	// Extraire le texte avec les options par défaut
	text, err := p2t.ExtractText(pdfFile, nil)
	if err != nil {
		fmt.Printf("Erreur lors de l'extraction: %v\n", err)
		os.Exit(1)
	}

	// Afficher le résultat
	fmt.Println("========== TEXTE EXTRAIT ==========")
	fmt.Println(text)
	fmt.Println("===================================")
	fmt.Printf("\nTotal: %d caractères extraits\n", len(text))
}
