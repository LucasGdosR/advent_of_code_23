package main

import (
	"aoc/2023/common"
	"fmt"
	"syscall"
)

const (
	LINE_LENGTH                = 117
	CARD_COUNT                 = 211
	WINNERS_START, WINNERS_END = 10, 40
	WINNERS_LENGTH             = WINNERS_END - WINNERS_START
	NUMBERS_START, NUMBERS_END = 42, 116
	NUMBERS_LENGTH             = NUMBERS_END - NUMBERS_START
)

func main() {
	mappedFile := common.Mmap("input")
	scratchcards := mappedFile.File
	defer syscall.Munmap(scratchcards)
	size := int(mappedFile.Size)

	var points = 0
	scratchcardsWinnings := make([]int, CARD_COUNT)
	// Forward would be:
	// 	i := 0; i < size; i += LINE_LENGTH
	// Backward allows for dynamic programming,
	// since earlier lines depend on later lines
	for i, card := size-LINE_LENGTH+1, byte(CARD_COUNT-1); i >= 0; i, card = i-LINE_LENGTH, card-1 {
		winners := scratchcards[i+WINNERS_START : i+WINNERS_END]
		numbers := scratchcards[i+NUMBERS_START : i+NUMBERS_END]

		winnerSet := makeWinnersSet(winners)
		matches := getMatches(winnerSet, numbers)
		points += 1 << (matches - 1)

		scratchcardsWon := 1
		for i := byte(1); i <= matches; i++ {
			scratchcardsWon += scratchcardsWinnings[card+i]
		}
		scratchcardsWinnings[card] = scratchcardsWon
	}
	totalScratchcards := 0
	for _, w := range scratchcardsWinnings {
		totalScratchcards += w
	}

	fmt.Println(points, totalScratchcards)
}

func makeWinnersSet(winners []byte) map[byte]struct{} {
	winnerSet := make(map[byte]struct{})
	var winner byte = 0
	for i := 0; i < WINNERS_LENGTH; i++ {
		n := winners[i]
		if n == ' ' {
			winnerSet[winner] = struct{}{}
			winner = 0
		} else {
			winner = winner*10 + n - '0'
		}
	}
	delete(winnerSet, 0)
	return winnerSet
}

func getMatches(set map[byte]struct{}, numbers []byte) byte {
	var number, matches byte = 0, 0
	for i := 0; i < NUMBERS_LENGTH; i++ {
		n := numbers[i]
		if n == ' ' {
			_, E := set[number]
			if E {
				matches++
			}
			number = 0
		} else {
			number = number*10 + n - '0'
		}
	}
	// Add the last number, that has '\n' after it instead of ' '
	_, E := set[number]
	if E {
		matches++
	}
	return matches
}
