package butcher

import (
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type ButcherRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *ButcherRepository {
	return &ButcherRepository{db}
}

type SaveValidation struct {
	NameExists    bool `db:"name_exists"`
	CnpjExists    bool `db:"cnpj_exists"`
	AddressExists bool `db:"address_exists"`
}

func (r *ButcherRepository) CheckSaveConflicts(entry *ButcherSave) (*SaveValidation, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM butchers
			WHERE name = :name
				AND id IS DISTINCT FROM :id
				AND user_id = :user_id
				AND deleted_at IS NULL
		) AS name_exists,
		EXISTS (
			SELECT 1
			FROM butchers
			WHERE cnpj = :cnpj
				AND id IS DISTINCT FROM :id
				AND user_id = :user_id
				AND deleted_at IS NULL
		) AS cnpj_exists,
		EXISTS (
			SELECT 1
			FROM butchers
			WHERE address = :address
				AND id IS DISTINCT FROM :id
				AND user_id = :user_id
				AND deleted_at IS NULL
		) AS address_exists
	`
	return util.NamedGet(r.DB, query, SaveValidation{}, entry)
}

func (r *ButcherRepository) CheckDelete(params *ButcherDelete) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM slaughter_entries
			WHERE butcher_id = :id
				AND user_id = :user_id
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.NamedPrimitive(r.DB, query, &exists, params)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *ButcherRepository) FindAll(userId string) (*[]ButcherEntry, error) {

	query := `
		WITH entries_stats AS (
			SELECT
				butcher_id,
				COUNT(*) AS animals_number,
				AVG(dead_weight) AS avg_weight,
				AVG(COALESCE(dead_weight / NULLIF(weight * (1 - discount_rate), 0) , 0)) AS avg_rate
			FROM slaughter_entries
			WHERE deleted_at IS NULL AND user_id = $1
			GROUP BY butcher_id
		)
		SELECT 
			b.id,
			b.name,
			b.cnpj,
			b.discount  AS discount,
			b.address,
			s.animals_number,
			s.avg_weight,
			s.avg_rate
		FROM butchers b
			JOIN entries_stats s ON s.butcher_id = b.id
		WHERE b.user_id = $1 AND b.deleted_at IS NULL
		ORDER BY name
	`
	return util.GetList[ButcherEntry](r.DB, query, userId)
}

func (r *ButcherRepository) FindById(id string, userId string) (*ButcherEntry, error) {

	query := `
		SELECT 
			id,
			name,
			cnpj,
			discount  AS discount,
			address
		FROM butchers 
		WHERE id = $1 
			AND user_id = $2 
			AND deleted_at IS NULL
	`
	return util.GetOne[ButcherEntry](r.DB, query, id, userId)
}

func (r *ButcherRepository) Add(entry *ButcherSave) error {
	query := `
		INSERT INTO butchers (name, cnpj, address, discount, user_id)
		VALUES (:name, :cnpj, :address, CAST(:discount AS float) / 100, :user_id)
	`
	return util.NamedExec(r.DB, query, entry)
}

func (r *ButcherRepository) Update(entry *ButcherSave) (*ButcherEntry, error) {

	query := `
		UPDATE butchers 
		SET name = :name, 
			cnpj = :cnpj, 
			discount = :discount,
			address = :address
		WHERE id = :id AND user_id = :user_id
	`

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return nil, err
	}

	selectQuery := `
		WITH entries_stats AS (
			SELECT
				butcher_id,
				COUNT(*) AS animals_number,
				AVG(dead_weight) AS avg_weight,
				AVG(COALESCE(dead_weight / NULLIF(weight * (1 - discount_rate), 0) , 0)) AS avg_rate
			FROM slaughter_entries
			GROUP BY butcher_id
		)
		SELECT 
			b.id,
			b.name,
			b.cnpj,
			b.discount  AS discount,
			b.address,
			COALESCE(s.animals_number, 0) AS animals_number,
			COALESCE(s.avg_weight, 0) AS avg_weight,
			COALESCE(s.avg_rate, 0) AS avg_rate
		FROM butchers b 
		LEFT JOIN entries_stats s ON s.butcher_id = b.id
		WHERE b.id = :id AND b.user_id = :user_id
	`
	return util.NamedGet(r.DB, selectQuery, ButcherEntry{}, entry)
}

func (r *ButcherRepository) Delete(params *ButcherDelete) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `
		UPDATE butchers
		SET deleted_at = NOW()
		WHERE id = :id AND user_id = :user_id
	`
	err = util.NamedExecTx(tx, query, params)
	if err != nil {
		return err
	}

	if !params.IgnoreDeaths {
		query = `
			UPDATE animals a
			SET death_date = NULL
			FROM slaughter_entries se
			WHERE se.animal_id = a.id
				AND se.butcher_id = :id
				AND se.user_id = :user_id
				AND se.deleted_at IS NULL
		`
		err = util.NamedExecTx(tx, query, params)
		if err != nil {
			return err
		}
	}

	query = `
		UPDATE slaughter_entries
		SET deleted_at = NOW()
		WHERE butcher_id = :id 
			AND user_id = :user_id
			AND deleted_at IS NULL
	`
	err = util.NamedExecTx(tx, query, params)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
