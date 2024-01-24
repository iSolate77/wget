package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"wget/helpers"
)

func main() {
	// Get flags
	output_name_arg_flag := flag.String("out", "", "a string")
	flag.Parse()
	// Parse Args after all flags usually it is the URL_PATH
	URL_PATH := flag.Args()[0]
	response, err := http.Get(URL_PATH)
	if err != nil {
		fmt.Println("error")
		return

	}

	defer response.Body.Close()

	fmt.Println("Output_name: ", *output_name_arg_flag)
	fmt.Println(flag.Args())

	// Get fileName
	var fileName string
	if *output_name_arg_flag != "" {
		fileName = *output_name_arg_flag
	} else if helpers.Get_filename(URL_PATH) != "" {
		fileName = helpers.Get_filename(URL_PATH)
	} else {
		fileName = "index.html"
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

	fmt.Println("Contents saved to", fileName)
}
