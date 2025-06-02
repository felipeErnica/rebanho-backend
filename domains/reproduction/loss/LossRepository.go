package loss

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type LossRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *LossRepository {
	selectQuery := `
        SELECT loss.*, 
            animals.name as animal_name, animals.number as animal_number, animals.order as animal_order
        FROM Pregnancy_losses as loss
            LEFT JOIN animals ON animals.id = loss.animal_id
    `
	return &LossRepository{selectQuery, "losses", db}
}

func (r *LossRepository) FindPage(pageProps repositoriesUtil.PageProps) (*entity.Page[PregnancyLoss], error) {

	nullFields := []string{"observation"}
	props := repositoriesUtil.PageBuilderProps{
		QueryBody:  r.SelectQuery,
		PageProps:  pageProps,
		NullFields: nullFields,
		TableName:  r.TableName,
		DbConn:     r.Db,
	}
	return repositoriesUtil.BuildPage[PregnancyLoss](props)
}

func (r *LossRepository) FindByAnimalId(animalId string) (*[]PregnancyLoss, error) {
	query := r.SelectQuery + " WHERE loss.animal_id = $1 and loss.deleted_at is null"
	return repositoriesUtil.GetList[PregnancyLoss](r.Db, query, animalId)
}

func (r *LossRepository) FindById(id string) (*PregnancyLoss, error) {
	query := r.SelectQuery + " WHERE loss.id = $1 AND loss.deleted_at is null"
	return repositoriesUtil.GetOne[PregnancyLoss](r.Db, query, id)
}

func (r *LossRepository) Add(newModel *PregnancyLossSave) (*PregnancyLossSave, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newModel)
}

func (r *LossRepository) Update(model *PregnancyLossSave) error {
	return repositoriesUtil.Update(r.Db, r.TableName, model)
}

func (r *LossRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
