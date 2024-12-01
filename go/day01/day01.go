package main

import (
	"fmt"
	"math"
	"slices"

	utils "example.com/aoc2024"
)

func Read() ([]int, []int) {
	reader := utils.Read(1, false)
	arrayReader := utils.ArrayReader(reader)

	arr1 := []int{}
	arr2 := []int{}

	for {
		nums, err := arrayReader()
		if err != nil {
			break
		}

		arr1 = append(arr1, nums[0])
		arr2 = append(arr2, nums[1])
	}

	return arr1, arr2
}

func Part1() {
	arr1, arr2 := Read()
	slices.Sort(arr1)
	slices.Sort(arr2)
	sum := 0
	for i := range len(arr1) {
		sum += int(math.Abs(float64(arr1[i] - arr2[i])))
	}
	fmt.Println(sum)
}

func Part2() {
	arr1, arr2 := Read()
	counter := make(map[int]int)
	for _, elem := range arr2 {
		counter[elem] += 1
	}

	sum := 0
	for _, elem := range arr1 {
		sum += elem * counter[elem]
	}

	fmt.Println(sum)
}

func main() {
	Part1()
	Part2()
}
