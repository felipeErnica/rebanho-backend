package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
	"github.com/felipeErnica/rebanho-backend/util"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Impl       HandlerImpl[entity.User]
	Repository repositories.UserRepository
}

func InitUserAuthentication(app *app.App) {
	repository := new(repositories.UserRepository)
	repository.Init()
	impl := HandlerImpl[entity.User]{
		Repository: repository.Impl,
	}
	handler := UserHandler{
		Repository: *repository,
		Impl:       impl,
	}

	app.HandleFuncNoMiddleware("POST /login", handler.Authenticate)
	app.HandleFuncNoMiddleware("POST /register", handler.Register)
    LogControllersInit("Usuário")
}

func (h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		JsonServerError(err, w)
        return
	}

	userDatabase, err := h.Repository.FindByName(user.Name)
	if err != nil {
		util.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
        return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDatabase.Password), []byte(user.Password))
	if err != nil {
		util.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
        return
	}

	tokenString, err := util.GenerateToken(userDatabase)
	if err != nil {
		util.LogError(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
        return
	}

	authToken := entity.AuthToken{
		Token: tokenString,
	}

	jsonToken, err := json.Marshal(authToken)
	if err != nil {
		JsonServerError(err, w)
        return
	}

    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type","application/json")
	w.Write(jsonToken)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var newUser *entity.User
    if err:= json.NewDecoder(r.Body).Decode(&newUser); err != nil {
        JsonServerError(err, w)
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
    model, err:= h.Repository.Impl.Add(newUser)
    if err != nil {
        DatabaseSendError(err, w)
        return
    }

    response, err:= json.Marshal(model)
    if err != nil {
        JsonServerError(err,w)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write(response)
}
