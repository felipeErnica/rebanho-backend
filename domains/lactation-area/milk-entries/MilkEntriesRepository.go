package milkEntries

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type MilkRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *MilkRepository {
	selectQuery := `
    SELECT milk_entries.*, 
        animals.name as animal_name, animals.order as animal_order, animals.number as animal_number
        pastures.name as pasture_name
    FROM milk_entries
        LEFT JOIN animals ON animals.id = milk_entries.animal_id 
        LEFT JOIN pastures ON pastures.id = milk_entries.pasture_id
    `
	return &MilkRepository{selectQuery, "milk_entries", db}
}

func (r *MilkRepository) FindPage(pageProps repositoriesUtil.PageProps) (*entity.Page[MilkEntry], error) {
	props := repositoriesUtil.PageBuilderProps{
		QueryBody:  r.SelectQuery,
		TableName:  r.TableName,
		DbConn:     r.Db,
        PageProps: pageProps,
	}
	return repositoriesUtil.BuildPage[MilkEntry](props)
}

func (r *MilkRepository) FindByAnimal(animalId string) (*[]MilkEntry, error) {
	query := r.SelectQuery + " WHERE milk_entries.deleted_at is null and milk_entries.animal_id = $1"
	return repositoriesUtil.GetList[MilkEntry](r.Db, query, animalId)
}

func (r *MilkRepository) Add(newMilk *MilkEntrySave) (*MilkEntrySave, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newMilk)
}

func (r *MilkRepository) Update(milk *MilkEntrySave) error {
	return repositoriesUtil.Update(r.Db, r.TableName, milk)
}

func (r *MilkRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
