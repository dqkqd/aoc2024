package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	utils "example.com/aoc2024"
)

type Operation int

const (
	Add Operation = iota
	Mul
	Concat
)

type Equation struct {
	rhs []int
	lhs int
}

func binary_operators(size int) func(mask int) []Operation {
	return func(mask int) []Operation {
		ops := []Operation{}
		for b := range size {
			var op Operation
			if (mask & (1 << b)) > 0 {
				op = Add
			} else {
				op = Mul
			}
			ops = append(ops, op)
		}
		return ops
	}
}

func ternary_operators(size int) func(mask int) []Operation {
	return func(mask int) []Operation {
		ops := []Operation{}
		for b := range size {
			pow := int(math.Pow(float64(3), float64(b)))
			s := (mask / pow) % 3
			var op Operation
			switch s {
			case 0:
				op = Add
			case 1:
				op = Mul
			default:
				op = Concat
			}
			ops = append(ops, op)
		}
		return ops
	}
}

func (op Operation) Eval(lhs int, rhs int) int {
	switch op {
	case Add:
		return lhs + rhs
	case Mul:
		return lhs * rhs
	case Concat:
		log := int(math.Floor(math.Log10(float64(rhs))))
		pow := int(math.Pow(10, float64(log)+1))
		return lhs*pow + rhs
	}
	panic("unreachable")
}

func EquationFromString(s string) (Equation, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return Equation{[]int{}, 0}, errors.New("invalid string")
	}
	splitted := strings.Split(s, ":")

	lhs, err := strconv.Atoi(splitted[0])
	utils.CheckError(err)

	rhs_str := strings.Split(splitted[1], " ")
	rhs := []int{}
	for _, item := range rhs_str {
		v, err := strconv.Atoi(item)
		if err == nil {
			rhs = append(rhs, v)
		}
	}

	return Equation{rhs, lhs}, nil
}

func (equation Equation) Eval(ops []Operation) int {
	value := equation.rhs[0]
	for i, op := range ops {
		value = op.Eval(value, equation.rhs[i+1])
	}
	return value
}

func ReadEquations() []Equation {
	reader := utils.Read(7, false)
	equations := []Equation{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		equation, err := EquationFromString(line)
		if err == nil {
			equations = append(equations, equation)
		}
	}
	return equations
}

func Part1() {
	equations := ReadEquations()
	ans := 0
	for _, equation := range equations {
		gen := binary_operators(len(equation.rhs) - 1)
		for mask := range 1<<len(equation.rhs) - 1 {
			ops := gen(mask)
			if equation.Eval(ops) == equation.lhs {
				ans += equation.lhs
				break
			}
		}
	}
	fmt.Println(ans)
}

func Part2() {
	equations := ReadEquations()
	ans := 0
	for _, equation := range equations {
		gen := ternary_operators(len(equation.rhs) - 1)
		allMask := int(math.Pow(float64(3), float64(len(equation.rhs)-1)))
		for mask := range allMask {
			ops := gen(mask)
			if equation.Eval(ops) == equation.lhs {
				ans += equation.lhs
				break
			}
		}
	}
	fmt.Println(ans)
}

func main() {
	Part1()
	Part2()
}
