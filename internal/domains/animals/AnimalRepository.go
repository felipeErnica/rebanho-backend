package animals

import (
	"fmt"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type AnimalRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *AnimalRepository {
	return &AnimalRepository{db}
}

type WarnRecordValidations struct {
	HasLactation          bool `db:"has_lactation"`
	HasSlaughter          bool `db:"has_slaughter"`
	HasInsemination       bool `db:"has_insemination"`
	HasBreeding           bool `db:"has_breeding"`
	HasTransfer           bool `db:"has_transfer"`
	HasBullPastureEntries bool `db:"has_bull_pastures_entries"`
}

type ErrorRecordValidations struct {
	HasChildren        bool `db:"has_children"`
	IsCalfInLactation  bool `db:"is_calf_lac"`
	IsEmbryoDonor      bool `db:"is_embryo_donor"`
	IsEmbryoBull       bool `db:"is_embryo_bull"`
	IsInseminationBull bool `db:"is_insemination_bull"`
	IsBreedingBull     bool `db:"is_breeding_bull"`
}

type ValidationResult struct {
	NumberExists     bool `db:"number_exists"`
	NameExists       bool `db:"name_exists"`
	DeadNumberExists bool `db:"dead_number_exists"`
	DeadNameExists   bool `db:"dead_name_exists"`
}

func (r *AnimalRepository) CheckSaveConflicts(entry *AnimalSave) (*ValidationResult, error) {
	query := `
		SELECT
			EXISTS (
				SELECT 1
				FROM animals
				WHERE tag = :tag
					AND death_date IS NULL
					AND name IS NOT NULL
					AND (:id IS NULL OR id <> :id)
					AND deleted_at IS NULL
					AND user_id = :user_id
			) AS number_exists,
			EXISTS (
				SELECT 1
				FROM animals
				WHERE name = :name
					AND death_date IS NULL
					AND (:id IS NULL OR id <> :id)
					AND deleted_at IS NULL
					AND user_id = :user_id
			) AS name_exists,
			EXISTS (
				SELECT 1
				FROM animals
				WHERE tag = :tag
					AND name IS NOT NULL
					AND death_date IS NOT NULL
					AND deleted_at IS NULL
					AND (:id IS NULL OR id <> :id)
					AND user_id = :user_id
			) AS dead_number_exists,
			EXISTS (
				SELECT 1
				FROM animals
				WHERE name = :name
					AND death_date IS NOT NULL
					AND deleted_at IS NULL
					AND (:id IS NULL OR id <> :id)
					AND user_id = :user_id
			) AS dead_name_exists
	`
	return util.NamedGet(r.DB, query, ValidationResult{}, entry)
}

func (r *AnimalRepository) CheckDeleteErrorConditions(id string, userId string) (*ErrorRecordValidations, error) {

	query := `
		SELECT 
			EXISTS (
				SELECT 1
				FROM animals
				WHERE deleted_at IS NULL
					AND (mother_id = $1 OR father_id = $1)
					AND user_id = $2
			) AS has_children,
			EXISTS (
				SELECT 1
				FROM lactations
				WHERE deleted_at IS NULL
					AND calf_id = $1
					AND user_id = $2
			) AS is_calf_lac,
			EXISTS (
				SELECT 1
				FROM embryo_transfer
				WHERE donor_id = $1
					AND user_id = $2
					AND deleted_at IS NULL
			) AS is_embryo_donor,
			EXISTS (
				SELECT 1
				FROM embryo_transfer
				WHERE bull_id = $1
					AND user_id = $2
					AND deleted_at IS NULL
			) AS is_embryo_bull,
			EXISTS (
				SELECT 1
				FROM insemination_entries
				WHERE bull_id = $1
					AND user_id = $2
					AND deleted_at IS NULL
			) AS is_insemination_bull,
			EXISTS (
				SELECT 1
				FROM breeding_entries
				WHERE bull_id = $1
					AND user_id = $2
					AND deleted_at IS NULL
			) AS is_breeding_bull
	`
	return util.GetOne[ErrorRecordValidations](r.DB, query, id, userId)
}

func (r *AnimalRepository) CheckDeleteWarningConditions(id string, userId string) (*WarnRecordValidations, error) {
	query := `
		SELECT 
			EXISTS (
				SELECT 1
				FROM lactations
				WHERE deleted_at IS NULL
					AND animal_id = $1
					AND user_id = $2
			) AS has_lactation,
			EXISTS (
				SELECT 1
				FROM breeding_entries
				WHERE deleted_at IS NULL
					AND animal_id = $1
					AND user_id = $2
			) AS has_breeding,
			EXISTS (
				SELECT 1
				FROM embryo_transfer
				WHERE deleted_at IS NULL
					AND receiver_id = $1
					AND user_id = $2
			) AS has_transfer,
			EXISTS (
				SELECT 1
				FROM slaughter_entries
				WHERE deleted_at IS NULL
					AND animal_id = $1
					AND user_id = $2
			) AS has_slaughter,
			EXISTS (
				SELECT 1
				FROM insemination_entries
				WHERE deleted_at IS NULL
					AND animal_id = $1
					AND user_id = $2
			) AS has_insemination,
			EXISTS (
				SELECT 1
				FROM pasture_entries pe
					JOIN animals a ON a.id = pe.animal_id
				WHERE pe.deleted_at IS NULL
					AND a.sex = 'M'
					AND a.animal_type = 'REPRODUCTION_ANIMAL'
					AND pe.animal_id = $1
					AND pe.user_id = $2
			) AS has_bull_pastures_entries
	`
	return util.GetOne[WarnRecordValidations](r.DB, query, id, userId)
}

func (r *AnimalRepository) GetBirthHist(userId string) (*[]util.GraphData, error) {

	query := `
		WITH cte AS (
			SELECT
				DATE_TRUNC('month', birth_date) AS date,
				COUNT(*) AS value
			FROM animals
			WHERE user_id = $1 
				AND deleted_at IS NULL
				AND birth_date IS NOT NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 12
		)
		SELECT * FROM cte ORDER BY date
	`

	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *AnimalRepository) GetDairyHist(userId string) (*[]util.GraphData, error) {

	query := `
		WITH calendar AS (
			SELECT GENERATE_SERIES(
				DATE_TRUNC('month', MAX(end_date) - INTERVAL '12 months'),
				DATE_TRUNC('month', MAX(end_date)),
				INTERVAL '1 month'
			) AS entry_date
			FROM lactations
		)
		SELECT
			c.entry_date AS date,
			COUNT(*) AS value
		FROM lactations 
			JOIN calendar c ON start_date <= c.entry_date
				AND c.entry_date <= COALESCE(end_date, NOW())
		WHERE user_id = $1 AND deleted_at IS NULL
		GROUP BY 1
		ORDER BY 1
	`

	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *AnimalRepository) GetDeathHist(userId string) (*[]util.GraphData, error) {

	query := `
		WITH cte AS (
			SELECT
				DATE_TRUNC('month', death_date) AS date,
				COUNT(*) AS value
			FROM animals a
			WHERE user_id = $1 
				AND deleted_at IS NULL
				AND death_date IS NOT NULL
				AND NOT EXISTS (
					SELECT 1
					FROM slaughter_entries s
					WHERE s.animal_id = a.id
						AND a.death_date = s.entry_date
						AND s.user_id = $1
				)
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 12
		)
		SELECT *
		FROM cte 
		ORDER BY date
	`

	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *AnimalRepository) GetSlaughterHist(userId string) (*[]util.GraphData, error) {

	query := `
		WITH cte AS (
			SELECT
				entry_date AS date,
				COUNT(*) AS value
			FROM slaughter_entries
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 12
		)
		SELECT *
		FROM cte
		ORDER BY date
	`

	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *AnimalRepository) GetAnimalTypes(userId string) (*AnimalByType, error) {

	query := `
		SELECT
			COUNT(*) FILTER (WHERE animal_type = 'DAIRY_ANIMAL') AS dairy_animals,
			COUNT(*) FILTER (WHERE animal_type = 'OFFSPRING') AS offspring,
			COUNT(*) FILTER (WHERE animal_type = 'BEEF_ANIMAL') AS beef_animals,
			COUNT(*) FILTER (WHERE animal_type = 'REPRODUCTION_ANIMAL') AS reproduction_animals
		FROM animals a
		WHERE user_id = $1 
			AND deleted_at IS NULL
			AND is_outside_animal = FALSE
			AND death_date IS NULL
	`
	result, err := util.GetOne[AnimalByType](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AnimalRepository) GetLastDeaths(userId string) (*[]AnimalDB, error) {

	query := `
		SELECT
			a.id,
			a.tag,
			a.name,
			a.sex,
			a.animal_type,
			a.birth_date,
			a.death_date,
			a.observation,

			a.father_id,
			f.tag,
			f.name,
			
			a.mother_id,
			m.tag,
			m.name
		FROM animals a
		LEFT JOIN animals f ON f.id = a.father_id
		LEFT JOIN animals m ON m.id = a.mother_id
		WHERE a.user_id = $1 
			AND a.deleted_at IS NULL
			AND a.is_outside_animal = FALSE
			AND a.death_date IS NOT NULL
			AND NOT EXISTS (
				SELECT 1
				FROM slaughter_entries s
				WHERE s.animal_id = a.id
					AND s.user_id = $1
					AND s.deleted_at IS NULL
			)
		ORDER BY a.death_date DESC
		LIMIT 20
	`
	result, err := util.GetList[AnimalDB](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AnimalRepository) GetAgeAndSex(userId string) (*[]AnimalsByAge, error) {
	query := `
        SELECT 
            CASE 
                WHEN age(birth_date) < INTERVAL '2 months' THEN '0-2 meses'
                WHEN age(birth_date) BETWEEN INTERVAL '2 months' AND INTERVAL '8 months' THEN '2-8 meses'
                WHEN age(birth_date) BETWEEN INTERVAL '8 months' AND INTERVAL '12 months' THEN '8-12 meses'
                WHEN age(birth_date) BETWEEN INTERVAL '12 months' AND INTERVAL '24 months' THEN '12-24 meses'
                WHEN age(birth_date) BETWEEN INTERVAL '24 months' AND INTERVAL '36 months' THEN '24-36 meses'
                WHEN age(birth_date) > INTERVAL '36 months' THEN '+36 meses'
				ELSE 'Sem Data'
            END AS category,
			COUNT(*) FILTER (WHERE sex = 'M') AS male,
			COUNT(*) FILTER (WHERE sex = 'F') AS female
        FROM animals
        WHERE user_id = $1
            AND deleted_at IS NULL
			AND is_outside_animal = FALSE
			AND birth_date IS NOT NULL
			AND death_date IS NULL
		GROUP BY category
		ORDER BY MIN(birth_date) DESC
    `
	return util.GetList[AnimalsByAge](r.DB, query, userId)
}

func (r *AnimalRepository) FindPage(
	userId string,
	cursor string,
	sort string,
	order string,
	limit int,
	filter *AnimalFilter,
) (page *[]AnimalDB, err error) {

	sortMap := map[string]util.SortField{
		"name":                   {Field: "coalesce(a.name, '')", Order: "asc"},
		"average_birth_interval": {Field: "coalesce(b.average_birth_interval, 0)", Order: "asc"},
		"average_lac_interval":   {Field: "coalesce(l.average_lac_interval, 0)", Order: "asc"},
		"average_prod":           {Field: "coalesce(l.average_prod, 0)", Order: "asc"},
		"average_peak":           {Field: "coalesce(l.average_peak, 0)", Order: "asc"},
		"death_date":             {Field: "coalesce(a.death_date, '-infinity')", Order: "asc"},
		"weaning_date":           {Field: "coalesce(a.weaning_date, '-infinity')", Order: "asc"},
		"birth_date":             {Field: "coalesce(a.birth_date, '-infinity')", Order: "desc"},
		"animal_order":           {Field: "coalesce(nullif(regexp_replace(a.tag, '[^0-9]', '', 'g'), '')::int, 0)", Order: "asc"},
		"created_at":             {Field: "a.created_at", Order: "asc"},
		"id":                     {Field: "a.id", Order: "asc"},
	}

	query := `
		WITH milk_stats AS (
			SELECT
				l.id AS lactation_id,
				AVG(m.quantity) AS avg_prod,
				MAX(m.quantity) AS peak
			FROM milk_entries m
				JOIN lactations l
					ON l.animal_id = m.animal_id
				   AND l.start_date <= m.entry_date
				   AND COALESCE(l.end_date, NOW()) >= m.entry_date
				   AND l.deleted_at IS NULL
			WHERE m.deleted_at IS NULL
			GROUP BY l.id
		),

		lac_interval_cte AS (
			SELECT
				l.animal_id,
				l.start_date,
				l.end_date,
				EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) AS lac_interval,
				p.avg_prod,
				p.peak
			FROM lactations l
				LEFT JOIN milk_stats p ON p.lactation_id = l.id
			WHERE l.user_id = $1 AND l.deleted_at IS NULL
		),

		lac_stats AS (
			SELECT
				animal_id,
				AVG(lac_interval) AS average_lac_interval,
				AVG(avg_prod) AS average_prod,
				AVG(peak) AS average_peak
			FROM lac_interval_cte
			GROUP BY animal_id
		),

		birth_interval_cte AS (
			SELECT
				mother_id,
				EXTRACT(days FROM birth_date - LAG(birth_date) OVER (PARTITION BY mother_id ORDER BY birth_date)) AS birth_interval
			FROM animals
			WHERE user_id = $1
			  AND deleted_at IS NULL
			  AND mother_id IS NOT NULL
		),

		birth_stats AS (
			SELECT
				mother_id,
				AVG(birth_interval) AS average_birth_interval
			FROM birth_interval_cte
			GROUP BY mother_id
		),

		pasture_cte AS (
			SELECT 
				pe.animal_id,
				pe.pasture_id,
				p.name as pasture_name,
				p.farm_id,
				f.name as farm_name
			FROM pasture_entries pe
			JOIN pastures p ON p.id = pe.pasture_id
			JOIN farms f ON f.id = p.farm_id
			WHERE pe.user_id = $1 
				and pe.deleted_at IS NULL
				and pe.exit_date IS NULL
		)

		SELECT
			a.id,
			a.tag,
			a.name,
			a.sex,
			a.birth_date,
			a.weight_birth,
			a.weaning_date,
			a.death_date,
			a.animal_type,
			l.average_prod,
			l.average_lac_interval,
			l.average_peak,
			b.average_birth_interval,
			a.observation,

			a.father_id,
			f.name as father_name,
			f.tag as father_tag,

			a.mother_id,
			m.name as mother_name,
			m.tag as mother_tag,

			p.pasture_id,
			p.pasture_name,
			p.farm_id,
			p.farm_name,

			EXISTS (
				SELECT 1
				FROM lactations
				WHERE user_id = $1
					AND deleted_at IS NULL
					AND animal_id = a.id
			) AS is_lactating,
			COALESCE(NULLIF(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g'), '')::int, 0) as animal_order,
			a.created_at
		FROM animals a
		LEFT JOIN animals m ON m.id = a.mother_id
		LEFT JOIN animals f ON f.id = a.father_id
		LEFT JOIN birth_stats b ON b.mother_id = a.id
		LEFT JOIN lac_stats l ON l.animal_id = a.id
		LEFT JOIN pasture_cte p ON p.animal_id = a.id
	`

	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "a", 2)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := util.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	mainExpression := "a.deleted_at IS NULL and a.user_id = $1"
	whereExpression := util.GetWhereExpression(mainExpression, filterExpression, cursorExpression)

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	orderExpression := " ORDER BY " + sortExpression
	query += whereExpression + orderExpression + fmt.Sprintf(" LIMIT %d", limit)

	return util.GetList[AnimalDB](r.DB, query, args...)
}

func (r *AnimalRepository) GetPageFoot(userId string, filter *AnimalFilter) (*AnimalFoot, error) {

	query := `
		WITH milk_stats AS (
			SELECT
				l.id AS lactation_id,
				AVG(m.quantity) AS avg_prod,
				MAX(m.quantity) AS peak
			FROM milk_entries m
			JOIN lactations l
				ON l.animal_id = m.animal_id
			   AND l.start_date <= m.entry_date
			   AND COALESCE(l.end_date, NOW()) >= m.entry_date
			   AND l.deleted_at IS NULL
			WHERE m.deleted_at IS NULL
			GROUP BY l.id
		),

		lac_interval_cte AS (
			SELECT
				l.animal_id,
				l.start_date,
				l.end_date,
				EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) AS lac_interval,
				p.peak,
				p.avg_prod
			FROM lactations l
			JOIN milk_stats p ON p.lactation_id = l.id
			WHERE l.user_id = $1 AND l.deleted_at IS NULL
		),

		lac_stats AS (
			SELECT
				animal_id,
				AVG(lac_interval) AS average_lac_interval,
				AVG(avg_prod) AS average_prod,
				AVG(peak) AS average_peak
			FROM lac_interval_cte
			GROUP BY animal_id
		),

		birth_interval_cte AS (
			SELECT
				mother_id,
				EXTRACT(days FROM birth_date - LAG(birth_date) OVER (PARTITION BY mother_id ORDER BY birth_date)) AS birth_interval
			FROM animals
			WHERE user_id = $1
			  AND deleted_at IS NULL
			  AND mother_id IS NOT NULL
		),

		birth_stats AS (
			SELECT
				mother_id,
				COUNT(*) AS children_number,
				AVG(birth_interval) AS average_birth_interval
			FROM birth_interval_cte
			GROUP BY mother_id
		),

		cte AS (
			SELECT
				a.father_id,
				a.mother_id,
				COALESCE(NULLIF(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g'), '')::int, 0) animal_order,
				a.tag,
				a.name,
				a.sex,
				a.birth_date,
				a.weight_birth,
				a.weaning_date,
				a.death_date,
				a.animal_type,
				l.average_prod,
				l.average_lac_interval,
				l.average_peak,
				b.average_birth_interval,
				a.observation,
				EXISTS (
					SELECT 1
					FROM lactations
					WHERE user_id = $1
						AND deleted_at IS NULL
						AND animal_id = a.id
				) AS is_lactating,
				a.is_outside_animal,
				a.is_embryo_donor,
				a.is_breeding_bull,
				a.is_insemination_bull,
				a.is_transfer_bull
			FROM animals a
				LEFT JOIN animals m ON m.id = a.mother_id
				LEFT JOIN animals f ON f.id = a.father_id
				LEFT JOIN birth_stats b ON b.mother_id = a.id
				LEFT JOIN lac_stats l ON l.animal_id = a.id
			WHERE a.user_id = $1 AND a.deleted_at IS NULL
		)

		SELECT 
			COUNT(*) AS total,
			AVG(cte.average_prod) AS average_prod,
			AVG(cte.average_lac_interval) AS average_lac_interval,
			AVG(cte.average_birth_interval) AS average_birth_interval,
			AVG(cte.average_peak) AS average_peak
		FROM cte
	`
	filterExpression, _, err := util.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExpression)
	query += whereExpression

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return util.GetOne[AnimalFoot](r.DB, query, args...)
}

func (r *AnimalRepository) FindById(id string, userId string) (*AnimalDB, error) {

	query := `
		WITH pasture_cte AS (
			SELECT 
				pe.pasture_id,
				pe.animal_id,
				p.name as pasture_name,
				p.farm_id,
				f.name as farm_name
			FROM pasture_entries pe
			JOIN pastures p ON p.id = pe.pasture_id
			JOIN farms f ON f.id = p.farm_id
			WHERE pe.user_id = $2 AND pe.animal_id = $1
			ORDER BY entry_date DESC
			LIMIT 1
		)

        SELECT 
			a.id,
			a.tag,
			a.name,
			a.sex,
			a.death_date,
			a.weaning_date,
			a.animal_type,
			a.is_embryo_donor,
			a.is_transfer_bull,
			a.is_breeding_bull,
			a.is_insemination_bull,
			a.is_outside_animal,
			a.observation,

			a.father_id,
			f.tag AS father_tag,
			f.name AS father_name,
			
			a.mother_id,
			m.name AS mother_name,
			m.tag AS mother_tag,

			p.pasture_id,
			p.pasture_name,
			p.farm_id,
			p.farm_name
        FROM animals a
		CROSS JOIN pasture_cte p
		LEFT JOIN animals f ON f.id = a.father_id
		LEFT JOIN animals m ON m.id = a.mother_id
		WHERE a.id = $1 AND a.user_id = $2
	`
	return util.GetOne[AnimalDB](r.DB, query, id, userId)
}

func (r *AnimalRepository) Search(
	sort string,
	order string,
	filter *AnimalFilter,
	userId string,
) (*[]AnimalDB, error) {

	sortMap := map[string]util.SortField{
		"name":         {Field: "COALESCE(cte.name, '')", Order: "ASC"},
		"birth_date":   {Field: "COALESCE(cte.birth_date, '-infinity')", Order: "DESC"},
		"animal_order": {Field: "COALESCE(NULLIF(REGEXP_REPLACE(cte.tag, '[^0-9]', '', 'g'), '')::int, 0)", Order: "ASC"},
	}

	query := `
		WITH pasture_cte AS (
			SELECT 
				pe.animal_id,
				pe.pasture_id,
				p.name AS pasture_name,
				p.farm_id,
				f.name AS farm_name
			FROM pasture_entries pe
			JOIN pastures p ON p.id = pe.pasture_id
			JOIN farms f ON f.id = p.farm_id
			WHERE pe.user_id = $1 
				AND pe.deleted_at IS NULL
				AND pe.exit_date IS NULL
		),
		
        cte AS (
			SELECT 
				a.id,
				a.tag,
				a.name,
				a.sex,
				a.animal_type,
				a.birth_date,
				a.death_date,
				a.weaning_date,
				a.is_embryo_donor,
				a.is_transfer_bull,
				a.is_breeding_bull,
				a.is_insemination_bull,
				a.is_outside_animal,
				a.observation,

				a.father_id,
				f.name AS father_name,
				f.tag AS father_tag,

				a.mother_id,
				m.tag AS mother_tag,
				m.name AS mother_name,

				p.pasture_id,
				p.pasture_name,
				p.farm_id,
				p.farm_name,

				(
					SELECT COUNT(*)
					FROM animals 
					WHERE mother_id = a.id 
						AND user_id = $1 
						AND deleted_at IS NULL
				) AS children_number
			FROM animals a
			LEFT JOIN animals f ON f.id = a.father_id
			LEFT JOIN animals m ON m.id = a.mother_id
			LEFT JOIN pasture_cte p ON p.animal_id = a.id
			WHERE a.user_id = $1 AND a.deleted_at IS NULL
		)
		SELECT 
			id,
			tag,
			name,
			sex,
			animal_type,
			birth_date,
			death_date,
			weaning_date,

			father_id,
			father_tag,
			father_name,

			mother_id,
			mother_tag,
			mother_name,

			pasture_id,
			pasture_name,
			farm_id,
			farm_name
		FROM cte
	`

	filterExpression, _, err := util.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	whereExp := util.GetWhereExpression(filterExpression)
	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	query += whereExp + " ORDER BY " + sortExpression
	return util.GetList[AnimalDB](r.DB, query, args...)
}

func (r *AnimalRepository) Delete(id string, userId string) *log.APIError {

	tx, err := r.DB.Beginx()
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	query := `
		UPDATE animals
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, query, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	lacQuery := `
		UPDATE lactations
		SET deleted_at = NOW()
		WHERE animal_id = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, lacQuery, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	milkQuery := `
		UPDATE milk_entries
		SET deleted_at = NOW()
		WHERE animal_id = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, milkQuery, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	embryoQuery := `
		UPDATE embryo_transfer
		SET deleted_at = NOW()
		WHERE receiver_id = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, embryoQuery, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	inseminationQuery := `
		UPDATE insemination_entries
		SET deleted_at = NOW()
		WHERE animal_id = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, inseminationQuery, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	breedingQuery := `
		UPDATE breeding_entries
		SET deleted_at = NOW()
		WHERE animal_id = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, breedingQuery, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	pastureQuery := `
		UPDATE pasture_entries
		SET deleted_at = NOW()
		WHERE animal_id = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, pastureQuery, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	pregnancyQuery := `
		UPDATE pregnancy_tests
		SET deleted_at = NOW()
		WHERE animal_id = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, pregnancyQuery, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	slaughterQuery := `
		UPDATE slaughter_entries
		SET deleted_at = NOW()
		WHERE animal_id = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, slaughterQuery, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	weightQuery := `
		UPDATE weight_entries
		SET deleted_at = NOW()
		WHERE animal_id = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, weightQuery, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	err = tx.Commit()
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *AnimalRepository) Update(newEntry *AnimalSave) (*AnimalDB, error) {

	updateQuery := `
		UPDATE animals
		SET tag = :tag,
			name = :name,
			death_date = :death_date,
			birth_date = :birth_date,
			sex = :sex,
			father_id = :father_id,
			mother_id = :mother_id,
			animal_type = :animal_type,
			weaning_date = :weaning_date,
			is_outside_animal = :is_outside_animal,
			is_insemination_bull = :is_insemination_bull,
			is_transfer_bull = :is_transfer_bull,
			is_breeding_bull = :is_breeding_bull,
			is_embryo_donor = :is_embryo_donor
		WHERE id = :id AND user_id = :user_id
	`

	err := util.NamedExec(r.DB, updateQuery, newEntry)
	if err != nil {
		return nil, err
	}

	selectQuery := `
		WITH pature_cte AS (
			SELECT
				pe.animal_id,
				pe.pasture_id,
				p.name AS pasture_name,
				p.farm_id,
				f.name AS farm_name
			FROM pasture_entries pe
			WHERE pe.animal_id = :id
				AND pe.exit_date IS NULL
				AND pe.user_id = :user_id
				AND pe.deleted_at IS NULL
		)	

        SELECT 
			a.id,
			a.tag,
			a.name,
			a.sex,
			a.animal_type,
			a.birth_date,
			a.death_date,
			a.weaning_date,
			a.observation,

			a.father_id,
			f.tag AS father_tag,
			f.name AS father_name,

			a.mother_id,
			m.tag AS mother_tag,
			m.name AS mother_name,
				
			p.pasture_id,
			p.pasture_name,
			p.farm_id,
			p.farm_name
        FROM animals a
		CROSS JOIN pasture_cte p
		LEFT JOIN animals f ON f.id = a.father_id
		LEFT JOIN animals m ON m.id = a.mother_id
		WHERE a.id = :id AND a.user_id = :user_id
    `

	response, err := util.NamedGet(r.DB, selectQuery, AnimalDB{}, newEntry)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *AnimalRepository) Add(entry *AnimalSave) error {

	query := `
		INSERT INTO animals (
			tag, 
			name, 
			sex, 
			father_id, 
			mother_id, 
			animal_type, 
			birth_date, 
			death_date, 
			weaning_date, 
			is_insemination_bull,
			is_transfer_bull,
			is_breeding_bull,
			is_embryo_donor,
			is_outside_animal,
			observation,
			user_id
		)
		VALUES (
			:tag, 
			:name, 
			:sex, 
			:father_id, 
			:mother_id, 
			:animal_type, 
			:birth_date, 
			:death_date, 
			:weaning_date, 
			:is_insemination_bull,
			:is_transfer_bull,
			:is_breeding_bull,
			:is_embryo_donor,
			:is_outside_animal,
			:observation,
			:user_id
		)
	`

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return err
	}

	return nil
}
