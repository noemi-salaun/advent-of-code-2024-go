package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

var registerA = 0
var registerB = 0
var registerC = 0

var pointer = 0

var output []int

var program = []int{2, 4, 1, 2, 7, 5, 1, 3, 4, 3, 5, 5, 0, 3, 3, 0}

//var program = []int{0, 3, 5, 4, 3, 0}

func main() {
	resolve()
}

func resolve() {
	var j = 0
	for i := len(program) - 1; i >= 0; i-- {
		j *= 8
		currTarget := program[i:]
		fmt.Println("NEW TARGET", i, currTarget)
		for {
			curr := runProgram(j)
			fmt.Println(j, curr)
			if slices.Equal(curr, currTarget) {
				break
			}
			j++
		}
	}
	fmt.Println(j)
}

func runProgram(i int) []int {
	registerA = i
	registerB = 0
	registerC = 0

	pointer = 0

	output = []int{}

	for pointer < len(program) {
		opcode := program[pointer]
		operand := program[pointer+1]

		switch opcode {
		case 0:
			adv(operand)
		case 1:
			bxl(operand)
		case 2:
			bst(operand)
		case 3:
			jnz(operand)
		case 4:
			bxc(operand)
		case 5:
			out(operand)
		case 6:
			bdv(operand)
		case 7:
			cdv(operand)
		default:
			panic(fmt.Sprintf("invalid opcode %d", opcode))
		}
	}

	return output
}

func getOutput(o []int) int {
	ss := make([]string, len(o))
	for i, n := range o {
		ss[i] = strconv.Itoa(n)
	}
	var ii, _ = strconv.Atoi(strings.Join(ss, ""))
	return ii
}

// adv instruction (opcode 0) performs division.
// The numerator is the value in the A register.
// The denominator is found by raising 2 to the power of the instruction's combo operand.
// (So, an operand of 2 would divide A by 4 (2^2); an operand of 5 would divide A by 2^B.)
// The result of the division operation is truncated to an integer and then written to the A register.
func adv(operand int) {
	registerA = registerA / powInt(2, combo(operand))
	pointer += 2
}

// bxl instruction (opcode 1) calculates the bitwise XOR of register B and the instruction's literal operand, then stores the result in register B.
func bxl(operand int) {
	registerB = registerB ^ operand
	pointer += 2
}

// bst instruction (opcode 2) calculates the value of its combo operand modulo 8 (thereby keeping only its lowest 3 bits), then writes that value to the B register.
func bst(operand int) {
	registerB = combo(operand) % 8
	pointer += 2
}

// jnz instruction (opcode 3) does nothing if the A register is 0.
// However, if the A register is not zero, it jumps by setting the instruction pointer to the value of its literal operand;
// if this instruction jumps, the instruction pointer is not increased by 2 after this instruction.
func jnz(operand int) {
	if registerA == 0 {
		pointer += 2
		return
	}

	pointer = operand
}

// bxc instruction (opcode 4) calculates the bitwise XOR of register B and register C,
// then stores the result in register B. (For legacy reasons, this instruction reads an operand but ignores it.)
func bxc(operand int) {
	registerB = registerB ^ registerC
	pointer += 2
}

// out instruction (opcode 5) calculates the value of its combo operand modulo 8, then outputs that value.
// (If a program outputs multiple values, they are separated by commas.)
func out(operand int) {
	output = append(output, combo(operand)%8)
	pointer += 2
}

// bdv instruction (opcode 6) works exactly like the adv instruction except that the result is stored in the B register. (The numerator is still read from the A register.)
func bdv(operand int) {
	registerB = registerA / powInt(2, combo(operand))
	pointer += 2
}

// bdv instruction (opcode 6) works exactly like the adv instruction except that the result is stored in the B register. (The numerator is still read from the A register.)
func cdv(operand int) {
	registerC = registerA / powInt(2, combo(operand))
	pointer += 2
}

func combo(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return registerA
	case 5:
		return registerB
	case 6:
		return registerC
	default:
		panic(fmt.Sprintf("Invalid operand %d", operand))
	}
}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
