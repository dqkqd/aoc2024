package main

import (
	"fmt"
	"strconv"
	"strings"

	utils "example.com/aoc2024"
)

const Inf = 1 << 32

type Object int

const (
	Wall Object = iota
	Empty
)

type Position struct {
	x, y int
}

func (p Position) nextPosition() []Position {
	return []Position{
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x, p.y - 1},
		{p.x, p.y + 1},
	}
}

func (p Position) inside(n int) bool {
	return p.x >= 0 && p.x < n && p.y >= 0 && p.y < n
}

func (p Position) to1d(n int) int {
	return p.x*n + p.y
}

func readPoints() []Position {
	reader := utils.Read(18, false)
	points := []Position{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			s := strings.Split(line, ",")
			a, err := strconv.Atoi(s[0])
			utils.CheckError(err)
			b, err := strconv.Atoi(s[1])
			utils.CheckError(err)
			points = append(points, Position{a, b})
		}
	}
	return points
}

func Part1() {
	steps := 1024
	n := 71

	scores := make([]int, n*n)
	for i := range n * n {
		scores[i] = Inf
	}

	objects := make([]Object, n*n)
	for i := range n * n {
		objects[i] = Empty
	}
	points := readPoints()

	for i := range steps {
		s := points[i].to1d(n)
		objects[s] = Wall
	}

	start := Position{0, 0}
	scores[start.to1d(n)] = 0

	var dfs func(Position)
	dfs = func(p Position) {
		sp := scores[p.to1d(n)]
		for _, np := range p.nextPosition() {
			if np.inside(n) {
				s := np.to1d(n)
				if objects[s] == Empty {
					if scores[s] > sp+1 {
						scores[s] = sp + 1
						dfs(np)
					}
				}
			}
		}
	}

	dfs(start)

	end := Position{n - 1, n - 1}

	fmt.Println(scores[end.to1d(n)])
}

type DSU []int

func newDSU(n int) DSU {
	dsu := make([]int, n)
	for i := range n {
		dsu[i] = i
	}
	return DSU(dsu)
}

func (d DSU) merge(a, b int) {
	pa := d.parent(a)
	pb := d.parent(b)
	if pa > pb {
		pa, pb = pb, pa
	}
	d[pb] = pa
}

func (d DSU) parent(a int) int {
	for {
		p := d[d[a]]
		if p == a {
			break
		}
		a = p
	}
	return a
}

func (d DSU) same(a, b int) bool {
	return d.parent(a) == d.parent(b)
}

func Part2() {
	n := 71

	objects := make([]Object, n*n)
	for i := range n * n {
		objects[i] = Empty
	}

	points := readPoints()
	for _, p := range points {
		objects[p.to1d(n)] = Wall
	}

	colors := make([]int, n*n)
	visited := map[int]bool{}

	dsu := newDSU(n*n + 10)

	color := 0

	var dfs func(Position)
	dfs = func(p Position) {
		sp := p.to1d(n)
		if objects[sp] != Empty || visited[sp] {
			return
		}
		visited[sp] = true
		colors[sp] = color

		for _, np := range p.nextPosition() {
			sn := np.to1d(n)
			if np.inside(n) && objects[sn] == Empty {
				if visited[sn] {
					dsu.merge(colors[sn], colors[sp])
				} else {
					colors[sn] = colors[sp]
					dfs(np)
				}
			}
		}
	}

	for x := range n {
		for y := range n {
			p := Position{x, y}
			if !visited[p.to1d(n)] && objects[p.to1d(n)] == Empty {
				color += 1
				dfs(p)
			}
		}
	}

	start := Position{0, 0}.to1d(n)
	end := Position{n - 1, n - 1}.to1d(n)
	if dsu.same(colors[start], colors[end]) {
		panic("...")
	}

	for i := len(points) - 1; i >= 0; i-- {
		p := points[i]
		sp := p.to1d(n)
		objects[sp] = Empty

		if !visited[sp] {
			color += 1
			dfs(p)
		}

		if dsu.same(colors[start], colors[end]) {
			fmt.Printf("%d,%d\n", p.x, p.y)
			return
		}
	}
}

func main() {
	Part1()
	Part2()
}
