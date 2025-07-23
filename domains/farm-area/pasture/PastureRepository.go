package pasture

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type PastureRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PastureRepository {
	selectQuery := `
        SELECT pastures.*, bull.name as bull_name, farms.name as farm_name
        FROM pastures
        LEFT JOIN animals as bull ON bull.id = pastures.animal_id
        LEFT JOIN farms ON farms.id = pastures.farm_id
    `
	return &PastureRepository{selectQuery, "pastures", db}
}

func (r *PastureRepository) SearchPasture(userId string, input string, farmsId []string) (*[]entity.SearchEntity, error) {
	args := []any{userId, input}

	arrayStatement := ""
	if len(farmsId) != 0 {
		var buff bytes.Buffer
		for i, id := range farmsId {
			placeholder := fmt.Sprintf("$%d, ", 3+i)
			buff.WriteString(placeholder)
			args = append(args, id)
		}
		placeholderArray := buff.String()
		placeholderArray, _ = strings.CutSuffix(placeholderArray, ", ")
		arrayStatement = fmt.Sprintf("and farm_id in (%s)", placeholderArray)
	}

	query := fmt.Sprintf(`
        select id, name as label 
        from pastures 
        where 
            user_id = $1 
            and name ilike $2 
            %s
            and deleted_at is null
        order by label
    `, arrayStatement)

	return repositoriesUtil.GetList[entity.SearchEntity](r.Db, query, args...)
}

func (r *PastureRepository) FindAnimalsByPasture(
    pastureId string, 
    userId string, 
    sort string, 
    order string,
) (*[]PastureAnimal, error) {
	sortMap := map[string]string{
		"ring_number": "coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0)",
		"name":        "coalesce(animals.name, '')",
		"birth_date":  "coalesce(animals.birth_date, '-infinity')",
		"death_date":  "coalesce(animals.death_date, '-infinity')",
	}

	expression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query := `
        select
            animals.id, 
            animals.name, 
            animals.ring_number, 
            animals.sex, 
            animals.father_id, 
            animals.mother_id, 
            animals.birth_date, 
            animals.death_date, 
            animals.animal_type,
            concat_ws(' - ', father.ring_number, father.name) as father_name, 
            concat_ws(' - ', mother.ring_number, mother.name) as mother_name
        from animals
            left join animals as father on father.id = animals.father_id
            left join animals as mother on mother.id = animals.mother_id
        where animals.pasture_id = $1 and animals.user_id = $2 and animals.deleted_at is null
    `
	orderExpression := " order by " + expression
	query = query + orderExpression

	return repositoriesUtil.GetList[PastureAnimal](r.Db, query, pastureId, userId)
}

func (r *PastureRepository) FindAll(userId string) (*[]Pasture, error) {
	query := r.SelectQuery + " where pastures.user_id = $1 AND pastures.deleted_at is null"
	return repositoriesUtil.GetList[Pasture](r.Db, query)
}

func (r *PastureRepository) FindById(id string) (*Pasture, error) {
	query := r.SelectQuery + " WHERE pastures.id = $1 AND pastures.deleted_at is null"
	return repositoriesUtil.GetOne[Pasture](r.Db, query, id)
}

func (r *PastureRepository) Add(newPasture *PastureSave) (*PastureSave, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newPasture)
}

func (r *PastureRepository) Update(pasture *PastureSave) error {
	return repositoriesUtil.Update(r.Db, r.TableName, pasture)
}

func (r *PastureRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
