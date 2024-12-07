package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	utils "example.com/aoc2024"
)

type (
	Row int
	Col int
)

type Direction struct {
	dx Row
	dy Col
}

var (
	Up    = Direction{-1, 0}
	Down  = Direction{1, 0}
	Left  = Direction{0, -1}
	Right = Direction{0, 1}
)

type Position struct {
	x         Row
	y         Col
	direction Direction
}

type Point struct {
	x Row
	y Col
}

type Map struct {
	verticalObstacles   map[Col][]Row
	horizontalObstacles map[Row][]Col
	data                [][]byte
	start               Position
	h                   Row
	w                   Col
}

func (direction Direction) Turn() Direction {
	switch direction {
	case Up:
		return Right
	case Right:
		return Down
	case Down:
		return Left
	case Left:
		return Up
	}
	// do not turn
	return direction
}

func (p Position) Move() Position {
	return Position{
		p.x + p.direction.dx,
		p.y + p.direction.dy,
		p.direction,
	}
}

func (p Position) Turn() Position {
	return Position{
		p.x,
		p.y,
		p.direction.Turn(),
	}
}

func (p Position) Point() Point {
	return Point{p.x, p.y}
}

func IsObstacle(c byte) bool {
	return c == '#'
}

func OutOfBound(p Position, m Map) bool {
	return p.x < 0 || p.y < 0 || p.x >= m.h || p.y >= m.w
}

func IsValid(p Position, m Map) bool {
	return !OutOfBound(p, m) && !IsObstacle(m.data[p.x][p.y])
}

func (m Map) VerticalObstacles(y Col) *[]Row {
	_, hasObstacles := m.verticalObstacles[y]
	if !hasObstacles {
		m.verticalObstacles[y] = []Row{}
	}
	v := m.verticalObstacles[y]
	return &v
}

func (m Map) HorizontalObstacles(x Row) *[]Col {
	_, hasObstacles := m.horizontalObstacles[x]
	if !hasObstacles {
		m.horizontalObstacles[x] = []Col{}
	}
	v := m.horizontalObstacles[x]
	return &v
}

func (m *Map) MarkObstacle(x Row, y Col) {
	if IsObstacle(m.data[x][y]) {
		panic("already obstacle")
	}
	m.data[x][y] = '#'

	verticalObstacles := m.VerticalObstacles(y)
	horizontalObstacles := m.HorizontalObstacles(x)

	*verticalObstacles = append(*verticalObstacles, x)
	*horizontalObstacles = append(*horizontalObstacles, y)

	sort.Slice(*verticalObstacles, func(i int, j int) bool {
		return (*verticalObstacles)[i] < (*verticalObstacles)[j]
	})
	sort.Slice(*horizontalObstacles, func(i int, j int) bool {
		return (*horizontalObstacles)[i] < (*horizontalObstacles)[j]
	})

	m.verticalObstacles[y] = *verticalObstacles
	m.horizontalObstacles[x] = *horizontalObstacles
}

func (m *Map) UnMarkObstacle(x Row, y Col) {
	if !IsObstacle(m.data[x][y]) {
		panic("not an obstacle")
	}
	m.data[x][y] = '.'

	verticalObstacles := m.VerticalObstacles(y)
	horizontalObstacles := m.HorizontalObstacles(x)

	newVerticleObstacles := []Row{}
	for _, e := range *verticalObstacles {
		if e != x {
			newVerticleObstacles = append(newVerticleObstacles, e)
		}
	}
	m.verticalObstacles[y] = newVerticleObstacles

	newHorizontalObstacles := []Col{}
	for _, e := range *horizontalObstacles {
		if e != y {
			newHorizontalObstacles = append(newHorizontalObstacles, e)
		}
	}
	m.horizontalObstacles[x] = newHorizontalObstacles
}

func (m Map) Run() []Point {
	pos := m.start
	positions := map[Position]bool{}
	for {
		positions[pos] = true

		nextPos := pos.Move()
		if OutOfBound(nextPos, m) {
			break
		}
		if IsValid(nextPos, m) {
			pos = nextPos
		} else {
			pos = pos.Turn()
		}
	}

	mapPoints := map[Point]bool{}
	for position := range positions {
		mapPoints[position.Point()] = true
	}

	points := []Point{}
	for p := range mapPoints {
		points = append(points, p)
	}

	return points
}

func (m Map) IsLoop() bool {
	pos := m.start
	traces := map[Position]bool{}

	for {
		_, exist := traces[pos]
		if exist {
			return true
		}
		traces[pos] = true
		nextPos, err := MoveUntilBlocked(pos, m)
		if err != nil {
			// outside
			return false
		}
		pos = nextPos.Turn()
	}
}

func MoveUntilBlocked(p Position, m Map) (Position, error) {
	switch p.direction {
	case Up:
		obstacles := m.VerticalObstacles(p.y)
		index := sort.Search(len(*obstacles), func(i int) bool {
			return (*obstacles)[i] >= p.x
		})
		index--
		if index >= 0 {
			row := (*obstacles)[index] + 1
			return Position{row, p.y, p.direction}, nil
		}
	case Right:
		obstacles := m.HorizontalObstacles(p.x)
		index := sort.Search(len(*obstacles), func(i int) bool {
			return (*obstacles)[i] >= p.y
		})
		if index < len(*obstacles) {
			col := (*obstacles)[index] - 1
			return Position{p.x, col, p.direction}, nil
		}
	case Down:
		obstacles := m.VerticalObstacles(p.y)
		index := sort.Search(len(*obstacles), func(i int) bool {
			return (*obstacles)[i] >= p.x
		})
		if index < len(*obstacles) {
			row := (*obstacles)[index] - 1
			return Position{row, p.y, p.direction}, nil
		}
	case Left:
		obstacles := m.HorizontalObstacles(p.x)
		index := sort.Search(len(*obstacles), func(i int) bool {
			return (*obstacles)[i] >= p.y
		})
		index--
		if index >= 0 {
			col := (*obstacles)[index] + 1
			return Position{p.x, col, p.direction}, nil
		}
	}
	return Position{-1, -1, p.direction}, errors.New("Out of bound")
}

func ReadMap() Map {
	reader := utils.Read(6, false)
	data := [][]byte{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		data = append(data, []byte(strings.TrimSpace(line)))
	}

	start := Position{-1, -1, Up}
	verticalObstacles := map[Col][]Row{}
	horizontalObstacles := map[Row][]Col{}

	for x, line := range data {
		for y, char := range line {
			r := Row(x)
			c := Col(y)

			switch char {
			case '^':
				start.x = r
				start.y = c
			case '#':
				verticalObstacles[c] = append(verticalObstacles[c], r)
				horizontalObstacles[r] = append(horizontalObstacles[r], c)
			}
		}
	}

	return Map{
		verticalObstacles,
		horizontalObstacles,
		data,
		start,
		Row(len(data)),
		Col(len(data[0])),
	}
}

func Part1() {
	myMap := ReadMap()
	points := myMap.Run()
	fmt.Println(len(points))
}

func Part2() {
	myMap := ReadMap()

	points := myMap.Run()
	ans := 0

	for _, p := range points {
		if p.x == myMap.start.x && p.y == myMap.start.y {
			continue
		}
		myMap.MarkObstacle(p.x, p.y)
		if myMap.IsLoop() {
			ans++
		}
		myMap.UnMarkObstacle(p.x, p.y)
	}

	fmt.Println(ans)
}

func main() {
	Part1()
	Part2()
}
