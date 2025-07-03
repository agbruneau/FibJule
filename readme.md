# Calculateur de Fibonacci Concurrent en Go

Ce projet est un outil en ligne de commande écrit en Go pour calculer le n-ième nombre de Fibonacci en utilisant l'algorithme de Doublage Rapide (Fast Doubling).

L'application est conçue comme un outil pratique et une démonstration de concepts Go tels que la gestion de contexte, l'optimisation de la mémoire avec `sync.Pool`, et les tests de performance (benchmarking).

✨ Fonctionnalités

*   **Calcul de Très Grands Nombres**: Utilise le paquet `math/big` pour calculer des nombres de Fibonacci bien au-delà des limites des types entiers standards.
*   **Algorithme Performant**: Implémente l'algorithme de Doublage Rapide (Fast Doubling), connu pour son efficacité.
*   **Affichage de la Progression**: Montre en temps réel la progression du calcul sur une seule ligne qui se met à jour.
*   **Gestion du Délai d'Attente (Timeout)**: Utilise `context.WithTimeout` pour assurer que le programme se termine proprement si le calcul prend trop de temps.
*   **Optimisation de la Mémoire**: Emploie un `sync.Pool` pour recycler les objets `*big.Int`, réduisant la pression sur le Ramasse-Miettes (Garbage Collector).
*   **Suite de Tests Complète**: Inclut des tests unitaires pour valider la correction de l'algorithme et un benchmark pour mesurer ses performances.

🛠️ Prérequis

*   Go (version 1.18+ recommandée, car le projet utilise `go fmt ./...` et certaines pratiques plus récentes. La logique principale pourrait fonctionner sur Go 1.16+ mais 1.18+ est conseillé pour une compatibilité totale avec l'outillage de l'environnement de développement.)
*   Git (pour cloner le projet)

🚀 Installation

1.  Clonez le dépôt sur votre machine locale :
    ```sh
    git clone https://github.com/votre-nom-utilisateur/votre-depot.git
    ```
    (Remplacez l'URL par l'URL réelle de votre dépôt)

2.  Naviguez vers le répertoire du projet :
    ```sh
    cd votre-depot
    ```

💻 Utilisation

L'outil est exécuté directement depuis la ligne de commande.

**Exécution Simple**

Pour exécuter le calcul avec les valeurs par défaut (n=100000, timeout=1m, tous les algorithmes) :
```sh
go run .
```

**Options en Ligne de Commande**

Vous pouvez personnaliser l'exécution avec les options suivantes :

*   `-n <nombre>` : Spécifie l'index `n` du nombre de Fibonacci à calculer (entier non-négatif). Défaut : `100000`.
*   `-timeout <durée>` : Spécifie le délai d'attente global pour l'exécution (ex: `30s`, `2m`, `1h`). Défaut : `1m`.

**Exemples**

Calculer F(500 000) avec un délai d'attente de 30 secondes :
```sh
go run . -n 500000 -timeout 30s
```

Calculer F(1 000 000) avec un délai d'attente de 5 minutes :
```sh
go run . -n 1000000 -timeout 5m
```

**Exemple de Sortie**
```
2023/10/27 10:30:00 Calculating F(200000) using Fast Doubling with a timeout of 1m...
2023/10/27 10:30:00 Launching calculation...
Fast Doubling:   100.00%
2023/10/27 10:30:01 Calculation finished.

--------------------------- RESULT ---------------------------
Fast Doubling    : 8.8475ms     [OK              ] Result: 25974...03125
------------------------------------------------------------------------

📊 Algorithm: Fast Doubling (8.848ms)
Number of digits in F(200000): 41798
Value (scientific notation) ≈ 2.59740692e+41797
2023/10/27 10:30:01 Program finished.
```

🧠 Algorithme Implémenté

1.  **Doublage Rapide (Fast Doubling)**
    L'un des algorithmes connus les plus rapides pour les grands entiers. Il utilise les identités :
    *   `F(2k) = F(k) * [2*F(k+1) – F(k)]`
    *   `F(2k+1) = F(k)² + F(k+1)²`
    pour réduire significativement le nombre d'opérations. Complexité : O(log n) opérations arithmétiques.

🏗️ Architecture du Code

La base de code est organisée en plusieurs fichiers Go pour une meilleure modularité :

*   `main.go`: Contient la logique principale de l'application, y compris l'analyse des options en ligne de commande, l'orchestration de l'exécution de l'algorithme via une goroutine, et l'affichage final du résultat.
*   `algorithms.go`: Abrite l'implémentation de l'algorithme de calcul de Fibonacci `fibFastDoubling` et la définition du type `fibFunc`.
*   `utils.go`: Fournit des fonctions utilitaires partagées à travers l'application. Les composants clés sont le `progressPrinter` pour l'affichage en temps réel de la progression et l'assistant `newIntPool` pour la gestion du `sync.Pool` d'objets `*big.Int`.
*   `main_test.go`: Contient des tests unitaires pour vérifier la correction de l'algorithme `fibFastDoubling` et un benchmark pour mesurer ses caractéristiques de performance.

L'exécution est gérée à l'aide d'un `sync.WaitGroup` pour s'assurer que la goroutine de calcul se termine avant que le programme ne procède à l'affichage du résultat. Les mises à jour de progression sont envoyées via un canal partagé (`progressAggregatorCh`) à la goroutine `progressPrinter`, qui les affiche sur une seule ligne dans la console.

✅ Tests

Le projet inclut une suite complète de tests pour assurer la correction et mesurer les performances.

**Exécuter les Tests Unitaires**

Pour vérifier que l'algorithme implémenté produit des nombres de Fibonacci corrects pour un ensemble de valeurs connues (y compris les cas limites) :
```sh
go test -v ./...
```
Cette commande exécute tous les tests dans le paquet courant.

**Exécuter les Benchmarks**

Pour mesurer les performances (temps d'exécution et allocations mémoire) de l'algorithme :
```sh
go test -bench . ./...
```
Cette commande exécute tous les benchmarks dans le paquet courant. Le `.` indique tous les benchmarks.
Pour exécuter le benchmark spécifique de l'algorithme Fast Doubling :
```sh
go test -bench=BenchmarkFibFastDoubling ./...
```

📜 Licence

Ce projet est distribué sous la Licence MIT. (Typiquement, un fichier `LICENSE` serait inclus dans le dépôt avec le texte intégral de la Licence MIT.)
