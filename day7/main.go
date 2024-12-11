package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	result := 0
	for _, test := range *input {
		if test.isPossible() {
			result += test.value
		}
	}
	fmt.Println(result)
}

type input []*test

type test struct {
	value   int
	numbers []int
}

func (t *test) isPossible() bool {
	return check(t.value, t.numbers[0], t.numbers[1:])
}

func check(expectedValue int, head int, tail []int) bool {
	if len(tail) == 0 {
		return head == expectedValue
	}

	if check(expectedValue, head+tail[0], tail[1:]) || check(expectedValue, head*tail[0], tail[1:]) {
		return true
	}

	concat, err := strconv.Atoi(fmt.Sprintf("%d%d", head, tail[0]))
	if err != nil {
		return false
	}

	return check(expectedValue, concat, tail[1:])
}

func newTest(line string) (*test, error) {
	index := strings.Index(line, ":")
	if index == -1 {
		return nil, errors.New("missing : in line")
	}

	value, err := strconv.Atoi(line[:index])
	if err != nil {
		return nil, errors.New("invalid test value")
	}

	var numbers []int
	tail := line[index+1:]
	for _, field := range strings.Fields(tail) {
		number, err := strconv.Atoi(field)
		if err != nil {
			return nil, errors.New("invalid test number")
		}

		numbers = append(numbers, number)
	}

	if len(numbers) == 0 {
		return nil, errors.New("empty list of numbers")
	}

	return &test{value, numbers}, nil
}

func loadInput() (*input, error) {

	readFile, err := os.Open("day7/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var input input

	for fileScanner.Scan() {
		line := fileScanner.Text()

		test, err := newTest(line)
		if err != nil {
			return nil, err
		}

		input = append(input, test)
	}

	return &input, nil
}
