// This is a second degree equation:
// y(x) = x (t - x) = -x² + tx
// -x² + tx - y = 0

// If y is a constant (like in the input), we have:
// x_min = (t - sqrt(t² - 4y)) / 2
// x_max = (t + sqrt(t² - 4y)) / 2
// range = trunc(x_max) - trunc(x_min)

// Technically, if x_min and x_max are integers, there are edge cases.
// However, none of them showed up in my input.
package main

import (
	"aoc/2023/common"
	"bufio"
	"fmt"
	"math"
	"strings"
)

type race struct {
	time, distance int
}

func main() {
	input := common.Open("input")
	races, theRace := parseTimeAndDist(bufio.NewScanner(input))
	input.Close()

	racesMargin := 1
	for _, race := range races {
		racesMargin *= marginOfError(race)
	}
	theMargin := marginOfError(theRace)

	fmt.Println("Margin of error (part 1):", racesMargin)
	fmt.Println("Margin of error (part 2):", theMargin)
}

func marginOfError(r race) int {
	xMin := (float64(r.time) - math.Sqrt(math.Pow(float64(r.time), 2)-4*float64(r.distance))) / 2
	xMax := (float64(r.time) + math.Sqrt(math.Pow(float64(r.time), 2)-4*float64(r.distance))) / 2
	return int(xMax) - int(xMin)
}

func parseTimeAndDist(scanner *bufio.Scanner) ([]race, race) {
	scanner.Scan()
	timesStr := strings.Fields(scanner.Text())[1:]
	races := make([]race, len(timesStr))
	var theRace race
	for i, t := range timesStr {
		races[i].time = common.Atoi(t)
	}
	theRace.time = common.Atoi(strings.Join(timesStr, ""))
	scanner.Scan()
	distStr := strings.Fields(scanner.Text())[1:]
	for i, d := range distStr {
		races[i].distance = common.Atoi(d)
	}
	theRace.distance = common.Atoi(strings.Join(distStr, ""))
	return races, theRace
}
