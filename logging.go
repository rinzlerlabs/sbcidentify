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
	source        string
	group         string
	args          []string
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
	args := make([]string, len(h.args))
	copy(args, h.args)
	r.Attrs(func(a slog.Attr) bool {
		args = append(args, fmt.Sprintf("%v", a.Value))
		return true
	})

	msg := r.Message
	if len(args) == 0 {
		_, err := h.writer.WriteString(fmt.Sprintf("%s: %s\n", h.source, msg))
		return err
	}
	argStr := fmt.Sprint(args)
	_, err := h.writer.WriteString(fmt.Sprintf("%s: %s %s\n", h.source, msg, argStr))

	return err
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.handlerConfig.Level.Level()
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handler := &handler{
		writer:        h.writer,
		handlerConfig: h.handlerConfig,
		args:          make([]string, 0),
	}
	for _, a := range attrs {
		if a.Key == "source" {
			handler.source = fmt.Sprintf("%v", a.Value)
		} else {
			handler.args = append(handler.args, fmt.Sprintf("%v", a.Value))
		}
	}
	return handler
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		writer:        h.writer,
		handlerConfig: h.handlerConfig,
		group:         name,
	}
}
