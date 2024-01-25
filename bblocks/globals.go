package bblocks

import (
	"flag"
	"os"
)


var (
	SilentMode = flag.Bool("B", false, "Silent Mode")
	LogFile, _ = os.Create("wget-log.txt")
	AsyncFileInput = flag.String("i", "", "Async file download from input txt source")
	Output_name_arg_flag = flag.String("O", "", "Output file name")
)