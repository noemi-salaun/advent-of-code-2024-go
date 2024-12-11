package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	part1()
	part2()
}

func part1() {
	inputs, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Ints(inputs.Left)
	sort.Ints(inputs.Right)

	total := 0

	for i := range inputs.Left {
		left := inputs.Left[i]
		right := inputs.Right[i]

		distance := left - right
		if distance < 0 {
			distance = -distance
		}

		total += distance
	}
	fmt.Println(total)
}

func part2() {
	inputs, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	rightMap := make(map[int]int)
	for _, v := range inputs.Right {
		rightMap[v]++
	}

	total := 0
	for _, v := range inputs.Left {
		mult := v * rightMap[v]
		total += mult
	}

	fmt.Println(total)
}

type inputs struct {
	Left  []int
	Right []int
}

func loadInput() (inputs, error) {
	readFile, err := os.Open("day1/input.txt")
	var inputs inputs
	if err != nil {
		return inputs, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	for fileScanner.Scan() {
		parts := strings.Split(fileScanner.Text(), "   ")

		left, err := strconv.Atoi(parts[0])
		if err != nil {
			return inputs, err
		}

		right, err := strconv.Atoi(parts[1])
		if err != nil {
			return inputs, err
		}

		inputs.Left = append(inputs.Left, left)
		inputs.Right = append(inputs.Right, right)
	}

	return inputs, nil
}
