package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

var gridScoreUp [][]int
var gridScoreRight [][]int
var gridScoreDown [][]int
var gridScoreLeft [][]int
var grid [][]string

func main() {
	in, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	grid = in.grid
	for _, line := range grid {
		gridScoreUp = append(gridScoreUp, make([]int, len(line)))
		gridScoreRight = append(gridScoreRight, make([]int, len(line)))
		gridScoreDown = append(gridScoreDown, make([]int, len(line)))
		gridScoreLeft = append(gridScoreLeft, make([]int, len(line)))
	}

	paths := walk(in.start, right, 0, []vector2{})

	var bestScore = -1
	for _, p := range paths {
		if p.score < bestScore || bestScore == -1 {
			bestScore = p.score
		}
	}

	var bestPaths []path
	for _, p := range paths {
		if p.score == bestScore {
			bestPaths = append(bestPaths, p)
		}
	}

	spots := drawPaths(bestPaths)
	fmt.Println(spots)
}

func walk(position vector2, direction vector2, score int, previousSteps []vector2) []path {
	mySteps := slices.Clone(previousSteps)
	mySteps = append(mySteps, position)

	if !saveScore(position, direction, score) {
		return []path{}
	}
	//drawGrid()

	if grid[position.y][position.x] == "E" {
		return []path{
			{score, mySteps},
		}
	}

	nextSteps := findNextSteps(position, direction, score)
	if len(nextSteps) == 0 {
		return []path{}
	}

	var paths []path

	for _, nextStep := range nextSteps {
		ps := walk(nextStep.position, nextStep.direction, nextStep.score, mySteps)
		paths = append(paths, ps...)
	}

	return paths
}

type path struct {
	score int
	steps []vector2
}

func saveScore(position vector2, direction vector2, score int) bool {
	switch direction {
	case up:
		minScore := gridScoreUp[position.y][position.x]
		if score <= minScore || minScore == 0 {
			gridScoreUp[position.y][position.x] = score
			return true
		} else {
			return false
		}

	case right:
		minScore := gridScoreRight[position.y][position.x]
		if score <= minScore || minScore == 0 {
			gridScoreRight[position.y][position.x] = score
			return true
		} else {
			return false
		}

	case down:
		minScore := gridScoreDown[position.y][position.x]
		if score <= minScore || minScore == 0 {
			gridScoreDown[position.y][position.x] = score
			return true
		} else {
			return false
		}

	case left:
		minScore := gridScoreLeft[position.y][position.x]
		if score <= minScore || minScore == 0 {
			gridScoreLeft[position.y][position.x] = score
			return true
		} else {
			return false
		}

	default:
		panic("invalid direction")
	}
}

func drawGrid() {
	fmt.Print("\033[H\033[2J")
	for y, line := range grid {
		for x, char := range line {
			if gridScoreUp[y][x] > 0 || gridScoreRight[y][x] > 0 || gridScoreDown[y][x] > 0 || gridScoreLeft[y][x] > 0 {
				fmt.Print("X")
			} else {
				fmt.Print(char)
			}
		}
		fmt.Print("\n")
	}
	time.Sleep(time.Millisecond * 500)
}

func drawPaths(paths []path) int {
	var spots = 0
	for _, p := range paths {
		for _, s := range p.steps {
			if grid[s.y][s.x] != "O" {
				grid[s.y][s.x] = "O"
				spots++
			}
		}
	}

	for _, line := range grid {
		for _, char := range line {
			fmt.Print(char)
		}
		fmt.Print("\n")
	}

	return spots
}

func drawPath(p path) {
	var g [][]string
	for _, line := range grid {
		g = append(g, slices.Clone(line))
	}

	for _, s := range p.steps {
		if g[s.y][s.x] != "O" {
			g[s.y][s.x] = "O"
		}
	}

	fmt.Println("\nscore: ", p.score)
	for _, line := range g {
		for _, char := range line {
			fmt.Print(char)
		}
		fmt.Print("\n")
	}
}

func findNextSteps(position vector2, direction vector2, score int) []step {
	var steps []step

	steps = append(steps, step{position: position.target(direction), direction: direction, score: score + 1})

	rgt := direction.turnRight()
	steps = append(steps, step{position: position.target(rgt), direction: rgt, score: score + 1001})

	lft := direction.turnLeft()
	steps = append(steps, step{position: position.target(lft), direction: lft, score: score + 1001})

	return slices.DeleteFunc(steps, func(s step) bool {
		char, minScore := s.getData()
		if char == "#" {
			return true
		}
		if s.score > minScore && minScore != 0 {
			return true
		}

		return false
	})
}

type vector2 struct {
	x, y int
}

func (v *vector2) target(m vector2) vector2 {
	return vector2{v.x + m.x, v.y + m.y}
}
func (v *vector2) turnRight() vector2 {
	switch *v {
	case up:
		return right
	case right:
		return down
	case down:
		return left
	case left:
		return up
	default:
		panic("invalid direction")
	}
}
func (v *vector2) turnLeft() vector2 {
	switch *v {
	case up:
		return left
	case left:
		return down
	case down:
		return right
	case right:
		return up
	default:
		panic("invalid direction")
	}
}

var up = vector2{0, -1}
var right = vector2{1, 0}
var down = vector2{0, 1}
var left = vector2{-1, 0}

type step struct {
	position  vector2
	direction vector2
	score     int
}

func (s *step) getData() (char string, minScore int) {
	char = grid[s.position.y][s.position.x]
	switch s.direction {
	case up:
		minScore = gridScoreUp[s.position.y][s.position.x]
	case left:
		minScore = gridScoreLeft[s.position.y][s.position.x]
	case down:
		minScore = gridScoreDown[s.position.y][s.position.x]
	case right:
		minScore = gridScoreRight[s.position.y][s.position.x]
	default:
		panic("invalid direction")
	}

	return
}

func (s *step) isEnd() bool {
	char := grid[s.position.y][s.position.x]
	return char == "E"
}

type input struct {
	grid  [][]string
	start vector2
}

func loadInput() (*input, error) {
	readFile, err := os.Open("day16/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var in input

	var y = 0
	for fileScanner.Scan() {
		var chars = strings.Split(fileScanner.Text(), "")
		var line []string
		for x, char := range chars {
			if char == "S" {
				in.start = vector2{x, y}
				line = append(line, ".")
			} else {
				line = append(line, char)
			}
		}

		in.grid = append(in.grid, line)
		y++
	}

	return &in, nil
}
