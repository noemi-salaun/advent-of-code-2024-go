package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	//part1()
	part2()
}

func part1() {
	input, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	countSafe := 0
	for _, r := range input {
		if r.isSafe() {
			countSafe++
		}
	}

	fmt.Println(countSafe)
}

func part2() {
	input, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	countSafe := 0
	for _, r := range input {
		if r.isAlmostSafe() {
			countSafe++
		}
	}

	fmt.Println(countSafe)
}

type input []report

type report []int

func (r report) isSafe() bool {
	if len(r) <= 1 {
		return true
	}

	increasing := r[1] > r[0]

	for i := 1; i < len(r); i++ {
		var distance int
		if increasing {
			distance = r[i] - r[i-1]
		} else {
			distance = r[i-1] - r[i]
		}

		if distance < 1 || distance > 3 {
			return false
		}
	}

	return true
}

func (r report) isAlmostSafe() bool {
	if r.isSafe() {
		return true
	}

	for i := 0; i < len(r); i++ {
		head := r[:i]
		tail := r[i+1:]
		subReport := slices.Concat(head, tail)
		if subReport.isSafe() {
			return true
		}
	}

	return false
}

func loadInput() (input, error) {
	var input input

	readFile, err := os.Open("day2/input.txt")
	if err != nil {
		return input, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	for fileScanner.Scan() {
		sLevels := strings.Fields(fileScanner.Text())

		iLevels := make([]int, len(sLevels))
		for i, s := range sLevels {
			val, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			iLevels[i] = val
		}

		input = append(input, iLevels)
	}

	return input, nil
}
