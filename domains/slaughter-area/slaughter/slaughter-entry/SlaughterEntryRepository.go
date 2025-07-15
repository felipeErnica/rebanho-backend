package slaughterEntry

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type SlaughterEntryRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *SlaughterEntryRepository {
	selectQuery := `
        SELECT slaughter_entries.*, 
            animals.name as animal_name, animals.number as animal_number, 
            animals.order as animal_order, animals.birth_date as animal_birth,
            group.date as group_date,
            slaughterhouses.name as slaughterhouse
        FROM slaughter_entries
            LEFT JOIN animals ON animals.id = slaughter_entries.animal_id
            LEFT JOIN slaighter_groups as group ON group.id = slaughter_entries.group_id
            LEFT JOIN slaughterhouses ON slaughterhouses.id = group.slaughterhouse_id
    `
	return &SlaughterEntryRepository{selectQuery, "slaughter_entries", db}
}

func (r *SlaughterEntryRepository) FindPage(pageProps repositoriesUtil.PageProps) (*entity.Page[SlaughterEntry], error) {
	props := repositoriesUtil.PageBuilderProps{
		QueryBody:  r.SelectQuery,
		TableName:  r.TableName,
		DbConn:     r.Db,
		PageProps:  pageProps,
	}
	return repositoriesUtil.BuildPage[SlaughterEntry](props)
}

func (r *SlaughterEntryRepository) FindByGroupId(groupId string) (*[]SlaughterEntry, error) {
	query := r.SelectQuery + " WHERE slaughter_entries.user_id = $1 and entry.deleted_at is null"
	return repositoriesUtil.GetList[SlaughterEntry](r.Db, query, groupId)
}

func (r *SlaughterEntryRepository) FindById(id string) (*SlaughterEntry, error) {
	query := r.SelectQuery + " WHERE slaughter_entries.id = $1 and slaughter_entries.deleted_at is null"
	return repositoriesUtil.GetOne[SlaughterEntry](r.Db, query, id)
}

func (r *SlaughterEntryRepository) Add(newModel *SlaughterEntrySave) (*SlaughterEntrySave, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newModel)
}

func (r *SlaughterEntryRepository) Update(model *SlaughterEntrySave) error {
	return repositoriesUtil.Update(r.Db, r.TableName, model)
}

func (r *SlaughterEntryRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
