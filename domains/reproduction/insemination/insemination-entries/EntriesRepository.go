package inseminationEntries

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type EntriesRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) *EntriesRepository {
	selectQuery := `
        SELECT entries.*, 
            animals.number as animal_number, animals.order as animal_order, animals.name as animal_name,
            group.date as group_date, bull.name as bull_name
        FROM insemination_entries as entries
            LEFT JOIN animals ON animals.id = entries.animal_id
            LEFT JOIN insemination_groups as group ON group.id = entries.group_id
            LEFT JOIN animals as bull ON bull.id = group.bull_id
    `
	return &EntriesRepository{selectQuery, "entries", db}
}

func (r *EntriesRepository) FindByGroupId(groupId string) (*[]InseminationEntry, error) {
	query := r.SelectQuery + " WHERE entries.group_id = $1"
	return repositoriesUtil.GetList[InseminationEntry](r.Db, query, groupId)
}

func (r *EntriesRepository) FindById(id string) (*InseminationEntry, error) {
	query := r.SelectQuery + " WHERE entries.id = $1"
	return repositoriesUtil.GetOne[InseminationEntry](r.Db, query, id)
}

func (r *EntriesRepository) Add(newModel *InseminationEntrySave) (*InseminationEntrySave, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newModel)
}

func (r *EntriesRepository) Update(model *InseminationEntrySave) error {
	return repositoriesUtil.Update(r.Db, r.TableName, model)
}

func (r *EntriesRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
