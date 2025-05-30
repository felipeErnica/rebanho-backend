package pastureEntries

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type PastureEntryRepository struct {
	SelectQuery string
    TableName string
    Db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PastureEntryRepository {
    selectQuery := `
        SELECT pastures_entries.*, 
            animals.name as animal_name, animals.number as animal_number,
            pastures.name as pasture_name
        FROM pastures_entries
            LEFT JOIN animals ON animals.id = pastures_entries.animal_id
            LEFT JOIN pastures ON pastures.id = pastures_entries.pasture_id
    `
    return &PastureEntryRepository{selectQuery, "pastures_entries", db}
}

func (r *PastureEntryRepository) FindByAnimalId(animalId string) (*[]PastureEntry, error) {
    query := r.SelectQuery + " WHERE pastures_entries.deleted_at is null and pastures_entries.animal_id = $1"
	return repositoriesUtil.GetList[PastureEntry](r.Db, query, animalId)
}

func (r *PastureEntryRepository) Add(newEntry *PastureEntrySave) (*PastureEntrySave, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newEntry)
}

func (r *PastureEntryRepository) Update(entry *PastureEntrySave) error {
	return repositoriesUtil.Update(r.Db, r.TableName, entry)
}

func (r *PastureEntryRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
