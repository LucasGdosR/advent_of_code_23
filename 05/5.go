package main

import (
	"aoc/2023/common"
	"bufio"
	"fmt"
	"math"
	"strings"
)

type interval struct {
	min, max int
}

type intervalTree struct {
	key         interval
	offset      int
	left, right *intervalTree
}

func main() {
	almanac := common.Open("input")
	defer almanac.Close()

	scanner := bufio.NewScanner(almanac)

	individualSeeds, seedRanges := makeSeeds(scanner)

	const seed2Soil2Fertilizer2Water2Light2Temperature2Humidity2Location = 7
	for i := 0; i < seed2Soil2Fertilizer2Water2Light2Temperature2Humidity2Location; i++ {
		tree := makeIntervalTree(scanner)
		individualSeeds = tree.Map(individualSeeds)
		seedRanges = tree.Map(seedRanges)
	}
	individualMin, rangeMin := math.MaxInt, math.MaxInt
	for _, seed := range individualSeeds {
		if seed.min < individualMin {
			individualMin = seed.min
		}
	}
	for _, seed := range seedRanges {
		if seed.min < rangeMin {
			rangeMin = seed.min
		}
	}
	fmt.Println("Closest seed (part 1):", individualMin)
	fmt.Println("Closest seed (pat 2):", rangeMin)
}

func makeSeeds(s *bufio.Scanner) ([]interval, []interval) {
	s.Scan()
	seedsStr := strings.Fields(s.Text())[1:]
	l := len(seedsStr)
	individualSeeds := make([]interval, l)
	seedRanges := make([]interval, l>>1)
	for i := 0; i < l; i += 2 {
		min, rng := common.Atoi(seedsStr[i]), common.Atoi(seedsStr[i+1])
		individualSeeds[i] = interval{min: min, max: min}
		individualSeeds[i+1] = interval{min: rng, max: rng}
		seedRanges[i>>1] = interval{min: min, max: min + rng - 1}
	}
	s.Scan()
	s.Scan()
	return individualSeeds, seedRanges
}

func makeIntervalTree(s *bufio.Scanner) *intervalTree {
	var root *intervalTree
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			break
		}
		nodeStr := strings.Fields(line)
		dstStart := common.Atoi(nodeStr[0])
		srcStart := common.Atoi(nodeStr[1])
		rng := common.Atoi(nodeStr[2])
		node := intervalTree{key: interval{min: srcStart, max: srcStart + rng}, offset: dstStart - srcStart}
		root = root.insert(&node)
	}
	s.Scan()
	return root
}

func (t *intervalTree) insert(node *intervalTree) *intervalTree {
	if t == nil {
		return node
	}
	if node.key.max <= t.key.min {
		t.left = t.left.insert(node)
	} else {
		t.right = t.right.insert(node)
	}
	return t
}

func (t *intervalTree) Map(seeds []interval) []interval {
	var mappedSeeds []interval
	for _, seed := range seeds {
		mappedSeeds = append(mappedSeeds, mapInterval(t, seed)...)
	}
	return mappedSeeds
}

func mapInterval(node *intervalTree, seed interval) []interval {
	if node == nil {
		return []interval{seed}
	} else if seed.max < node.key.min {
		return mapInterval(node.left, seed)
	} else if seed.min > node.key.max {
		return mapInterval(node.right, seed)
	}

	var result []interval

	// Handle portion before mapping if it exists
	if seed.min < node.key.min {
		result = append(result, interval{min: seed.min, max: node.key.min - 1})
	}

	// Handle mapped portion
	mappedMin := max(seed.min, node.key.min)
	mappedMax := min(seed.max, node.key.max)
	result = append(result, interval{min: mappedMin + node.offset, max: mappedMax + node.offset})

	// Handle portion after mapping if it exists
	if seed.max > node.key.max {
		result = append(result, mapInterval(node.right, interval{min: node.key.max + 1, max: seed.max})...)
	}

	return result
}
