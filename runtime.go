package mist

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

var wmSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	wmSourceDir = regexp.MustCompile(`mist.mist\.go`).ReplaceAllString(file, "")
}
func WriteRuntimeMsg() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, wmSourceDir) || strings.HasSuffix(file, "_test.go")) {
			return fmt.Sprintf("\033[%d;1m%s:%d\033[0m", 36, file, line)
		}
	}
	return ""
}
