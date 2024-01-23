package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// Parse Args

	args := os.Args[1:]
	URL_PATH := args[0]

	response, err := http.Get(URL_PATH)
	if err != nil {
		fmt.Println("error")
		return

	}

	defer response.Body.Close()

	fileName := "output.html"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = io.Copy(writer, response.Body)
	if err != nil {
		fmt.Println("Error copying content to file:", err)
		return
	}

	fmt.Println("HTML content saved to", fileName)
}
