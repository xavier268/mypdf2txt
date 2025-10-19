# mypdf2txt

Bibliothèque Go pour extraire du texte de fichiers PDF, **y compris les PDF scannés** (avec OCR), en utilisant les API natives de Windows 11.

## Caractéristiques

- **100% natif Windows 11** : Utilise les API Windows intégrées (aucune installation externe requise)
- **Support OCR** : Extrait le texte même des PDF scannés (images)
- **Pas de dépendances Go externes** : Utilise uniquement la bibliothèque standard de Go
- **Multi-langues** : Support du français, anglais, et autres langues disponibles dans Windows
- **Simple à utiliser** : API Go simple et intuitive

## Prérequis

- **Windows 10/11** (les API UWP sont natives)
- **Go 1.20+**
- **.NET SDK 8.0** (installation automatisée)

## Installation rapide

```powershell
# 1. Compiler l'extracteur C# (installe .NET SDK si nécessaire)
powershell.exe -ExecutionPolicy Bypass -File tools/build.ps1

# 2. (Optionnel) Compiler l'exemple Go
go build -o example.exe ./cmd/example
```

**Note** : La première installation peut prendre 5-10 minutes (téléchargement et installation de .NET SDK ~209 MB).

Pour plus de détails, voir [INSTALL.md](INSTALL.md)

## Architecture

Le projet utilise une architecture hybride :

```
Go Application
    ↓
    Appelle PdfTextExtractor.exe (C#)
    ↓
Exécutable C#
    1. Convertit PDF → PNG (Windows.Data.Pdf API)
    2. Fait OCR sur PNG (Windows.Media.Ocr API)
    3. Retourne le texte extrait
    ↓
Go récupère le résultat
```

## Installation

```bash
go get github.com/xavier268/mypdf2txt
```

## Utilisation

### Exemple basique

```go
package main

import (
    "fmt"
    "github.com/xavier268/mypdf2txt/p2t"
)

func main() {
    // Extraire le texte avec les options par défaut
    text, err := p2t.ExtractText("document.pdf", nil)
    if err != nil {
        panic(err)
    }

    fmt.Println(text)
}
```

### Avec options personnalisées

```go
options := &p2t.ExtractTextOptions{
    Language: "en-US",  // Langue pour l'OCR
    DPI:      300,      // Résolution pour le rendu (plus haut = meilleur OCR mais plus lent)
}

text, err := p2t.ExtractText("document.pdf", options)
if err != nil {
    panic(err)
}

fmt.Println(text)
```

### Langues supportées

Les langues disponibles dépendent de votre installation Windows. Les plus courantes :

- `fr-FR` : Français (par défaut)
- `en-US` : Anglais
- `de-DE` : Allemand
- `es-ES` : Espagnol
- `it-IT` : Italien

## Exemple en ligne de commande

Un exemple complet est fourni dans `cmd/example/` :

```bash
# Construire l'exemple
go build -o example.exe ./cmd/example

# Utiliser
./example.exe document.pdf
```

## Tests

```bash
# Placer des fichiers PDF de test dans le dossier testFiles/
mkdir testFiles
# Copier vos PDF de test dans testFiles/

# Lancer les tests
go test ./p2t -v
```

## Structure du projet

```
mypdf2txt/
├── p2t/                    # Package principal
│   ├── extract.go          # Fonction d'extraction de texte
│   └── extract_test.go     # Tests
├── scripts/                # Scripts PowerShell
│   └── Extract-PdfText.ps1 # Script PowerShell natif Windows
├── cmd/
│   └── example/            # Exemple d'utilisation
│       └── main.go
├── testFiles/              # Fichiers PDF de test (non versionnés)
├── go.mod
└── README.md
```

## Comment ça marche ?

### 1. Script PowerShell (`scripts/Extract-PdfText.ps1`)

Le script utilise deux API Windows natives :

- **Windows.Data.Pdf.PdfDocument** : Convertit chaque page PDF en image
- **Windows.Media.Ocr.OcrEngine** : Fait l'OCR sur chaque image

### 2. Wrapper Go (`p2t/extract.go`)

Le code Go :
1. Trouve le script PowerShell
2. Appelle PowerShell avec les bons paramètres
3. Parse la sortie et retourne le texte extrait

## Avantages de cette approche

✅ **Pas d'installation externe** : Tesseract, Poppler, etc. ne sont pas nécessaires
✅ **Fonctionne sur tous les PC Windows** : Les API sont natives
✅ **Simple à maintenir** : Pas de CGO, pas de DLL natives complexes
✅ **Performant** : Les API Windows sont optimisées et maintenues par Microsoft

## Limitations

- **Windows uniquement** : Cette solution ne fonctionne que sur Windows 10/11
- **Langues OCR** : Dépend des packs de langues installés dans Windows
- **Vitesse** : L'OCR peut être lent pour les gros documents (plusieurs pages)

## Licence

MIT

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou une pull request.
