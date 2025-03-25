package lactation

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type LactationRepository struct {
	Base *PageRepositoryImpl[entity.Lactation]
}

func (r *LactationRepository) Init() {

	dateFields := []string{
		"start_date",
		"end_date",
		"created_at",
		"deleted_at",
	}

	selectQuery := `
        SELECT lactations.id, lactations.start_date, lactations.end_date, lactations.production_period, lactations.production_total, lactations.average_production
            lactations.peak_production, lactations.isr, lactations.observation,
            animal.id as animal_id, animal.identificantion_number as animal_number, animal.name as animal_name, animal.status as animal_status,
            cow_pasture.id, cow_pasture.name,
            calf.id as calf_id, calf.sex as calf_sex, calf.birth_date as calf_birth
            FROM lactations_active as lactations
        LEFT JOIN animals as animal ON animal.id = lactations.animal_id
        LEFT JOIN animals as calf ON calf.id = lactations.calf_id
        LEFT JOIN pastures as cow_pasture ON cow_pasture.id = animal.pasture_id
    `

	insertQuery := `
        INSERT INTO lactations(id, start_date, end_date, production_period, production_total, average_production, 
            peak_production, isr, observation, animal_id, calf_id, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `

	updateQuery := `
        UPDATE lactations
        SET start_date = $2, end_date = $3, production_period = $4, production_total = $5, average_production = $6, 
            peak_production = $7, isr = $8, observation = $9, animal_id = $10, calf_id = $11)
        WHERE id = $1
    `

	base := RepositoryImpl[entity.Lactation]{
		Repository:      r,
		TableName:       "lactations",
		SelectQueryBody: selectQuery,
		InsertQuery:     insertQuery,
		UpdateQuery:     updateQuery,
	}

	r.Base = &PageRepositoryImpl[entity.Lactation]{
		Base:           &base,
		PageRepository: r,
		DateFields:     dateFields,
	}

}

func (l *LactationRepository) SetNewEntity(model *entity.Lactation, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
}

func (l *LactationRepository) BuildEntity(row *sql.Row) (model *entity.Lactation, err error) {
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

func (l *LactationRepository) BuildListEntity(rows *sql.Rows) (list *[]entity.Lactation, err error) {
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

func (l *LactationRepository) getFields(sort string) (firstField string, secondField string) {
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

func (r *LactationRepository) createKey(sort string, lactation *entity.Lactation)  string  {
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
        key:=fmt.Sprintf("%s,%s", "null", lactation.Id)
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

func (l *LactationRepository) SaveOrUpdateScan(query string, lactation *entity.Lactation) error {
	return repositories.ExecQuery(query, lactation.Id, lactation.Cow.Id, lactation.Calf.Id, lactation.StartDate, lactation.EndDate,
		lactation.ProductionPeriod, lactation.ProductionTotal, lactation.AverageProduction, lactation.PeakProduction,
		lactation.Isr, lactation.Observation, lactation.CreatedAt, lactation.DeletedAt)
}

func (l *LactationRepository) FindPage(sort string, direction string, cursor string) (page *entity.Page[entity.Lactation], err error) {
    return l.Base.FindPage(sort, direction, cursor)
}

func (l *LactationRepository) FindByAnimal(animalId string) (arr *[]entity.Lactation, err error) {
    query:="WHERE lactations.animal_id = $1 ORDER BY lactations.birth_date"
    return l.Base.FindListByQuery(query, animalId)
}

func (l *LactationRepository) Add(newLactation entity.Lactation) (*entity.Lactation, error) {
    return l.Base.Add(newLactation)
}

func (l *LactationRepository) Save(lactation *entity.Lactation) error {
    return l.Base.Save(lactation)
}

func (l *LactationRepository) Delete(id string) error {
    return l.Base.SoftDelete(id)
}
