package auth

import (
	"encoding/json"
	"net/http"

	authConfig "github.com/felipeErnica/rebanho-backend/internal/config/auth-config"
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repository *UserRepository
}

func (h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	userRequest, ok := util.DecodeEntity(w, r, &User{})
	if !ok {
		return
	}

	user, err := h.Repository.FindByName(userRequest.Name)
	if err != nil {
		log.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		log.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, err := authConfig.GenerateToken(user.Id)
	if err != nil {
		log.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authToken := AuthToken{tokenString}

	jsonToken, err := json.Marshal(authToken)
	if err != nil {
		log.JsonServerError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonToken)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	newUser, ok := util.DecodeEntity(w, r, &User{})
	if !ok {
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.LogError("Erro na geração de senha")
		log.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newUser.Password = string(hashedPassword)
	// util.Add(w, r, h.Repository)
}
