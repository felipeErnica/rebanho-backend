package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type SlaughterhouseRepository struct {
    Impl RepositoryImpl[entity.Slaughterhouse]
}

func (r *SlaughterhouseRepository) Init() {
    selectQuery:=new(util.QueryConstructor).Select("", "id", "name", "tax_number").From("slaughterhouses","")
    updateQuery:=new(util.QueryConstructor).Update("slaughterhouses", "name", "tax_number", "created_at", "user_id")
    InsertQuery:=new(util.QueryConstructor).Insert("slaughterhouses", "id", "name", "tax_number", "created_at", "user_id")
    r.Impl = RepositoryImpl[entity.Slaughterhouse]{
        Repository: r,
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: InsertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
        TableName: "slaughterhouses",
    }
}

func (r *SlaughterhouseRepository) setNewEntity(model *entity.Slaughterhouse, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
}

func (r *SlaughterhouseRepository) buildEntity(row *sql.Row) (model *entity.Slaughterhouse, err error) {
    var slaughterhouse entity.Slaughterhouse
    err = row.Scan(&slaughterhouse.Id, &slaughterhouse.Name, &slaughterhouse.TaxNumber)
    return &slaughterhouse, err
}

func (r *SlaughterhouseRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.Slaughterhouse, err error) {
    var slaughterhouses []entity.Slaughterhouse
    for rows.Next() {
        var slaughterhouse entity.Slaughterhouse
        err = rows.Scan(&slaughterhouse.Id, &slaughterhouse.Name, &slaughterhouse.TaxNumber)
        if err != nil {
            return
        }
        slaughterhouses = append(slaughterhouses, slaughterhouse)
    }
    return &slaughterhouses, err
}

func (r *SlaughterhouseRepository) saveOrUpdateScan(query string, model *entity.Slaughterhouse) error {
    return execQuery(query, model.Id, model.Name, model.TaxNumber, model.CreatedAt, model.UserId)
}

func (r *SlaughterhouseRepository) FindAll() (*[]entity.Slaughterhouse, error) {
    return r.Impl.FindAll()
}

func (r *SlaughterhouseRepository) FindById(id string) (*entity.Slaughterhouse, error) {
    return r.Impl.FindById(id)
}

func (r *SlaughterhouseRepository) Add(newModel entity.Slaughterhouse) (*entity.Slaughterhouse, error) {
    return r.Impl.Add(newModel)
}

func (r *SlaughterhouseRepository) Save(model *entity.Slaughterhouse) error {
    return r.Impl.Save(model)
}

func (r *SlaughterhouseRepository) Delete(id string) error {
    return r.Impl.HardDelete(id)
}
