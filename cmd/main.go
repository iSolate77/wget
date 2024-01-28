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
	flag.Parse()

	if *bblocks.MirrorMode {
		var resources []string
		if len(flag.Args()) == 0 {
			return
		}
		fmt.Println(bblocks.MirrorWebsite(flag.Args()[0], &resources))
	} else {
		if *bblocks.SilentMode {
			if len(flag.Args()) == 0 {
				return
			}
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
				go bblocks.DownloadFileWithRateLimitAndProgressBar(link, &wg)
			}
			wg.Wait()
		} else {
			if len(flag.Args()) == 0 {
				return
			}
			URL_PATH := flag.Args()[0]
			bblocks.DownloadFileWithRateLimitAndProgressBar(URL_PATH, nil)
		}
	}
}
