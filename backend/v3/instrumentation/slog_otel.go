package instrumentation

import (
	"context"
	"log/slog"
	"slices"
)

// otelLogHandler wraps the handler that sends records to the OTEL [log.LoggerProvider],
// so that it honors the configured log level and masking like the standard error handler.
//
// The OTEL log processors report all levels as enabled
// and the OTEL slog bridge does not support a [slog.HandlerOptions.ReplaceAttr] function.
type otelLogHandler struct {
	next    slog.Handler
	level   slog.Leveler
	replace replacer
	groups  []string
}

// newOTELLogHandler returns a handler which only passes records with at least the given level
// and applies the replacer to all attributes before passing them to next.
// The replacer may be nil, in which case attributes are passed unchanged.
func newOTELLogHandler(next slog.Handler, level slog.Leveler, replace replacer) slog.Handler {
	return &otelLogHandler{
		next:    next,
		level:   level,
		replace: replace,
	}
}

func (h *otelLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level() && h.next.Enabled(ctx, level)
}

func (h *otelLogHandler) Handle(ctx context.Context, record slog.Record) error {
	if h.replace == nil || record.NumAttrs() == 0 {
		return h.next.Handle(ctx, record)
	}
	attrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, h.replaceAttr(h.groups, a))
		return true
	})
	replaced := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	replaced.AddAttrs(attrs...)
	return h.next.Handle(ctx, replaced)
}

func (h *otelLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h.replace != nil {
		replaced := make([]slog.Attr, len(attrs))
		for i, a := range attrs {
			replaced[i] = h.replaceAttr(h.groups, a)
		}
		attrs = replaced
	}
	return &otelLogHandler{
		next:    h.next.WithAttrs(attrs),
		level:   h.level,
		replace: h.replace,
		groups:  h.groups,
	}
}

func (h *otelLogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &otelLogHandler{
		next:    h.next.WithGroup(name),
		level:   h.level,
		replace: h.replace,
		groups:  append(slices.Clip(h.groups), name),
	}
}

// replaceAttr applies the replacer the same way as the [slog.HandlerOptions.ReplaceAttr] of the standard library handlers:
// it is called for every non-group attribute, with the groups it is nested in.
func (h *otelLogHandler) replaceAttr(groups []string, a slog.Attr) slog.Attr {
	a.Value = a.Value.Resolve()
	if a.Value.Kind() != slog.KindGroup {
		return h.replace(groups, a)
	}
	// attributes of a group without key are inlined into the parent
	if a.Key != "" {
		groups = append(slices.Clip(groups), a.Key)
	}
	members := a.Value.Group()
	replaced := make([]slog.Attr, len(members))
	for i, member := range members {
		replaced[i] = h.replaceAttr(groups, member)
	}
	a.Value = slog.GroupValue(replaced...)
	return a
}
