package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var grid [][]string
var bot vector2

func main() {
	in, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	grid = in.grid
	bot = in.bot

	for _, move := range in.moves {
		canMove := moveBot(move)
		drawGrid(move, canMove)
		time.Sleep(time.Millisecond * 200)
	}

	total := 0
	for y, line := range grid {
		for x, char := range line {
			if char == "[" {
				gps := 100*y + x

				fmt.Println(gps)

				total += gps
			}
		}
	}

	fmt.Println(total)
}

func drawGrid(lastMove vector2, canMove bool) {
	fmt.Print("\033[H\033[2J")
	for y, line := range grid {
		for x, char := range line {
			if bot.x == x && bot.y == y {
				if canMove {
					fmt.Printf("\033[1;33m%s\033[0m", lastMove.getBotChar())
				} else {
					fmt.Printf("\033[1;31m%s\033[0m", lastMove.getBotChar())
				}
			} else {
				if char == "." {
					char = " "
				}
				fmt.Print(char)
			}
		}
		fmt.Print("\n")
	}
}

func moveBot(move vector2) bool {
	tar := bot.target(move)

	if tar.isWall() {
		return false
	}

	if tar.isEmpty() {
		bot = tar
		return true
	}

	if tar.isBox() {
		osb := newSuperBox(tar)
		if osb.canMove(move) {
			osb.move(move)
			bot = tar
			return true
		} else {
			return false
		}
	}

	panic(fmt.Sprintf("Cannnot reach me. Target %d , %d = %s", tar.x, tar.y, grid[tar.y][tar.x]))
}

type superBox struct {
	l vector2
	r vector2
}

func newSuperBox(box vector2) superBox {
	if grid[box.y][box.x] == "[" {
		return superBox{box, box.target(right)}
	} else {
		return superBox{box.target(left), box}
	}
}

func (sb *superBox) targets(m vector2) []vector2 {
	if m == right {
		return []vector2{sb.r.target(right)}
	}
	if m == left {
		return []vector2{sb.l.target(left)}
	}
	if m == up {
		return []vector2{
			sb.l.target(up),
			sb.r.target(up),
		}
	}

	if m == down {
		return []vector2{
			sb.l.target(down),
			sb.r.target(down),
		}
	}

	panic(fmt.Sprintf("Invalid move %d, %d", m.x, m.y))
}

func (sb *superBox) canMove(move vector2) bool {
	tars := sb.targets(move)

	for _, tar := range tars {
		if tar.isWall() {
			return false
		}
	}

	var canMove = true
	for _, tar := range tars {
		if tar.isBox() {
			osb := newSuperBox(tar)
			if !osb.canMove(move) {
				canMove = false
				break
			}
		}
	}

	return canMove
}

func (sb *superBox) move(move vector2) {
	tars := sb.targets(move)

	var prevSb superBox

	for _, tar := range tars {
		if tar.isBox() {
			osb := newSuperBox(tar)
			if osb != prevSb {
				osb.move(move)

				prevSb = osb
			}
		}
	}

	tl := sb.l.target(move)
	tr := sb.r.target(move)

	grid[sb.l.y][sb.l.x] = "."
	grid[sb.r.y][sb.r.x] = "."

	grid[tl.y][tl.x] = "["
	grid[tr.y][tr.x] = "]"
}

type vector2 struct {
	x, y int
}

var up = vector2{0, -1}
var right = vector2{1, 0}
var down = vector2{0, 1}
var left = vector2{-1, 0}

func (v *vector2) target(m vector2) vector2 {
	return vector2{v.x + m.x, v.y + m.y}
}

func (v *vector2) isBox() bool {
	return grid[v.y][v.x] == "[" || grid[v.y][v.x] == "]"
}

func (v *vector2) isWall() bool {
	return grid[v.y][v.x] == "#"
}

func (v *vector2) isEmpty() bool {
	return grid[v.y][v.x] == "."
}

func (v *vector2) getBotChar() string {
	switch *v {
	case up:
		return "^"
	case right:
		return ">"
	case down:
		return "v"
	case left:
		return "<"
	default:
		panic("invalid bot direction")
	}
}

type input struct {
	grid  [][]string
	bot   vector2
	moves []vector2
}

func newVector2(char string) *vector2 {
	switch char {
	case "^":
		return &up
	case "<":
		return &left
	case ">":
		return &right
	case "v":
		return &down
	default:
		return nil
	}
}

func loadInput() (*input, error) {
	readFile, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var in input
	var buildGrid = true

	var y = 0
	for fileScanner.Scan() {
		if buildGrid {
			var chars = strings.Split(fileScanner.Text(), "")
			if len(chars) == 0 {
				buildGrid = false
				continue
			}
			var line []string
			for x, char := range chars {
				if char == "@" {
					in.bot = vector2{x * 2, y}
					line = append(line, ".", ".")
				} else {
					if char == "O" {
						line = append(line, "[", "]")
					} else {
						line = append(line, char, char)
					}
				}
			}

			in.grid = append(in.grid, line)
			y++
		} else {
			for _, char := range strings.Split(fileScanner.Text(), "") {

				var move = newVector2(char)
				if move != nil {
					in.moves = append(in.moves, *move)
				}
			}
		}
	}

	return &in, nil
}
