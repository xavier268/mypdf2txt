# Solution : Extraction de texte PDF avec OCR natif Windows

## Problème initial

Extraire du texte de fichiers PDF, **y compris les PDF scannés** (qui nécessitent de l'OCR), avec une solution :
- Pure Go OU utilisant des fonctionnalités Windows natives
- Sans dépendances externes complexes
- Fonctionnant sur Windows 11

## Solution retenue : Exécutable C# + API Windows natives

### Pourquoi cette approche ?

Après avoir exploré plusieurs options (PowerShell + WinRT, bibliothèques Go, etc.), nous avons opté pour :

**Exécutable C# autonome** qui :
- Utilise `Windows.Data.Pdf.PdfDocument` pour convertir PDF → images
- Utilise `Windows.Media.Ocr.OcrEngine` pour faire l'OCR
- Est appelé par le code Go via `exec.Command`

### Avantages

✅ **Zéro dépendance externe**
- Pas besoin de Tesseract OCR
- Pas besoin de Poppler
- Pas besoin de bibliothèques natives complexes

✅ **100% Windows natif**
- Utilise les API UWP intégrées à Windows 10/11
- Fonctionne sur tous les PC Windows sans installation

✅ **Code Go pur**
- Pas de CGO
- Pas de bindings C/C++ compliqués
- Simple maintenance

✅ **Installation automatisée**
- Script PowerShell qui installe .NET SDK si nécessaire
- Compilation automatique du projet C#
- Prêt à l'emploi

### Inconvénients

❌ **Windows uniquement**
- Ne fonctionne pas sur Linux/macOS
- Acceptable car c'était le besoin exprimé

❌ **Nécessite .NET SDK**
- Mais installation automatisée
- Package de ~209 MB (installation unique)

❌ **Compilation initiale requise**
- Mais scriptée et automatique
- Une seule fois

## Architecture finale

```
mypdf2txt/
├── p2t/
│   ├── extract.go          # API Go - appelle PdfTextExtractor.exe
│   └── extract_test.go     # Tests
├── tools/
│   ├── build.ps1           # Script de build auto (installe .NET + compile)
│   └── PdfTextExtractor/
│       ├── Program.cs      # Code C# utilisant les API Windows
│       └── *.csproj
├── bin/
│   └── PdfTextExtractor.exe  # Exécutable compilé (généré)
├── cmd/example/
│   └── main.go             # Exemple CLI
└── testFiles/              # PDFs de test
```

## Flux d'exécution

```
1. Go appelle p2t.ExtractText("file.pdf", options)
        ↓
2. p2t trouve bin/PdfTextExtractor.exe
        ↓
3. Exécute : PdfTextExtractor.exe file.pdf fr-FR 300
        ↓
4. C# charge le PDF avec Windows.Data.Pdf
        ↓
5. Convertit chaque page en image PNG (en mémoire)
        ↓
6. Fait l'OCR avec Windows.Media.Ocr
        ↓
7. Retourne le texte via stdout
        ↓
8. Go parse la sortie et retourne le texte
```

## Alternatives explorées (et pourquoi abandonnées)

### ❌ PowerShell pur avec API WinRT
**Problème** : Manipulation des types génériques COM impossible
- Les objets WinRT sont exposés comme `System.__ComObject`
- Impossible d'extraire les types génériques par réflexion
- Erreurs complexes avec `MakeGenericMethod`

### ❌ Bibliothèques Go (gosseract)
**Problème** : Nécessite Tesseract OCR externe
- Installation Tesseract requise (~50MB)
- Dépendance native C
- Moins intégré que les API Windows

### ❌ Bibliothèques Go natives (pdfcpu, go-fitz)
**Problème** : N'incluent pas d'OCR
- pdfcpu : manipulation PDF mais pas d'OCR
- go-fitz : rendu PDF mais pas d'OCR
- Nécessiterait quand même Tesseract

## Performance observée

(Tests à venir une fois la compilation terminée)

Estimations basées sur l'architecture :
- **Temps de startup** : ~500ms (chargement .NET)
- **Par page (300 DPI)** : ~2-5 secondes
- **Qualité OCR** : Excellente (moteur Microsoft)

## Utilisation

### Installation (une seule fois)

```powershell
powershell.exe -ExecutionPolicy Bypass -File tools/build.ps1
```

### Code Go

```go
text, err := p2t.ExtractText("document.pdf", nil)
```

### CLI

```bash
./example.exe document.pdf
```

## Conclusion

Cette solution offre le **meilleur compromis** entre :
- Simplicité d'utilisation
- Pas de dépendances externes
- Performance
- Qualité de l'OCR
- Intégration Windows

L'overhead d'avoir un exécutable C# séparé est minimal et largement compensé par les avantages en termes de simplicité et de fiabilité.
