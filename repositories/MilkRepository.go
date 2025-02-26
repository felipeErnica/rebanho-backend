package repositories

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type MilkRepository struct{}

func (m *MilkRepository) GetByAnimal(animalId string) (*[]entity.MilkEntry, error) {
	query := "SELECT * FROM milk_entries WHERE animal_id = $1"
	sqlStatement, err := SelectQueryList(query)
	defer sqlStatement.Close()
	if err != nil {
		return nil, err
	}

	var entries []entity.MilkEntry

	for sqlStatement.Next() {
		var entry entity.MilkEntry
		err := sqlStatement.Scan(&entry.Id, &entry.AnimalId, &entry.PastureId, &entry.LactationId, &entry.EntryDate, &entry.MilkQuantity)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return &entries, err
}

func (m *MilkRepository) GetByDate(entryDate time.Time) (*[]entity.MilkEntry, error) {
	query := "SELECT * FROM milk_entries WHERE entry_date = $1"
	sqlStatement, err := SelectQueryList(query)
	defer sqlStatement.Close()
	if err != nil {
		return nil, err
	}

	var entries []entity.MilkEntry

	for sqlStatement.Next() {
		var entry entity.MilkEntry
		err := sqlStatement.Scan(&entry.Id, &entry.AnimalId, &entry.PastureId, &entry.LactationId, &entry.EntryDate, &entry.MilkQuantity)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return &entries, err
}

func (m *MilkRepository) GetByLactation(lactationId string) (*[]entity.MilkEntry, error) {
	query := "SELECT * FROM milk_entries WHERE entry_date = $1"
	sqlStatement, err := SelectQueryList(query)
	defer sqlStatement.Close()
	if err != nil {
		return nil, err
	}

	var entries []entity.MilkEntry

	for sqlStatement.Next() {
		var entry entity.MilkEntry
		err := sqlStatement.Scan(&entry.Id, &entry.AnimalId, &entry.PastureId, &entry.LactationId, &entry.EntryDate, &entry.MilkQuantity)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return &entries, err
}
