package main

import (
	"fmt"
	"strings"

	utils "example.com/aoc2024"
)

const Inf = 1 << 32

type (
	Object byte
	Map    [][]Object
)

const (
	Wall  Object = '#'
	Start Object = 'S'
	End   Object = 'E'
	Empty Object = '.'
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

var (
	MoveUp    = []int{-1, 0}
	MoveRight = []int{0, 1}
	MoveDown  = []int{1, 0}
	MoveLeft  = []int{0, -1}
)

type Position struct {
	x int
	y int
	d Direction
}

func readMap() Map {
	reader := utils.Read(16, false)
	mapObjects := [][]Object{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if len(line) > 0 {
			objects := []Object{}
			for _, b := range line {
				objects = append(objects, Object(b))
			}
			mapObjects = append(mapObjects, objects)
		}
	}
	return Map(mapObjects)
}

func (m Map) height() int {
	return len(m)
}

func (m Map) width() int {
	return len(m[0])
}

func (m Map) startPoint() Position {
	for x, line := range m {
		for y, b := range line {
			if b == Start {
				return Position{x, y, Right}
			}
		}
	}
	panic("Unreachable")
}

func (m Map) endPoints() []Position {
	for x, line := range m {
		for y, b := range line {
			if b == End {
				return []Position{
					{x, y, Up},
					{x, y, Right},
					{x, y, Down},
					{x, y, Left},
				}
			}
		}
	}
	panic("Unreachable")
}

func (p Position) next() Position {
	switch p.d {
	case Down:
		return Position{p.x + MoveDown[0], p.y + MoveDown[1], p.d}
	case Left:
		return Position{p.x + MoveLeft[0], p.y + MoveLeft[1], p.d}
	case Right:
		return Position{p.x + MoveRight[0], p.y + MoveRight[1], p.d}
	case Up:
		return Position{p.x + MoveUp[0], p.y + MoveUp[1], p.d}
	default:
		panic(fmt.Sprintf("unexpected main.Direction: %#v", p.d))
	}
}

func (p Position) turnRight() Position {
	d := (int(p.d) + 1) % 4
	return Position{p.x, p.y, Direction(d)}
}

func (p Position) turnLeft() Position {
	d := (int(p.d) + 3) % 4
	return Position{p.x, p.y, Direction(d)}
}

func (p Position) inside(m Map) bool {
	return p.x >= 0 && p.x < m.height() && p.y >= 0 && p.y < m.width()
}

func (p Position) encode(m Map) int {
	return m.width()*p.x + p.y
}

func (m Map) initScoreMap() [][][]int {
	scores := [][][]int{}
	for range 4 {
		score := [][]int{}
		for range m.height() {
			ss := []int{}
			for range m.width() {
				ss = append(ss, Inf)
			}
			score = append(score, ss)
		}
		scores = append(scores, score)
	}
	return scores
}

func (m Map) scoreMap(scores *[][][]int, startPositions []Position) {
	var dfs func(Position)

	dfs = func(p Position) {
		sp := (*scores)[int(p.d)][p.x][p.y]

		var np Position

		np = p.next()
		if np.inside(m) && m[np.x][np.y] != Wall && sp+1 < (*scores)[int(np.d)][np.x][np.y] {
			(*scores)[int(np.d)][np.x][np.y] = sp + 1
			dfs(np)
		}

		np = p.turnRight()
		if np.inside(m) && m[np.x][np.y] != Wall && sp+1000 < (*scores)[int(np.d)][np.x][np.y] {
			(*scores)[int(np.d)][np.x][np.y] = sp + 1000
			dfs(np)
		}

		np = p.turnLeft()
		if np.inside(m) && m[np.x][np.y] != Wall && sp+1000 < (*scores)[int(np.d)][np.x][np.y] {
			(*scores)[int(np.d)][np.x][np.y] = sp + 1000
			dfs(np)
		}
	}

	for _, p := range startPositions {
		dfs(p)
	}
}

func Part1() {
	m := readMap()
	scores := m.initScoreMap()

	s := m.startPoint()
	scores[int(s.d)][s.x][s.y] = 0
	m.scoreMap(&scores, []Position{s})

	ans := 1 << 32
	for _, e := range m.endPoints() {
		ans = min(ans, scores[int(e.d)][e.x][e.y])
	}

	fmt.Println(ans)
}

func Part2() {
	m := readMap()

	startScores := m.initScoreMap()

	s := m.startPoint()
	startScores[int(s.d)][s.x][s.y] = 0
	m.scoreMap(&startScores, []Position{s})
	minScore := 1 << 32
	for _, e := range m.endPoints() {
		minScore = min(minScore, startScores[int(e.d)][e.x][e.y])
	}

	endScores := m.initScoreMap()
	for _, e := range m.endPoints() {
		endScores[int(e.d)][e.x][e.y] = 0
	}
	m.scoreMap(&endScores, m.endPoints())

	ans := 0
	for x := range m.height() {
		for y := range m.width() {
			for d := range 4 {
				if startScores[d][x][y]+endScores[(d+2)%4][x][y] <= minScore {
					ans += 1
					break
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
