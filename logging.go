package sbcidentify

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type handler struct {
	writer        *os.File
	handlerConfig *HandlerConfig
}

type HandlerConfig struct {
	Level *slog.LevelVar
}

func NewLogHandler(w *os.File, handlerConfig *HandlerConfig) *handler {
	return &handler{
		writer:        w,
		handlerConfig: handlerConfig,
	}
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	args := make([]string, 0)
	r.Attrs(func(a slog.Attr) bool {
		args = append(args, fmt.Sprintf("%v", a.Value))
		return true
	})

	msg := r.Message
	if len(args) == 0 {
		_, err := h.writer.WriteString(fmt.Sprintf("%s\n", msg))
		return err
	}
	argStr := fmt.Sprint(args)
	_, err := h.writer.WriteString(fmt.Sprintf("%s %s\n", msg, argStr))

	return err
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.handlerConfig.Level.Level()
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *handler) WithGroup(name string) slog.Handler {
	return h
}
