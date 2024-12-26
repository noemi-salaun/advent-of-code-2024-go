package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var gatesMap = make(map[string][]*gate)
var wiresMap = make(map[string]int)

var outGatesMap = make(map[string]*gate)

var remainingZ int

func mdfgdfgain() {
	in, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, g := range in.gates {
		outGatesMap[g.out] = g
	}

	for _, w := range in.wires {
		wiresMap[w.name] = w.val
	}

	var zWires []*wire

	for _, g := range in.gates {
		if string(g.out[0]) == "z" {
			fmt.Println(g.out, g.findDeps())

		}
	}

	fmt.Println(getWiresNumberFromSlices("z", zWires))
}

func uniq(list []string) []string {
	if list == nil {
		return nil
	}
	out := make([]string, len(list))
	copy(out, list)
	slices.Sort(out)
	uniq := out[:0]
	for _, x := range out {
		if len(uniq) == 0 || uniq[len(uniq)-1] != x {
			uniq = append(uniq, x)
		}
	}
	return uniq
}

func (g *gate) findDeps() []string {
	var deps []string

	if g.gA == nil {
		gA, ok := outGatesMap[g.a]
		if ok {
			g.gA = gA
			deps = append(deps, g.gA.out)
			deps = append(deps, g.gA.findDeps()...)
		} else {
			deps = append(deps, g.a)
		}
	}

	if g.gB == nil {
		gB, ok := outGatesMap[g.b]
		if ok {
			g.gB = gB
			deps = append(deps, g.gB.out)
			deps = append(deps, g.gB.findDeps()...)
		} else {
			deps = append(deps, g.b)
		}
	}

	return deps
}

func (g *gate) getValue() int {
	if g.value != -1 {
		return g.value
	}

	var a int
	if g.gA != nil {
		a = g.gA.getValue()
	} else {
		a = wiresMap[g.a]
	}

	var b int
	if g.gB != nil {
		b = g.gB.getValue()
	} else {
		b = wiresMap[g.b]
	}

	g.value = g.op.apply(a, b)
	return g.value
}

func (g *gate) simplify() bool {
	if (g.gA == nil && g.gB != nil) || (g.gA != nil && g.gB == nil) {
		panic("incoherent state")
	}

	if g.value != -1 {
		return false
	}

	var a = -1
	var b = -1

	if g.gA == nil && g.gB == nil {
		a = wiresMap[g.a]
		b = wiresMap[g.b]

		val := g.op.apply(a, b)
		g.value = val

		return true
	}

	if g.gA == nil {
		a = wiresMap[g.a]
	} else if g.gA.value != -1 {
		a = g.gA.value
	}

	if g.gB == nil {
		b = wiresMap[g.b]
	} else if g.gB.value != -1 {
		b = g.gB.value
	}

	if a != -1 && b == -1 {
		val := g.op.apply(a, b)
		g.value = val

		return true
	}

	var sa bool
	if a == -1 {
		sa = g.gA.simplify()
	}
	var sb bool
	if b == -1 {
		sb = g.gB.simplify()
	}

	return sa || sb
}

func (g *gate) print() string {
	if g.value != -1 {
		return fmt.Sprintf("%d", g.value)
	}

	var gA string
	if g.gA == nil {
		gA = g.a
	} else {
		gA = g.gA.print()
	}

	var gB string
	if g.gB == nil {
		gB = g.b
	} else {
		gB = g.gB.print()
	}

	return fmt.Sprintf("(%s %s %s)", gA, g.op.print(), gB)
}

func main2() {
	in, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, g := range in.gates {
		gatesMap[g.a] = append(gatesMap[g.a], g)
		gatesMap[g.b] = append(gatesMap[g.b], g)

		if string(g.out[0]) == "z" {
			remainingZ++
		}
	}

	x := getWiresNumberFromSlices("x", in.wires)
	y := getWiresNumberFromSlices("y", in.wires)

	fmt.Println(x, y)

	fmt.Println(add(x, y, remainingZ-1))

	fmt.Println(x + y)
}

