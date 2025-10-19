# Guide d'installation - mypdf2txt

Ce guide explique comment installer et utiliser mypdf2txt pour extraire du texte de fichiers PDF, y compris les PDF scannés (avec OCR).

## Prérequis

- **Windows 10/11** (obligatoire - utilise les API natives Windows)
- **Go 1.20+** (pour compiler le code Go)
- **.NET SDK 8.0** (sera installé automatiquement si absent)

## Installation rapide

### Étape 1 : Compiler l'extracteur PDF

L'extracteur C# utilise les API Windows natives et doit être compilé avant utilisation.

```powershell
# Depuis la racine du projet
powershell.exe -ExecutionPolicy Bypass -File tools/build.ps1
```

Ce script va :
1. Vérifier si .NET SDK est installé
2. L'installer automatiquement via winget si nécessaire (peut prendre 5-10 minutes)
3. Compiler le projet C# PdfTextExtractor
4. Copier l'exécutable dans `bin/PdfTextExtractor.exe`

**Note** : La première installation peut prendre du temps car .NET SDK fait ~209 MB.

### Étape 2 : Vérifier l'installation

```bash
# Vérifier que l'exécutable a été créé
ls bin/PdfTextExtractor.exe
```

### Étape 3 : Compiler le code Go (optionnel)

```bash
# Compiler l'exemple
go build -o example.exe ./cmd/example
```

## Utilisation

### Depuis Go

```go
package main

import (
    "fmt"
    "log"
    "github.com/xavier268/mypdf2txt/p2t"
)

func main() {
    // Extraction avec options par défaut (français, 300 DPI)
    text, err := p2t.ExtractText("document.pdf", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(text)

    // Extraction avec options personnalisées
    options := &p2t.ExtractTextOptions{
        Language: "en-US",  // Langue OCR
        DPI:      200,      // Résolution (plus bas = plus rapide)
    }
    text, err = p2t.ExtractText("document.pdf", options)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(text)
}
```

### Depuis la ligne de commande

```bash
# Utiliser l'exemple compilé
./example.exe document.pdf
```

### Directement avec l'exécutable C#

```powershell
# Syntaxe : PdfTextExtractor.exe <fichier.pdf> [langue] [dpi]
.\bin\PdfTextExtractor.exe document.pdf fr-FR 300
```

## Langues supportées

Les langues disponibles dépendent des packs de langues installés dans Windows. Les plus courantes :

- `fr-FR` : Français (par défaut)
- `en-US` : Anglais
- `de-DE` : Allemand
- `es-ES` : Espagnol
- `it-IT` : Italien
- `pt-PT` : Portugais
- `ru-RU` : Russe
- `zh-CN` : Chinois simplifié
- `ja-JP` : Japonais
- `ko-KR` : Coréen

## Dépannage

### Le build échoue

**Problème** : Le script de build échoue à installer .NET SDK

**Solution** : Installez .NET SDK 8.0 manuellement
1. Téléchargez depuis : https://dotnet.microsoft.com/download/dotnet/8.0
2. Installez le SDK (pas juste le runtime)
3. Relancez : `powershell.exe -ExecutionPolicy Bypass -File tools/build.ps1`

### L'exécutable n'est pas trouvé

**Problème** : `impossible de trouver l'exécutable PdfTextExtractor.exe`

**Solution** : Recompilez l'extracteur
```powershell
powershell.exe -ExecutionPolicy Bypass -File tools/build.ps1
```

### Langue OCR non disponible

**Problème** : `Langue XX-XX non disponible, utilisation de l'anglais`

**Solution** : Le pack de langue n'est pas installé dans Windows
1. Allez dans Paramètres > Heure et langue > Langue
2. Ajoutez la langue souhaitée
3. Téléchargez les fonctionnalités facultatives (dont OCR)

### Erreur de compilation C#

**Problème** : Erreur lors de la compilation du projet C#

**Solution** : Vérifiez la version de .NET SDK
```powershell
dotnet --version
# Devrait afficher 8.0.xxx
```

Si la version est incorrecte, réinstallez .NET SDK 8.0.

## Architecture technique

```
┌─────────────────┐
│  Application Go │
└────────┬────────┘
         │ appelle
         ▼
┌──────────────────────┐
│ PdfTextExtractor.exe │ (exécutable C#)
└──────────┬───────────┘
           │ utilise
           ▼
┌──────────────────────────┐
│  API Windows natives     │
│  - Windows.Data.Pdf      │ (conversion PDF → PNG)
│  - Windows.Media.Ocr     │ (OCR sur images)
└──────────────────────────┘
```

### Avantages de cette approche

✅ **Aucune dépendance externe** (Tesseract, Poppler, etc.)
✅ **100% natif Windows** - utilise les API système
✅ **Pas de CGO** - code Go pur
✅ **Performant** - API optimisées par Microsoft
✅ **Multi-langues** - support OCR intégré

### Fichiers importants

- `tools/PdfTextExtractor/Program.cs` - Code C# de l'extracteur
- `tools/build.ps1` - Script de compilation automatique
- `p2t/extract.go` - Wrapper Go
- `bin/PdfTextExtractor.exe` - Exécutable compilé (généré)

## Performance

- **PDF avec texte sélectionnable** : Non applicable (cette solution fait uniquement de l'OCR)
- **PDF scanné (OCR)** : ~2-5 secondes par page (dépend du DPI et de la complexité)
- **Résolution recommandée** : 300 DPI (bon compromis qualité/vitesse)

## Licence

MIT
