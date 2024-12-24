package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	in, err := loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	bestSequenceMap := make(map[sequence]int)
	for _, secret := range in {
		bestSequences := nextNth(secret, 2000)
		for _, seq := range bestSequences {
			bestSequenceMap[seq]++
		}
	}

	type kv struct {
		Key   sequence
		Value int
	}

	var ss []kv
	for k, v := range bestSequenceMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	var bestBananas int

	for _, kv := range ss {
		count := 0
		for _, secret := range in {
			bestPrice := findPrice(secret, kv.Key, 2000)
			if bestPrice > 0 {
				count += bestPrice
			}
		}
		if count > bestBananas {
			bestBananas = count
			fmt.Println(bestBananas)
		}
	}

	fmt.Println(bestBananas)
}

func lastDigit(secret int) int {
	return secret % 10
}

type sequence struct {
	c1 int
	c2 int
	c3 int
	c4 int
}

type step struct {
	i      int
	secret int
	price  int
	change int
}

func findPrice(secret int, seq sequence, maxRound int) int {
	var current sequence

	for i := range maxRound {
		s := next(secret, i)
		secret = s.secret

		current.c1, current.c2, current.c3, current.c4 = current.c2, current.c3, current.c4, s.change

		if i >= 3 && current == seq {
			return s.price
		}
	}

	return -1
}

func nextNth(secret int, round int) []sequence {
	var steps []step
	var priceMap = make(map[int][]step)

	for i := range round {
		s := next(secret, i)
		secret = s.secret

		steps = append(steps, s)
		priceMap[s.price] = append(priceMap[s.price], s)
	}

	keys := make([]int, 0, len(priceMap))
	for k := range priceMap {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	var bestSequences []sequence
	var count = 0
out:
	for _, k := range keys {
		for _, stp := range priceMap[k] {
			if stp.i < 3 {
				continue
			}

			bestSequences = append(bestSequences, sequence{
				steps[stp.i-3].change,
				steps[stp.i-2].change,
				steps[stp.i-1].change,
				stp.change,
			})
			count++

			if count >= 1000 {
				break out
			}
		}
	}

	return bestSequences
}

func next(secret, i int) step {
	oldPrice := lastDigit(secret)

	tmp := secret * 64
	secret = mix(tmp, secret)
	secret = prune(secret)

	tmp = secret / 32
	secret = mix(tmp, secret)
	secret = prune(secret)

	tmp = secret * 2048
	secret = mix(tmp, secret)
	secret = prune(secret)

	newPrice := lastDigit(secret)

	return step{i, secret, newPrice, newPrice - oldPrice}
}

func mix(val int, secret int) int {
	return val ^ secret
}

func prune(secret int) int {
	return secret % 16777216
}

type input []int

func loadInput() (input, error) {
	var in input

	readFile, err := os.Open("day22/input.txt")
	if err != nil {
		return in, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	for fileScanner.Scan() {
		var line = fileScanner.Text()
		in = append(in, atoi(line))
	}

	return in, nil
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
