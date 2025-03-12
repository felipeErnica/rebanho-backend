package repositories

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/serverErrors"
	"github.com/felipeErnica/rebanho-backend/util"
)

type LactationRepository struct{}

func (l *LactationRepository) createCriteriaFirstPage(sort string, direction string) string {

	var criteriaOrder string
	direction = strings.ToUpper(direction)
    criteriaWhere:="WHERE lac.deleted_at IS NULL"

	switch sort {
	case "name":
		criteriaOrder = fmt.Sprintf("ORDER BY animal.name %[1]s, lac.id %[1]s", direction)
	case "identification_number":
		criteriaOrder = fmt.Sprintf("ORDER BY animal.ring_order %[1]s, lac.id %[1]s", direction)
	case "birth_date":
		criteriaOrder = fmt.Sprintf("ORDER BY calf.birth_date %[1]s, lac.id %[1]s", direction)
	case "start_date":
		criteriaOrder = fmt.Sprintf("ORDER BY animal.start_date %[1]s, lac.id %[1]s", direction)
	case "end_date":
		criteriaOrder = fmt.Sprintf("ORDER BY lac.end_date %[1]s, lac.id %[1]s", direction)
	case "production_period":
		criteriaOrder = fmt.Sprintf("ORDER BY lac.production_period %[1]s, lac.id %[1]s", direction)
	case "production_total":
		criteriaOrder = fmt.Sprintf("ORDER BY lac.production_total %[1]s, lac.id %[1]s", direction)
	case "average_production":
		criteriaOrder = fmt.Sprintf("ORDER BY lac.average_production %[1]s, lac.id %[1]s", direction)
	case "peak_production":
		criteriaOrder = fmt.Sprintf("ORDER BY lac.peak_production %[1]s, lac.id %[1]s", direction)
	case "isr":
		criteriaOrder = fmt.Sprintf("ORDER BY lac.isr %[1]s, lac.id %[1]s", direction)
	default:
		criteriaOrder = "ORDER BY animal.created_at, animal.id"
	}

	return fmt.Sprintf(`%s %s`, criteriaWhere, criteriaOrder)
}

func (l *LactationRepository) createCriteriaNextPage(sort string, direction string) string {

	var signal string
	switch direction {
	case "asc":
		signal = ">"
	case "desc":
		signal = "<"
	}

	var criteria string
	direction = strings.ToUpper(direction)

	switch sort {
	case "name":
		criteria = fmt.Sprintf(`
            WHERE (animal.name, lac.id) %s ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY animal.name %[2]s, lac.id %[2]s`, signal, direction)
	case "identification_number":
		criteria = fmt.Sprintf(`
            WHERE (animal.ring_order, lac.id) %s ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY animal.ring_order %[2]s, lac.id %[2]s`, signal, direction)
	case "birth_date":
		criteria = fmt.Sprintf(`
            WHERE (animal.birth_date, lac.id) %s ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY calf.birth_date %[2]s, lac.id %[2]s`, signal, direction)
	case "start_date":
		criteria = fmt.Sprintf(`
            WHERE (lac.start_date, lac.id) %s ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY animal.start_date %[2]s, lac.id %[2]s`, signal, direction)
	case "end_date":
		criteria = fmt.Sprintf(`
            WHERE (lac.end_date, lac.id) %s ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY lac.end_date %[2]s, lac.id %[2]s`, signal, direction)
	case "production_period":
		criteria = fmt.Sprintf(`
            WHERE (lac.production_period, lac.id) %s ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY lac.production_period %[2]s, lac.id %[2]s`, signal, direction)
	case "production_total":
		criteria = fmt.Sprintf(`
            WHERE (lac.production_total, lac.id) %s ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY lac.production_total %[2]s, lac.id %[2]s`, signal, direction)
	case "average_production":
		criteria = fmt.Sprintf(`
            WHERE (lac.average_production, lac.id) %s ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY lac.average_production %[2]s, lac.id %[2]s`, signal, direction)
	case "peak_production":
		criteria = fmt.Sprintf(`
            WHERE (lac.peak_production, lac.id) %s ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY lac.peak_production %[2]s, lac.id %[2]s`, signal, direction)
	case "isr":
		criteria = fmt.Sprintf(`
            WHERE (lac.isr, lac.id) %s ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY lac.isr %[2]s, lac.id %[2]s`, signal, direction)
	default:
		criteria = `
            WHERE (lac.created_at, lac.id) > ($1, $2) AND lac.deleted_at IS NULL
            ORDER BY lac.created_at, lac.id
        `
	}

	return criteria
}

