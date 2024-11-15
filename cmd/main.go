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
	board := flag.String("b", "", "Specify the board type")
	output := flag.String("o", "StdOut", "Specify the log output, accept StdOut, StdErr, or a file path")
	flag.Parse()

	logLevel := new(slog.LevelVar)
	if *debug {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}

	handlerConfig := &sbcidentify.HandlerConfig{Level: logLevel}

	switch *output {
	case "StdOut":
		sbcidentify.SetLogger(slog.New(sbcidentify.NewLogHandler(os.Stdout, handlerConfig)))
	case "StdErr":
		sbcidentify.SetLogger(slog.New(sbcidentify.NewLogHandler(os.Stderr, handlerConfig)))
	default:
		file, err := os.OpenFile(*output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		sbcidentify.SetLogger(slog.New(sbcidentify.NewLogHandler(file, handlerConfig)))
	}

	if board != nil && *board != "" {
		fmt.Println(sbcidentify.IsBoardType(sbcidentify.BoardType(*board)))
		return
	} else {
		board, err := sbcidentify.GetBoardType()
		if err != nil {
			panic(err)
		}
		fmt.Println(board)
	}
}
