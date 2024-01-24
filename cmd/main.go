package main

import (
	"flag"

	"wget/bblocks"
)

func main() {
	// Get flags
	output_name_arg_flag := flag.String("O", "", "Output file name")
	// background_flag := flag.String("B", "", "Run in background")
	flag.Parse()
	// Parse Args after all flags usually it is the URL_PATH
	URL_PATH := flag.Args()[0]

	bblocks.DownloadFile(URL_PATH, *output_name_arg_flag)
}
