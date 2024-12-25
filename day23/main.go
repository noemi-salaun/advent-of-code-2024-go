package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
)

var theGraph map[string][]string
var theCouples map[couple]struct{}

type triplet struct {
	c1 string
	c2 string
	c3 string
}

func newTriplet(c1, c2, c3 string) triplet {
	list := []string{c1, c2, c3}
	sort.Strings(list)
	return triplet{list[0], list[1], list[2]}
}

func (t *triplet) hasChief() bool {
	return string(t.c1[0]) == "t" || string(t.c2[0]) == "t" || string(t.c3[0]) == "t"
}

type couple struct {
	c1 string
	c2 string
}

func newCouple(c1, c2 string) couple {
	list := []string{c1, c2}
	sort.Strings(list)
	return couple{list[0], list[1]}
}

func main() {
	in, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	theGraph = in.links
	theCouples = in.couples

	var pcs []string
	for pc := range theGraph {
		pcs = append(pcs, pc)
	}

	R := []string{}
	P := pcs
	X := []string{}
	result := make(map[string]struct{})

	bronKerbosch(R, P, X, result)

	var best string
	var bestSize int
	for r := range result {
		if len(r) > bestSize {
			best = r
			bestSize = len(r)
		}
	}

	fmt.Println(best)
}

// Fonction pour convertir une liste triée en une clé unique
func toKey(clique []string) string {
	sort.Strings(clique)             // Trie les éléments de la clique
	return strings.Join(clique, ",") // Crée une chaîne unique
}

// R: Une liste contenant les nœuds de la clique actuelle.
// P: Une liste des candidats potentiels à ajouter à la clique.
// X: Une liste des nœuds qui ne peuvent pas être dans la clique actuelle.
func bronKerbosch(R []string, P []string, X []string, resultSet map[string]struct{}) {
	if len(P) == 0 && len(X) == 0 {
		key := toKey(R) // Convertit la clique triée en une clé unique
		if _, ok := resultSet[key]; !ok {
			resultSet[key] = struct{}{} // Stocke la clique unique dans le set
		}
		return
	}
	// Choisir un pivot
	u := choosePivot(union(P, X))
	neighbors := make(map[string]struct{})
	for _, n := range theGraph[u] {
		neighbors[n] = struct{}{}
	}
	var PNu []string
	for _, p := range P {
		if _, ok := neighbors[p]; !ok {
			PNu = append(PNu, p)
		}
	}

	// Explorer les nœuds de P qui ne sont pas voisins du pivot
	for _, node := range PNu {
		bronKerbosch(
			append(R, node),
			intersect(P, theGraph[node]),
			intersect(X, theGraph[node]),
			resultSet,
		)
		P = remove(P, node)
		X = append(X, node)
	}
}

func remove(slice []string, element string) []string {
	var result []string
	for _, v := range slice {
		if v != element {
			result = append(result, v)
		}
	}
	return result
}

func choosePivot(P []string) string {
	return P[rand.Intn(len(P))]
}

func intersect(list1 []string, list2 []string) []string {
	set := make(map[string]struct{})
	for _, v := range list2 {
		set[v] = struct{}{}
	}

	var intersection []string
	for _, v := range list1 {
		if _, ok := set[v]; ok {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

func union(a, b []string) []string {
	set := make(map[string]bool)
	for _, v := range a {
		set[v] = true
	}
	for _, v := range b {
		set[v] = true
	}

	var uni []string
	for k := range set {
		uni = append(uni, k)
	}
	return uni
}

type input struct {
	couples map[couple]struct{}
	links   map[string][]string
}

func loadInput() (input, error) {
	var in input
	in.couples = make(map[couple]struct{})
	in.links = make(map[string][]string)

	readFile, err := os.Open("day23/input.txt")
	if err != nil {
		return in, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	for fileScanner.Scan() {
		var line = fileScanner.Text()

		comp := strings.Split(line, "-")
		in.couples[newCouple(comp[0], comp[1])] = struct{}{}

		in.links[comp[0]] = append(in.links[comp[0]], comp[1])
		in.links[comp[1]] = append(in.links[comp[1]], comp[0])
	}

	return in, nil
}
