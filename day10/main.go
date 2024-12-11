package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var theGrid grid
var gridHeight int
var gridWidth int

func main() {
	grid, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	theGrid = grid
	gridHeight = len(theGrid)
	gridWidth = len(theGrid[0])

	sum := 0
	for y, line := range theGrid {
		for x, n := range line {
			if n == 0 {
				sum += discoverTrailheadRating(vector2{x, y})
			}
		}
	}

	fmt.Println(sum)
}

type grid [][]int

type vector2 struct {
	x int
	y int
}

func (v *vector2) isInsideGrid() bool {
	return v.x >= 0 && v.x < gridWidth && v.y >= 0 && v.y < gridHeight
}

func (v *vector2) getGridElevation() int {
	if !v.isInsideGrid() {
		return -1
	}

	return theGrid[v.y][v.x]
}

var up = vector2{0, -1}
var right = vector2{1, 0}
var down = vector2{0, 1}
var left = vector2{-1, 0}

func discoverTrailheadScore(start vector2) int {
	ends := advance(start)

	score := 0
	set := make(map[vector2]struct{})
	for _, e := range ends {
		_, has := set[e]
		if !has {
			set[e] = struct{}{}
			score++
		}
	}

	return score
}

func discoverTrailheadRating(start vector2) int {
	ends := advance(start)

	return len(ends)
}

func advance(step vector2) []vector2 {

	if step.getGridElevation() == 9 {
		return []vector2{step}
	}

	nextSteps := findNextSteps(step)

	var ends []vector2
	for _, s := range nextSteps {
		ends = append(ends, advance(s)...)
	}

	return ends
}

func findNextSteps(step vector2) []vector2 {
	currentElevation := step.getGridElevation()
	nextElevation := currentElevation + 1

	var nextSteps []vector2
	for _, dir := range []vector2{up, right, down, left} {
		nextStep := vector2{step.x + dir.x, step.y + dir.y}
		if nextStep.getGridElevation() == nextElevation {
			nextSteps = append(nextSteps, nextStep)
		}
	}

	return nextSteps
}

func loadInput() (grid, error) {
	readFile, err := os.Open("day10/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var input grid
	for fileScanner.Scan() {
		chars := strings.Split(fileScanner.Text(), "")

		line := make([]int, len(chars))
		for i, c := range chars {
			n, err := strconv.Atoi(c)
			if err != nil {
				return nil, err
			}
			line[i] = n
		}

		input = append(input, line)
	}
	return input, nil
}
