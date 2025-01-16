package main

import (
	"aoc/2023/common"
	"bufio"
	"bytes"
	"runtime"
	"slices"
	"strings"
	"sync"
	"syscall"
)

type handAndBid struct {
	hand string
	bid  int
}

var cardToScore = map[byte]byte{'2': 1, '3': 2, '4': 3, '5': 4, '6': 5, '7': 6, '8': 7, '9': 8, 'T': 9, 'J': 10, 'Q': 11, 'K': 12, 'A': 13}

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  rankHandsGetWinningsST,
		MT_Impl:  rankHandsGetWinningsMT,
		Part1Str: "Total winnings",
		Part2Str: "Total winnings",
	}
	common.Benchmark(thisProgram, 1000)
}

func rankHandsGetWinningsST() common.Results[int, int] {
	var winnings common.Results[int, int]
	handsAndBids := parseInput()
	slices.SortFunc(handsAndBids, func(a handAndBid, b handAndBid) int {
		return sortingFunc(a, b, countSameOfAKindRegular)
	})
	winnings.Part1 = countWinnings(handsAndBids)

	cardToScore['J'] = 0
	slices.SortFunc(handsAndBids, func(a handAndBid, b handAndBid) int {
		return sortingFunc(a, b, countSameOfAKindJoker)
	})
	cardToScore['J'] = 10 // So this can be benchmarked multiple times
	winnings.Part2 = countWinnings(handsAndBids)
	return winnings
}

func rankHandsGetWinningsMT() common.Results[int, int] {
	// Variables
	mappedFile := common.Mmap("input")
	file := mappedFile.File
	defer syscall.Munmap(file)
	size := mappedFile.Size

	numWorkers := runtime.GOMAXPROCS(0)
	bytesPerWorker := size/int64(numWorkers) + 1
	var wg sync.WaitGroup
	partialSortedHands := make(chan []handAndBid, numWorkers)
	partialWinnings := make(chan int, numWorkers)
	var winnings common.Results[int, int]

	// Actual solution starts here
	for i := 0; i < numWorkers; i++ {
		start := int64(i) * bytesPerWorker
		end := start + bytesPerWorker
		if i == numWorkers-1 {
			end = size
		}
		// Backtracking for getting lines out of Mmap
		go func(start, end int64) {
			if start > 0 {
				for file[start-1] != '\n' {
					start--
				}
			}
			if end < int64(len(file)) {
				for file[end-1] != '\n' {
					end--
				}
			}
			partialSortedHands <- parseAndSort(file, start, end)
		}(start, end)
	}
	handsAndBids := merge(&wg, numWorkers, partialSortedHands, countSameOfAKindRegular)
	handsPerWorker := len(handsAndBids) / numWorkers
	winnings.Part1 = getWinnings(numWorkers, handsPerWorker, handsAndBids, partialWinnings)

	// Part 2
	cardToScore['J'] = 0
	for i := 0; i < numWorkers; i++ {
		start := i * handsPerWorker
		end := start + handsPerWorker
		if i == numWorkers-1 {
			end = len(handsAndBids)
		}
		go func(HABs []handAndBid) {
			slices.SortFunc(HABs, func(a handAndBid, b handAndBid) int {
				return sortingFunc(a, b, countSameOfAKindJoker)
			})
			partialSortedHands <- HABs
		}(handsAndBids[start:end])
	}
	handsAndBids = merge(&wg, numWorkers, partialSortedHands, countSameOfAKindJoker)
	close(partialSortedHands)
	winnings.Part2 = getWinnings(numWorkers, handsPerWorker, handsAndBids, partialWinnings)
	close(partialWinnings)

	cardToScore['J'] = 10 // So this can be benchmarked multiple times
	return winnings
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

func parseAndSort(file []byte, start, end int64) []handAndBid {
	handsAndBids := make([]handAndBid, 0, 256) // tailored for 4 workers. could receive numworkers as parameter for initial capacity
	for ; start < end; start++ {
		var HAB handAndBid
		HAB.hand = string(file[start : start+5])
		start += 6
		for ; start < end && file[start] != '\n'; start++ {
			HAB.bid = HAB.bid*10 + int(file[start]-'0')
		}
		handsAndBids = append(handsAndBids, HAB)
	}
	slices.SortFunc(handsAndBids, func(a handAndBid, b handAndBid) int {
		return sortingFunc(a, b, countSameOfAKindRegular)
	})
	return handsAndBids
}

func merge(
	wg *sync.WaitGroup,
	numWorkers int,
	partialSortedHands chan []handAndBid,
	countSameOfAKind func([]byte) (byte, byte)) []handAndBid {

	for i := 0; i < numWorkers-1; i++ {
		wg.Add(1)
		go func() {
			partialSortedHands <- mergeSort(partialSortedHands, countSameOfAKind)
			wg.Done()
		}()
	}
	wg.Wait()
	return <-partialSortedHands
}

func mergeSort(partialSortedHands chan []handAndBid, countSameOfAKind func([]byte) (byte, byte)) []handAndBid {
	a := <-partialSortedHands
	b := <-partialSortedHands
	i, j, la, lb := 0, 0, len(a), len(b)
	merged := make([]handAndBid, 0, la+lb)
	for i < la && j < lb {
		if sortingFunc(a[i], b[j], countSameOfAKind) < 0 {
			merged = append(merged, a[i])
			i++
		} else {
			merged = append(merged, b[j])
			j++
		}
	}
	merged = append(merged, a[i:]...)
	merged = append(merged, b[j:]...)
	return merged
}

func getWinnings(numWorkers, handsPerWorker int,
	HABs []handAndBid, partialWinnings chan int) int {
	for i := 0; i < numWorkers; i++ {
		start := i * handsPerWorker
		end := start + handsPerWorker
		if i == numWorkers-1 {
			end = len(HABs)
		}
		go func(start, end int) {
			totalWinnings := 0
			for ; start < end; start++ {
				totalWinnings += (start + 1) * HABs[start].bid
			}
			partialWinnings <- totalWinnings
		}(start, end)
	}
	var winnings int
	for i := 0; i < numWorkers; i++ {
		winnings += <-partialWinnings
	}
	return winnings
}
