package bblocks

import "github.com/schollz/progressbar/v3"

func CreateProgressBar(totalSize int64) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(
		int(totalSize),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetWidth(35),
		progressbar.OptionSetDescription(FormatSize(totalSize)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	return bar
}
