package main

import (
	"fmt"
	"strings"

	utils "example.com/aoc2024"
)

func read() (int, map[string]bool, []string) {
	reader := utils.Read(19, false)

	patterns := map[string]bool{}
	maximum := 0

	line, err := reader.ReadString('\n')
	utils.CheckError(err)
	if err != nil {
		panic("...")
	}

	line = strings.TrimSpace(line)
	for _, p := range strings.Split(line, ",") {
		p = strings.TrimSpace(p)
		maximum = max(maximum, len(p))
		patterns[p] = true
	}

	reader.ReadString('\n')

	designs := []string{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			designs = append(designs, line)
		}
	}

	return maximum, patterns, designs
}

func Part1() {
	maximum, patterns, designs := read()

	good := func(line string) bool {
		goodPos := make([]bool, len(line)+1)
		goodPos[0] = true

		for i := range line {
			c := ""
			for j := i; j >= 0 && len(c) <= maximum; j-- {
				c = string(line[j]) + c
				if goodPos[j] && patterns[c] {
					goodPos[i+1] = true
					break
				}
			}
		}

		return goodPos[len(line)]
	}

	ans := 0
	for _, d := range designs {
		if good(d) {
			ans++
		}
	}

	fmt.Println(ans)
}

func Part2() {
	maximum, patterns, designs := read()

	ans := 0
	search := func(line string) {
		cnt := make([]int, len(line)+1)
		cnt[0] = 1

		for i := range line {
			c := ""
			for j := i; j >= 0 && len(c) <= maximum; j-- {
				c = string(line[j]) + c
				if patterns[c] {
					cnt[i+1] += cnt[j]
				}
			}
		}

		ans += cnt[len(line)]
	}

	for _, d := range designs {
		search(d)
	}
	fmt.Println(ans)
}

func main() {
	Part1()
	Part2()
}
