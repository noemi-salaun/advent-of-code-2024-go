package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var grid [][]string
var path []step
var stepMap = make(map[vector2]step)
var num = 0

var cheats []cheat
var cheatsTemplate []jump

func main() {
	in, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	grid = in.grid

	preparePath(in)
	prepareCheatsTemplate(20)

	part2()
}

func part2() {

	var totalCount = 0
	var cheatCountIndexedByGain = make(map[int]int)

	for _, s := range path {
		for _, j := range cheatsTemplate {
			tar := s.pos.target(j.vec)
			dest, ok := stepMap[tar]
			if !ok {
				continue
			}
			gain := dest.num - s.num - j.len
			if gain >= 100 {
				cheatCountIndexedByGain[gain]++
				totalCount++
			}
		}
	}

	fmt.Println(totalCount)
}

func part1() {
	for _, s := range path {
		for _, d := range []vector2{up, right, down, left} {
			wall := s.pos.target(d)
			if wall.getGridChar() != "#" {
				continue
			}

			tar := wall.target(d)
			dest, ok := stepMap[tar]
			if !ok || dest.num <= s.num {
				continue
			}

			cheats = append(cheats, cheat{wall, s, dest, dest.num - s.num - 2})
		}
	}

	var cheatMap = make(map[int]int)

	var count int
	for _, c := range cheats {
		cheatMap[c.gain]++

		if c.gain >= 100 {
			count++
		}
	}
	fmt.Println(count)
}

func preparePath(in input) {
	var s = newStep(in.start)

out:
	for {
		for _, d := range []vector2{up, right, down, left} {
			tar := s.pos.target(d)
			_, done := stepMap[tar]
			var gridChar = tar.getGridChar()
			if gridChar == "." && !done {
				s = newStep(tar)
				break
			}
			if gridChar == "E" {
				s = newStep(tar)
				break out
			}
		}
	}
}

// Generate all reachable positions starting from (0, 0) with up to n moves
func prepareCheatsTemplate(n int) {
	var x, y int
	for d := 1; d <= n; d++ {
		x = d
		y = 0
		for x > 0 {
			cheatsTemplate = append(cheatsTemplate, jump{vector2{x, y}, d})
			x--
			y++
		}

		x = 0
		y = d
		for y > 0 {
			cheatsTemplate = append(cheatsTemplate, jump{vector2{x, y}, d})
			x--
			y--
		}

		x = -d
		y = 0
		for x < 0 {
			cheatsTemplate = append(cheatsTemplate, jump{vector2{x, y}, d})
			x++
			y--
		}

		x = 0
		y = -d
		for y < 0 {
			cheatsTemplate = append(cheatsTemplate, jump{vector2{x, y}, d})
			x++
			y++
		}
	}
}

func newStep(pos vector2) step {
	var s = step{pos, num}
	num++
	path = append(path, s)
	stepMap[pos] = s

	return s
}

type input struct {
	grid  [][]string
	start vector2
}

type cheat struct {
	pos  vector2
	from step
	to   step
	gain int
}

type jump struct {
	vec vector2
	len int
}

type vector2 struct {
	x, y int
}

func (v *vector2) getGridChar() string {
	return grid[v.y][v.x]
}

type step struct {
	pos vector2
	num int
}

func (v *vector2) target(m vector2) vector2 {
	return vector2{v.x + m.x, v.y + m.y}
}

var up = vector2{0, -1}
var right = vector2{1, 0}
var down = vector2{0, 1}
var left = vector2{-1, 0}

func loadInput() (input, error) {
	var in input

	readFile, err := os.Open("day20/input.txt")
	if err != nil {
		return in, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var y = 0
	for fileScanner.Scan() {
		var line = fileScanner.Text()

		var l []string
		var chars = strings.Split(line, "")
		for x, char := range chars {
			if char == "S" {
				in.start = vector2{x, y}
			}
			l = append(l, char)
		}
		in.grid = append(in.grid, l)

		y++
	}

	return in, nil
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
