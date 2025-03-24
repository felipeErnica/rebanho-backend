package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("dfc4f9218a732c96ac527031a2c19b64202ce119793ade2f9e9fb814d7b235535743bf9a2a3079abcc61ab69c1d6b0b082bb26e57f3d7907623d82464be65ebc9147b65243dd7f37a74737fdb9eb68513f9c86d9734f5f2de205a7e521f561e431d6c88833e566a78669add8f2135664dba9aadaecd9a2f55b3c29a3d1ae096e")

type Middleware func(http.HandlerFunc) http.HandlerFunc

func reclaimToken(authorizationString string) (token string, err error) {
    if authorizationString == "" {
        err = errors.New("nenhuma autenticação fornecida!")
        return
    }

    tokenStrings:= strings.Split(authorizationString, " ")
    if tokenStrings[0] != "Bearer" {
        err = errors.New("tipo de autenticação inválido")
        return
    } 

    if tokenStrings[1] == "" {
        err = errors.New("autenticação inválida")
        return
    }

    return tokenStrings[0], err
}

func AuthenticationMiddleware(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authorizationString:=r.Header.Get("Authorization")
        tokenString, err:=reclaimToken(authorizationString)
        if err != nil {
            serverErrors.AuthenticationError(w, err)
            return
        }

        credential, err:=jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
            return secretKey, nil
        })
        if err != nil {
            err:=errors.New("autenticação inválida")
            serverErrors.AuthenticationError(w, err)
            return
        }

        println(credential.Claims)
        handler.ServeHTTP(w,r)
    }
}
