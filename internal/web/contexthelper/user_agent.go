package contexthelper

import (
	"context"

	"github.com/ernestngugi/todo/internal/entities"
)

func UserAgent(ctx context.Context) string {
	existing := ctx.Value(entities.ContextKeyUserAgent)
	if existing == nil {
		return ""
	}
	return existing.(string)
}

func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, entities.ContextKeyUserAgent, userAgent)
}
