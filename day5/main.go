package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	input, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	goodMiddleSum := 0
	badMiddleSum := 0
	for _, update := range input.updates {
		if input.checkRules(update) {
			goodMiddleSum += update.middle
		} else {
			badMiddleSum += update.getSortedMiddle(input)
		}
	}

	fmt.Printf("Part 1 : %d\n", goodMiddleSum)
	fmt.Printf("Part 2 : %d", badMiddleSum)
}

type rule struct {
	first int
	last  int
}

type update struct {
	pages  []int
	index  map[int]int
	middle int
}

func (u *update) getSortedMiddle(i *input) int {
	pages := make([]int, len(u.pages))
	copy(pages, u.pages)
	slices.SortFunc(pages, func(a, b int) int {
		for _, rule := range i.rules {
			if a == rule.first && b == rule.last {
				return -1
			}
			if b == rule.first && a == rule.last {
				return 1
			}
		}
		return 0
	})

	middleIndex := (len(pages) - 1) / 2
	return pages[middleIndex]
}

type input struct {
	rules   []rule
	updates []update
}

func (input *input) checkRules(update update) bool {
	for _, rule := range input.rules {
		firstPos, ok := update.index[rule.first]
		if !ok {
			continue
		}

		lastPos, ok := update.index[rule.last]
		if !ok {
			continue
		}

		if firstPos > lastPos {
			return false
		}
	}

	return true
}

func loadInput() (*input, error) {

	readFile, err := os.Open("day5/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var rules []rule
	var updates []update

	rulesPart := true
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if rulesPart && line == "" {
			rulesPart = false
			continue
		}

		if rulesPart {
			parts := strings.Split(line, "|")
			first, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, err
			}
			last, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}

			rules = append(rules, rule{first, last})
		} else {
			parts := strings.Split(line, ",")
			length := len(parts)
			if length%2 == 0 {
				return nil, errors.New("update with even number of pages")
			}
			middleIndex := (length - 1) / 2
			var update update
			update.index = make(map[int]int)
			for i, sVal := range parts {
				iVal, err := strconv.Atoi(sVal)
				if err != nil {
					return nil, err
				}
				update.pages = append(update.pages, iVal)
				update.index[iVal] = i
				if i == middleIndex {
					update.middle = iVal
				}
			}

			updates = append(updates, update)
		}

	}

	return &input{rules, updates}, nil
}
