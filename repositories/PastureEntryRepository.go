package repositories

import "github.com/felipeErnica/rebanho-backend/entity"

type PastureEntryRepository struct {}

func (p *PastureEntryRepository) GetByPastureId(pastureId string) (*[]entity.PastureEntry, error) {
	query := "SELECT * FROM pasture_entries WHERE pasture_id = $1"
	sqlStatement, err := selectQueryList(query)
	defer sqlStatement.Close()
	if err != nil {
		return nil, err
	}

	var entries []entity.PastureEntry

	for sqlStatement.Next() {
		var entry entity.PastureEntry
		err := sqlStatement.Scan(&entry.Id, &entry.AnimalId, &entry.PastureId, &entry.EntryDate, &entry.ExitDate)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return &entries, err
}

func (p *PastureEntryRepository) GetByAnimalId(animalId string) (*entity.PastureEntry, error) {
	query := "SELECT * FROM pasture_entries WHERE animal_id = $1"
    var entry entity.PastureEntry
	sqlStatement:= selectQueryOne(query, animalId)
    err:= sqlStatement.Scan(&entry.Id, &entry.AnimalId, &entry.PastureId, &entry.EntryDate, &entry.ExitDate)
    return &entry, err
}

func (p *PastureEntryRepository) Add(newEntry *entity.CreatePastureEntry) (*entity.PastureEntry, error) {
    query:="INSERT INTO pasture_entries(id, animal_id, pasture_id, entry_date, exit_date) VALUES ($1, $2, $3)"
    entry:= new(entity.PastureEntry).New(newEntry)
    err:= execQuery(query, entry.Id, entry.AnimalId, entry.PastureId, entry.EntryDate, entry.ExitDate)
    return entry, err
}

func (p *PastureEntryRepository) Save(entry *entity.PastureEntry) (*entity.PastureEntry, error) {
    query:="UPDATE pasture_entries SET animal_id = $1, pasture_id = $2, entry_date = $3, exit_date =$4 WHERE id = $5"
    err:= execQuery(query, entry.AnimalId, entry.PastureId, entry.EntryDate, entry.ExitDate, entry.Id)
    return entry, err
}
