package middlewares

import (
	"net/http"

	authConfig "github.com/felipeErnica/rebanho-backend/config/auth-config"
	"github.com/felipeErnica/rebanho-backend/serverErrors"
)

func AuthenticationMiddleware(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authorizationString:=r.Header.Get("Authorization")
        userId, err:= authConfig.VerifyToken(authorizationString)
        if err != nil {
            serverErrors.AuthenticationError(w, err)
            return
        }
        r = authConfig.RegisterUserId(r, userId)
        handler.ServeHTTP(w,r)
    }
}
