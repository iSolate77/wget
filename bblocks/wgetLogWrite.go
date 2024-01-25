package bblocks

import "fmt"

func WriteToWgetLog(a ...any) (n int, e error) {
	str := ""
	for _, v := range a {
		str += fmt.Sprintf("%s", v)
	}
	LogFile.WriteString(str)
	return 0, nil
}
