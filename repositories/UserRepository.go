package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type UserRepository struct {
	Impl RepositoryImpl[entity.User]
}

func (r *UserRepository) Init() {
    selectQuery:=new(util.QueryConstructor).Select("", "id", "username")
        selectQuery.From("users_active", "")
    insertQuery:=new(util.QueryConstructor).Insert("users", "id", "username", "password", "created_at")
    updateQuery:=new(util.QueryConstructor).Update("users", "username", "password", "created_at")
    r.Impl = RepositoryImpl[entity.User]{
        TableName: "users",
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
        Repository: r,
    }
}

func (r *UserRepository) setNewEntity(model *entity.User, id string, createdAt time.Time) {
    model.Id= id
    model.CreatedAt = createdAt
}

func (r *UserRepository) buildEntity(row *sql.Row) (model *entity.User, err error) {
    var user entity.User
    err = row.Scan(&user.Id, &user.Username)
    return &user, err
}

func (r *UserRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.User, err error) {
    var users []entity.User
    for rows.Next() {
        var user entity.User
        err = rows.Scan(&user.Id, &user.Username)
        if err != nil {
            return
        }
        users = append(users, user)
    }
    return &users, err
}

func (r *UserRepository) saveOrUpdateScan(query string, model *entity.User) error {
    return ExecQuery(query, model.Id, model.Username, model.Password, model.CreatedAt)
}

func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
    query:="WHERE username = $1"
    return r.Impl.FindByQuery(query, username)
}

func (r *UserRepository) FindById(id string) (*entity.User, error) {
    return r.Impl.FindById(id)
}

func (r *UserRepository) ValidateUser(user entity.User) (*entity.User, error) {
    query:="WHERE username = $1 AND password = $2"
    return r.Impl.FindByQuery(query, user.Username, user.Password)
}

func (r *UserRepository) Add(user entity.User) (*entity.User, error) {
    return r.Impl.Add(user)
}

func (r *UserRepository) Delete(id string) error {
    return r.Impl.SoftDelete(id)
}
