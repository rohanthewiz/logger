package logger

import (
	"fmt"
	"path/filepath"
	"runtime/debug"
	"strings"
)

type PrintStackTraceOpts struct {
	WithoutHeading bool
}

func PrintStackTrace(options ...PrintStackTraceOpts) {
	var opts PrintStackTraceOpts

	if len(options) > 0 {
		opts = options[0]
	}

	if !opts.WithoutHeading {
		fmt.Println(strings.Repeat("-", 20), "USER PRINTED STACK TRACE", strings.Repeat("-", 20))
	}

	var out []string
	newLine := string(filepath.Separator)

	arr := strings.Split(string(debug.Stack()), newLine)
	for _, item := range arr {
		if trimmed := strings.TrimSpace(item); trimmed != "" { // don't add an empty item
			out = append(out, trimmed)
		}
	}

	fmt.Println(strings.Join(out, newLine))

	if !opts.WithoutHeading {
		fmt.Println(strings.Repeat("-", 60))
	}
}
