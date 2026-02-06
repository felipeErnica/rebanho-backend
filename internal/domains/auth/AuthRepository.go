package auth

import (
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepostory(db *sqlx.DB) *UserRepository {
	selectQuery := "SELECT users.* FROM users"
	return &UserRepository{selectQuery, "users", db}
}

func (r *UserRepository) FindByName(name string) (*User, error) {
	query := r.SelectQuery + " WHERE users.name = $1 AND users.deleted_at IS NULL"
	return util.GetOne[User](r.Db, query, name)
}

func (r *UserRepository) FindByEmailAddress(email string) (*[]User, error) {
	query := r.SelectQuery + " WHERE users.email_address = $1 AND users.deleted_at IS NULL"
	return util.GetList[User](r.Db, query, email)
}

func (r *UserRepository) ValidateUser(user User) (*User, error) {
	query := r.SelectQuery + "users.name = $1 AND users.password = $2 AND users.deleted_at IS NULL"
	return util.GetOne[User](r.Db, query, user.Name, user.Password)
}
