package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
