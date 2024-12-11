package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	globalStart := time.Now()

	var array, err = loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	var count = 0
	for _, num := range *array {

		count += blink(num, 1000)
	}

	globalElapsed := time.Since(globalStart)
	fmt.Printf("Nb stones = %d - Total time %s\n", count, globalElapsed)
}

type blinkCacheKey struct {
	num   int
	depth int
}

var blinkCache = make(map[blinkCacheKey]int)

func blink(num int, depth int) int {
	key := blinkCacheKey{num, depth}
	res, ok := blinkCache[key]
	if ok {
		return res
	}

	if num == -1 {
		return 0
	}
	if depth == 0 {
		return 1
	}

	var left = -1
	var right = -1

	if num == 0 {
		left = 1
	} else {
		txt := strconv.Itoa(num)
		if len(txt)%2 == 0 {
			left, _ = strconv.Atoi(txt[:len(txt)/2])
			right, _ = strconv.Atoi(txt[len(txt)/2:])
		} else {
			left = num * 2024
		}
	}

	res = blink(left, depth-1) + blink(right, depth-1)
	blinkCache[key] = res

	return res
}

func loadInput() (*[]int, error) {
	readFile, err := os.Open("day11/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var array []int

	fileScanner.Scan()
	for _, field := range strings.Fields(fileScanner.Text()) {
		num, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}

		array = append(array, num)
	}

	return &array, nil
}
