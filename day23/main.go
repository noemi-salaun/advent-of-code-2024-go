package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

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

	var triplets = make(map[triplet]struct{})

	for c1, c1Links := range in.links {
		if len(c1Links) < 2 {
			continue
		}

		for _, c2 := range c1Links {
			c2Links := in.links[c2]
			if len(c2Links) < 2 {
				continue
			}

			inter := intersect(c1Links, c2, in.couples)
			for _, c3 := range inter {
				triplets[newTriplet(c1, c2, c3)] = struct{}{}
			}
		}
	}

	count := 0
	for k := range triplets {
		if k.hasChief() {
			count++
		}
	}

	fmt.Println(count)
}

func intersect(c1Links []string, c2 string, coupleMap map[couple]struct{}) []string {
	var result []string

	for _, c1Link := range c1Links {
		cpl := newCouple(c1Link, c2)
		if _, ok := coupleMap[cpl]; ok {
			result = append(result, c1Link)
		}
	}

	return result
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
