package main

import (
	"fmt"
	"strings"

	utils "example.com/aoc2024"
)

func ReadMap() []string {
	words := []string{}
	buf := utils.Read(4, false)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		if len(line) > 0 {
			words = append(words, line)
		}
	}
	return words
}

func ReadSquare(words []string, x int, y int, size int) []string {
	h := len(words)
	w := len(words[0])

	out := []string{}

	ly, ry := y, min(w, y+size)
	lx, rx := x, min(h, x+size)

	for i := lx; i < rx; i++ {
		out = append(out, words[i][ly:ry])
	}
	return out
}

func XMasCount(square []string) int {
	count := 0

	if len(square) == 0 {
		return count
	}

	if square[0] == "XMAS" {
		count++
	}
	if square[0] == "SAMX" {
		count++
	}

	if len(square) >= 4 {
		if square[0][0] == 'X' && square[1][0] == 'M' && square[2][0] == 'A' && square[3][0] == 'S' {
			count++
		}
		if square[3][0] == 'X' && square[2][0] == 'M' && square[1][0] == 'A' && square[0][0] == 'S' {
			count++
		}
	}

	if len(square) >= 4 && len(square[0]) >= 4 {
		if square[0][0] == 'X' && square[1][1] == 'M' && square[2][2] == 'A' && square[3][3] == 'S' {
			count++
		}
		if square[3][3] == 'X' && square[2][2] == 'M' && square[1][1] == 'A' && square[0][0] == 'S' {
			count++
		}
		if square[3][0] == 'X' && square[2][1] == 'M' && square[1][2] == 'A' && square[0][3] == 'S' {
			count++
		}
		if square[0][3] == 'X' && square[1][2] == 'M' && square[2][1] == 'A' && square[3][0] == 'S' {
			count++
		}
	}

	return count
}

func IsMas(square []string) bool {
	if len(square) != 3 {
		return false
	}

	for _, s := range square {
		if len(s) != 3 {
			return false
		}
	}

	if square[1][1] != 'A' {
		return false
	}

	positions := [][]int{
		{0, 0},
		{0, 2},
		{2, 2},
		{2, 0},
		{0, 0},
		{0, 2},
		{2, 2},
		{2, 0},
	}

	for offset := 0; offset < 4; offset++ {
		acc := []byte{}
		for i := 0; i < 4; i++ {
			x, y := positions[offset+i][0], positions[offset+i][1]
			acc = append(acc, square[x][y])
		}
		if string(acc) == "MMSS" {
			return true
		}
	}

	return false
}

func Part1() {
	words := ReadMap()
	h := len(words)
	w := len(words[0])

	cnt := 0
	for x := 0; x < h; x++ {
		for y := 0; y < w; y++ {
			square := ReadSquare(words, x, y, 4)
			cnt += XMasCount(square)
		}
	}
	fmt.Println(cnt)
}

func Part2() {
	words := ReadMap()
	h := len(words)
	w := len(words[0])

	cnt := 0
	for x := 0; x < h; x++ {
		for y := 0; y < w; y++ {
			square := ReadSquare(words, x, y, 3)
			if IsMas(square) {
				cnt++
			}
		}
	}
	fmt.Println(cnt)
}

func main() {
	Part1()
	Part2()
}
