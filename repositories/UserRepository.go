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
	selectQuery := util.NewSelectQuery(util.SELECT, *util.NewGroup("id", "name", "password", "email_address", "phone_number")).
		From("users")

	insertQuery := util.NewInsertQuery("users", "id", "name", "password", "email_address", "phone_number", "created_at")
	updateQuery := util.NewUpdateQuery("users", "name", "password", "email_address", "phone_number", "created_at")
	r.Impl = RepositoryImpl[entity.User]{
		TableName:       "users",
		SelectQueryBody: *selectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
		Repository:      r,
	}
}

func (r *UserRepository) setNewEntity(model *entity.User, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
}

func (r *UserRepository) buildEntity(row *sql.Row) (model *entity.User, err error) {
	var user entity.User
	err = row.Scan(&user.Id, &user.Name, &user.Password, &user.EmailAddress, &user.PhoneNumber)
	return &user, err
}

func (r *UserRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.User, err error) {
	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.EmailAddress, &user.PhoneNumber)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	return &users, err
}

func (r *UserRepository) saveOrUpdateScan(query string, model *entity.User) error {
	return execQuery(query, model.Id, model.Name, model.Password, model.EmailAddress, model.PhoneNumber, model.CreatedAt)
}

func (r *UserRepository) FindByName(name string) (*entity.User, error) {
	query := r.Impl.SelectQueryBody
	query.Where("users.name = $1 and users.deleted_at is null")
	return r.Impl.FindByQuery(&query, name)
}

func (r *UserRepository) FindByEmailAddress(email string) (*[]entity.User, error) {
	query := r.Impl.SelectQueryBody
	query.Where("users.email_address = $1 and users.deleted_at is null")
	return r.Impl.FindListByQuery(&query, email)
}

func (r *UserRepository) ValidateUser(user entity.User) (*entity.User, error) {
	query := r.Impl.SelectQueryBody
	query.Where("users.name = $1 and users.password = $2 and users.deleted_at is null")
	return r.Impl.FindByQuery(&query, user.Name, user.Password)
}
