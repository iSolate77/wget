package bblocks

import (
	"flag"
	"os"
)


var (
	SilentMode = flag.Bool("B", false, "Silent Mode")
	Reject = flag.String("R", "", "Reject specific file types")
	Exclude = flag.String("X", "", "Exclude specific directories")
	MirrorMode = flag.Bool("mirror",false,"Mirror Website")
	LogFile, _ = os.Create("wget-log.txt")
	Output_name_arg_flag = flag.String("O", "", "Output file name")
	New_file_path = flag.String("P", "", "File path")
	File *os.File
	Any_error error
	FilePath string
	AsyncFileInput = flag.String("i", "", "Async file download from input txt source")
	RateLimit = flag.String("rate-limit", "", "Speed limit for download (e.g., 400k, 2M)")
)

type RateLimiter struct {
	limit  float64
	burst  float64
	tokens float64
}

type CustomWriter interface {
	Write([]byte) (int, error)
}
