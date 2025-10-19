# Extract-PdfText.ps1
# Script PowerShell pour extraire le texte d'un PDF (y compris PDF scannés)
# en utilisant les API natives Windows 11
# Nécessite : Windows 10/11 (API UWP natives)

param(
    [Parameter(Mandatory=$true)]
    [string]$PdfPath,

    [Parameter(Mandatory=$false)]
    [string]$Language = "fr-FR",

    [Parameter(Mandatory=$false)]
    [int]$DPI = 300
)

# Vérifier que le fichier existe
if (-not (Test-Path $PdfPath)) {
    Write-Error "Le fichier PDF n'existe pas: $PdfPath"
    exit 1
}

# Charger les assemblies Windows Runtime nécessaires
Add-Type -AssemblyName System.Runtime.WindowsRuntime

# Charger les types WinRT nécessaires
$null = [Windows.Storage.StorageFile,Windows.Storage,ContentType=WindowsRuntime]
$null = [Windows.Storage.StorageFolder,Windows.Storage,ContentType=WindowsRuntime]
$null = [Windows.Data.Pdf.PdfDocument,Windows.Data.Pdf,ContentType=WindowsRuntime]
$null = [Windows.Media.Ocr.OcrEngine,Windows.Foundation,ContentType=WindowsRuntime]
$null = [Windows.Globalization.Language,Windows.Foundation,ContentType=WindowsRuntime]
$null = [Windows.Graphics.Imaging.BitmapDecoder,Windows.Graphics,ContentType=WindowsRuntime]
$null = [Windows.Graphics.Imaging.SoftwareBitmap,Windows.Graphics,ContentType=WindowsRuntime]

# Fonction pour attendre une opération async WinRT de manière simple
Function Await-Task {
    param([Parameter(Mandatory=$true)]$WinRtTask)

    # Obtenir le type de l'opération
    $taskType = $WinRtTask.GetType()

    # Obtenir le type générique - essayer différentes méthodes
    $resultType = $null
    if ($taskType.GenericTypeArguments -and $taskType.GenericTypeArguments.Count -gt 0) {
        $resultType = $taskType.GenericTypeArguments[0]
    } elseif ($taskType.GetGenericArguments().Count -gt 0) {
        $resultType = $taskType.GetGenericArguments()[0]
    }

    if ($null -eq $resultType) {
        throw "Impossible de déterminer le type de retour pour: $($taskType.FullName)"
    }

    # Trouver la méthode AsTask appropriée
    $asTaskMethod = ([System.WindowsRuntimeSystemExtensions].GetMethods() | Where-Object {
        $_.Name -eq 'AsTask' -and
        $_.GetParameters().Count -eq 1 -and
        $_.IsGenericMethod
    })[0]

    if ($null -eq $asTaskMethod) {
        throw "Méthode AsTask introuvable"
    }

    # Créer la méthode générique
    $genericAsTask = $asTaskMethod.MakeGenericMethod($resultType)

    # Convertir en Task et attendre
    $netTask = $genericAsTask.Invoke($null, @($WinRtTask))
    $netTask.GetAwaiter().GetResult()
}

