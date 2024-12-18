package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var gridScore [][]int
var grid [][]string
var end vector2

func main() {
	bytes, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	grid = *buildGrid(71, 71, &bytes, 3040)

	end = vector2{70, 70}

	for _, line := range grid {
		var l []int
		for range line {
			l = append(l, -1)
		}
		gridScore = append(gridScore, l)
	}

	fmt.Println(walk(vector2{0, 0}, 0))
}

func walk(position vector2, score int) int {

	if !saveScore(position, score) {
		return -1
	}

	if position == end {
		return score
	}

	nextSteps := findNextSteps(position, score)
	if len(nextSteps) == 0 {
		return -1
	}

	var bestScore = -1
	for _, nextStep := range nextSteps {
		s := walk(nextStep.position, nextStep.score)
		if s == -1 {
			continue
		}
		if bestScore == -1 || s < bestScore {
			bestScore = s
		}
	}

	return bestScore
}

func findNextSteps(position vector2, score int) []step {
	var nextScore = score + 1
	var steps []step

	for _, dir := range []vector2{up, right, down, left} {
		tar := position.target(dir)

		if tar.x < 0 || tar.y < 0 || tar.x > end.x || tar.y > end.y {
			continue
		}

		if grid[tar.y][tar.x] == "#" {
			continue
		}

		if (gridScore[tar.y][tar.x] != -1) && (nextScore >= gridScore[tar.y][tar.x]) {
			continue
		}

		steps = append(steps, step{position: tar, score: nextScore})
	}

	return steps
}

func saveScore(position vector2, score int) bool {
	minScore := gridScore[position.y][position.x]
	if score < minScore || minScore == -1 {
		gridScore[position.y][position.x] = score
		return true
	} else {
		return false
	}
}

func buildGrid(width int, height int, bytes *[]vector2, iterations int) *[][]string {
	var g [][]string
	for range height {
		var line = make([]string, width)
		for x := range width {
			line[x] = "."
		}
		g = append(g, line)
	}

	for i := range iterations {
		b := (*bytes)[i]
		g[b.y][b.x] = "#"
	}

	return &g
}

type vector2 struct {
	x, y int
}

func (v *vector2) target(m vector2) vector2 {
	return vector2{v.x + m.x, v.y + m.y}
}

var up = vector2{0, -1}
var right = vector2{1, 0}
var down = vector2{0, 1}
var left = vector2{-1, 0}

type step struct {
	position vector2
	score    int
}

func printGrid() {
	fmt.Print("\033[H\033[2J")
	for _, line := range grid {
		for _, char := range line {
			fmt.Print(char)
		}
		fmt.Print("\n")
	}
}

func loadInput() ([]vector2, error) {
	readFile, err := os.Open("day18/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var bytes []vector2

	for fileScanner.Scan() {
		var parts = strings.Split(fileScanner.Text(), ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])

		bytes = append(bytes, vector2{x, y})
	}

	return bytes, nil
}