func initWires(prefix string, number int64, nbWires int) []*wire {
	bin := reverse(strconv.FormatInt(number, 2))
	var wires []*wire

	var length = len(bin)
	for i := 0; i < nbWires; i++ {
		var val int
		if i < length {
			val, _ = strconv.Atoi(string(bin[i]))
		} else {
			val = 0
		}

		wires = append(wires, &wire{name: fmt.Sprintf("%s%02d", prefix, i), val: val})
	}

	return wires
}

func getWiresNumberFromMap(prefix string, m map[string]int) int64 {
	var names []string
	for w := range m {
		if string(w[0]) == prefix {
			names = append(names, w)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(names)))

	var outputS string
	for _, name := range names {
		outputS += strconv.Itoa(m[name])
	}

	outputI, _ := strconv.ParseInt(outputS, 2, 64)
	return outputI
}

func getWiresNumberFromSlices(prefix string, wires []*wire) int64 {
	var m = make(map[string]int)

	for _, w := range wires {
		m[w.name] = w.val
	}

	return getWiresNumberFromMap(prefix, m)
}

func add(x int64, y int64, nbWires int) int64 {

	wires := slices.Concat(initWires("x", x, nbWires), initWires("y", y, nbWires))

	for _, w := range wires {
		w.emit()
	}

	return getWiresNumberFromMap("z", wiresMap)
}

type wire struct {
	name string
	val  int
}

func (w *wire) emit() {
	wiresMap[w.name] = w.val

	if string(w.name[0]) == "z" {
		remainingZ--
		if remainingZ == 0 {
			return
		}
	}

	gates := getOpenedGatesWithInput(w.name)

	for _, g := range gates {
		g.emitIfAllInputs()
	}
}

func getOpenedGatesWithInput(wireName string) []*gate {
	gates, ok := gatesMap[wireName]
	if !ok {
		return []*gate{}
	}

	return slices.DeleteFunc(gates, func(g *gate) bool {
		return g.closed
	})
}

func newWire(s string) *wire {
	parts := strings.Split(s, ": ")
	i, _ := strconv.Atoi(parts[1])
	return &wire{parts[0], i}
}

type gate struct {
	a      string
	b      string
	op     operator
	out    string
	closed bool

	gA *gate
	gB *gate

	value int
}

func newGate(s string) *gate {

	parts := strings.Split(s, " -> ")
	subParts := strings.Split(parts[0], " ")

	return &gate{
		a:      subParts[0],
		b:      subParts[2],
		op:     newOperator(subParts[1]),
		out:    parts[1],
		closed: false,
		value:  -1,
	}
}

func (g *gate) emitIfAllInputs() {
	if g.closed {
		return
	}

	a, ok := wiresMap[g.a]
	if !ok {
		return
	}

	b, ok := wiresMap[g.b]
	if !ok {
		return
	}

	out := g.op.apply(a, b)
	w := wire{
		name: g.out,
		val:  out,
	}

	w.emit()
}

type operator int

const (
	AND operator = iota
	OR
	XOR
)

func newOperator(s string) operator {
	switch s {
	case "AND":
		return AND
	case "OR":
		return OR
	case "XOR":
		return XOR
	}
	panic("invalid new operator")
}

func (op operator) apply(a, b int) int {
	switch op {
	case AND:
		return a & b
	case OR:
		return a | b
	case XOR:
		return a ^ b
	}
	panic("invalid operator")
}

func (op operator) print() string {
	switch op {
	case AND:
		return "AND"
	case OR:
		return "OR"
	case XOR:
		return "XOR"
	}
	panic("invalid new operator")
}

type input struct {
	wires []*wire
	gates []*gate
}

func loadInput() (input, error) {
	var in input

	readFile, err := os.Open("day24/input.txt")
	if err != nil {
		return in, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var wiresDone bool
	for fileScanner.Scan() {
		var line = fileScanner.Text()

		if line == "" {
			wiresDone = true
			continue
		}

		if !wiresDone {
			in.wires = append(in.wires, newWire(line))
		} else {
			in.gates = append(in.gates, newGate(line))
		}
	}

	return in, nil
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