try {
    # Convertir le chemin en chemin absolu
    $AbsolutePdfPath = (Resolve-Path $PdfPath).Path

    Write-Host "Traitement du fichier: $AbsolutePdfPath" -ForegroundColor Cyan

    # Ouvrir le fichier PDF
    $pdfFile = Await-Task -WinRtTask ([Windows.Storage.StorageFile]::GetFileFromPathAsync($AbsolutePdfPath))
    $pdfDocument = Await-Task -WinRtTask ([Windows.Data.Pdf.PdfDocument]::LoadFromFileAsync($pdfFile))

    $pageCount = $pdfDocument.PageCount
    Write-Host "Nombre de pages: $pageCount" -ForegroundColor Cyan

    # Initialiser l'engine OCR
    $ocrLanguage = [Windows.Globalization.Language]::new($Language)
    $ocrEngine = [Windows.Media.Ocr.OcrEngine]::TryCreateFromLanguage($ocrLanguage)

    if ($null -eq $ocrEngine) {
        # Essayer avec anglais par défaut si la langue demandée n'est pas disponible
        Write-Warning "Langue $Language non disponible, utilisation de l'anglais"
        $ocrLanguage = [Windows.Globalization.Language]::new("en-US")
        $ocrEngine = [Windows.Media.Ocr.OcrEngine]::TryCreateFromLanguage($ocrLanguage)
    }

    if ($null -eq $ocrEngine) {
        Write-Error "Impossible d'initialiser l'engine OCR"
        exit 1
    }

    Write-Host "Engine OCR initialisé: $($ocrEngine.RecognizerLanguage.DisplayName)" -ForegroundColor Green

    # Créer un dossier temporaire pour les images
    $tempFolder = Join-Path $env:TEMP "pdf2txt_$(Get-Random)"
    New-Item -ItemType Directory -Path $tempFolder -Force | Out-Null

    # StringBuilder pour accumuler tout le texte
    $allText = [System.Text.StringBuilder]::new()

    # Traiter chaque page
    for ($pageIndex = 0; $pageIndex -lt $pageCount; $pageIndex++) {
        $pageNumber = $pageIndex + 1
        Write-Host "Traitement de la page $pageNumber/$pageCount..." -ForegroundColor Yellow

        # Obtenir la page
        $pdfPage = $pdfDocument.GetPage($pageIndex)

        # Créer un fichier temporaire pour l'image
        $imagePath = Join-Path $tempFolder "page_$pageNumber.png"
        $tempStorageFolder = Await-Task -WinRtTask ([Windows.Storage.StorageFolder]::GetFolderFromPathAsync($tempFolder))
        $imageFile = Await-Task -WinRtTask ($tempStorageFolder.CreateFileAsync("page_$pageNumber.png", [Windows.Storage.CreationCollisionOption]::ReplaceExisting))

        # Rendre la page en image avec le DPI spécifié
        $renderOptions = [Windows.Data.Pdf.PdfPageRenderOptions]::new()
        $renderOptions.DestinationWidth = [uint32]($pdfPage.Size.Width * $DPI / 72)
        $renderOptions.DestinationHeight = [uint32]($pdfPage.Size.Height * $DPI / 72)

        $stream = Await-Task -WinRtTask ($imageFile.OpenAsync([Windows.Storage.FileAccessMode]::ReadWrite))
        Await-Task -WinRtTask ($pdfPage.RenderToStreamAsync($stream, $renderOptions)) | Out-Null
        $stream.Dispose()
        $pdfPage.Dispose()

        # Ouvrir l'image pour l'OCR
        $imageFile = Await-Task -WinRtTask ([Windows.Storage.StorageFile]::GetFileFromPathAsync($imagePath))
        $imageStream = Await-Task -WinRtTask ($imageFile.OpenAsync([Windows.Storage.FileAccessMode]::Read))

        # Créer un décodeur pour l'image
        $decoder = Await-Task -WinRtTask ([Windows.Graphics.Imaging.BitmapDecoder]::CreateAsync($imageStream))
        $softwareBitmap = Await-Task -WinRtTask ($decoder.GetSoftwareBitmapAsync())

        # Faire l'OCR
        $ocrResult = Await-Task -WinRtTask ($ocrEngine.RecognizeAsync($softwareBitmap))

        # Extraire le texte
        $pageText = $ocrResult.Text

        if ($pageText.Trim().Length -gt 0) {
            [void]$allText.AppendLine("=== Page $pageNumber ===")
            [void]$allText.AppendLine($pageText)
            [void]$allText.AppendLine("")
            Write-Host "  Texte extrait: $($pageText.Length) caractères" -ForegroundColor Green
        } else {
            Write-Host "  Aucun texte trouvé sur cette page" -ForegroundColor Gray
        }

        # Nettoyer
        $softwareBitmap.Dispose()
        $imageStream.Dispose()
    }

    # Nettoyer le dossier temporaire
    Remove-Item -Path $tempFolder -Recurse -Force -ErrorAction SilentlyContinue

    # Fermer le document PDF
    $pdfDocument.Dispose()

    # Afficher le résultat sur stdout
    Write-Host "`n========== TEXTE EXTRAIT ==========" -ForegroundColor Cyan
    Write-Output $allText.ToString()
    Write-Host "===================================`n" -ForegroundColor Cyan

    exit 0

} catch {
    Write-Error "Erreur lors du traitement: $_"
    Write-Error $_.Exception.Message
    Write-Error $_.ScriptStackTrace
    exit 1
}
