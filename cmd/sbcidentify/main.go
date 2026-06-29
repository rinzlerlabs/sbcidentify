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

// cliHandler formats log records as: <time> <LEVEL> <message> key=value ...
type cliHandler struct {
	w        io.Writer
	level    slog.Level
	mu       sync.Mutex
	prefix   string      // built from WithGroup calls
	preAttrs []slog.Attr // built from WithAttrs calls
}

func (h *cliHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *cliHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	buf := fmt.Sprintf("%s %s %s",
		r.Time.Format("2006-01-02T15:04:05.999Z07:00"),
		r.Level,
		r.Message,
	)
	for _, a := range h.preAttrs {
		buf += " " + h.fmtAttr(a)
	}
	r.Attrs(func(a slog.Attr) bool {
		buf += " " + h.fmtAttr(a)
		return true
	})
	_, err := fmt.Fprintln(h.w, buf)
	return err
}

func (h *cliHandler) fmtAttr(a slog.Attr) string {
	if h.prefix != "" {
		return h.prefix + "." + a.Key + "=" + fmt.Sprintf("%v", a.Value.Any())
	}
	return a.Key + "=" + fmt.Sprintf("%v", a.Value.Any())
}

func (h *cliHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	c := *h
	c.preAttrs = make([]slog.Attr, len(h.preAttrs)+len(attrs))
	copy(c.preAttrs, h.preAttrs)
	copy(c.preAttrs[len(h.preAttrs):], attrs)
	return &c
}

func (h *cliHandler) WithGroup(name string) slog.Handler {
	c := *h
	if h.prefix != "" {
		c.prefix = h.prefix + "." + name
	} else {
		c.prefix = name
	}
	return &c
}

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
