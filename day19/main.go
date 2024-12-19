package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var towels []string

var solved = make(map[string]int)

func main() {
	in, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	towels = in.towels

	var totalCount int
	for _, pattern := range in.patterns {
		count := solve(pattern)
		totalCount += count
		fmt.Println(pattern, count)
	}
	fmt.Println(totalCount)
}

func solve(p string) int {
	if p == "" {
		return 1
	}

	res, ok := solved[p]
	if ok {
		return res
	}

	var matches = 0
	for _, towel := range towels {
		after, found := strings.CutPrefix(p, towel)
		if !found {
			continue
		}

		matches += solve(after)
	}

	solved[p] = matches
	return matches
}

func part1(in input) {
	var count int
	var re = regexp.MustCompile(fmt.Sprintf(`^((%s))+$`, strings.Join(in.towels, ")|(")))
	for _, pattern := range in.patterns {
		if re.MatchString(pattern) {
			count++
		}
	}

	fmt.Println(count)
}

type input struct {
	towels   []string
	patterns []string
}

func loadInput() (input, error) {
	var in input

	readFile, err := os.Open("day19/input.txt")
	if err != nil {
		return in, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var readTowels = true

	for fileScanner.Scan() {
		var line = fileScanner.Text()

		if readTowels {
			in.towels = strings.Split(line, ", ")
			readTowels = false
			continue
		}

		if line != "" {
			in.patterns = append(in.patterns, line)
		}
	}

	return in, nil
}
