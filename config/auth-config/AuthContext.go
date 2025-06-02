package authConfig

import (
	"context"
	"errors"
	"net/http"
)

func RegisterUserId(r *http.Request, userId string) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "userId", userId)
	r = r.WithContext(ctx)
}

func GetUserId(r *http.Request) (userId string, err error) {
	ctx := r.Context()
	userId, ok := ctx.Value("userId").(string)
    if !ok {
        err := errors.New("Falha na leitura do ID de usuário")
        return userId, err
    }
	return userId, err
}
