package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	in, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	var count int
	for _, k := range in.keys {
		for _, l := range in.locks {
			if k.fitIn(l) {
				count++
			}
		}
	}
	fmt.Println(count)
}

type key [5]int
type lock [5]int

func (k *key) fitIn(l lock) bool {
	for i := range *k {
		if k[i]+l[i] > 5 {
			return false
		}
	}

	return true
}

type input struct {
	keys  []key
	locks []lock
}

func loadInput() (input, error) {
	var in input

	b, err := os.ReadFile("day25/input.txt")
	if err != nil {
		return in, err
	}

	txt := string(b)

	patterns := strings.Split(txt, "\n\n")

	for _, p := range patterns {
		lines := strings.Split(p, "\n")

		if len(lines) == 0 {
			continue
		}

		var item [5]int
		for _, line := range lines {
			if len(line) != 5 {
				panic("invalid line")
			}

			for i, c := range line {
				if rune(c) == '#' {
					item[i]++
				}
			}
		}

		for i := range item {
			item[i]--
		}

		if rune(lines[0][0]) == '#' {
			in.locks = append(in.locks, item)
		} else {
			in.keys = append(in.keys, item)
		}

	}

	return in, nil
}
