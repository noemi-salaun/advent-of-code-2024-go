package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	operations, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	var total int
	for _, op := range operations {
		total += op.left * op.right
	}
	fmt.Println(total)
}

type operationType struct {
	left  int
	right int
}

const token_start = "mul("
const token_comma = ","
const token_end = ")"
const max_digit = 3
const do = "do()"
const dont = "don't()"

const (
	START = iota
	LEFT_DIGITS
	RIGHT_DIGITS
)

func loadInput() ([]operationType, error) {
	readFile, err := os.Open("day3/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	return parseOperations(readFile)
}

func parseOperations(reader io.Reader) ([]operationType, error) {
	var operations []operationType

	fileScanner := bufio.NewScanner(reader)
	fileScanner.Split(bufio.ScanRunes)

	var enabled = true
	var operation operationType
	var search = START
	var bufferOp = ""
	var bufferDo = ""
	for fileScanner.Scan() {
		char := fileScanner.Text()

		if enabled {
			// Search for don't()
			bufferDo += char
			if bufferDo == dont[:len(bufferDo)] {
				if len(bufferDo) == len(dont) {
					bufferDo = ""
					bufferOp = ""
					search = START
					enabled = false
					continue
				}
			} else {
				bufferDo = ""
			}

		} else {
			// Search for do()
			bufferDo += char
			if bufferDo == do[:len(bufferDo)] {
				if len(bufferDo) == len(do) {
					bufferDo = ""
					bufferOp = ""
					search = START
					enabled = true
					continue
				}
			} else {
				bufferDo = ""
			}
		}

		if !enabled {
			continue
		}

		switch search {
		case START:
			bufferOp += char
			if bufferOp == token_start[:len(bufferOp)] {
				if len(bufferOp) == len(token_start) {
					bufferOp = ""
					search = LEFT_DIGITS
				}
				continue
			}

		case LEFT_DIGITS:
			if _, err := strconv.Atoi(char); err == nil {
				bufferOp += char
				if len(bufferOp) <= max_digit {
					continue
				}
			} else if (len(bufferOp) > 0) && (char == token_comma) {
				leftVal, err := strconv.Atoi(bufferOp)
				if err == nil {
					operation.left = leftVal
					bufferOp = ""
					search = RIGHT_DIGITS
					continue
				}
			}

		case RIGHT_DIGITS:
			if _, err := strconv.Atoi(char); err == nil {
				bufferOp += char
				if len(bufferOp) <= max_digit {
					continue
				}
			} else if (len(bufferOp) > 0) && (char == token_end) {
				rightVal, err := strconv.Atoi(bufferOp)
				if err == nil {
					operation.right = rightVal
					operations = append(operations, operation)
				}
			}

		}

		// Bad or complete match, reset the parsing
		bufferOp = ""
		search = START
		operation = operationType{}
	}

	return operations, nil
}
