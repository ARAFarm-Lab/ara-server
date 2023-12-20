package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	colorFatal = color("\033[1;31m%s\033[0m")
	colorWarn  = color("\033[1;33m%s\033[0m")
	colorInfo  = color("\033[1;36m%s\033[0m")

	timeFormat = "2006-01-02 15:04:05"
)

func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func Fatal(err error) {
	fmt.Printf("%s %s\n", colorFatal(fmt.Sprintf("[%s] [FATAL]", time.Now().Format(timeFormat))), err.Error())
	os.Exit(1)
}

func Info(message string) {
	fmt.Printf("%s %s\n", colorInfo(fmt.Sprintf("[%s] [INFO]", time.Now().Format(timeFormat))), message)
}

func Error(err error, message string, metadata ...interface{}) {
	logTemplate := fmt.Sprintf("%s %s\n", colorWarn(fmt.Sprintf("[%s] [ERROR]", time.Now().Format(timeFormat))), "%s: %s")
	if len(metadata) > 0 {
		payload, _ := json.Marshal(metadata[0])
		logTemplate += " [Metadata] %v\n"
		fmt.Printf(logTemplate, message, err.Error(), string(payload))
		return
	}

	if err == nil {
		err = errors.New("")
	}

	fmt.Printf(logTemplate, message, err.Error())
}
