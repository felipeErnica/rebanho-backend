package auth

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
    SelectQuery string
    TableName string
    Db *sqlx.DB
}

func NewRepostory(db *sqlx.DB) *UserRepository {
    selectQuery := "SELECT users.* FROM users"
    return &UserRepository{selectQuery, "users", db}
}

func (r *UserRepository) FindByName(name string) (*User, error) {
	query := r.SelectQuery + " WHERE users.name = $1 and users.deleted_at is null"
	return repositoriesUtil.GetOne[User](r.Db, query, name)
}

func (r *UserRepository) FindByEmailAddress(email string) (*[]User, error) {
	query := r.SelectQuery + " WHERE users.email_address = $1 and users.deleted_at is null"
	return repositoriesUtil.GetList[User](r.Db, query, email)
}

func (r *UserRepository) ValidateUser(user User) (*User, error) {
	query := r.SelectQuery + "users.name = $1 and users.password = $2 and users.deleted_at is null"
	return repositoriesUtil.GetOne[User](r.Db, query, user.Name, user.Password)
}

func (r *UserRepository) Add(user *User) (*User, error) {
    return repositoriesUtil.Add(r.Db, r.TableName, user)
}

func (r *UserRepository) Update(user *User) error {
    return repositoriesUtil.Update(r.Db, r.TableName, user)
}
