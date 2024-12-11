package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Cache global pour mémoriser les combinaisons déjà générées
var operatorCombinationCache = map[int][][]string{}

// Génère toutes les combinaisons possibles d'opérateurs '+' et '*' pour n-1 positions
func generateOperatorCombinations(n int) [][]string {
	// Vérifie si la combinaison existe déjà dans le cache
	if cached, exists := operatorCombinationCache[n]; exists {
		return cached
	}

	// Sinon, génère les combinaisons
	var combinations [][]string
	var generate func([]string, int)
	generate = func(current []string, depth int) {
		if depth == n {
			combination := make([]string, len(current))
			copy(combination, current)
			combinations = append(combinations, combination)
			return
		}
		for _, op := range []string{"+", "*"} {
			generate(append(current, op), depth+1)
		}
	}
	generate([]string{}, 0)

	// Stocke les combinaisons générées dans le cache
	operatorCombinationCache[n] = combinations
	return combinations
}

// Évalue une expression donnée sous forme de nombres et d'opérateurs
func evaluateExpression(numbers []int, operators []string) int {
	result := numbers[0]
	for i, op := range operators {
		if op == "+" {
			result += numbers[i+1]
		} else if op == "*" {
			result *= numbers[i+1]
		}
	}
	return result
}

// Vérifie si une équation est valide
func isEquationValid(numbers []int, targetValue int, operatorCombinations [][]string) bool {
	for _, operators := range operatorCombinations {
		if evaluateExpression(numbers, operators) == targetValue {
			return true
		}
	}
	return false
}

func chatGpt() {
	readFile, err := os.Open("day7/input.txt")
	if err != nil {
		return
	}
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)
	var wg sync.WaitGroup
	results := make(chan int)

	for scanner.Scan() {
		line := scanner.Text()

		// Crée une goroutine par équation
		wg.Add(1)
		go func(line string) {
			defer wg.Done()

			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				return
			}

			// Valeur cible
			targetValue, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return
			}

			// Nombres de l'équation
			numberStrings := strings.Fields(strings.TrimSpace(parts[1]))
			var numbers []int
			for _, numStr := range numberStrings {
				num, err := strconv.Atoi(numStr)
				if err != nil {
					return
				}
				numbers = append(numbers, num)
			}

			// Générer toutes les combinaisons d'opérateurs (avec mise en cache)
			operatorCombinations := generateOperatorCombinations(len(numbers) - 1)

			// Vérifier si une combinaison d'opérateurs donne la valeur cible
			if isEquationValid(numbers, targetValue, operatorCombinations) {
				results <- targetValue // Envoie le résultat valide dans le canal
			}
		}(line)
	}

	// Lancer une goroutine pour fermer le canal après que toutes les tâches soient terminées
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collecter les résultats
	totalCalibrationResult := 0
	for result := range results {
		totalCalibrationResult += result
	}

	fmt.Println(totalCalibrationResult)
}
