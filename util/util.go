//TODO use logger
package util

import (
	"errors"
	"strings"

	"github.com/fatih/color"
)

func GetMap(flags []string) (map[string]string, error) {
	args := make(map[string]string)
	if len(flags) == 0 {
		return args, nil
	}
	for i := range flags {
		if !strings.Contains(flags[i], "=") {
			return args, errors.New("bad format of labels (expected name=value)")
		} else {
			parts := strings.SplitN(flags[i], "=", 2)
			name := strings.ToLower(strings.TrimSpace(parts[0]))
			value := strings.TrimSpace(parts[1])
			args[name] = value
		}
	}
	return args, nil
}

func Err(format string, args ...interface{}) {
	color.Red(format, args...)
}

func Info(format string, args ...interface{}) {
	color.Blue(format, args...)
}

func Warn(format string, args ...interface{}) {
	color.Yellow(format, args...)
}
