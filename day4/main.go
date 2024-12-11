package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	rows, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	height := len(rows)
	width := len(rows[0])

	count := 0
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			if checkX_MAS(rows, x, y) {
				count++
			}
		}
	}
	fmt.Println(count)
}

func checkX_MAS(rows rows, x int, y int) bool {
	return rows[y][x] == 'A' &&
		((rows[y-1][x-1] == 'M' && rows[y+1][x+1] == 'S') || (rows[y-1][x-1] == 'S' && rows[y+1][x+1] == 'M')) &&
		((rows[y+1][x-1] == 'M' && rows[y-1][x+1] == 'S') || (rows[y+1][x-1] == 'S' && rows[y-1][x+1] == 'M'))
}

func countXMAS(rows rows, x int, y int) int {
	height := len(rows)
	width := len(rows[0])
	matches := 0

	// North
	if y >= 3 && rows[y][x] == 'X' && rows[y-1][x] == 'M' && rows[y-2][x] == 'A' && rows[y-3][x] == 'S' {
		matches++
	}
	// NorthEast
	if y >= 3 && x <= (width-4) && rows[y][x] == 'X' && rows[y-1][x+1] == 'M' && rows[y-2][x+2] == 'A' && rows[y-3][x+3] == 'S' {
		matches++
	}
	// East
	if x <= (width-4) && rows[y][x] == 'X' && rows[y][x+1] == 'M' && rows[y][x+2] == 'A' && rows[y][x+3] == 'S' {
		matches++
	}
	// SouthEast
	if y <= (height-4) && x <= (width-4) && rows[y][x] == 'X' && rows[y+1][x+1] == 'M' && rows[y+2][x+2] == 'A' && rows[y+3][x+3] == 'S' {
		matches++
	}
	// South
	if y <= (height-4) && rows[y][x] == 'X' && rows[y+1][x] == 'M' && rows[y+2][x] == 'A' && rows[y+3][x] == 'S' {
		matches++
	}
	// SouthWest
	if y <= (height-4) && x >= 3 && rows[y][x] == 'X' && rows[y+1][x-1] == 'M' && rows[y+2][x-2] == 'A' && rows[y+3][x-3] == 'S' {
		matches++
	}
	// West
	if x >= 3 && rows[y][x] == 'X' && rows[y][x-1] == 'M' && rows[y][x-2] == 'A' && rows[y][x-3] == 'S' {
		matches++
	}
	// NorthWest
	if y >= 3 && x >= 3 && rows[y][x] == 'X' && rows[y-1][x-1] == 'M' && rows[y-2][x-2] == 'A' && rows[y-3][x-3] == 'S' {
		matches++
	}

	return matches
}

type rows [][]rune

func loadInput() (rows, error) {
	var rows rows

	readFile, err := os.Open("day4/input.txt")
	if err != nil {
		return rows, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	for fileScanner.Scan() {
		rows = append(rows, []rune(fileScanner.Text()))
	}

	return rows, nil
}
