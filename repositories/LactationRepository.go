package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type LactationRepository struct {
	Impl        *PageRepositoryImpl[entity.Lactation]
	SelectQuery *util.QueryConstructor
}

func (r *LactationRepository) Init() {

	dateFields := []string{
		"start_date",
		"end_date",
		"created_at",
		"deleted_at",
	}

	selectQuery := new(util.QueryConstructor).Select("lactations", "id", "start_date", "end_date", "production_period",
		"production_total", "average_production", "peak_production", "isr", "observation")
	selectQuery.AndSelect("pastures", "id", "name")
	selectQuery.AndSelect("cow", "id", "identificantion_number", "name", "status")
	selectQuery.AndSelect("calf", "id", "sex", "birth_date")
	selectQuery.From("lactations", "")
	selectQuery.LeftJoin("pastures", "").On("pastures.id", "lactations.pasture_id")
	selectQuery.LeftJoin("animals", "cow").On("cow.id", "lactations.animal_id")
	selectQuery.LeftJoin("animals", "calf").On("calf.id", "lactations.calf_id")

	insertQuery := new(util.QueryConstructor).Insert("lactations", "id", "start_date", "end_date", "production_period",
		"production_total", "average_production", "peak_production", "isr", "observation", "animal_id", "calf_id", "created_at",
		"user_id")

	updateQuery := new(util.QueryConstructor).Update("lactations", "id", "start_date", "end_date", "production_period",
		"production_total", "average_production", "peak_production", "isr", "observation", "animal_id", "calf_id", "created_at",
		"user_id")

	base := RepositoryImpl[entity.Lactation]{
		Repository:      r,
		TableName:       "lactations",
		SelectQueryBody: *selectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
	}

	r.Impl = &PageRepositoryImpl[entity.Lactation]{
		Base:           &base,
		PageRepository: r,
		DateFields:     dateFields,
	}

}

func (r *LactationRepository) setNewEntity(model *entity.Lactation, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
	model.UserId = GetUserId()
}

func (r *LactationRepository) buildEntity(row *sql.Row) (model *entity.Lactation, err error) {
	var lactation entity.Lactation
	err = row.Scan(&lactation.Id, &lactation.StartDate, &lactation.EndDate, &lactation.ProductionPeriod, &lactation.ProductionTotal,
		&lactation.AverageProduction, &lactation.PeakProduction, &lactation.Isr, &lactation.Observation,
		&lactation.Cow.Id, &lactation.Cow.IdentificationNumber, &lactation.Cow.Name, &lactation.Cow.Status,
		&lactation.Cow.Pasture.Id, &lactation.Cow.Pasture.Name,
		&lactation.Calf.Id, &lactation.Calf.Sex, &lactation.Calf.BirthDate)
	if err != nil {
		return
	}
	return &lactation, err
}

func (r *LactationRepository) buildListEntity(rows *sql.Rows) (list *[]entity.Lactation, err error) {
	var lactations []entity.Lactation
	for rows.Next() {
		var lactation entity.Lactation
		err = rows.Scan(&lactation.Id, &lactation.StartDate, &lactation.EndDate, &lactation.ProductionPeriod, &lactation.ProductionTotal,
			&lactation.AverageProduction, &lactation.PeakProduction, &lactation.Isr, &lactation.Observation,
			&lactation.Cow.Id, &lactation.Cow.IdentificationNumber, &lactation.Cow.Name, &lactation.Cow.Status,
			&lactation.Cow.Pasture.Id, &lactation.Cow.Pasture.Name,
			&lactation.Calf.Id, &lactation.Calf.Sex, &lactation.Calf.BirthDate)
		if err != nil {
			return
		}
		lactations = append(lactations, lactation)
	}
	return &lactations, err
}

func (r *LactationRepository) getFields(sort string) (firstField string, secondField string) {
	switch sort {
	case "name":
		return "lactations.name", "lactations.id"
	case "identification_number":
		return "lactations.identification_number", "lactations.id"
	case "birth_date":
		return "lactations.birth_date", "lactations.id"
	case "start_date":
		return "lactations.start_date", "lactations.id"
	case "end_date":
		return "lactations.end_date", "lactations.id"
	case "production_period":
		return "lactations.production_period", "lactations.id"
	case "production_total":
		return "lactations.production_total", "lactations.id"
	case "average_production":
		return "lactations.average_production", "lactations.id"
	case "peak_production":
		return "lactations.peak_production", "lactations.id"
	case "isr":
		return "lactations.isr", "lactations.id"
	default:
		return "lactations.created_at", "lactations.id"
	}
}

func (r *LactationRepository) createKey(sort string, lactation *entity.Lactation) string {
	switch sort {
	case "name":
		return fmt.Sprintf("%s,%s", *lactation.Cow.Name, lactation.Id)
	case "identification_number":
		return fmt.Sprintf("%s,%s", *lactation.Cow.IdentificationNumber, lactation.Id)
	case "birth_date":
		return fmt.Sprintf("%s,%s", lactation.Calf.BirthDate, lactation.Id)
	case "start_date":
		return fmt.Sprintf("%s,%s", lactation.StartDate, lactation.Id)
	case "end_date":
		key := fmt.Sprintf("%s,%s", "null", lactation.Id)
		if lactation.EndDate != nil {
			key = fmt.Sprintf("%s,%s", lactation.EndDate, lactation.Id)
		}
		return key
	case "production_period":
		return fmt.Sprintf("%d,%s", lactation.ProductionPeriod, lactation.Id)
	case "production_total":
		return fmt.Sprintf("%f,%s", lactation.ProductionTotal, lactation.Id)
	case "average_production":
		return fmt.Sprintf("%f,%s", lactation.AverageProduction, lactation.Id)
	case "peak_production":
		return fmt.Sprintf("%f,%s", lactation.PeakProduction, lactation.Id)
	case "isr":
		return fmt.Sprintf("%f,%s", lactation.Isr, lactation.Id)
	default:
		return fmt.Sprintf("%s,%s", lactation.CreatedAt, lactation.Id)
	}
}

func (r *LactationRepository) saveOrUpdateScan(query string, lactation *entity.Lactation) error {
	return execQuery(query, lactation.Id, lactation.Cow.Id, lactation.Calf.Id, lactation.StartDate, lactation.EndDate,
		lactation.ProductionPeriod, lactation.ProductionTotal, lactation.AverageProduction, lactation.PeakProduction,
		lactation.Isr, lactation.Observation, lactation.CreatedAt, lactation.DeletedAt)
}

func (r *LactationRepository) FindPage(sort string, direction string, cursor string) (page *entity.Page[entity.Lactation], err error) {
	query := r.SelectQuery.Where("lactations.deleted_at is null").And("lactations.user_id = $1")
	return r.Impl.FindRandomQueryPage(query, sort, direction, cursor, GetUserId())
}

func (r *LactationRepository) FindByAnimal(animalId string) (arr *[]entity.Lactation, err error) {
	query := r.SelectQuery.Where("lactations.deleted_at is null").And("lactations.animal_id = $1")
	return r.Impl.FindListByQuery(query, animalId)
}

func (r *LactationRepository) Add(newLactation *entity.Lactation) (*entity.Lactation, error) {
	return r.Impl.Add(newLactation)
}

func (r *LactationRepository) Save(lactation *entity.Lactation) error {
	return r.Impl.Save(lactation)
}

func (r *LactationRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
