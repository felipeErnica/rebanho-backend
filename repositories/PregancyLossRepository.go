package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type PregnancyLossRepository struct {
	Impl PageRepositoryImpl[entity.PregnancyLoss]
}

func (r *PregnancyLossRepository) Init() {

	dateFields := []string{
		"loss_date",
	}

	selectQuery := util.NewSelectQuery(util.SELECT, 
        *util.NewNamedGroup("loss", "id", "loss_type", "loss_date", "created_at"),
	    *util.NewNamedGroup("animal", "id", "identification_number", "name", "animal_order")).
	    From("pregnancy_losses as loss").
	    Joins("left join animals as animal on animal.id = loss.animal_id")

	updateQuery := util.NewUpdateQuery("pregnancy_losses", "animal_id", "loss_type",
		"loss_date", "created_at", "user_id")
	insertQuery := util.NewInsertQuery("pregnancy_losses", "id", "animal_id", "loss_type",
		"loss_date", "created_at", "user_id")

	base := RepositoryImpl[entity.PregnancyLoss]{
		Repository:      r,
		SelectQueryBody: *selectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
		TableName:       "pregnancy_losses",
	}

	r.Impl = PageRepositoryImpl[entity.PregnancyLoss]{
		Base:           &base,
		PageRepository: r,
		DateFields:     dateFields,
	}

}

func (r *PregnancyLossRepository) setNewEntity(model *entity.PregnancyLoss, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
	model.UserId = GetUserId()
}

func (r *PregnancyLossRepository) buildEntity(row *sql.Row) (model *entity.PregnancyLoss, err error) {
	var loss entity.PregnancyLoss
	err = row.Scan(&loss.Id, &loss.LossType, &loss.LossDate, &loss.CreatedAt,
		&loss.AnimalId, &loss.AnimalNumber, &loss.AnimalOrder)
	return &loss, err
}

func (r *PregnancyLossRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.PregnancyLoss, err error) {
	var losses []entity.PregnancyLoss
	for rows.Next() {
		var loss entity.PregnancyLoss
		err = rows.Scan(&loss.Id, &loss.LossType, &loss.LossDate, &loss.CreatedAt,
			&loss.AnimalId, &loss.AnimalNumber, &loss.AnimalOrder)
		if err != nil {
			return
		}
		losses = append(losses, loss)
	}
	return &losses, err
}

func (r *PregnancyLossRepository) saveOrUpdateScan(query string, model *entity.PregnancyLoss) error {
	return execQuery(query, model.Id, model.AnimalId, model.LossType, model.LossDate, model.CreatedAt, model.UserId)
}

func (r *PregnancyLossRepository) getFields(sort string) (firstField string, secondField string) {
	switch sort {
	case "name":
		return "animals.name", "animals.id"
	case "identification_number":
		return "animals.animal_order", "animals.id"
	case "loss_date":
		return "animals.birth_date", "animals.id"
	default:
		return "animals.created_at", "animals.id"
	}
}

func (r *PregnancyLossRepository) createKey(sort string, lastEntry *entity.PregnancyLoss) (key string) {
	switch sort {
	case "name":
		return fmt.Sprintf("%s,%s", lastEntry.AnimalName, lastEntry.Id)
	case "identification_number":
		return fmt.Sprintf("%d,%s", lastEntry.AnimalOrder, lastEntry.Id)
	case "loss_date":
		return fmt.Sprintf("%s,%s", lastEntry.LossDate, lastEntry.Id)
	default:
		return fmt.Sprintf("%s,%s", lastEntry.CreatedAt, lastEntry.Id)
	}
}

func (r *PregnancyLossRepository) FindPage(cursor string, sort string, order string) (*entity.Page[entity.PregnancyLoss], error) {
	return r.Impl.FindPage(cursor, sort, order)
}

func (r *PregnancyLossRepository) FindByAnimalId(animalId string) (*[]entity.PregnancyLoss, error) {
	query := r.Impl.Base.SelectQueryBody.Where("loss.animal_id = $1 and loss.deleted_at is null")
	return r.Impl.FindListByQuery(query)
}

func (r *PregnancyLossRepository) FindById(id string) (*entity.PregnancyLoss, error) {
	return r.Impl.FindById(id)
}

func (r *PregnancyLossRepository) Add(newModel *entity.PregnancyLoss) (*entity.PregnancyLoss, error) {
	return r.Impl.Add(newModel)
}

func (r *PregnancyLossRepository) Save(model *entity.PregnancyLoss) error {
	return r.Impl.Save(model)
}

func (r *PregnancyLossRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
