package goproxy

import (
	"fmt"
	"strings"
	"sync"
)

var Pool = &sync.Pool{
	New: func() interface{} {
		return make([]string, 0, 10)
	},
}

func Getload(up uint64, down uint64) string {
	ss := Pool.Get().([]string)
	getLoadStr(&ss, Blue, up, "Upload:")
	ss = append(ss, " ")
	getLoadStr(&ss, Green, down, "Download:")
	s := strings.Join(ss, "")
	Pool.Put(ss[:0])
	return s
}

func getLoadStr(ss *[]string, color string, num uint64, prefix string) {
	u := float64(num)
	setColor(ss, color, prefix)
	if mod := num / 1024; mod > 0 {
		// KB or High
		if mod > 900 {
			// MB
			*ss = append(*ss, fmt.Sprintf("%6.2fMB/s", u/(1024*1024)))
		} else {
			// KB
			*ss = append(*ss, fmt.Sprintf("%6.2fKB/s", u/1024))
		}
	} else {
		// B
		*ss = append(*ss, fmt.Sprintf("%7dB/s", num))
	}
	// Set the default color.
	*ss = append(*ss, End)
}

func setColor(ss *[]string, color string, data string) {
	*ss = append(*ss, Start)
	*ss = append(*ss, color)
	*ss = append(*ss, Mid)
	*ss = append(*ss, data)
}

func GetColorError(err error) string {
	return Start + Red + Mid + End
}
