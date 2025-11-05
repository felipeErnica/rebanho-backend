package middlewares

import (
	"net/http"

	authConfig "github.com/felipeErnica/rebanho-backend/config/auth-config"
	"github.com/felipeErnica/rebanho-backend/apiError"
)

func AuthenticationMiddleware(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authorizationString:=r.Header.Get("Authorization")
        userId, err:= authConfig.VerifyToken(authorizationString)
        if err != nil {
            apiError.AuthenticationError(w, err)
            return
        }
        r = authConfig.RegisterUserId(r, userId)
        handler.ServeHTTP(w,r)
    }
}
