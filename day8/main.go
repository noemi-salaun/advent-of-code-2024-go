package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	var antinodes [][]bool
	for _, line := range input.grid {
		antinodes = append(antinodes, make([]bool, len(line)))
	}

	boundaries := vector2{len(input.grid[0]), len(input.grid)}

	antinodesCount := 0
	for y, line := range input.grid {
		for x, char := range line {
			if char != '.' {
				currentAntenna := vector2{x, y}
				nextAntennas := findNextAntennas(&input.grid, currentAntenna, char)
				for _, otherAntenna := range nextAntennas {
					for _, an := range getAntinodes(currentAntenna, otherAntenna, boundaries) {
						if antinodes[an.y][an.x] == false {
							antinodesCount++
							antinodes[an.y][an.x] = true
						}
					}
				}
			}
		}
	}

	fmt.Println(antinodesCount)
}

func findNextAntennas(grid *grid, startAt vector2, frequency rune) []vector2 {
	var antennas []vector2
	for y := startAt.y; y < len(*grid); y++ {
		line := (*grid)[y]
		for x, char := range line {
			if y == startAt.y && x <= startAt.x {
				continue
			}

			if char == frequency {
				antennas = append(antennas, vector2{x, y})
			}
		}
	}
	return antennas
}

func getAntinodes(t1 vector2, t2 vector2, boundaries vector2) []vector2 {

	var newAntinodes []vector2
	newAntinodes = append(newAntinodes, t1, t2)

	antinode1 := t1
	antinode2 := t2
	for {
		antinode1, antinode2 = antinode2, vector2{antinode2.x + (antinode2.x - antinode1.x), antinode2.y + (antinode2.y - antinode1.y)}

		if antinode2.y < 0 || antinode2.y >= boundaries.y || antinode2.x < 0 || antinode2.x >= boundaries.x {
			break
		}

		newAntinodes = append(newAntinodes, antinode2)
	}

	antinode1 = t1
	antinode2 = t2
	for {
		antinode1, antinode2 = vector2{antinode1.x + (antinode1.x - antinode2.x), antinode1.y + (antinode1.y - antinode2.y)}, antinode1

		if antinode1.y < 0 || antinode1.y >= boundaries.y || antinode1.x < 0 || antinode1.x >= boundaries.x {
			break
		}

		newAntinodes = append(newAntinodes, antinode1)
	}

	return newAntinodes
}

type vector2 struct {
	x int
	y int
}

type input struct {
	grid grid
}

type grid [][]rune

func loadInput() (*input, error) {

	readFile, err := os.Open("day8/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var input input

	for fileScanner.Scan() {
		line := fileScanner.Text()

		input.grid = append(input.grid, []rune(line))
	}

	return &input, nil
}
