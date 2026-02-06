package birth

import (
	"database/sql"
	"fmt"

	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type BirthRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *BirthRepository {
	return &BirthRepository{db}
}

type BirthValidation struct {
	BirthExists          bool `db:"birth_exists"`
	RingExists           bool `db:"ring_exists"`
	InvalidPreviousBirth bool `db:"invalid_previous_birth"`
	InvalidNextBirth     bool `db:"invalid_next_birth"`
}

func (r *BirthRepository) CheckBirthConflicts(entry *BirthEntrySave) (*BirthValidation, error) {

	query := fmt.Sprintf(`
		SELECT 
			EXISTS (
				SELECT 1
				FROM animals
				WHERE mother_id = :mother_id
					AND birth_date = :birth_date
					AND id IS DISTINCT FROM :id
					AND user_id = :user_id
					AND deleted_at IS NULL
			) AS birth_exists,
			EXISTS (
				SELECT 1
				FROM animals
				WHERE tag = :tag 
					AND id IS DISTINCT FROM :id
					AND user_id = :user_id
					AND deleted_at IS NULL
			) AS ring_exists,
			EXISTS (
				SELECT 1
				FROM animals
				WHERE mother_id = :mother_id
					AND birth_date < :birth_date
					AND age(birth_date, :birth_date) <= 'interval %[1]d days'
					AND user_id = :user_id
					AND id IS DISTINCT FROM :id
					AND deleted_at IS NULL
				ORDER BY birth_date DESC
				LIMIT 1
			) AS invalid_previous_birth,
			EXISTS (
				SELECT 1
				FROM animals
				WHERE mother_id = :mother_id
					AND birth_date > :birth_date
					AND age(:birth_date, birth_date) <= 'interval %[1]d days'
					AND user_id = :user_id
					AND id IS DISTINCT FROM :id
					AND deleted_at IS NULL
				ORDER BY birth_date DESC
				LIMIT 1
			) AS invalid_next_birth
	`, util.MinGestantionDays)

	return util.NamedGet(r.DB, query, BirthValidation{}, entry)
}

func (r *BirthRepository) GetPotentialFather(entry *BirthEntrySave) (string, error) {
	query := fmt.Sprintf(`
		WITH insemination_cte AS (
			SELECT bull_id
			FROM insemination_entries i
			WHERE i.deleted_at IS NULL
				AND i.user_id = :user_id
				AND i.animal_id = :mother_id
				AND i.insemination_date < :birth_date
				AND age(:birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND '%[2]d days'
				AND NOT EXISTS (
					SELECT 1 FROM pregnancy_tests t
					WHERE t.deleted_at IS NULL
						AND t.pregnancy_status = 'FAILED'
						AND t.test_date BETWEEN i.insemination_date AND :birth_date
						AND t.animal_id = :mother_id
				)
			ORDER BY i.insemination_date DESC
			LIMIT 1
		),

		transfer_cte AS (
			SELECT bull_id
			FROM embryo_transfer t
			WHERE t.deleted_at IS NULL
				AND t.user_id = :user_id
				AND t.receiver_id = :mother_id
				AND t.transfer_date < :birth_date
				AND age(:birth_date, t.transfer_date) BETWEEN INTERVAL '%[1]d days' AND '%[2]d days'
				AND NOT EXISTS (
					SELECT 1 FROM pregnancy_tests pt
					WHERE pt.deleted_at IS NULL
						AND pt.pregnancy_status = 'FAILED'
						AND pt.test_date BETWEEN t.transfer_date AND :birth_date
						AND pt.animal_id = :mother_id
				)
			ORDER BY t.transfer_date DESC
			LIMIT 1
		),

		breeding_cte AS (
			SELECT bull_id
			FROM breeding_entries b
			WHERE b.deleted_at IS NULL
				AND b.user_id = :user_id
				AND b.animal_id = :mother_id
				AND b.breeding_date < :birth_date
				AND age(:birth_date, b.breeding_date) BETWEEN INTERVAL '%[1]d days' AND '%[2]d days'
				AND NOT EXISTS (
					SELECT 1 FROM pregnancy_tests pt
					WHERE pt.deleted_at IS NULL
						AND pt.pregnancy_status = 'FAILED'
						AND pt.test_date BETWEEN b.breeding_date AND :birth_date
						AND pt.animal_id = :mother_id
				)
			ORDER BY b.breeding_date DESC
			LIMIT 1
		),

		mother_entries_cte AS (
			SELECT pasture_id
			FROM pasture_entries
			WHERE deleted_at IS NULL
				AND user_id = :user_id
				AND animal_id = :mother_id
				AND (CAST(:birth_date AS timestamptz) - INTERVAL '%[3]d days') BETWEEN entry_date AND COALESCE(exit_date, NOW())
		),
		
		pasture_cte AS (
			SELECT pe.animal_id
			FROM pasture_entries pe
			CROSS JOIN mother_entries_cte me
			JOIN animals bull ON bull.id = pe.animal_id
				AND bull.animal_type = 'REPRODUCTION_ANIMAL'
				AND bull.sex = 'M'
			WHERE pe.deleted_at IS NULL
				AND pe.user_id = :user_id
				AND (CAST(:birth_date AS timestamptz) - INTERVAL '%[3]d days') BETWEEN pe.entry_date AND COALESCE(pe.exit_date, NOW())
				AND pe.pasture_id = me.pasture_id
			ORDER BY (COALESCE(pe.exit_date, NOW()) - pe.entry_date) DESC
			LIMIT 1
		)

		SELECT COALESCE(
			(SELECT bull_id FROM insemination_cte),
			(SELECT bull_id FROM transfer_cte),
			(SELECT bull_id FROM breeding_cte),
			(SELECT animal_id FROM pasture_cte)
		)
	`, util.MinGestantionDays, util.MaxGestationDays, util.AverageGestationDays)

	var fatherId sql.NullString
	err := util.NamedPrimitive(r.DB, query, &fatherId, entry)
	if err != nil {
		return "", err
	}

	if fatherId.Valid {
		return fatherId.String, nil
	}

	return "", nil
}

func (r *BirthRepository) GetBestIntervals(userId string) (*[]IntervalAnimal, error) {
	query := `
        WITH interval_list AS (
			SELECT 
				a.mother_id,
				EXTRACT(days FROM a.birth_date - LAG(a.birth_date) OVER win) AS birth_interval
			FROM animals a
				JOIN animals m ON m.id = a.mother_id AND m.is_outside_animal = FALSE
			WHERE 
				a.user_id = $1 
				AND a.deleted_at IS NULL
				AND a.birth_date IS NOT NULL
				AND a.is_outside_animal = FALSE
			WINDOW win AS (PARTITION BY a.mother_id ORDER BY a.birth_date)
		),
		average_list AS (
            SELECT
                CONCAT_WS(' - ', a.tag, a.name) animal_name,
                COUNT(l.*) birth_numbers,
                AVG(l.birth_interval) interval_average
            FROM interval_list l
                JOIN animals a ON 
					a.id = l.mother_id
					AND a.death_date IS NULL
            GROUP BY animal_name
        ),
        birth_stats AS (
            SELECT 
				AVG(interval_average) gn_interval_average,
				STDDEV(interval_average) dev_interval
            FROM average_list
			WHERE birth_numbers >= 3
        ),
        scores AS (
            SELECT 
                a.animal_name,
                a.interval_average,
                a.birth_numbers,
                (a.interval_average - b.gn_interval_average) / NULLIF(b.dev_interval, 0) AS reproductive_score,
                ((a.interval_average / NULLIF(b.gn_interval_average, 0)) - 1) * 100 average_rate
            FROM average_list a, birth_stats b
			WHERE a.birth_numbers >= 3
        )
        SELECT *
        FROM scores
		WHERE -reproductive_score > 0
        ORDER BY -reproductive_score DESC
		LIMIT 10
    `
	return util.GetList[IntervalAnimal](r.DB, query, userId)
}

func (r *BirthRepository) GetWorstIntervals(userId string) (*[]IntervalAnimal, error) {
	query := `
        WITH interval_list AS (
			SELECT 
				a.mother_id,
				EXTRACT(days FROM a.birth_date - LAG(a.birth_date) OVER win) AS birth_interval
			FROM animals a
				JOIN animals m ON m.id = a.mother_id AND m.is_outside_animal = FALSE
			WHERE 
				a.user_id = $1 
				AND a.deleted_at IS NULL
				AND a.birth_date IS NOT NULL
				AND a.is_outside_animal = FALSE
			WINDOW win AS (PARTITION BY a.mother_id ORDER BY a.birth_date)
		),
		average_list AS (
            SELECT
                CONCAT_WS(' - ', a.tag, a.name) animal_name,
                COUNT(l.*) birth_numbers,
                AVG(l.birth_interval) interval_average
            FROM interval_list l
                JOIN animals a ON 
					a.id = l.mother_id
					AND a.death_date IS NULL
            GROUP BY animal_name
        ),
        birth_stats AS (
            SELECT 
				AVG(interval_average) gn_interval_average,
				STDDEV(interval_average) dev_interval
            FROM average_list
			WHERE birth_numbers >= 3
        ),
        scores AS (
            SELECT 
                a.animal_name,
                a.interval_average,
                a.birth_numbers,
                (a.interval_average - b.gn_interval_average) / NULLIF(b.dev_interval, 0) AS reproductive_score,
                ((a.interval_average / NULLIF(b.gn_interval_average, 0)) - 1) * 100 average_rate
            FROM average_list a, birth_stats b
			WHERE a.birth_numbers >= 3
        )
        SELECT *
        FROM scores
		WHERE reproductive_score > 0
        ORDER BY reproductive_score DESC
		LIMIT 10
    `
	return util.GetList[IntervalAnimal](r.DB, query, userId)
}

func (r *BirthRepository) GetBirthIntervalHistory(userId string) (*[]util.GraphData, error) {
	query := `
		WITH year_birth_interval AS (
			SELECT 
				DATE_TRUNC('year', a.birth_date) birth_date,
				EXTRACT(days FROM a.birth_date - LAG(a.birth_date) OVER win) birth_interval
			FROM animals a
				JOIN animals m ON m.id = a.mother_id AND m.is_outside_animal = FALSE
			WHERE 
				a.user_id = $1 
				AND a.deleted_at IS NULL
				AND a.birth_date IS NOT NULL
				AND a.is_outside_animal = FALSE
			WINDOW win AS (PARTITION BY a.mother_id ORDER BY a.birth_date)
		),
		cte AS (
			SELECT 
				birth_date AS date,
				AVG(birth_interval) AS value
			FROM year_birth_interval
			WHERE birth_interval <> 0
			GROUP BY 1
			ORDER BY birth_date DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY birth_date
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *BirthRepository) GetLastBirthsNumber(userId string) (*[]util.GraphData, error) {
	query := `
		WITH cte AS (
			SELECT
				DATE_TRUNC('month', a.birth_date) as date,
				COUNT(a.*) value
			FROM animals a
				JOIN animals m ON m.id = a.mother_id AND m.is_outside_animal = FALSE
			WHERE 
				a.user_id = $1 
				AND a.deleted_at IS NULL
				AND a.birth_date IS NOT NULL
				AND a.is_outside_animal = FALSE
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY 1
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *BirthRepository) GetYearBirthsNumber(userId string) (*[]util.GraphData, error) {
	query := `
		WITH cte AS (
			SELECT
				DATE_TRUNC('year', a.birth_date) date,
				COUNT(a.*) value
			FROM animals a
				JOIN animals m ON m.id = a.mother_id 
					AND m.is_outside_animal = FALSE
			WHERE 
				a.user_id = $1 
				AND a.deleted_at IS NULL
				AND a.birth_date IS NOT NULL
				AND a.is_outside_animal = FALSE
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 20
		)
		SELECT * FROM cte ORDER BY entry_date
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *BirthRepository) GetYearDeathsNumber(userId string) (*[]util.GraphData, error) {
	query := `
		WITH cte AS (
			SELECT
				DATE_TRUNC('year', a.death_date) entry_date,
				COUNT(a.*) deaths_total
			FROM animals a
				JOIN animals m ON m.id = a.mother_id 
					AND m.is_outside_animal = FALSE
			WHERE 
				a.user_id = $1 
				AND a.deleted_at IS NULL
                AND age(a.death_date, a.birth_date) < INTERVAL '1 year'
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 20
		)
		SELECT * FROM cte ORDER BY entry_date
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *BirthRepository) GetDeathIndex(userId string) (*[]util.GraphData, error) {
	query := `
        WITH death_tbl AS (
            SELECT 
                DATE_TRUNC('year', a.death_date) date,
                COUNT(a.*) deaths
            FROM animals a
				JOIN animals m ON m.id = a.mother_id AND m.is_outside_animal = FALSE
            WHERE
                a.user_id = $1
                AND a.death_date IS NOT NULL
                AND age(a.death_date, a.birth_date) < INTERVAL '1 year'
                AND a.deleted_at IS NULL
            GROUP BY 1
        ),
        birth_tbl AS (
            SELECT
                DATE_TRUNC('year', a.birth_date) date,
                COUNT(a.*) births
            FROM animals a
            WHERE 
				a.user_id = $1 
				AND a.deleted_at IS NULL
				AND a.birth_date IS NOT NULL
				AND a.is_outside_animal = FALSE
            GROUP BY 1
        ),
        cte AS (
			SELECT
				date,
				COALESCE((deaths::float / NULLIF(births, 0)::float)*100, 0) death_index
			FROM birth_tbl FULL JOIN death_tbl USING(date)
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY date
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *BirthRepository) GetBirthHistory(userId string) (*[]BirthsByDate, error) {
	query := `
        WITH death_data AS ( 
            SELECT 
                DATE_TRUNC('month', death_date) date,
                COUNT(*) death_total 
			FROM animals  
            WHERE 
                user_id = $1
                AND deleted_at IS NULL 
                AND death_date IS NOT NULL
                AND age(death_date, birth_date) < INTERVAL '1 year'
            GROUP BY 1
        ), 
        birth_data AS (
            SELECT 
                DATE_TRUNC('month', a.birth_date) date,
                COUNT(a.*) AS birth_total
            FROM animals a
				JOIN animals m ON m.id = a.mother_id AND m.is_outside_animal = FALSE
            WHERE 
                a.user_id = $1
                AND a.deleted_at IS NULL 
				AND a.is_outside_animal = FALSE
				AND a.birth_date IS NOT NULL
            GROUP BY 1
        ),
        cte AS (
			SELECT
				date,
				COALESCE(birth_data.birth_total,0) birth_total,
				COALESCE(death_data.death_total, 0) death_total
			FROM birth_data FULL JOIN death_data USING(date)
			ORDER BY date DESC
			LIMIT 60
		) 
		SELECT * FROM cte ORDER BY date
    `
	return util.GetList[BirthsByDate](r.DB, query, userId)
}

func (r *BirthRepository) TotalBySex(userId string) (*[]TotalBirthsBySex, error) {
	query := `
        WITH cte AS (
			SELECT 
				DATE_TRUNC('month', a.birth_date) birth_month,
				COUNT(a.*) FILTER (WHERE a.sex = 'M') males,
				COUNT(a.*) FILTER (WHERE a.sex = 'F') females
			FROM animals a
				JOIN animals m ON m.id = a.mother_id 
					AND m.is_outside_animal = FALSE
			WHERE 
				a.user_id = $1 
				AND a.deleted_at IS NULL
				AND a.is_outside_animal = FALSE
				AND a.birth_date IS NOT NULL
			GROUP BY birth_month
			ORDER BY birth_month DESC
			LIMIT 24
		)
		SELECT * FROM cte ORDER BY birth_month
    `
	return util.GetList[TotalBirthsBySex](r.DB, query, userId)
}

func (r *BirthRepository) GetYearBySex(userId string) (*[]TotalBirthsBySex, error) {
	query := `
		SELECT 
			DATE_TRUNC('year', a.birth_date) birth_month,
			COUNT(a.*) FILTER (WHERE a.sex = 'M') males,
			COUNT(a.*) FILTER (WHERE a.sex = 'F') females
		FROM animals a
			JOIN animals m ON m.id = a.mother_id 
				AND m.is_outside_animal = FALSE
		WHERE 
			a.user_id = $1 
			AND a.deleted_at IS NULL
			AND a.birth_date IS NOT NULL
			AND a.is_outside_animal = FALSE
		GROUP BY birth_month
		ORDER BY birth_month DESC
		LIMIT 10
    `
	return util.GetList[TotalBirthsBySex](r.DB, query, userId)
}

func (r *BirthRepository) GetLastBirths(userId string) (*[]BirthDB, error) {
	query := `
        SELECT 
            CONCAT_WS(' - ', m.tag, m.name) mother_info,
            a.birth_date AS calf_birth_date,
            a.sex AS calf_sex,
            CONCAT_WS(' - ', f.tag, f.name) AS calf_father,
			EXTRACT(days FROM a.birth_date - LAG(a.birth_date) OVER win) AS birth_interval
        FROM animals a
            JOIN animals m ON m.id = a.mother_id AND m.is_outside_animal = FALSE
            LEFT JOIN animals f ON f.id = a.father_id
		WHERE 
			a.user_id = $1 
			AND a.deleted_at IS NULL
			AND a.birth_date IS NOT NULL
			AND a.is_outside_animal = FALSE
		WINDOW win AS (PARTITION BY a.mother_id ORDER BY a.birth_date)
		ORDER BY a.birth_date DESC, COALESCE(REGEXP_REPLACE(m.tag, '[^0-9]', '', 'g')::int, 0)
		LIMIT 15
    `
	return util.GetList[BirthDB](r.DB, query, userId)
}

func (r *BirthRepository) FindPage(
	userId string,
	sort string,
	order string,
	filter *BirthEntryFilter,
	cursor string,
	limit int,
) (*[]BirthDB, error) {
	sortMap := map[string]util.SortField{
		"calf_birth_date": {Field: "cte.calf_birth_date", Order: "asc"},
		"mother_order":    {Field: "cte.mother_order", Order: "asc"},
		"mother_name":     {Field: "cte.mother_name", Order: "asc"},
		"birth_interval":  {Field: "coalesce(cte.birth_interval, 0)", Order: "asc"},
		"id":              {Field: "cte.calf_id", Order: "asc"},
	}

	query := `
		WITH cte AS (
			SELECT 
				a.id AS calf_id,
				a.sex AS calf_sex,
				a.name AS calf_name,
				a.tag AS calf_tag,
				a.birth_date AS calf_birth_date,
				a.observation AS calf_observation,

				a.mother_id,
				m.name AS mother_name,
				m.tag AS mother_tag,
				COALESCE(REGEXP_REPLACE(m.tag, '[^0-9]', '', 'g')::int, 0) AS mother_order,

				a.father_id AS father_id,
				f.name AS father_name,
				f.tag AS father_tag,

				EXTRACT(days FROM a.birth_date - LAG(a.birth_date) OVER win) AS birth_interval
			FROM animals a
				JOIN animals m ON m.id = a.mother_id AND m.is_outside_animal = FALSE
				LEFT JOIN animals f ON f.id = a.father_id
			WHERE 
				a.user_id = $1 
				AND a.deleted_at IS NULL 
				AND m.is_outside_animal = FALSE
				AND a.birth_date IS NOT NULL
				AND a.mother_id IS NOT NULL
			WINDOW win AS (PARTITION BY a.mother_id ORDER BY a.birth_date)
		)
		SELECT * FROM cte
    `
	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "cte", 2)
	cursorExpression, nextParam, err := util.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExpression, cursorExpression)
	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	orderExpression := " ORDER BY " + sortExpression + fmt.Sprintf(" LIMIT %d", limit)
	query += whereExpression + orderExpression
	return util.GetList[BirthDB](r.DB, query, args...)
}

func (r *BirthRepository) GetPageFoot(userId string, filter *BirthEntryFilter) (*BirthFooter, error) {

	animalQuery := `
		SELECT
			a.id,
			EXTRACT(days FROM a.birth_date - LAG(a.birth_date) OVER win) AS birth_interval
		FROM animals a JOIN animals m ON m.id = a.mother_id 
    `
	whereExpression := `
		a.user_id = $1 
		AND a.deleted_at IS NULL
		AND a.birth_date IS NOT NULL
		AND m.is_outside_animal = FALSE
	`
	windowExp := "WINDOW win AS (PARTITION BY a.mother_id ORDER BY a.birth_date)"

	filterExpression, _, err := util.GetFilterExpressions(filter, "a", 2)
	if err != nil {
		return nil, err
	}

	whereExpression = util.GetWhereExpression(whereExpression, filterExpression)

	animalQuery += whereExpression + " " + windowExp
	query := fmt.Sprintf(` 
		WITH animal_cte AS (%s)
		SELECT
			COUNT(a.*) AS total,
			AVG(a.birth_interval) AS interval_average
		FROM animal_cte a
	`, animalQuery)

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return util.GetOne[BirthFooter](r.DB, query, args...)
}

func (r *BirthRepository) UpdateBirth(entry *BirthEntrySave) (*BirthDB, error) {

	tx, err := r.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `
		UPDATE animals
		SET birth_date = :birth_date,
			sex = :sex,
			father_id = :father_id,
			observation = :observation
			WHERE id = :id AND user_id = :user_id
	`
	err = util.NamedExecTx(tx, query, entry)
	if err != nil {
		return nil, err
	}

	selectQuery := `
		SELECT 
			a.id,
			a.mother_id,
			m.name AS mother_name,
			CONCAT_WS(' - ', m.tag, m.name) AS mother_info,
			a.birth_date AS calf_birth_date,
			a.sex AS calf_sex,
			CASE 
				WHEN a.name IS NULL THEN ''
				ELSE CONCAT_WS(' - ', a.tag, a.name)
			END AS calf_name,
			a.father_id AS calf_father_id,
			CONCAT_WS(' - ', f.tag, f.name) calf_father,
			bi.birth_interval
		FROM animals a
			JOIN animals m ON m.id = a.mother_id 
			LEFT JOIN animals f ON f.id = a.father_id
			JOIN (
				SELECT
					id,
					EXTRACT(days FROM bi.birth_date - LAG(bi.birth_date) OVER win) AS birth_interval
				FROM animals bi
				WHERE bi.mother_id = :mother_id
					AND bi.deleted_at IS NULL
					AND bi.user_id = :user_id
				WINDOW win AS (PARTITION BY bi.mother_id ORDER BY bi.birth_date)
			) bi ON bi.id = a.id
			WHERE a.id = :id 
				AND a.user_id = :user_id
				AND m.is_outside_animal = FALSE
	`

	result, err := util.NamedGet(r.DB, selectQuery, BirthDB{}, entry)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *BirthRepository) AddBirth(entry *BirthEntrySave) error {

	if entry.Overwrite {
		query := `
			UPDATE animals
			SET tag = :tag,
				sex = :sex,
				father_id = :father_id,
				observation = :observation
			WHERE user_id = :user_id
				AND deleted_at IS NULL
				AND birth_date = :birth_date
				AND mother_id = :mother_id
		`
		err := util.NamedExec(r.DB, query, entry)
		if err != nil {
			return err
		}

		return nil
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	birthQuery := `
		INSERT INTO animals (tag, sex, birth_date, father_id, mother_id, animal_type, observation, user_id)
		VALUES (:tag, :sex, :birth_date, :father_id, :mother_id, 'OFFSPRING', :observation, :user_id)
	`

	newId, err := util.NamedExecReturningIdTx(tx, birthQuery, entry)
	if err != nil {
		return err
	}

	if entry.PastureId != nil {
		entry.Id = newId

		pastureEntryQuery := `
			INSERT INTO pasture_entries (animal_id, pasture_id, entry_date, user_id)
			VALUES (:id, :pasture_id, :birth_date, :user_id)
		`
		err = util.NamedExecTx(tx, pastureEntryQuery, entry)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
