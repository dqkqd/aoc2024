package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	utils "example.com/aoc2024"
)

type InstructionType int

const (
	Do InstructionType = iota
	Dont
	Mul
	Invalid
)

type Instruction struct {
	instructionType InstructionType
	value           int
}

var (
	doInstruction      = Instruction{Do, -1}
	dontInstruction    = Instruction{Dont, -1}
	invalidInstruction = Instruction{Invalid, -1}
)

func ReadWhile(buf *bufio.Reader, pred func(byte) bool) []byte {
	out := []byte{}
	for {
		c, err := buf.ReadByte()
		if err != nil {
			break
		}
		if !pred(c) {
			buf.UnreadByte()
			break
		}
		out = append(out, c)
	}
	return out
}

func Expect(buf *bufio.Reader, expect string) bool {
	expectedByte := []byte(expect)
	out, err := buf.Peek(len(expectedByte))
	if err != nil {
		return false
	}

	isEqual := bytes.Equal(out, expectedByte)
	if isEqual {
		// advance
		io.ReadFull(buf, out)
		return true
	}

	return false
}

func ReadInt(buf *bufio.Reader) (int, error) {
	isDigit := func(b byte) bool {
		_, err := strconv.Atoi(string(b))
		return err == nil
	}
	outBytes := ReadWhile(buf, isDigit)
	value, err := strconv.Atoi(string(outBytes))
	if err != nil {
		return 0, err
	}
	return value, nil
}

func ReadMulInstruction(buf *bufio.Reader) Instruction {
	notM := func(b byte) bool {
		return b != 'm'
	}

	ReadWhile(buf, notM)

	if !Expect(buf, "mul(") {
		return invalidInstruction
	}

	lhs, err := ReadInt(buf)
	if err != nil {
		return invalidInstruction
	}

	if !Expect(buf, ",") {
		return invalidInstruction
	}

	rhs, err := ReadInt(buf)
	if err != nil {
		return invalidInstruction
	}

	if !Expect(buf, ")") {
		return invalidInstruction
	}

	return Instruction{Mul, lhs * rhs}
}

func ReadInstruction(buf *bufio.Reader) Instruction {
	notMorD := func(b byte) bool {
		return b != 'm' && b != 'd'
	}

	ReadWhile(buf, notMorD)

	if Expect(buf, "do()") {
		return doInstruction
	}
	if Expect(buf, "don't()") {
		return dontInstruction
	}

	mul := ReadMulInstruction(buf)
	if mul != invalidInstruction {
		return mul
	}

	// read at least 1 byte to avoid infinite loop
	buf.ReadByte()
	return invalidInstruction
}

func ReadInstructions() []Instruction {
	buf := utils.Read(3, false)

	instructions := []Instruction{}
	for {
		_, err := buf.Peek(1)
		if err != nil {
			break
		}
		ins := ReadInstruction(buf)
		if ins != invalidInstruction {
			instructions = append(instructions, ins)
		}
	}
	return instructions
}

func Part1() {
	ans := 0
	for _, ins := range ReadInstructions() {
		if ins.instructionType == Mul {
			ans += ins.value
		}
	}
	fmt.Println(ans)
}

func Part2() {
	ans := 0

	latestInstruction := doInstruction
	for _, ins := range ReadInstructions() {
		if ins.instructionType == Mul && latestInstruction == doInstruction {
			ans += ins.value
		}
		if ins == doInstruction || ins == dontInstruction {
			latestInstruction = ins
		}
	}
	fmt.Println(ans)
}

func main() {
	Part1()
	Part2()
}
