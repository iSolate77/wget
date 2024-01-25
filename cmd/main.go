package main

import (
	"flag"
	"fmt"

	"wget/bblocks"
)

func main() {
	// Print start time
	flag.Parse()

	if *bblocks.SilentMode {
		fmt.Println("output will be written to wget-log")
	}

	// Get flags

	// background_flag := flag.String("B", "", "Run in background")
	// Parse Args after all flags usually it is the URL_PATH
	URL_PATH := flag.Args()[0]

	bblocks.DownloadFile(URL_PATH, *bblocks.Output_name_arg_flag)
	
}
