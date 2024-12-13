package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var machines, err = loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	var tokensCount int

	for _, m := range *machines {
		combo := m.solve()

		if combo != nil {
			tokensCount += combo.getTokenPrice()
		}

	}

	fmt.Println(tokensCount)
}

type machine struct {
	buttonA button
	buttonB button
	prize   prize
}

type buttonCombo struct {
	a int
	b int
}

func (bc *buttonCombo) getTokenPrice() int {
	return bc.a*3 + bc.b*1
}

func (m *machine) solve() *buttonCombo {
	ax, ay := m.buttonA.x, m.buttonA.y
	bx, by := m.buttonB.x, m.buttonB.y
	tx, ty := m.prize.x, m.prize.y

	det := float64(by*ax - bx*ay)
	if det == 0 {
		panic("AAAAAAAAAAHHHHHhhhhhh")
	}

	b := -float64(tx*ay-ty*ax) / det
	a := float64(tx*by-ty*bx) / det

	if math.Ceil(b) == b && math.Ceil(a) == a {
		return &buttonCombo{int(a), int(b)}
	}
	return nil
}

func (m *machine) findGoodCombos() []buttonCombo {
	maxB := min(m.prize.x/m.buttonB.x, m.prize.y/m.buttonB.y)

	var goodCombos []buttonCombo

	for countB := 0; countB <= maxB; countB++ {

		if countB%100000 == 0 {
			fmt.Println(countB)
			fmt.Println(maxB)
		}

		res := calculateButtonCount(m.buttonA.x, m.buttonB.x, countB, m.prize.x)
		if math.Ceil(res) != res {
			continue
		}

		countA := int(res)

		if m.buttonA.y*countA+m.buttonB.y*countB != m.prize.y {

			continue
		}

		fmt.Printf("X and Y ok : %d %d\n", countA, countB)

		goodCombos = append(goodCombos, buttonCombo{countA, countB})
	}

	return goodCombos
}

type button vector2

func newButton(line string) button {
	var re = regexp.MustCompile(`X\+(\d+), Y\+(\d+)`)

	var b button

	matches := re.FindAllStringSubmatch(line, -1)
	b.x, _ = strconv.Atoi(matches[0][1])
	b.y, _ = strconv.Atoi(matches[0][2])

	return b
}

type prize vector2

func newPrize(line string) prize {
	var re = regexp.MustCompile(`X=(\d+), Y=(\d+)`)

	var p prize

	matches := re.FindAllStringSubmatch(line, -1)
	x, _ := strconv.Atoi(matches[0][1])
	y, _ := strconv.Atoi(matches[0][2])

	p.x = 10000000000000 + x
	p.y = 10000000000000 + y

	return p
}

type vector2 struct {
	x, y int
}

func calculateButtonCount(theButton, otherButton, countOtherButton, target int) float64 {
	return (float64(target) - (float64(otherButton) * float64(countOtherButton))) / float64(theButton)
}

func loadInput() (*[]machine, error) {
	readFile, err := os.Open("day13/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var machines []machine
	var line string
	var index int

	var currentMachine machine

	for fileScanner.Scan() {
		line = fileScanner.Text()

		index = strings.Index(line, "Button A: ")
		if index != -1 {
			currentMachine.buttonA = newButton(line[len("Button A: "):])
			continue
		}

		index = strings.Index(line, "Button B: ")
		if index != -1 {
			currentMachine.buttonB = newButton(line[len("Button B: "):])
			continue
		}

		index = strings.Index(line, "Prize: ")
		if index != -1 {
			currentMachine.prize = newPrize(line[len("Prize: "):])
			machines = append(machines, currentMachine)
			currentMachine = machine{}
			continue
		}
	}

	return &machines, nil
}
