package main

import (
	"fmt"
	"math"

	utils "example.com/aoc2024"
)

type IsSafe func([]int) bool

func eval(f func(lhs int, rhs int) bool) func([]int) bool {
	return func(levels []int) bool {
		for i := 0; i+1 < len(levels); i++ {
			if !f(levels[i], levels[i+1]) {
				return false
			}
		}
		return true
	}
}

func IsDiff(a int, b int) bool {
	diff := math.Abs(float64(a - b))
	return diff >= 1 && diff <= 3
}

func IsIncr(a int, b int) bool {
	return a < b
}

func IsDecr(a int, b int) bool {
	return a > b
}

func isSafe(levels []int) bool {
	isDiff := eval(IsDiff)
	isIncr := eval(IsIncr)
	isDecr := eval(IsDecr)
	return isDiff(levels) && (isIncr(levels) || isDecr(levels))
}

func isSafeRm(levels []int) bool {
	if isSafe(levels) {
		return true
	}

	for i := 0; i < len(levels); i++ {
		arr := []int{}

		for j := 0; j < len(levels); j++ {
			if i != j {
				arr = append(arr, levels[j])
			}
		}

		if isSafe(arr) {
			return true
		}
	}

	return false
}

func Run(f IsSafe) {
	reader := utils.Read(2, false)
	arrayReader := utils.ArrayReader(reader)

	safeCount := 0
	for {
		levels, err := arrayReader()
		if err != nil {
			break
		}
		if f(levels) {
			safeCount += 1
		}
	}
	fmt.Println(safeCount)
}

func Part1() {
	Run(isSafe)
}

func Part2() {
	Run(isSafeRm)
}

func main() {
	Part1()
	Part2()
}
