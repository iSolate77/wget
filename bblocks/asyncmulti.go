package bblocks

import (
	"bufio"
	"fmt"
	"os"
)

func GetLinksFromFile() ([]string, error) {
	pathArr := []string{}
	f, err := os.Open(*AsyncFileInput)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		pathArr = append(pathArr, scanner.Text())
	}
	fmt.Println(pathArr)
	return pathArr, nil
}
