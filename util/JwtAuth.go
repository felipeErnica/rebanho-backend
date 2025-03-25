package util

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("dfc4f9218a732c96ac527031a2c19b64202ce119793ade2f9e9fb814d7b235535743bf9a2a3079abcc61ab69c1d6b0b082bb26e57f3d7907623d82464be65ebc9147b65243dd7f37a74737fdb9eb68513f9c86d9734f5f2de205a7e521f561e431d6c88833e566a78669add8f2135664dba9aadaecd9a2f55b3c29a3d1ae096e")

func reclaimToken(authorizationString string) (token string, err error) {
	if authorizationString == "" {
		err = errors.New("nenhuma autenticação fornecida!")
		return
	}
	tokenParts := strings.Split(authorizationString, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		err = errors.New("tipo de autenticação inválido")
		return
	}
	if tokenParts[1] == "" {
		err = errors.New("autenticação inválida")
		return
	}
	return tokenParts[1], err
}

func VerifyToken(authorizationString string) (userId string, err error) {
    tokenString, err := reclaimToken(authorizationString)
    if err != nil {
        return
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })

    if err != nil {
        err = errors.New("autenticação inválida")
        return
    }

    if claims, ok:=token.Claims.(jwt.MapClaims); ok && token.Valid {
        return fmt.Sprint(claims["user_id"]), nil
    } else {
        return
    }
}

func GenerateToken(user *entity.User) (tokenString string, err error) {
    claims:=jwt.MapClaims{}
    claims["user_id"] = user.Id
    claims["exp"] = time.Now().Add(time.Hour * 5).Unix()
    token:=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}
