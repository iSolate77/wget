package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"

	"wget/helpers"
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

	var fileName string

	// // Get flags
	// Output_name_arg_flag := flag.String("O", "foo", "a string")
	// flag.Parse()
	//
	// fmt.Println("Output_name: ", *Output_name_arg_flag)

	// Get path_extension for the file name
	path_extension := helpers.Get_extension(URL_PATH)
	if path_extension == "" {
		fileName = "index.html"
	} else {
		fileName = path_extension
	}

	// Create file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write to file
	writer := bufio.NewWriter(file)
	_, err = io.Copy(writer, response.Body)
	if err != nil {
		fmt.Println("Error copying content to file:", err)
		return
	}

	fmt.Println("HTML content saved to", fileName)
}
