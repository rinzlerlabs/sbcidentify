package main

import (
	"flag"
	"log/slog"
	"os"

	"fmt"

	"github.com/thegreatco/sbcidentify"
)

func main() {
	debug := flag.Bool("d", false, "Enable debug logging")
	output := flag.String("o", "StdOut", "Specify the log output, accept StdOut, StdErr, or a file path")
	flag.Parse()

	logLevel := new(slog.LevelVar)
	if *debug {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}

	handlerConfig := &sbcidentify.HandlerConfig{Level: logLevel}

	var logger *slog.Logger
	switch *output {
	case "StdOut":
		logger = slog.New(sbcidentify.NewLogHandler(os.Stdout, handlerConfig))
	case "StdErr":
		logger = slog.New(sbcidentify.NewLogHandler(os.Stderr, handlerConfig))
	default:
		file, err := os.OpenFile(*output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		logger = slog.New(sbcidentify.NewLogHandler(file, handlerConfig))
	}

	sbcidentify.SetLogger(logger.With("source", "sbcidentify"))

	board, err := sbcidentify.GetBoardType()
	if err != nil {
		if errList, ok := err.(interface{ Unwrap() []error }); ok {
			// Access the slice of errors
			errs := errList.Unwrap()
			for _, e := range errs {
				fmt.Printf("Error: %v\n", e)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		fmt.Println(board.GetPrettyName())
	}
}
