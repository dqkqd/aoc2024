package main

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	utils "example.com/aoc2024"
)

const (
	A = 0
	B = 1
	C = 2
)

type Opcode int

const (
	Adv Opcode = iota
	Bxl
	Bst
	Jnz
	Bxc
	Out
	Bdv
	Cdv
)

type Program struct {
	instructions []int
	outs         []int
	regs         [3]int
	pc           int
}

func readProgram() Program {
	reader := utils.Read(17, false)

	toInt := func(s string) int {
		v, e := strconv.Atoi(strings.TrimSpace(s))
		utils.CheckError(e)
		return v
	}

	program := Program{pc: 0, outs: []int{}}

	// read register
	for i := range program.regs {
		r, e := reader.ReadString('\n')
		utils.CheckError(e)

		s := strings.Split(r, ":")
		program.regs[i] = toInt(s[1])
	}
	// read empty string
	reader.ReadString('\n')

	line, e := reader.ReadString('\n')
	utils.CheckError(e)

	ins := strings.Split(line, ":")[1]
	ins = strings.TrimSpace(ins)

	instructions := strings.Split(ins, ",")
	program.instructions = make([]int, len(instructions))
	for i, ins := range instructions {
		program.instructions[i] = toInt(ins)
	}

	return program
}

func (p Program) opcode() Opcode {
	return Opcode(p.instructions[p.pc])
}

func (p Program) operand() int {
	return p.instructions[p.pc+1]
}

func (p Program) literal() int {
	return p.operand()
}

func (p Program) combo() int {
	switch p.operand() {
	case 0, 1, 2, 3:
		return p.operand()
	case 4:
		return p.regs[A]
	case 5:
		return p.regs[B]
	case 6:
		return p.regs[C]
	default:
		panic("unreachable")
	}
}

func (p *Program) run() bool {
	if p.pc+1 >= len(p.instructions) {
		return false
	}
	switch p.opcode() {
	case Adv:
		p.regs[A] = p.regs[A] >> p.combo()
		p.pc += 2
	case Bdv:
		p.regs[B] = p.regs[A] >> p.combo()
		p.pc += 2
	case Cdv:
		p.regs[C] = p.regs[A] >> p.combo()
		p.pc += 2

	case Bxl:
		p.regs[B] ^= p.literal()
		p.pc += 2
	case Bst:
		p.regs[B] = p.combo() % 8
		p.pc += 2
	case Jnz:
		if p.regs[A] == 0 {
			p.pc += 2
		} else {
			p.pc = p.literal()
		}
	case Bxc:
		p.regs[B] ^= p.regs[C]
		p.pc += 2
	case Out:
		p.outs = append(p.outs, p.combo()%8)
		p.pc += 2
	default:
		panic("unreachable")
	}

	return true
}

func Part1() {
	program := readProgram()
	for program.run() {
	}

	outs := []string{}
	for _, o := range program.outs {
		outs = append(outs, strconv.Itoa(o))
	}
	out := strings.Join(outs, ",")
	fmt.Println(out)
}

func Part2() {
	program := readProgram()

	expected := program.instructions

	on := [48]bool{}
	good := []int{}

	bitAt := func(num, pos int) int {
		return (num >> pos) & 1
	}

	setBitAt := func(num, pos, b int) int {
		if bitAt(num, pos) != b {
			num ^= (1 << pos)
		}
		return num
	}

	var dfs func(int, int)

	dfs = func(i, a int) {
		if i >= 48 {
			b := a >> (i - 3)
			if ((b&7)^(b>>(7-(b&7))))&7 == expected[(i-1)/3] {
				good = append(good, a)
			}
			return
		}

		// b = (a % 8), c = 0
		// b = 7 - (a % 8), c = 0
		// b = 7 - (a % 8), c = a >> (7 - (a % 8))
		// b = (a % 8), c = a >> (7 - (a % 8))
		// b = (a % 8) ^ (a >> (7 - (a % 8))), c = a >> (7 - (a % 8))
		// a = a >> 3

		turned := [3]bool{}

		var b, m int

		// -------------------------i-xx
		//                    7 - m  i - 3
		if i%3 == 0 && i > 0 {
			// check for (i - 3, i - 2, i - 1)
			b = a >> (i - 3)
			m = b & 7
			// um := 7 - m // (i - 3 + 7 - m) = (i + 4 - m)
			// (m ^ (b >> um)) & 7 == expected[(i - 1) / 3]

			// these should match
			// (i + 4 - m, i + 5 - m, i + 6 - m)
			// (i - 3, i - 2, i - 1)

			need := m ^ expected[(i-1)/3]

			for c := 0; c <= 2; c++ {
				index := i + 4 + c - m
				if (index >= 48 || on[index]) && bitAt(need, c) != bitAt(a, index) {
					return
				}
			}

			for c := 0; c <= 2; c++ {
				index := i + 4 + c - m
				if index < 48 && !on[index] {
					on[index] = true
					a = setBitAt(a, index, bitAt(need, c))
					turned[c] = true
				}
			}
		}

		if on[i] {
			dfs(i+1, a)
		} else {
			on[i] = true
			dfs(i+1, setBitAt(a, i, 0))
			dfs(i+1, setBitAt(a, i, 1))
			on[i] = false
		}

		for c := range turned {
			if turned[c] {
				index := i + 4 + c - m
				on[index] = false
			}
		}
	}

	dfs(0, 0)

	runWithA := func(a int) Program {
		p := Program{instructions: program.instructions, regs: program.regs, pc: program.pc}
		p.regs[A] = a
		for p.run() {
		}
		return p
	}

	verify := func(a int) {
		p := runWithA(a)
		if !reflect.DeepEqual(p.instructions, p.outs) {
			panic("...")
		}
	}

	for _, g := range good {
		verify(g)
	}

	slices.Sort(good)
	fmt.Println(good[0])
}

func main() {
	Part1()
	Part2()
}