func (l *LactationRepository) createNextCursor(sort string, arr []entity.LactationComplete) (cursor string, err error) {

	if len(arr) == 0 {
		err = serverErrors.EmptyList()
		return
	}

	lastEntry := arr[len(arr)-1]

	switch sort {
	case "name":
		cursor = encodeCursor(lastEntry.AnimalName, lastEntry.Id)
	case "identification_number":
		cursor = encodeCursor(strconv.Itoa(lastEntry.AnimalOrder), lastEntry.Id)
	case "birth_date":
		cursor = encodeCursor(lastEntry.CalfBirthDate.String(), lastEntry.Id)
	case "start_date":
		cursor = encodeCursor(lastEntry.StartDate.String(), lastEntry.Id)
	case "end_date":
		cursor = encodeCursor(lastEntry.EndDate.String(), lastEntry.Id)
	case "production_period":
		cursor = encodeCursor(strconv.Itoa(int(lastEntry.ProductionPeriod)), lastEntry.Id)
	case "production_total":
		cursor = encodeCursor(util.Float32ToString(lastEntry.ProductionTotal), lastEntry.Id)
	case "average_production":
		cursor = encodeCursor(util.Float32ToString(lastEntry.AverageProduction), lastEntry.Id)
	case "peak_production":
		cursor = encodeCursor(util.Float32ToString(lastEntry.PeakProduction), lastEntry.Id)
	case "isr":
		cursor = encodeCursor(lastEntry.AnimalName, lastEntry.Id)
	default:
		cursor = encodeCursor(lastEntry.CreatedAt.String(), lastEntry.Id)
	}

	return cursor, err
}

func (l *LactationRepository) saveOrUpdateScan(query string, lactation *entity.Lactation) error {
    return execQuery(query, lactation.Id, lactation.AnimalId, lactation.CalfId, lactation.StartDate, lactation.EndDate,
        lactation.ProductionPeriod, lactation.ProductionTotal, lactation.AvarageProduction, lactation.PeakProduction,
        lactation.Isr, lactation.Observation, lactation.CreatedAt, lactation.DeletedAt)
}

func (l *LactationRepository) hasNextPage(arr []entity.LactationComplete) bool {
	return len(arr) == PAGE_LIMIT
}


func (l *LactationRepository) GetFirstPage(sort string, direction string) (page *entity.LactationPage, err error) {
    criteria:= l.createCriteriaFirstPage(sort, direction)
	query := fmt.Sprintf(`
        SELECT lac.id, lac.start_date, lac.end_date, lac.production_period, lac.production_total, lac.average_production
        lac.peak_production, lac.isr, lac.observation,
        animal.id as animal_id, animal.identificantion_number as animal_number, animal.name as animal_name,
        animal.ring_order as animal_order, animal.pasture_id as animal_pasture, animal.status as animal_status,
        calf.id as calf_id, calf.sex as calf_sex, calf.birth_date as calf_birth
        FROM lactations as lac
        LEFT JOIN animals as animal ON animal.id = lac.animal_id
        LEFT JOIN animals as calf ON calf.id = lac.animal_id
        %s
        LIMIT %d
        `, criteria, PAGE_LIMIT)
	sqlStatement, err := selectQueryList(query)
	defer sqlStatement.Close()
	
    if err != nil {
		return 
	}

    var lactations []entity.LactationComplete

    for sqlStatement.Next() {
        var lactation entity.LactationComplete
        err = sqlStatement.Scan(&lactation.Id, &lactation.StartDate, &lactation.EndDate, &lactation.ProductionPeriod, &lactation.ProductionTotal,
            &lactation.AverageProduction, &lactation.PeakProduction, &lactation.Isr, &lactation.Observation, &lactation.AnimalId, &lactation.AnimalNumber,
            &lactation.AnimalName, &lactation.AnimalOrder, &lactation.AnimalPasture, &lactation.AnimalPasture, 
            &lactation.CalfId, &lactation.CalfSex, &lactation.CalfBirthDate)
        if err != nil {
            return
        }
        lactations = append(lactations, lactation)
    }
    
    nextCursor, err:= l.createNextCursor(sort, lactations)
    if err != nil {
        return
    }

    page = &entity.LactationPage{
        HasNextPage: l.hasNextPage(lactations),
        NextCursor: nextCursor,
        List: &lactations,
    }

    return page, err

}

