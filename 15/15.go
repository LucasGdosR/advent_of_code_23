package main

import (
	"aoc/2023/common"
	"syscall"
)

type lens struct {
	label       string
	focalLength byte
}

func main() {
	mappedFile := common.Mmap("input")
	initializationSequence := mappedFile.File
	size := int(mappedFile.Size)
	defer syscall.Munmap(initializationSequence)

	var results common.Results[int, int]

	var HASH byte
	boxes := make([][]lens, 256)
	for i := range boxes {
		boxes[i] = make([]lens, 0)
	}

	for l, r := 0, 0; r < size; r++ {
		b := initializationSequence[r]
		if b == ',' {
			results.Part1 += int(HASH)
			HASH = 0
			l = r + 1
		} else {
			switch b {
			case '-':
				box := boxes[HASH]
				label := string(initializationSequence[l:r])
				for i, lens := range box {
					if lens.label == label {
						boxes[HASH] = append(box[:i], box[i+1:]...)
						break
					}
				}
			case '=':
				box := boxes[HASH]
				label := string(initializationSequence[l:r])
				focalLength := initializationSequence[r+1] - '0'
				found := false
				for i, lens := range box {
					if lens.label == label {
						box[i].focalLength = focalLength
						found = true
						break
					}
				}
				if !found {
					boxes[HASH] = append(box, lens{label: label, focalLength: focalLength})
				}
			}
			HASH += b
			HASH *= 17
		}

	}
	results.Part1 += int(HASH)

	for i, box := range boxes {
		for j, lens := range box {
			results.Part2 += (i + 1) * (j + 1) * int(lens.focalLength)
		}
	}
	println(results.Part1)
	println(results.Part2)
}
