package utils

import (
	"bufio"
	"errors"
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

func ArrayReader(reader *bufio.Reader, delim string) func() (arr []int, err error) {
	return func() (arr []int, err error) {
		s, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		s = strings.TrimSpace(s)
		if len(s) == 0 {
			return nil, errors.New("empty string")
		}

		arr = []int{}

		items := strings.Split(s, delim)

		for _, item := range items {
			num, err := strconv.Atoi(strings.TrimSpace(item))
			if err != nil {
				continue
			}
			arr = append(arr, num)
		}

		return arr, nil
	}
}
