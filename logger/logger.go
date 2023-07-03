package logger

import "fmt"

func Logger(logChan chan string) {
	for message := range logChan {
		// Print the log message
		fmt.Println(message)
	}
}

func ErrLogger(logChan chan error) {
	for err := range logChan {
		// Print the log message
		fmt.Println("error while parsing: ", err.Error())
	}
}
