package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type FarmRepository struct {
	Impl        PageRepositoryImpl[entity.Farm]
	SelectQuery *util.SelectConstructor
}

func (r *FarmRepository) Init() {
	r.SelectQuery = util.NewSelectQuery(util.SELECT, 
        *util.NewNamedGroup("farms", "id", "name", "state", "city", "tax_number", "status"),
	    *util.NewNamedGroup("owner", "id", "name")).
        From("farms").
        Joins("left join users as owner on owner.id = farms.owner_id")

	insertQuery := util.NewInsertQuery("farms", "id", "name", "state", "city", "tax_number", "created_at")
	updateQuery := util.NewUpdateQuery("farms", "name", "state", "city", "tax_number", "owner_id", "created_at")
	base := RepositoryImpl[entity.Farm]{
		Repository:      r,
		TableName:       "farms",
		SelectQueryBody: *r.SelectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
	}

	r.Impl = PageRepositoryImpl[entity.Farm]{
		Base:       &base,
		DateFields: []string{"created_at"},
	}
}

func (r *FarmRepository) setNewEntity(model *entity.Farm, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
	model.Owner.Id = GetUserId()
}

func (r *FarmRepository) buildEntity(row *sql.Row) (model *entity.Farm, err error) {
	var farm entity.Farm
	err = row.Scan(&farm.Id, &farm.Name, &farm.State, &farm.City, &farm.TaxNumber,
		&farm.Owner.Id, &farm.Owner.Name)
	return &farm, err
}

func (r *FarmRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.Farm, err error) {
	var farms []entity.Farm
	for rows.Next() {
		var farm entity.Farm
		err = rows.Scan(&farm.Id, &farm.Name, &farm.State, &farm.City, &farm.TaxNumber,
			&farm.Owner.Id, &farm.Owner.Name)
		if err != nil {
			return
		}
		farms = append(farms, farm)
	}
	return &farms, err
}

func (r *FarmRepository) saveOrUpdateScan(query string, model *entity.Farm) error {
	return execQuery(query, model.Id, model.Name, model.State, model.City, model.TaxNumber, model.Owner.Id, model.CreatedAt)
}

func (r *FarmRepository) getFields(sort string) (firstField string, secondField string) {
	return "name", "id"
}

func (r *FarmRepository) createKey(sort string, lastEntry *entity.Farm) (key string) {
	return fmt.Sprintf("%s,%s", lastEntry.Name, lastEntry.Id)
}

func (r *FarmRepository) buildFilterQuery(filter *entity.FarmFilter) (*util.SelectConstructor, *[]any) {
    query:=r.SelectQuery.Where("farms.deleted_at is null and farms.status = $1")
    args:=[]any{ "ACTIVE", } 

    if filter == nil {
        return r.SelectQuery, &args
    }
    numParam:=2
    if filter.Name != nil {
        condition:=fmt.Sprintf("and farms.name LIKE $%d", numParam)
        query.AppendWhere(condition)
        args = append(args, filter.Name)
        numParam++
    }
    if filter.TaxNumber != nil {
        condition:=fmt.Sprintf("and farms.tax_number LIKE $%d", numParam)
        query.AppendWhere(condition)
        args = append(args, filter.TaxNumber)
        numParam++
    }
    if filter.City != nil {
        condition:=fmt.Sprintf("and farms.city LIKE $%d", numParam)
        query.AppendWhere(condition)
        args = append(args, filter.City)
        numParam++
    }
    if filter.State != nil {
        condition:=fmt.Sprintf("and farms.state LIKE $%d", numParam)
        query.AppendWhere(condition)
        args = append(args, filter.State)
        numParam++
    }
    if filter.OwnerName != nil {
        condition:=fmt.Sprintf("and owner.name LIKE $%d", numParam)
        query.AppendWhere(condition)
        args = append(args, filter.OwnerName)
        numParam++
    }
    return query, &args
}

func (r *FarmRepository) FindAllPage(cursor string, sort string, order string,
    filter *entity.FarmFilter) (*entity.Page[entity.Farm], error) {
    query, args:=r.buildFilterQuery(filter)
    return r.Impl.FindRandomQueryPage(query, sort, order, cursor, *args...)
}

func (r *FarmRepository) FindById(id string) (*entity.Farm, error) {
    return r.Impl.FindById(id)
}

func (r *FarmRepository) FindByOwner() (*[]entity.Farm, error) {
    query:=r.SelectQuery.Where("owner.id = $1")
    return r.Impl.FindListByQuery(query, GetUserId())
}

func (r *FarmRepository) FindByNameAndOwner(farm entity.Farm) (*[]entity.Farm, error) {
    query:=r.SelectQuery.Where("owner.id = $1 and farms.name = $2 and farms.city = $3 and farms.state = $3")
    return r.Impl.FindListByQuery(query, GetUserId(), farm.Name, farm.City, farm.State)
}

func (r *FarmRepository) FindByTaxNumberAndOwner(taxNumber string) (*[]entity.Farm, error) {
    query:=r.SelectQuery.Where("owner.id = $1 and farms.tax_number = $2")
    return r.Impl.FindListByQuery(query, GetUserId(), taxNumber)
}

func (r *FarmRepository) Add(newModel *entity.Farm) (*entity.Farm, error) {
    return r.Impl.Add(newModel)
}

func (r *FarmRepository) Save(model *entity.Farm) error {
    return r.Impl.Save(model)
}

func (r *FarmRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
