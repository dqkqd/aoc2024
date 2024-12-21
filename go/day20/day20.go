package main

import (
	"fmt"
	"math"
	"strings"

	utils "example.com/aoc2024"
)

const Inf = 1 << 32

type Object rune

const (
	Wall  Object = '#'
	Empty Object = '.'
	Start Object = 'S'
	End   Object = 'E'
)

type Position struct {
	x, y int
}

type Map [][]Object

func (p Position) nextPosition() []Position {
	return []Position{
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x, p.y - 1},
		{p.x, p.y + 1},
	}
}

func (p Position) inside(m Map) bool {
	return p.x >= 0 && p.x < m.height() && p.y >= 0 && p.y < m.width()
}

func readMap() Map {
	reader := utils.Read(20, false)
	mapObjects := [][]Object{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			objects := make([]Object, len(line))
			for i := range line {
				objects[i] = Object(line[i])
			}
			mapObjects = append(mapObjects, objects)
		}
	}
	return mapObjects
}

func (m Map) height() int {
	return len(m)
}

func (m Map) width() int {
	return len(m[0])
}

func (m Map) scoreMapInit() [][]int {
	score := make([][]int, m.height())
	scoreFlat := make([]int, m.height()*m.width())
	for i := range scoreFlat {
		scoreFlat[i] = Inf
	}
	for i := range score {
		score[i], scoreFlat = scoreFlat[:m.width()], scoreFlat[m.width():]
	}
	return score
}

func (m Map) shortestPath(score [][]int, to Position) [][]int {
	var dfs func(Position)
	dfs = func(p Position) {
		for _, np := range p.nextPosition() {
			if np.inside(m) && m[np.x][np.y] != Wall && score[p.x][p.y]+1 < score[np.x][np.y] {
				score[np.x][np.y] = score[p.x][p.y] + 1
				dfs(np)
			}
		}
	}
	score[to.x][to.y] = 0
	dfs(to)
	return score
}

func (m Map) findObject(o Object) Position {
	for x, line := range m {
		for y, c := range line {
			if c == o {
				return Position{x, y}
			}
		}
	}
	panic("unreachable")
}

func Part1() {
	m := readMap()

	start := m.findObject(Start)
	end := m.findObject(End)

	shortestToStart := m.shortestPath(m.scoreMapInit(), start)
	shortestToEnd := m.shortestPath(m.scoreMapInit(), end)

	before := shortestToEnd[start.x][start.y]
	ans := 0
	for x := range m.height() {
		for y := range m.width() {
			p := Position{x, y}
			if m[p.x][p.y] != Wall {
				for _, np1 := range p.nextPosition() {
					if np1.inside(m) {
						for _, np2 := range np1.nextPosition() {
							if np2.inside(m) && m[np2.x][np2.y] != Wall {
								cheatScore := shortestToStart[p.x][p.y] + shortestToEnd[np2.x][np2.y] + 2
								if cheatScore+100 <= before {
									ans++
								}
							}
						}
					}
				}
			}
		}
	}

	fmt.Println(ans)
}

func Part2() {
	m := readMap()

	start := m.findObject(Start)
	end := m.findObject(End)

	shortestToStart := m.shortestPath(m.scoreMapInit(), start)
	shortestToEnd := m.shortestPath(m.scoreMapInit(), end)

	before := shortestToEnd[start.x][start.y]
	ans := 0
	for x := range m.height() {
		for y := range m.width() {
			if m[x][y] != Wall {
				for x1 := range m.height() {
					for y1 := range m.width() {
						d := int(math.Abs(float64(x1-x)) + math.Abs(float64(y1-y)))
						if d <= 20 {
							if m[x1][y1] != Wall {
								cheatScore := shortestToStart[x][y] + shortestToEnd[x1][y1] + d
								if cheatScore+100 <= before {
									ans++
								}
							}
						}
					}
				}
			}
		}
	}

	fmt.Println(ans)
}

func main() {
	Part1()
	Part2()
}
