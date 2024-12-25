package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var numericalKeypad keypad
var directionalKeypad keypad

var getAllPossibleSequencesForButtonCache = make(map[bool]map[vector2]map[vector2][]code)

var chainCodeCache = make(map[int]map[string][]code)

func main() {
	initKeypads()

	codes, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	var total = 0
	for _, c := range codes {
		complexity := calculateComplexity(c)

		return
		total += complexity
	}

	fmt.Println(total)
}

func calculateComplexity(c code) int {
	var bestStep code
	var bestStepLen = -1

	for _, step1 := range numericalKeypad.getAllPossibleSequencesForCode(vector2{0, 0}, c) {

		for _, stepEnd := range chainCode(step1, 25) {
			sLen := len(stepEnd)
			if bestStepLen == -1 || sLen < bestStepLen {
				bestStepLen = sLen
				bestStep = stepEnd
			}
		}
	}

	fmt.Println(strings.Join(bestStep, ""))
	fmt.Println(len(bestStep), c.getNumericPart())

	complexity := len(bestStep) * c.getNumericPart()
	return complexity
}

func chainCode(c code, steps int) []code {
	if steps == 0 {
		return []code{c}
	}

	res := get_chainCodeCache(steps, c)
	if res != nil {
		return res
	}

	var result [][]code

	for _, possibleCode := range directionalKeypad.getAllPossibleSequencesForCode(vector2{0, 0}, c) {
		result = append(result, chainCode(possibleCode, steps-1))
	}

	pouet := slices.Concat(result...)
	save_chainCodeCache(steps, c, pouet)
	return pouet
}

func (pad *keypad) get(key string) vector2 {
	v, ok := pad.layout[key]
	if !ok {
		panic(fmt.Sprintf("Invalid key %s", key))
	}

	return v
}

func (pad *keypad) getSequenceForCode(c code) code {
	var sequence code
	var pos = pad.get("A")

	for _, k := range c {
		tar := pad.get(k)
		sequence = append(sequence, getSequenceForButton(pos, tar, pad.inverse)...)
		pos = tar
	}

	return sequence
}

func (pad *keypad) getAllPossibleSequencesForCode(start vector2, c code) []code {

	if len(c) == 0 {
		return nil
	}

	res := pad.get_getAllPossibleSequencesForCodeCache(start, c)
	if res != nil {
		return res
	}

	head := c[0]
	tail := c[1:]

	var result []code

	tar := pad.get(head)
	for _, seqBut := range getAllPossibleSequencesForButton(start, tar, pad.inverse) {

		seqs := pad.getAllPossibleSequencesForCode(tar, tail)

		if len(seqs) == 0 {
			result = append(result, seqBut)
		} else {
			for _, seqCode := range seqs {
				result = append(result, append(seqBut, seqCode...))
			}
		}

	}

	pad.save_getAllPossibleSequencesForCodeCache(start, c, result)

	return result
}

func getAllPossibleSequencesForButton(from vector2, to vector2, inverse bool) []code {
	res := get_getAllPossibleSequencesForButtonCache(inverse, from, to)
	if res != nil {
		return res
	}

	var codes []code

	var steps = getAllPossibleStepsForButton(from, to, inverse)

	if len(steps) == 0 {
		codes = []code{{"A"}}
		save_getAllPossibleSequencesForButtonCache(inverse, from, to, codes)

		return codes
	}

	for _, s := range steps {
		for _, seq := range getAllPossibleSequencesForButton(s.pos, to, inverse) {
			var c = append(code{s.move}, seq...)
			codes = append(codes, c)
		}
	}

	save_getAllPossibleSequencesForButtonCache(inverse, from, to, codes)
	return codes
}

type step struct {
	move string
	pos  vector2
}

func getAllPossibleStepsForButton(from vector2, to vector2, inverse bool) []step {
	var steps []step
	if to.y > from.y {
		var char = "^"
		if inverse {
			char = "v"
		}
		steps = append(steps, step{char, vector2{from.x, from.y + 1}})
	}
	if to.x < from.x {
		steps = append(steps, step{">", vector2{from.x - 1, from.y}})
	}
	if to.y < from.y && !(from.x == 2 && from.y == 1) {
		var char = "v"
		if inverse {
			char = "^"
		}
		steps = append(steps, step{char, vector2{from.x, from.y - 1}})
	}
	if to.x > from.x && !(from.x == 1 && from.y == 0) {
		steps = append(steps, step{"<", vector2{from.x + 1, from.y}})
	}

	return steps
}

