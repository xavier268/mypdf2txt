# Script de build pour PdfTextExtractor
# Installe automatiquement .NET SDK si necessaire et compile l'executable

Write-Host "=== Build de PdfTextExtractor ===" -ForegroundColor Cyan

# Verifier si dotnet est installe
$dotnetInstalled = $false
try {
    $dotnetVersion = dotnet --version 2>$null
    if ($LASTEXITCODE -eq 0) {
        $dotnetInstalled = $true
        Write-Host "[OK] .NET SDK detecte: version $dotnetVersion" -ForegroundColor Green
    }
} catch {
    $dotnetInstalled = $false
}

# Installer .NET SDK si necessaire
if (-not $dotnetInstalled) {
    Write-Host "[!] .NET SDK n'est pas installe" -ForegroundColor Yellow
    Write-Host "Installation de .NET SDK 8.0 via winget..." -ForegroundColor Yellow

    # Verifier que winget est disponible
    try {
        $wingetVersion = winget --version 2>$null
        if ($LASTEXITCODE -ne 0) {
            throw "winget non disponible"
        }
    } catch {
        Write-Error "winget n'est pas disponible. Veuillez installer .NET SDK manuellement depuis: https://dotnet.microsoft.com/download"
        Write-Host "`nTelechargez et installez .NET SDK 8.0, puis relancez ce script." -ForegroundColor Yellow
        exit 1
    }

    # Installer .NET SDK via winget
    Write-Host "Installation en cours (cela peut prendre quelques minutes)..." -ForegroundColor Yellow
    winget install Microsoft.DotNet.SDK.8 --silent --accept-source-agreements --accept-package-agreements

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Echec de l'installation de .NET SDK"
        Write-Host "Veuillez installer manuellement depuis: https://dotnet.microsoft.com/download" -ForegroundColor Yellow
        exit 1
    }

    # Rafraichir les variables d'environnement
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

    Write-Host "[OK] .NET SDK installe avec succes" -ForegroundColor Green

    # Verifier l'installation
    $dotnetVersion = dotnet --version
    Write-Host "Version installee: $dotnetVersion" -ForegroundColor Green
}

# Aller dans le dossier du projet
$projectDir = Join-Path $PSScriptRoot "PdfTextExtractor"
Set-Location $projectDir

Write-Host "`nCompilation de PdfTextExtractor..." -ForegroundColor Cyan

# Compiler le projet en mode Release
dotnet publish -c Release -r win-x64 --self-contained -p:PublishSingleFile=true -p:PublishReadyToRun=true

if ($LASTEXITCODE -ne 0) {
    Write-Error "Echec de la compilation"
    exit 1
}

# Trouver le dossier de compilation (win-x64)
$buildDir = Get-ChildItem -Path "bin\Release" -Directory -Recurse | Where-Object { $_.Name -eq "win-x64" } | Select-Object -First 1

if ($null -eq $buildDir) {
    Write-Error "Impossible de trouver le dossier de compilation win-x64"
    exit 1
}

Write-Host "[OK] Compilation reussie" -ForegroundColor Green
Write-Host "Dossier de compilation: $($buildDir.FullName)" -ForegroundColor Green

# Creer le dossier bin a la racine du projet si necessaire
$rootDir = Split-Path -Parent $PSScriptRoot
$binDir = Join-Path $rootDir "bin"
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir | Out-Null
}

# Copier TOUS les fichiers du dossier de compilation dans bin (sauf le sous-dossier publish)
Get-ChildItem -Path $buildDir.FullName -File | ForEach-Object {
    Copy-Item -Path $_.FullName -Destination $binDir -Force
}

Write-Host "`n[OK] Fichiers copies dans: $binDir" -ForegroundColor Green
Write-Host "`n=== Build termine avec succes ===" -ForegroundColor Cyan
Write-Host "`nVous pouvez maintenant utiliser l'extracteur:" -ForegroundColor White
Write-Host "  .\bin\PdfTextExtractor.exe document.pdf" -ForegroundColor Gray
