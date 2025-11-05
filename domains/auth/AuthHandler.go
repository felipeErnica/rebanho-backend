package auth

import (
	"encoding/json"
	"net/http"

	authConfig "github.com/felipeErnica/rebanho-backend/config/auth-config"
	"github.com/felipeErnica/rebanho-backend/apiError"
	"github.com/felipeErnica/rebanho-backend/util"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repository *UserRepository
}

func (h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	userRequest, ok := handlersUtil.DecodeEntity(w, r, &User{}); if !ok {
        return
    }

	user, err := h.Repository.FindByName(userRequest.Name)
	if err != nil {
		util.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		util.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, err := authConfig.GenerateToken(user.Id)
	if err != nil {
		util.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authToken := AuthToken{tokenString}

	jsonToken, err := json.Marshal(authToken)
	if err != nil {
		apiError.JsonServerError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonToken)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	newUser, ok := handlersUtil.DecodeEntity(w, r, &User{}); if !ok {
        return
    }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		util.LogError("Erro na geração de senha")
		util.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newUser.Password = string(hashedPassword)
    // handlersUtil.Add(w, r, h.Repository)
}