func getSequenceForButton(from vector2, to vector2, inverse bool) code {
	var sequence code
	if to.y > from.y {
		var char = "^"
		if inverse {
			char = "v"
		}
		sequence = append(sequence, strings.Split(strings.Repeat(char, to.y-from.y), "")...)
	}
	if to.x < from.x {
		sequence = append(sequence, strings.Split(strings.Repeat(">", from.x-to.x), "")...)
	}
	if to.y < from.y {
		var char = "v"
		if inverse {
			char = "^"
		}
		sequence = append(sequence, strings.Split(strings.Repeat(char, from.y-to.y), "")...)
	}
	if to.x > from.x {
		sequence = append(sequence, strings.Split(strings.Repeat("<", to.x-from.x), "")...)
	}

	sequence = append(sequence, "A")

	return sequence
}

type vector2 struct {
	x, y int
}

type code []string

func (c *code) toString() string {
	return strings.Join(*c, "")
}

func (c *code) getNumericPart() int {
	var number string
	for _, char := range *c {
		if isNumeric(char) {
			number += char
		}
	}

	res, err := strconv.Atoi(number)
	if err != nil {
		return 0
	}
	return res
}

type keypad struct {
	layout                              map[string]vector2
	inverse                             bool
	getAllPossibleSequencesForCodeCache map[vector2]map[string][]code
}

func (pad *keypad) get_getAllPossibleSequencesForCodeCache(v vector2, c code) []code {
	res, ok := pad.getAllPossibleSequencesForCodeCache[v]
	if !ok {
		return nil
	}

	res2, ok := res[c.toString()]
	if !ok {
		return nil
	}

	return res2
}

func (pad *keypad) save_getAllPossibleSequencesForCodeCache(v vector2, c code, res []code) {
	subMap := pad.getAllPossibleSequencesForCodeCache[v]
	if subMap == nil {
		pad.getAllPossibleSequencesForCodeCache[v] = make(map[string][]code)
	}
	pad.getAllPossibleSequencesForCodeCache[v][c.toString()] = res
}

func get_getAllPossibleSequencesForButtonCache(inverse bool, from vector2, to vector2) []code {
	res, ok := getAllPossibleSequencesForButtonCache[inverse]
	if !ok {
		return nil
	}

	res2, ok := res[from]
	if !ok {
		return nil
	}

	res3, ok := res2[to]
	if !ok {
		return nil
	}

	return res3
}

func save_getAllPossibleSequencesForButtonCache(inverse bool, from vector2, to vector2, res []code) {
	subMap := getAllPossibleSequencesForButtonCache[inverse]
	if subMap == nil {
		getAllPossibleSequencesForButtonCache[inverse] = make(map[vector2]map[vector2][]code)
	}

	subSubMap := getAllPossibleSequencesForButtonCache[inverse][from]
	if subSubMap == nil {
		getAllPossibleSequencesForButtonCache[inverse][from] = make(map[vector2][]code)
	}

	getAllPossibleSequencesForButtonCache[inverse][from][to] = res
}

func get_chainCodeCache(steps int, c code) []code {
	res, ok := chainCodeCache[steps]
	if !ok {
		return nil
	}

	res2, ok := res[c.toString()]
	if !ok {
		return nil
	}

	return res2
}

func save_chainCodeCache(steps int, c code, res []code) {
	subMap := chainCodeCache[steps]
	if subMap == nil {
		chainCodeCache[steps] = make(map[string][]code)
	}
	chainCodeCache[steps][c.toString()] = res
}

func loadInput() ([]code, error) {

	readFile, err := os.Open("day21/demo.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var codes []code
	for fileScanner.Scan() {
		var line = fileScanner.Text()
		var chars = strings.Split(line, "")

		codes = append(codes, chars)
	}

	return codes, nil
}

func initKeypads() {
	numericalKeypad = keypad{layout: make(map[string]vector2), inverse: false}
	numericalKeypad.layout["A"] = vector2{0, 0}
	numericalKeypad.layout["0"] = vector2{1, 0}
	numericalKeypad.layout["3"] = vector2{0, 1}
	numericalKeypad.layout["2"] = vector2{1, 1}
	numericalKeypad.layout["1"] = vector2{2, 1}
	numericalKeypad.layout["6"] = vector2{0, 2}
	numericalKeypad.layout["5"] = vector2{1, 2}
	numericalKeypad.layout["4"] = vector2{2, 2}
	numericalKeypad.layout["9"] = vector2{0, 3}
	numericalKeypad.layout["8"] = vector2{1, 3}
	numericalKeypad.layout["7"] = vector2{2, 3}

	numericalKeypad.getAllPossibleSequencesForCodeCache = make(map[vector2]map[string][]code)

	directionalKeypad = keypad{layout: make(map[string]vector2), inverse: true}
	directionalKeypad.layout["A"] = vector2{0, 0}
	directionalKeypad.layout["^"] = vector2{1, 0}
	directionalKeypad.layout[">"] = vector2{0, 1}
	directionalKeypad.layout["v"] = vector2{1, 1}
	directionalKeypad.layout["<"] = vector2{2, 1}

	directionalKeypad.getAllPossibleSequencesForCodeCache = make(map[vector2]map[string][]code)
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
