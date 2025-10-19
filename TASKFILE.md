# Guide d'utilisation du Taskfile

Ce projet utilise [Task](https://taskfile.dev) pour automatiser toutes les opérations de build, test et développement.

## Installation de Task

```bash
# Installer Task (une seule fois)
go install github.com/go-task/task/v3/cmd/task@latest

# Ou exécuter (Task s'installera automatiquement)
task
```

## Commandes principales

### Build

```bash
# Compiler tout (C# + Go + exemple)
task build

# Compiler uniquement l'extracteur C# (installe .NET SDK si nécessaire)
task build-extractor

# Compiler uniquement le code Go
task build-go

# Compiler uniquement l'exemple CLI
task build-example

# Recompiler tout depuis zéro
task rebuild

# Forcer la recompilation de l'extracteur C#
task rebuild-extractor
```

### Test

```bash
# Lancer tous les tests (verbose)
task test

# Tests rapides (sans verbose)
task test-quick

# Tests avec rapport de couverture
task test-coverage
```

### Run

```bash
# Exécuter l'exemple avec un PDF
task run -- testFiles/test1.pdf

# Syntaxe générale
task run -- <chemin-vers-pdf>
```

### Clean

```bash
# Nettoyer les fichiers de build
task clean

# Nettoyage complet (build + fichiers temporaires)
task clean-all
```

### Développement

```bash
# Formater le code Go
task fmt

# Vérifier le code (go vet)
task vet

# Linter (nécessite golangci-lint)
task lint

# Nettoyer go.mod
task tidy

# Pipeline CI complète (format + vet + build + test)
task ci
```

### Documentation

```bash
# Lancer godoc et ouvrir le navigateur
task godoc

# Générer la documentation API en texte
task docs
```

### Information

```bash
# Afficher toutes les tâches disponibles
task --list-all

# Afficher les informations du projet
task info

# Afficher les versions
task version
```

### Installation complète

```bash
# Installer et builder tout en une commande
task install
```

## Workflow recommandé

### Première utilisation

```bash
# 1. Cloner le repo
git clone <repo-url>
cd mypdf2txt

# 2. Installer et builder
task install

# 3. Tester
task test

# 4. Essayer l'exemple
task run -- testFiles/test1.pdf
```

### Développement quotidien

```bash
# Après avoir modifié du code
task ci              # Vérifier que tout compile et les tests passent

# Avant un commit
task fmt            # Formater le code
task test           # Vérifier les tests
```

### Après avoir modifié le C#

```bash
# Recompiler l'extracteur
task rebuild-extractor

# Ou rebuild complet
task rebuild
```

## Tâches disponibles

| Tâche | Description |
|-------|-------------|
| `task` | Affiche le menu des tâches |
| `task build` | Compile tout |
| `task build-extractor` | Compile l'extracteur C# |
| `task build-go` | Compile la bibliothèque Go |
| `task build-example` | Compile l'exemple CLI |
| `task test` | Lance les tests (verbose) |
| `task test-quick` | Tests rapides |
| `task test-coverage` | Tests avec couverture |
| `task run -- <pdf>` | Exécute l'exemple |
| `task clean` | Nettoie les builds |
| `task clean-all` | Nettoyage complet |
| `task fmt` | Formate le code Go |
| `task vet` | Vérifie le code |
| `task lint` | Lance le linter |
| `task tidy` | Nettoie go.mod |
| `task godoc` | Lance godoc |
| `task docs` | Génère la doc API |
| `task version` | Affiche les versions |
| `task info` | Info sur le projet |
| `task ci` | Pipeline CI |
| `task install` | Installation complète |
| `task rebuild` | Rebuild complet |
| `task rebuild-extractor` | Force rebuild C# |

## Variables d'environnement

Le Taskfile définit automatiquement :

- `GOOS` - Système d'exploitation (windows/linux)
- `GOARCH` - Architecture (amd64/arm64)
- `VERSION` - Version du projet
- `EXE` - Extension exécutable (.exe sur Windows)
- `DATE` - Date de build

## Dépendances entre tâches

Le Taskfile gère automatiquement les dépendances :

```
build
├── build-extractor (compile C#)
├── build-go (dépend de build-extractor)
└── build-example (dépend de build-extractor)

test
└── build-extractor (s'assure que l'extracteur existe)

run
└── build-example (s'assure que l'exemple est compilé)
```

## Personnalisation

Pour modifier une tâche, éditez `Taskfile.yaml` :

```yaml
tasks:
  ma-tache:
    desc: Description de ma tâche
    cmds:
      - echo "Commande 1"
      - echo "Commande 2"
```

## Aide

```bash
# Aide générale de Task
task --help

# Liste de toutes les tâches avec description
task --list-all

# Voir le contenu d'une tâche
task --summary <nom-tache>
```

## Troubleshooting

### Task n'est pas reconnu

```bash
# Vérifier que Go bin est dans le PATH
go env GOPATH

# Ajouter au PATH si nécessaire (Windows PowerShell)
$env:Path += ";$(go env GOPATH)\bin"
```

### Les tâches ne fonctionnent pas

```bash
# Vérifier la syntaxe du Taskfile
task --dry

# Voir les commandes exactes qui seront exécutées
task --dry <nom-tache>
```

### .NET SDK ne s'installe pas

```bash
# Installer manuellement .NET SDK 8.0
# Puis relancer
task build-extractor
```

## Liens utiles

- [Task Documentation](https://taskfile.dev)
- [Task GitHub](https://github.com/go-task/task)
- [Task Installation](https://taskfile.dev/installation/)