func (l *LactationRepository) GetNextPage(cursor string, sort string, direction string) (page *entity.LactationPage, err error) {
    param, id, err:= decodeCursor(cursor)
    if err != nil {
        return
    }
    criteria:= l.createCriteriaNextPage(sort, direction)
	query := fmt.Sprintf(`
        SELECT lac.id, lac.start_date, lac.end_date, lac.production_period, lac.production_total, lac.average_production
            lac.peak_production, lac.isr, lac.observation,
            animal.id as animal_id, animal.identificantion_number as animal_number, animal.name as animal_name,
            animal.ring_order as animal_order, animal.pasture_id as animal_pasture, animal.status as animal_status,
            calf.id as calf_id, calf.sex as calf_sex, calf.birth_date as calf_birth
        FROM lactations as lac
        LEFT JOIN animals as animal ON animal.id = lac.animal_id
        LEFT JOIN animals as calf ON calf.id = lac.animal_id
        %s
        LIMIT %d
        `, criteria, PAGE_LIMIT)
	sqlStatement, err := selectQueryList(query, param, id)
	defer sqlStatement.Close()
	
    if err != nil {
		return 
	}

    var lactations []entity.LactationComplete

    for sqlStatement.Next() {
        var lactation entity.LactationComplete
        err = sqlStatement.Scan(lactation.Id, lactation.StartDate, &lactation.EndDate, &lactation.ProductionPeriod, &lactation.ProductionTotal,
            &lactation.AverageProduction, &lactation.PeakProduction, &lactation.Isr, &lactation.Observation, &lactation.AnimalId, &lactation.AnimalNumber,
            &lactation.AnimalName, &lactation.AnimalOrder, &lactation.AnimalPasture, &lactation.AnimalPasture, 
            &lactation.CalfId, &lactation.CalfSex, &lactation.CalfBirthDate)
        if err != nil {
            return
        }
        lactations = append(lactations, lactation)
    }
    
    nextCursor, err:= l.createNextCursor(sort, lactations)
    if err != nil {
        return
    }

    page = &entity.LactationPage{
        HasNextPage: l.hasNextPage(lactations),
        NextCursor: nextCursor,
        List: &lactations,
    }

    return page, err

}

func (l *LactationRepository) GetByAnimal(animalId string) (arr *[]entity.LactationComplete, err error) {
	query := `
        SELECT lac.id, lac.start_date, lac.end_date, lac.production_period, lac.production_total, lac.average_production
            lac.peak_production, lac.isr, lac.observation,
            calf.id as calf_id, calf.sex as calf_sex, calf.birth_date as calf_birth
        FROM lactation as lac
        LEFT JOIN animals as calf ON calf.id = lac.animal_id
        WHERE lac.animal_id = $1
    `
	sqlStatement, err := selectQueryList(query, animalId)
	defer sqlStatement.Close()
	if err != nil {
		return nil, err
	}

	var entries []entity.LactationComplete

    for sqlStatement.Next() {
        var lactation entity.LactationComplete
        err = sqlStatement.Scan(lactation.Id, lactation.StartDate, lactation.EndDate, lactation.ProductionPeriod, lactation.ProductionTotal,
            lactation.AverageProduction, lactation.PeakProduction, lactation.Isr, lactation.Observation, 
            lactation.CalfId, lactation.CalfSex, lactation.CalfBirthDate)
        if err != nil {
            return
        }
        entries = append(entries, lactation)
    }
    
	return &entries, err
}

func (l *LactationRepository) Add(newLactation *entity.CreateLactation) (*entity.Lactation, error) {
    query:= 
        `INSERT INTO lactations (id, animal_id, calf_id, start_date, end_date, 
            production_period, production_total, avarage_production, peak_production, 
            isr, observation, created_at, deleted_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
    lactation:= new(entity.Lactation).New(newLactation)
    err:= l.saveOrUpdateScan(query, lactation)
    return lactation, err
}

func (l *LactationRepository) Save(lactation *entity.Lactation) error {
    query:= 
        `UPDATE lactations
        SET animal_id = $2, calf_id = $3, start_date = $4, end_date = $5, 
            production_period = $6, production_total = $7, avarage_production = $8, peak_production = $9, 
            isr = $10, observation = $11, created_at = $12, deleted_at = $13)
        WHERE id = $1`
    return l.saveOrUpdateScan(query, lactation)
}

func (l *LactationRepository) Delete(id string) error {
    query:=
        `UPDATE lactations
        SET deleted_at = $1)
        WHERE id = $2`
    return execQuery(query, time.Now(), id)
}
