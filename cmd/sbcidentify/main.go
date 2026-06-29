package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/rinzlerlabs/sbcidentify"
)

// cliHandler formats log records as: <time> <LEVEL> <message>
type cliHandler struct {
	w     io.Writer
	level slog.Level
	mu    sync.Mutex
}

func (h *cliHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *cliHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := fmt.Fprintf(h.w, "%s %s %s\n",
		r.Time.Format("2006-01-02T15:04:05.999Z07:00"),
		r.Level,
		r.Message,
	)
	return err
}

func (h *cliHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *cliHandler) WithGroup(string) slog.Handler       { return h }

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

	slog.SetDefault(slog.New(&cliHandler{w: w, level: level}))

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
