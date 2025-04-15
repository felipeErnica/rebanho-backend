package middlewares

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/repositories"
	"github.com/felipeErnica/rebanho-backend/serverErrors"
	"github.com/felipeErnica/rebanho-backend/util"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func AuthenticationMiddleware(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authorizationString:=r.Header.Get("Authorization")
        userId, err:= util.VerifyToken(authorizationString)
        if err != nil {
            serverErrors.AuthenticationError(w, err)
            return
        }
        repositories.SetUserId(&userId)
        handler.ServeHTTP(w,r)
    }
}
