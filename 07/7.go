package main

import (
	"aoc/2023/common"
	"bufio"
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type handAndBid struct {
	hand string
	bid  int
}

var cardToScore = map[byte]byte{'2': 1, '3': 2, '4': 3, '5': 4, '6': 5, '7': 6, '8': 7, '9': 8, 'T': 9, 'J': 10, 'Q': 11, 'K': 12, 'A': 13}

func main() {
	handsAndBids := parseInput()
	slices.SortFunc(handsAndBids, func(a handAndBid, b handAndBid) int {
		return sortingFunc(a, b, countSameOfAKindRegular)
	})
	fmt.Println(countWinnings(handsAndBids))

	// Part 2:
	cardToScore['J'] = 0
	slices.SortFunc(handsAndBids, func(a handAndBid, b handAndBid) int {
		return sortingFunc(a, b, countSameOfAKindJoker)
	})
	fmt.Println(countWinnings(handsAndBids))
}

func sortingFunc(a handAndBid, b handAndBid, countSameOfAKind func([]byte) (byte, byte)) int {
	aScore := getHandScore(a.hand, countSameOfAKind)
	bScore := getHandScore(b.hand, countSameOfAKind)
	if aScore < bScore {
		return -1
	} else if bScore < aScore {
		return 1
	} else {
		return cmpIndividualCards(a.hand, b.hand)
	}
}

func getHandScore(hand string, countSameOfAKind func([]byte) (byte, byte)) byte {
	h := []byte(hand)
	slices.Sort(h)

	big, small := countSameOfAKind(h)
	switch big {
	case 5:
		return 6
	case 4:
		return 5
	case 3:
		if small == 2 {
			return 4
		} else {
			return 3
		}
	case 2:
		if small == 2 {
			return 2
		} else {
			return 1
		}
	default:
		return 0
	}
}

func countSameOfAKindRegular(h []byte) (byte, byte) {
	var big, small, streak, curr byte
	for _, c := range h {
		if c != curr {
			curr = c
			big, small = endStreak(big, small, streak)
			streak = 0
		}
		streak++
	}
	return endStreak(big, small, streak)
}

func countSameOfAKindJoker(h []byte) (byte, byte) {
	var i int
	j := 4
	i = bytes.IndexByte(h, 'J')
	if i == -1 {
		i = 5
	} else {
		j = bytes.LastIndexByte(h, 'J')
	}
	jokers := j - i + 1

	startBig, startSmall := countSameOfAKindRegular(h[:i])
	endBig, endSmall := countSameOfAKindRegular(h[j+1:])

	var theBig, theSmall byte
	if startBig > endBig {
		theBig = startBig + byte(jokers)
		if endBig > startSmall {
			theSmall = endBig
		} else {
			theSmall = startSmall
		}
	} else {
		theBig = endBig + byte(jokers)
		if startBig > endSmall {
			theSmall = startBig
		} else {
			theSmall = endSmall
		}
	}

	return theBig, theSmall
}

func endStreak(big, small, streak byte) (byte, byte) {
	if streak > big {
		small = big
		big = streak
	} else if streak > small {
		small = streak
	}
	return big, small
}

func cmpIndividualCards(a string, b string) int {
	aScore, bScore := cardToScore[a[0]], cardToScore[b[0]]
	if aScore < bScore {
		return -1
	} else if bScore < aScore {
		return 1
	} else {
		return cmpIndividualCards(a[1:], b[1:])
	}
}

func countWinnings(HABs []handAndBid) int {
	totalWinnings := 0
	for rankMinusOne, HAB := range HABs {
		totalWinnings += (rankMinusOne + 1) * HAB.bid
	}
	return totalWinnings
}

func parseInput() []handAndBid {
	input := common.Open("input")
	defer input.Close()
	scanner := bufio.NewScanner(input)
	handsAndBids := make([]handAndBid, 0, 1000)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		handsAndBids = append(handsAndBids, handAndBid{hand: line[0], bid: common.Atoi(line[1])})
	}
	return handsAndBids
}
