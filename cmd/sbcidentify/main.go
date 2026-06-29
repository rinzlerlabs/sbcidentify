package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/rinzlerlabs/sbcidentify"
)

func main() {
	verbose := flag.Bool("v", false, "Enable verbose logging")
	output := flag.String("o", "StdErr", "Log output: StdOut, StdErr, or a file path")
	flag.Parse()

	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}

	var w io.Writer
	switch *output {
	case "StdOut":
		w = os.Stdout
	case "StdErr":
		w = os.Stderr
	default:
		file, err := os.OpenFile(*output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot open log file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		w = file
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})))

	board, err := sbcidentify.GetBoardType()
	if err != nil {
		if errList, ok := err.(interface{ Unwrap() []error }); ok {
			for _, e := range errList.Unwrap() {
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
	fmt.Println(board.GetPrettyName())
}
