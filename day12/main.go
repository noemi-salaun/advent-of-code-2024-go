package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var grid [][]string
var gridHeight int
var gridWidth int
var processed = make(map[vector2]struct{})
var regions []region
var currentRegion region

func main() {
	var err error

	grid, err = loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	gridHeight = len(grid)
	gridWidth = len(grid[0])

	for y, line := range grid {
		for x, label := range line {
			var plot = vector2{x, y}

			_, alreadyDone := processed[plot]
			if alreadyDone {
				continue
			}

			processRegion(label, plot)
		}
	}

	var price int
	for _, r := range regions {
		price += r.getFencePriceWithDiscount()
		r.printRegionStats()
	}
	fmt.Println(price)
}

func processRegion(label string, startingPlot vector2) {
	currentRegion = region{
		label:   label,
		min:     vector2{-1, -1},
		max:     vector2{-1, -1},
		tFences: make(map[vector2]struct{}),
		rFences: make(map[vector2]struct{}),
		bFences: make(map[vector2]struct{}),
		lFences: make(map[vector2]struct{}),
	}

	processPlot(label, startingPlot)

	regions = append(regions, currentRegion)
}

func processPlot(label string, plot vector2) {
	if _, alreadyDone := processed[plot]; alreadyDone {
		return
	}

	currentRegion.plots = append(currentRegion.plots, plot)
	if currentRegion.min.x == -1 || plot.x < currentRegion.min.x {
		currentRegion.min.x = plot.x
	}
	if currentRegion.min.y == -1 || plot.y < currentRegion.min.y {
		currentRegion.min.y = plot.y
	}
	if currentRegion.max.x == -1 || plot.x > currentRegion.max.x {
		currentRegion.max.x = plot.x
	}
	if currentRegion.max.y == -1 || plot.y > currentRegion.max.y {
		currentRegion.max.y = plot.y
	}

	currentRegion.perimeter += calculatePerimeter(label, plot)
	processed[plot] = struct{}{}

	neighbours := findNeighbours(label, plot)
	for _, neighbour := range neighbours {
		processPlot(label, neighbour)
	}
}

func findNeighbours(label string, plot vector2) []vector2 {

	var neighbours []vector2

	for _, dir := range []vector2{up, right, down, left} {
		neighbour := vector2{plot.x + dir.x, plot.y + dir.y}

		if neighbour.x < 0 || neighbour.x >= gridWidth || neighbour.y < 0 || neighbour.y >= gridHeight {
			continue
		}

		if grid[neighbour.y][neighbour.x] != label {
			continue
		}

		neighbours = append(neighbours, neighbour)
	}

	return neighbours
}

func calculatePerimeter(label string, plot vector2) int {
	var fences int

	for _, dir := range []vector2{up, right, down, left} {
		pos := vector2{plot.x + dir.x, plot.y + dir.y}

		if pos.x < 0 || pos.x >= gridWidth || pos.y < 0 || pos.y >= gridHeight || grid[pos.y][pos.x] != label {
			fences++

			switch dir {
			case left:
				currentRegion.lFences[vector2{plot.x, plot.y}] = struct{}{}
			case right:
				currentRegion.rFences[vector2{plot.x + 1, plot.y}] = struct{}{}
			case up:
				currentRegion.tFences[vector2{plot.x, plot.y}] = struct{}{}
			case down:
				currentRegion.bFences[vector2{plot.x, plot.y + 1}] = struct{}{}
			}
		}
	}

	return fences
}

var up = vector2{0, -1}
var right = vector2{1, 0}
var down = vector2{0, 1}
var left = vector2{-1, 0}

type vector2 struct {
	x int
	y int
}

type region struct {
	label     string
	plots     []vector2
	tFences   map[vector2]struct{}
	rFences   map[vector2]struct{}
	bFences   map[vector2]struct{}
	lFences   map[vector2]struct{}
	min       vector2
	max       vector2
	perimeter int
}

func (r *region) getFencePrice() int {
	return len(r.plots) * r.perimeter
}

func (r *region) getFencePriceWithDiscount() int {
	return len(r.plots) * r.countSides()
}

func (r *region) printRegionStats() {
	fmt.Printf("A region of %s plants with price %d * %d = %d.\n", r.label, len(r.plots), r.countSides(), r.getFencePriceWithDiscount())
}

func (r *region) countSides() int {
	var count = 0

	var lPrevious = make(map[int]struct{})
	var rPrevious = make(map[int]struct{})

	for y := r.min.y; y <= r.max.y+1; y++ {
		var lCurrent = make(map[int]struct{})
		var rCurrent = make(map[int]struct{})

		for x := r.min.x; x <= r.max.x+1; x++ {
			pos := vector2{x, y}
			if _, ok := r.lFences[pos]; ok {
				lCurrent[pos.x] = struct{}{}
				if _, p := lPrevious[pos.x]; !p {
					count++
				}
			}
			if _, ok := r.rFences[pos]; ok {
				rCurrent[pos.x] = struct{}{}
				if _, p := rPrevious[pos.x]; !p {
					count++
				}
			}
		}

		lPrevious = lCurrent
		rPrevious = rCurrent
	}

	var tPrevious = make(map[int]struct{})
	var bPrevious = make(map[int]struct{})

	for x := r.min.x; x <= r.max.x+1; x++ {
		var tCurrent = make(map[int]struct{})
		var bCurrent = make(map[int]struct{})

		for y := r.min.y; y <= r.max.y+1; y++ {
			pos := vector2{x, y}
			if _, ok := r.tFences[pos]; ok {
				tCurrent[pos.y] = struct{}{}
				if _, p := tPrevious[pos.y]; !p {
					count++
				}
			}
			if _, ok := r.bFences[pos]; ok {
				bCurrent[pos.y] = struct{}{}
				if _, p := bPrevious[pos.y]; !p {
					count++
				}
			}
		}

		tPrevious = tCurrent
		bPrevious = bCurrent
	}

	return count
}

func loadInput() ([][]string, error) {
	readFile, err := os.Open("day12/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var input [][]string

	for fileScanner.Scan() {
		line := fileScanner.Text()
		input = append(input, strings.Split(line, ""))
	}

	return input, nil
}
