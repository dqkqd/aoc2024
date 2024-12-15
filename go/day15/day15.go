package main

import (
	"errors"
	"fmt"
	"strings"

	utils "example.com/aoc2024"
)

type Move int

const (
	Up    Move = '^'
	Down  Move = 'v'
	Left  Move = '<'
	Right Move = '>'
)

var (
	MoveUp    []int = []int{-1, 0}
	MoveDown  []int = []int{1, 0}
	MoveLeft  []int = []int{0, -1}
	MoveRight []int = []int{0, 1}
)

func direction(move Move) []int {
	switch move {
	case Up:
		return MoveUp
	case Down:
		return MoveDown
	case Left:
		return MoveLeft
	case Right:
		return MoveRight
	}
	panic("Unreachable")
}

func nextPosition(x int, y int, move Move) (int, int) {
	d := direction(move)
	return x + d[0], y + d[1]
}

func prevPosition(x int, y int, move Move) (int, int) {
	d := direction(move)
	return x - d[0], y - d[1]
}

type Object rune

const (
	Robot Object = '@'
	Dot   Object = '.'
	Wall  Object = '#'
	Box   Object = 'O'
	LBox  Object = '['
	RBox  Object = ']'
)

type Map [][]Object

type Moves []Move

func readData() (Map, Moves) {
	reader := utils.Read(15, false)

	mapObjects := [][]Object{}

	// read map
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			break
		}

		objects := []Object{}

		for _, c := range line {
			objects = append(objects, Object(c))
		}

		mapObjects = append(mapObjects, objects)
	}

	moves := Moves{}
	// read moves
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		for _, c := range strings.TrimSpace(line) {
			moves = append(moves, Move(c))
		}
	}

	return mapObjects, Moves(moves)
}

func (m Map) toString() string {
	s := ""
	for _, line := range m {
		for _, c := range line {
			s += string(c)
		}
		s += "\n"
	}
	return s
}

func (m Map) robotLocation() (int, int) {
	for x, objects := range m {
		for y, o := range objects {
			if o == Robot {
				return x, y
			}
		}
	}
	panic("Unreachable")
}

func (m Map) get(x int, y int) (Object, error) {
	if x < 0 || x >= len(m) || y < 0 || y > len(m[0]) {
		return Wall, errors.New("Out of bound")
	}
	return m[x][y], nil
}

func (m Map) movedPositions(robotX int, robotY int, move Move) (movedPositions [][]int, movable bool) {
	positions := [][]int{}
	x, y := robotX, robotY
	for {
		x, y = nextPosition(x, y, move)
		positions = append(positions, []int{x, y})

		o, err := m.get(x, y)
		if err != nil || o == Wall {
			return [][]int{}, false
		}

		if o == Dot {
			return positions, true
		}
	}
}

func (m *Map) move(positions [][]int, move Move) {
	w := len((*m)[0])

	// store old objects and remove old object in the real map
	oldData := map[int]Object{}
	for _, p := range positions {
		x, y := prevPosition(p[0], p[1], move)

		k := x*w + y
		o, e := m.get(x, y)
		utils.CheckError(e)
		oldData[k] = o

		(*m)[x][y] = Dot
	}

	for _, p := range positions {
		x, y := prevPosition(p[0], p[1], move)
		k := x*w + y
		(*m)[p[0]][p[1]] = oldData[k]
	}
}

func (m *Map) move1(robotX *int, robotY *int, move Move) {
	positions, movable := m.movedPositions(*robotX, *robotY, move)
	if movable {
		m.move(positions, move)
		*robotX, *robotY = nextPosition(*robotX, *robotY, move)
	}
}

func Part1() {
	mapObjects, moves := readData()
	x, y := mapObjects.robotLocation()
	for _, move := range moves {
		mapObjects.move1(&x, &y, move)
	}

	ans := 0
	for i, objects := range mapObjects {
		for j, o := range objects {
			if o == Box {
				ans += 100*i + j
			}
		}
	}
	fmt.Println(ans)
}

func (m *Map) extend() {
	h := len(*m)

	extend := map[Object][]Object{
		Robot: {Robot, Dot},
		Wall:  {Wall, Wall},
		Dot:   {Dot, Dot},
		Box:   {LBox, RBox},
	}

	for x := range h {
		line := []Object{}
		for _, o := range (*m)[x] {
			no := extend[o]
			line = append(line, no[0])
			line = append(line, no[1])
		}
		(*m)[x] = line
	}
}

func (m Map) movedPositions2(robotX int, robotY int, move Move) (movedPositions [][]int, movable bool) {
	positions := [][]int{}
	visited := map[int]int{}

	w := len(m[0])

	movable = true

	var dfs func(int, int)
	dfs = func(x int, y int) {
		if !movable {
			return
		}

		k := x*w + y
		_, ok := visited[k]
		if ok {
			return
		}
		visited[k] = 1
		positions = append(positions, []int{x, y})

		o, err := m.get(x, y)
		if err != nil || o == Wall {
			// cannot move
			movable = false
			return
		}

		switch move {
		case Up, Down:
			if o != Dot {
				nx, ny := nextPosition(x, y, move)
				dfs(nx, ny)
				switch o {
				case LBox:
					dfs(nx, ny+1)
				case RBox:
					dfs(nx, ny-1)
				default:
					panic("Unreachable")
				}
			}
		case Left, Right:
			if o != Dot {
				nx, ny := nextPosition(x, y, move)
				dfs(nx, ny)
			}
		default:
			panic(fmt.Sprintf("unexpected main.Move: %#v", move))
		}
	}

	x, y := nextPosition(robotX, robotY, move)
	dfs(x, y)

	return positions, movable
}

func (m *Map) move2(robotX *int, robotY *int, move Move) {
	positions, movable := m.movedPositions2(*robotX, *robotY, move)
	if movable {
		m.move(positions, move)
		*robotX, *robotY = nextPosition(*robotX, *robotY, move)
	}
}

func Part2() {
	mapObjects, moves := readData()
	mapObjects.extend()

	x, y := mapObjects.robotLocation()
	for _, move := range moves {
		mapObjects.move2(&x, &y, move)
	}

	ans := 0
	for i, objects := range mapObjects {
		for j, o := range objects {
			if o == LBox {
				ans += 100*i + j
			}
		}
	}
	fmt.Println(ans)
}

func main() {
	Part1()
	Part2()
}
