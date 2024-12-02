package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func Read(day int, sample bool) (reader *bufio.Reader) {
	wd, err := os.Getwd()
	CheckError(err)

	inputFolder := filepath.Dir(wd)
	var fileName string
	if sample {
		fileName = "sample.txt"
	} else {
		fileName = "input.txt"
	}
	filePath := filepath.Join(inputFolder, "input", fmt.Sprintf("day%02d", day), fileName)

	f, err := os.Open(filePath)
	CheckError(err)

	return bufio.NewReader(f)
}

func ArrayReader(reader *bufio.Reader) func() (arr []int, err error) {
	return func() (arr []int, err error) {
		s, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		arr = []int{}

		items := strings.Fields(s)
		for _, item := range items {
			num, err := strconv.Atoi(item)
			if err != nil {
				return nil, err
			}
			arr = append(arr, num)
		}

		return arr, nil
	}
}
