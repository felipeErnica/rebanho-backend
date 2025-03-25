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
	Impl HandlerImpl[entity.User]
    Repository repositories.UserRepository
}

func InitUserAuthentication(app *app.App) {
    repository:=new(repositories.UserRepository)
    repository.Init()
    impl:=HandlerImpl[entity.User]{
        Repository: repository.Impl,
    }
    handler:=UserHandler{
        Repository: *repository,
        Impl: impl,
    }

    app.HandleFuncNoMiddleware("POST /login", handler.Authenticate)
    app.HandleFuncNoMiddleware("POST /register", handler.Register)
}

func (h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
    var user entity.User
    err:=json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        JsonServerError(err, w)
    }
    
    userDatabase, err:= h.Repository.FindByUsername(user.Username)
    if err != nil {
        util.LogError(err.Error())
        w.WriteHeader(http.StatusUnauthorized)
    }

    err = bcrypt.CompareHashAndPassword([]byte(userDatabase.Password), []byte(user.Password))
    if err != nil {
        util.LogError(err.Error())
        w.WriteHeader(http.StatusUnauthorized)
    }

    tokenString, err:=util.GenerateToken(userDatabase) 
    if err != nil {
        util.LogError(err.Error())
        w.WriteHeader(http.StatusUnauthorized)
    }
    
    authToken:=entity.AuthToken{
        Token: tokenString,
    }

    jsonToken, err:=json.Marshal(authToken)
    if err != nil {
        JsonServerError(err, w)
    }
    w.Write(jsonToken)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

}
