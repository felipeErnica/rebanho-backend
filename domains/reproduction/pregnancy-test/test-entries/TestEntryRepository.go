package testEntries

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type TestEntryRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *TestEntryRepository {
	selectQuery := `
		SELECT pregnancy_test_entries.*, 
			group.date as group_date,
			animals.name as animal_name, animals.number as animal_number, animals.order as animal_order
		FROM pregnancy_test_entries 
			LEFT JOIN pregnancy_test_groups as group ON group.id = pregnancy_test_entries.group_id
			LEFT JOIN animals ON animals.id = pregnancy_test_entries.animal_id
	`
	return &TestEntryRepository{selectQuery, "pregnancy_test_entries", db}
}

func (r *TestEntryRepository) FindByGroupId(groupId string) (*[]TestEntry, error) {
	query := r.SelectQuery + " WHERE entries.group_id = $1 AND entries.deleted_at is null"
	return repositoriesUtil.GetList[TestEntry](r.Db, query, groupId)
}

func (r *TestEntryRepository) FindByAnimalId(animalId string) (*[]TestEntry, error) {
	query := r.SelectQuery + " WHERE entries.animal_id = $1 AND entries.deleted_at is null"
	return repositoriesUtil.GetList[TestEntry](r.Db, query, animalId)
}

func (r *TestEntryRepository) FindById(id string) (*TestEntry, error) {
	query := r.SelectQuery + " WHERE entries.id = $1 AND entries.deleted_at is null"
	return repositoriesUtil.GetOne[TestEntry](r.Db, query, id)
}

func (r *TestEntryRepository) Add(newModel *TestEntrySave) (*TestEntrySave, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newModel)
}

func (r *TestEntryRepository) Update(model *TestEntrySave) error {
	return repositoriesUtil.Update(r.Db, r.TableName, model)
}

func (r *TestEntryRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
