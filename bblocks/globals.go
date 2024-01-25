package bblocks

import (
	"flag"
	"os"
)


var (
	SilentMode = flag.Bool("B", false, "Silent Mode")
	LogFile, _ = os.Create("wget-log.txt")
	Output_name_arg_flag = flag.String("O", "", "Output file name")
)