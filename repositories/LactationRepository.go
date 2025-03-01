package repositories

import (
	"database/sql"
	"fmt"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type LactationRepository struct{}

func (l *LactationRepository) returnSimpleQuery(criteria string) string {
    return fmt.Sprintf(`SELECT id, animal_id, calf_id, start_date, end_date, 
            production_period, production_total, avarage_production, peak_production, 
            isr, observation, created_at, deleted_at
        FROM lactations
        %s`, criteria)
}

func (l *LactationRepository) returnFirstPageQuery(criteria string) string {
    return fmt.Sprintf(`SELECT *
        FROM (%s)
        ORDER BY created_at, id
        LIMIT %d`, l.returnSimpleQuery(criteria), PAGE_LIMIT)
}

func (l *LactationRepository) returnPageQuery(criteria string) string {
    return fmt.Sprintf(`SELECT *
        FROM (%s)
        WHERE (created_at, id) < ($1, $2)
        ORDER BY created_at, id
        LIMIT %d`, l.returnSimpleQuery(criteria), PAGE_LIMIT)
}
 func (l *LactationRepository) createPage(arr []entity.Lactation) *entity.LactationPage {
    lastEntry:=arr[len(arr) - 1]
    return &entity.LactationPage{
        List: &arr,
        HasNextPage: len(arr) == PAGE_LIMIT,
        NextCursor: encodeCursor(lastEntry.CreatedAt, lastEntry.Id),
    }
} 


func (l *LactationRepository) ScanQueryRows(sqlStatement *sql.Rows) (entity.Lactation, error) {
    var entry entity.Lactation
    err := sqlStatement.Scan(&entry.Id, &entry.AnimalId, &entry.StartDate, &entry.EndDate, &entry.ProductionPeriod, &entry.ProductionTotal,
        &entry.AvarageProduction, &entry.PeakProduction, &entry.Isr, &entry.Observation)
    return entry, err
}


func (l *LactationRepository) GetFirstPage() (*entity.LactationPage, error) {
	query := l.returnFirstPageQuery("ORDER BY start_date")
	sqlStatement, err := selectQueryList(query)
	defer sqlStatement.Close()
	if err != nil {
		return nil, err
	}

	var entries []entity.Lactation

	for sqlStatement.Next() {
        entry, err:= l.ScanQueryRows(sqlStatement)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

    page:= l.createPage(entries)
	return page, err
}

func (l *LactationRepository) GetNextPage(cursor string) (*entity.LactationPage, error) {
    createdAt, id, err:= decodeCursor(cursor)
    if err != nil {
        return nil, err
    }

    query:=l.returnPageQuery("ORDER BY start_date")
	sqlStatement, err := selectQueryList(query, createdAt, id)
	defer sqlStatement.Close()
	if err != nil {
		return nil, err
	}

	var entries []entity.Lactation

	for sqlStatement.Next() {
        entry, err:= l.ScanQueryRows(sqlStatement)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
    
    page:=l.createPage(entries)
	return page, err
}

func (l *LactationRepository) GetByAnimal(animalId string) (*[]entity.Lactation, error) {
	query := l.returnSimpleQuery("WHERE animal_id = $1")
	sqlStatement, err := selectQueryList(query, animalId)
	defer sqlStatement.Close()
	if err != nil {
		return nil, err
	}

	var entries []entity.Lactation

	for sqlStatement.Next() {
        entry, err:= l.ScanQueryRows(sqlStatement)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return &entries, err
}
