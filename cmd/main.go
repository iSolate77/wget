package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"wget/bblocks"
)

var wg sync.WaitGroup

func main() {
	// Print start time
	flag.Parse()

	if *bblocks.SilentMode {
		fmt.Println("output will be written to wget-log")
	} else {
		os.Remove("wget-log.txt")
	}

	if *bblocks.AsyncFileInput != "" {
		links, err := bblocks.GetLinksFromFile()
		if err != nil {
			log.Fatal(err)
		}
		for _, link := range links {
			wg.Add(1)
			go bblocks.DownloadFile(link, *bblocks.Output_name_arg_flag, &wg)
		}
		wg.Wait()
	} else {
		if len(flag.Args()) == 0 {
			return
		}
		URL_PATH := flag.Args()[0]
		bblocks.DownloadFile(URL_PATH, *bblocks.Output_name_arg_flag, nil)
	}
}
