package authConfig

import (
	"context"
	"errors"
	"net/http"
)

type contextKey string

const userIDKey = contextKey("userId")

func RegisterUserId(r *http.Request, userId string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, userIDKey, userId)
    return r.WithContext(ctx)
}

func GetUserId(r *http.Request) (string, error) {
	ctx := r.Context()
	userId, ok := ctx.Value(userIDKey).(string)
    if !ok {
        err := errors.New("Falha na leitura do ID de usuário")
        return userId, err
    }
	return userId, nil
}
