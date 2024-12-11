package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

func main() {
	input, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	expanded := expandDiskMap(input)
	compacted := compact(expanded)
	sum := checksum(compacted)
	fmt.Println(sum)
}

func loadInput() ([]int, error) {
	readFile, err := os.Open("day9/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanBytes)

	var input []int

	for fileScanner.Scan() {
		char := fileScanner.Text()

		num, err := strconv.Atoi(char)
		if err != nil {
			return nil, err
		}

		input = append(input, num)
	}
	return input, nil
}

func expandDiskMap(dm []int) []int {
	var result []int
	id := 0
	for i, count := range dm {
		even := i%2 == 0

		var val int
		if even {
			val = id
			id++
		} else {
			val = -1
		}

		result = append(result, repeatInts(val, count)...)
	}

	return result
}

func repeatInts(val, count int) []int {
	result := make([]int, count)
	for i := 0; i < count; i++ {
		result[i] = val
	}
	return result
}

func defrag(input []int) []int {
	firstEmpty := slices.Index(input, -1)
	if firstEmpty == -1 || firstEmpty >= len(input)-1 {
		return input
	}

	for i := len(input) - 1; i > firstEmpty; i-- {
		if input[i] != -1 {
			input[firstEmpty], input[i] = input[i], input[firstEmpty]
			firstEmpty = indexAt(input, -1, firstEmpty)
		}
		if firstEmpty > i {
			break
		}
	}

	return input
}

func indexAt(s []int, v int, n int) int {
	idx := slices.Index(s[n:], v)
	if idx > -1 {
		idx += n
	}
	return idx
}

func compact(input []int) []int {
	for i := len(input) - 1; i > 0; i-- {
		num := input[i]
		if num == -1 {
			continue
		}
		numGrp := findGroup(input, num, 0)
		if numGrp.index == -1 {
			continue
		}

		emptyGrp := findEmptyGroup(input, numGrp.length)
		if emptyGrp.index == -1 || emptyGrp.index > numGrp.index {
			continue
		}

		for j := 0; j < numGrp.length; j++ {
			input[numGrp.index+j], input[emptyGrp.index+j] = input[emptyGrp.index+j], input[numGrp.index+j]
		}
	}
	return input
}

type groupPos struct {
	index  int
	length int
}

func findGroup(input []int, id int, startAt int) groupPos {
	index := indexAt(input, id, startAt)
	if index == -1 {
		return groupPos{-1, 0}
	}

	length := 1
	for i := index + 1; i < len(input); i++ {
		if input[i] != id {
			break
		}
		length++
	}

	return groupPos{index, length}
}

func findEmptyGroup(input []int, minLength int) groupPos {
	startAt := 0
	for startAt < len(input) {
		grp := findGroup(input, -1, startAt)
		if grp.index == -1 || grp.length >= minLength {
			return grp
		}
		startAt = grp.index + grp.length
	}
	return groupPos{-1, 0}
}

func checksum(input []int) int {
	sum := 0
	for i, v := range input {
		if v != -1 {
			sum += i * v
		}
	}
	return sum
}
