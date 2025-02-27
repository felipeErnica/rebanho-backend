package repositories

import (
	"database/sql"
	"fmt"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type LactationRepository struct{}

func (l *LactationRepository) ScanQueryRows(sqlStatement *sql.Rows) (entity.Lactation, error) {
    var entry entity.Lactation
    err := sqlStatement.Scan(&entry.Id, &entry.AnimalId, &entry.StartDate, &entry.EndDate, &entry.ProductionPeriod, &entry.ProductionTotal,
        &entry.AvarageProduction, &entry.PeakProduction, &entry.Isr, &entry.Observation)
    return entry, err
}


func (l *LactationRepository) GetFirstPage() (*[]entity.Lactation, error) {
	query := fmt.Sprintf(`SELECT * FROM lactations ORDER BY start_date DESC LIMIT %d`, PAGE_LIMIT)
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

	return &entries, err
}

func (l *LactationRepository) GetNextPage(cursor string) (*entity.LactationPage, error) {
    
    createdAt, id, err:= decodeCursor(cursor)
    if err != nil {
        return nil, err
    }

	query := fmt.Sprintf(`SELECT * 
        FROM lactations 
        WHERE (id, created_at) < (%s,%s)
        ORDER BY start_date DESC, id DESC, created_at DESC
        DESC LIMIT %d`, id, createdAt, PAGE_LIMIT)

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
    
    lastEntry:=entries[len(entries) - 1]

    page:=&entity.LactationPage{
        NextCursor: encodeCursor(lastEntry.CreatedAt, lastEntry.Id),
        HasNextPage: len(entries) < int(PAGE_LIMIT),
        List: &entries,
    }

	return page, err
}

func (l *LactationRepository) GetByAnimal(animalId string) (*[]entity.Lactation, error) {
	query := "SELECT * FROM lactations WHERE animal_id = $1"
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

	return &entries, err
}
