package middlewares

import (
	"context"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func RegisterUserId(ctx context.Context, userId string) context.Context {
    return context.WithValue(ctx, "userId", userId)
}

func GetUserId(ctx context.Context) (string, bool) {
    userId, ok := ctx.Value("userId").(string)
    return userId, ok
}
