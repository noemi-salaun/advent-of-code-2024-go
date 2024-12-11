package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	theMap, guard, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	tracePatrol(theMap, *guard)
	count := 0
	for y, line := range *theMap {
		for x, char := range line {
			if string(char) == "X" {
				(*theMap)[y] = replaceAtIndex((*theMap)[y], '#', x)
				if !tryPatrol(theMap, *guard) {
					count++
				}
				(*theMap)[y] = replaceAtIndex((*theMap)[y], 'X', x)
			}
		}
	}
	fmt.Println(count)
}

func tracePatrol(theMap *[]string, guard guard) {
	(*theMap)[guard.position.y] = replaceAtIndex((*theMap)[guard.position.y], 'X', guard.position.x)

	for {
		moved, out := guard.move(theMap)
		if out {
			return
		}

		if moved {
			(*theMap)[guard.position.y] = replaceAtIndex((*theMap)[guard.position.y], 'X', guard.position.x)
		} else {
			guard.direction = guard.direction.turn()
		}
	}
}

func tryPatrol(theMap *[]string, guard guard) bool {
	var visited [][]byte
	for _, line := range *theMap {
		visited = append(visited, make([]byte, len(line)))
	}

	visited[guard.position.y][guard.position.x] |= guard.direction.getVisitByte()

	for {
		moved, out := guard.move(theMap)
		if out {
			return true
		}

		if moved {
			if visited[guard.position.y][guard.position.x]&guard.direction.getVisitByte() != 0 {
				return false
			}
			visited[guard.position.y][guard.position.x] |= guard.direction.getVisitByte()
		} else {
			guard.direction = guard.direction.turn()
		}
	}
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

type vector2 struct {
	x int
	y int
}

type guard struct {
	position  vector2
	direction vector2
}

func (g *guard) move(theMap *[]string) (bool, bool) {
	x := g.position.x + g.direction.x
	y := g.position.y + g.direction.y

	height := len(*theMap)
	width := len((*theMap)[0])

	if x < 0 || x >= width || y < 0 || y >= height {
		g.position.x = x
		g.position.y = y
		return true, true
	}

	if string((*theMap)[y][x]) == "#" {
		return false, false
	}

	g.position.x = x
	g.position.y = y
	return true, false
}

func newGuard(char string, x int, y int) guard {
	var direction vector2
	switch char {
	case "^":
		direction = up
	case ">":
		direction = right
	case "v":
		direction = down
	case "<":
		direction = left
	}
	return guard{
		position:  vector2{x, y},
		direction: direction,
	}
}

var up = vector2{0, -1}
var right = vector2{1, 0}
var down = vector2{0, 1}
var left = vector2{-1, 0}

var bUp byte = 1
var bRight byte = 2
var bDown byte = 4
var bLeft byte = 8

func (v vector2) turn() vector2 {
	switch v {
	case up:
		return right
	case right:
		return down
	case down:
		return left
	case left:
		return up
	default:
		panic("invalid vector")
	}
}

func (v vector2) getVisitByte() byte {
	switch v {
	case up:
		return bUp
	case right:
		return bRight
	case down:
		return bDown
	case left:
		return bLeft
	default:
		panic("invalid vector")
	}
}

func loadInput() (*[]string, *guard, error) {

	readFile, err := os.Open("day6/input.txt")
	if err != nil {
		return nil, nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var theMap []string
	var guard guard
	y := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()

		theMap = append(theMap, line)
		guardIndex := strings.IndexAny(line, "^><v")
		if guardIndex != -1 {
			guard = newGuard(string(line[guardIndex]), guardIndex, y)
		}
		y++
	}

	return &theMap, &guard, nil
}
