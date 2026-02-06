package breeding

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/domains/animals"
	"github.com/felipeErnica/rebanho-backend/internal/entity"
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type BreedingRepository struct {
	DB *sqlx.DB
}

type SaveValidation struct {
	Repeated     bool `db:"repeated"`
	HasPregnancy bool `db:"has_pregnancy"`
	HasChildren  bool `db:"has_children"`
}

type DeleteValidation struct {
	HasPregnancy bool `db:"has_pregnancy"`
	HasChildren  bool `db:"has_children"`
}

func NewRepository(db *sqlx.DB) *BreedingRepository {
	return &BreedingRepository{db}
}

func (r *BreedingRepository) CheckEntriesOnDate(date time.Time, userId string) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM breeding_entries
			WHERE breeding_date = $1
				AND user_id = $2
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.GetPrimitive(r.DB, query, &exists, date, userId)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *BreedingRepository) CheckBreedingDelete(id string, userId string) (*DeleteValidation, error) {
	query := fmt.Sprintf(`
		WITH old_entry AS (
			SELECT 
				breeding_date,
				bull_id,
				animal_id,
				user_id
			FROM breeding_entries
			WHERE id = $1
				AND user_id = $2
				AND deleted_at IS NULL
		)

		SELECT EXISTS (
			SELECT 1
			FROM pregnancy_tests t
				CROSS JOIN old_entry o 
			WHERE t.pregnancy_status = 'SUCCESS'
				AND t.animal_id = o.animal_id
				AND t.test_date > o.breeding_date
				AND age(t.test_date, o.breeding_date) <= INTERVAL '%[2]d days'
				AND t.user_id = o.user_id
				AND t.deleted_at IS NULL
		) AND NOT EXISTS (
			SELECT 1
			FROM pregnancy_tests t
				CROSS JOIN old_entry o 
			WHERE t.pregnancy_status = 'FAILED'
				AND t.animal_id = o.animal_id
				AND t.test_date > o.breeding_date
				AND age(t.test_date, o.breeding_date) <= INTERVAL '%[2]d days'
				AND t.user_id = o.user_id
				AND t.deleted_at IS NULL
		) AS has_pregnancy,
		EXISTS (
			SELECT 1
			FROM animals a, old_entry o
			WHERE a.mother_id = o.animal_id
				AND a.birth_date > o.breeding_date
				AND age(a.birth_date, o.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				AND a.father_id =  o.bull_id
				AND a.deleted_at IS NULL 
				AND a.user_id = o.user_id
		) AS has_children
	`, util.MinGestantionDays, util.MaxGestationDays)

	validate, err := util.GetOne[DeleteValidation](r.DB, query, id, userId)
	if err != nil {
		return nil, err
	}

	return validate, nil
}

func (r *BreedingRepository) CheckBreedingAdd(entry *BreedingEntrySave) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM breeding_entries
			WHERE animal_id = :animal_id
			AND breeding_date = :breeding_date
				AND id IS DISTINCT FROM :id
				AND user_id = :user_id
				AND deleted_at IS NULL
		)
	`
	var exists bool
	err := util.NamedPrimitive(r.DB, query, exists, entry)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *BreedingRepository) CheckBreedingSave(entry *BreedingEntrySave) (*SaveValidation, error) {
	query := fmt.Sprintf(`
		WITH old_entry AS (
			SELECT 
				breeding_date,
				bull_id
			FROM breeding_entries
			WHERE id = :id
				AND user_id = :user_id
				AND deleted_at IS NULL
		)

		SELECT EXISTS (
			SELECT 1
			FROM breeding_entries
			WHERE animal_id = :animal_id
			AND breeding_date = :breeding_date
				AND id IS DISTINCT FROM :id
				AND user_id = :user_id
				AND deleted_at IS NULL
		) AS repeated,
		EXISTS (
			SELECT 1
			FROM pregnancy_tests t
				CROSS JOIN old_entry o 
			WHERE t.pregnancy_status = 'SUCCESS'
				AND t.animal_id = :animal_id
				AND t.test_date > o.breeding_date
				AND age(t.test_date, o.breeding_date) <= INTERVAL '%[2]d days'
				AND t.user_id = :user_id
				AND t.deleted_at IS NULL
		) AND NOT EXISTS (
			SELECT 1
			FROM pregnancy_tests t
				CROSS JOIN old_entry o 
			WHERE t.pregnancy_status = 'FAILED'
				AND t.animal_id = :animal_id
				AND t.test_date > o.breeding_date
				AND age(t.test_date, o.breeding_date) <= INTERVAL '%[2]d days'
				AND t.user_id = :user_id
				AND t.deleted_at IS NULL
		) AS has_pregnancy,
		EXISTS (
			SELECT 1
			FROM animals a
			CROSS JOIN old_entry o
			WHERE a.mother_id = :animal_id
				AND a.birth_date > o.breeding_date
				AND age(a.birth_date, o.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				AND a.father_id =  o.bull_id
				AND a.deleted_at IS NULL 
				AND a.user_id = :user_id
		) AS has_children
	`, util.MinGestantionDays, util.MaxGestationDays)

	validate, err := util.NamedGet(r.DB, query, SaveValidation{}, entry)
	if err != nil {
		return nil, err
	}

	return validate, nil
}

func (r *BreedingRepository) GetBirthRateStats(userId string) (*[]util.GraphData, error) {
	query := fmt.Sprintf(`
        WITH totals AS (
            SELECT 
                breeding_date,
                COUNT(*) AS total,
                COUNT(*) FILTER (WHERE EXISTS (
					SELECT 1 
					FROM animals a
					WHERE a.mother_id = i.animal_id
						AND a.birth_date > i.breeding_date
						AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
						AND NOT EXISTS (
							SELECT 1
							FROM pregnancy_tests t
							WHERE t.animal_id = a.mother_id
								AND t.test_date BETWEEN i.breeding_date AND a.birth_date
								AND t.pregnancy_status = 'FAILED'
						)
				)) birth_success
            FROM breeding_entries i
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
			GROUP BY 1
            ORDER BY 1 DESC
            LIMIT 10
        )
        SELECT 
            breeding_date,
			(birth_success::float / NULLIF(total, 0)::float) * 100 AS birth_rate
        FROM totals
		ORDER BY 1
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *BreedingRepository) GetPregnancyRateStats(userId string) (*[]util.GraphData, error) {
	query := fmt.Sprintf(`
		WITH insemination_status AS (
			SELECT
				i.breeding_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.breeding_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status
			FROM breeding_entries i
			WHERE i.user_id = $1 AND i.deleted_at IS NULL 
		),
		cte AS (
			SELECT
				t.breeding_date,
				COUNT(t.*) AS total,
				COUNT(t.*) FILTER (WHERE t.pregnancy_status = 'SUCCESS') AS pregnancy_success
			FROM insemination_status t
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT 
			breeding_date AS date,
			(pregnancy_success::float / NULLIF(total, 0)) * 100 AS value
		FROM cte
		ORDER BY 1
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *BreedingRepository) GetAnimalsNumber(userId string) (*[]util.GraphData, error) {
	query := `
		WITH cte AS (
			SELECT
				breeding_date AS date,
				COUNT(*) AS value
			FROM breeding_entries
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT *
		FROM cte
		ORDER BY 1
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *BreedingRepository) GetBreedingStats(userId string) (*[]BreedingHist, error) {
	query := fmt.Sprintf(`
        WITH cte AS (
			SELECT 
				i.breeding_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.breeding_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
								AND t.test_date BETWEEN i.breeding_date AND a.birth_date
								AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM breeding_entries i
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
		),
        totals AS (
            SELECT
                breeding_date,
                COUNT(*) animals_number,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') births_number,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancies_number
            FROM cte
            GROUP BY 1
            ORDER BY 1 DESC
            LIMIT 30
        )
        SELECT * FROM totals ORDER BY breeding_date
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[BreedingHist](r.DB, query, userId)
}

func (r *BreedingRepository) GetFutureBirths(userId string) (*[]util.GraphData, error) {
	query := fmt.Sprintf(`
		WITH upcoming_births AS (
			SELECT t.test_date + (%[3]d - t.pregnancy_time) AS birth_forecast
			FROM breeding_entries i
			JOIN pregnancy_tests t
				ON t.animal_id = i.animal_id
				AND t.test_date > i.breeding_date
				AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
				AND t.pregnancy_status = 'SUCCESS'
			WHERE i.user_id = $1
			  AND i.deleted_at IS NULL
			  AND NOT EXISTS (
				  SELECT 1
				  FROM animals a
				  WHERE a.mother_id = i.animal_id
					AND a.birth_date > i.breeding_date
					AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					AND NOT EXISTS (
						SELECT 1
						FROM pregnancy_tests t
						WHERE t.animal_id = a.mother_id
							AND t.test_date BETWEEN i.breeding_date AND a.birth_date
							AND t.pregnancy_status = 'FAILED'
					)
			  )
			  AND t.test_date::date + (%[3]d - t.pregnancy_time) >= NOW()  
		)
		SELECT
			DATE_TRUNC('month', birth_forecast) AS date,
			COUNT(*) AS value
		FROM upcoming_births
		GROUP BY 1
		ORDER BY 1;
	`, util.MinGestantionDays, util.MaxGestationDays, util.AverageGestationDays)
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *BreedingRepository) GetBestBull(userId string) (*[]BestBulls, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				CONCAT_WS(' - ', b.tag, b.name) bull_name,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.breeding_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
								AND t.test_date BETWEEN i.breeding_date AND a.birth_date
								AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM breeding_entries i
                LEFT JOIN animals b ON i.bull_id = b.id 
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
		),
        totals AS (
            SELECT
                s.bull_name,
                COUNT(s.*) total,
                COUNT(s.*) FILTER (WHERE s.birth_status = 'SUCCESS') birth_success,
                COUNT(s.*) FILTER (WHERE s.pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status s
            GROUP BY 1
        ),
        rates AS (
            SELECT 
                bull_name,
                total,
                (birth_success::float / NULLIF(total, 0)::float) * 100 birth_rate,
                (pregnancy_success::float / NULLIF(total, 0)::float) * 100 pregnancy_rate
            FROM totals
        )
        SELECT
			bull_name,
			total,
			birth_rate,
			pregnancy_rate,
			(birth_rate / NULLIF(AVG(birth_rate) OVER (), 0) - 1) * 100 AS birth_comparison_rate,
			(pregnancy_rate / NULLIF(AVG(pregnancy_rate) OVER (), 0) - 1) * 100 AS pregnancy_comparison_rate
		FROM rates
		ORDER BY birth_rate DESC;
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[BestBulls](r.DB, query, userId)
}

func (r *BreedingRepository) GetLastGroups(userId string) (*[]BreedingGroup, error) {
	query := fmt.Sprintf(`
		WITH insemination_data AS (
			SELECT
				i.breeding_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.breeding_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
								AND t.test_date BETWEEN i.breeding_date AND a.birth_date
								AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM breeding_entries i
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
		),
		daily_stats AS (
			SELECT
				breeding_date,
				COUNT(*) AS cow_number,
				COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') AS birth_success,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') AS pregnancy_success
			FROM insemination_data
			GROUP BY breeding_date
		),
		rates AS (
			SELECT
				breeding_date,
				cow_number,
				(birth_success::float * 100 / NULLIF(cow_number, 0)) AS birth_rate,
				(pregnancy_success::float * 100 / NULLIF(cow_number, 0)) AS pregnancy_rate
			FROM daily_stats
		)
		SELECT
			breeding_date,
			cow_number,
			birth_rate,
			pregnancy_rate,
			COALESCE(
				(birth_rate / NULLIF(LAG(birth_rate) OVER win, 0) - 1) * 100, 0
			) AS birth_comparison_rate,
			COALESCE(
				(pregnancy_rate / NULLIF(LAG(pregnancy_rate) OVER win, 0) - 1) * 100, 0
			) AS pregnancy_comparison_rate
		FROM rates
		WINDOW win AS (ORDER BY breeding_date)
		ORDER BY breeding_date DESC
		LIMIT 5;
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[BreedingGroup](r.DB, query, userId)
}

func (r *BreedingRepository) GetLastEntries(userId string) (*[]BreedingDB, error) {
	query := fmt.Sprintf(`
		WITH last_date_cte AS (
			SELECT max(breeding_date) AS last_date
			FROM breeding_entries
			WHERE user_id = $1 and deleted_at IS NULL
		),

		breeding_cte AS (
			SELECT
				id,
				breeding_date,
				animal_id,
				bull_id,
				observation
			FROM breeding_entries
			CROSS JOIN last_date_cte AS l
			WHERE user_id = $1 
				AND deleted_at IS NULL
				AND breeding_date = l.last_date
		),

		child_cte AS (
			SELECT 
				a.id AS child_id,
				a.mother_id,
				a.name AS child_name,
				a.tag AS child_tag,
				a.sex AS child_sex,
				a.birth_date AS child_birth_date
			FROM animals a
		)


		SELECT 
			i.id,
			i.breeding_date,

			i.animal_id,
			a.name AS animal_name,
			a.tag AS animal_tag,

			i.bull_id,
			b.name AS bull_name,
			b.tag AS bull_tag,

			c.child_id,
			c.child_tag,
			c.child_name,
			c.child_sex,
			c.child_birth_date,

			CASE
				WHEN c.child_id IS NOT NULL THEN 'SUCCESS'
				WHEN p.pregnancy_status IS NULL AND age(i.breeding_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE p.pregnancy_status
			END AS pregnancy_status,
			CASE
				WHEN c.child_id IS NOT NULL THEN 'SUCCESS'
				WHEN age(i.breeding_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status
		FROM breeding_entries i
		CROSS JOIN last_date_cte l
		LEFT JOIN animals a ON a.id = i.animal_id
		LEFT JOIN animals b ON b.id = i.bull_id
		LEFT JOIN pregnancy_tests p ON p.animal_id = i.animal_id
			AND p.test_date > i.breeding_date
			AND age(p.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
		LEFT JOIN child_cte c ON c.mother_id = i.animal_id
				AND c.child_birth_date > i.breeding_date
				AND age(c.child_birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				AND NOT EXISTS (
					SELECT 1
					FROM pregnancy_tests t
					WHERE t.animal_id = c.mother_id
						AND t.test_date BETWEEN i.breeding_date AND a.birth_date
						AND t.pregnancy_status = 'FAILED'
				)
		WHERE i.user_id = $1 
			AND i.breeding_date = l.last_date
			AND i.deleted_at IS NULL
		ORDER BY COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0);
    `, util.MinGestantionDays, util.MaxGestationDays)
	return  util.GetList[BreedingDB](r.DB, query, userId)
}

func (r *BreedingRepository) FindEntriesPage(
	userId string,
	filter *BreedingEntryFilter,
	sort string,
	order string,
	cursor string,
) (*entity.Page[BreedingDB], error) {

	sort = util.AddCommonFields(sort)
	sortMap := map[string]util.SortField{
		"animal_order":  {Field: "i.animal_order", Order: "asc"},
		"animal_name":   {Field: "i.animal_name", Order: "asc"},
		"breeding_date": {Field: "coalesce(i.breeding_date, '-infinity')", Order: "asc"},
		"id":            {Field: "i.id", Order: "asc"},
		"created_at":    {Field: "i.created_at", Order: "asc"},
	}

	query := fmt.Sprintf(`
        WITH cte AS (
			SELECT 
				i.id,
				COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0) animal_order,
				a.name AS animal_name,
				CONCAT_WS(' - ', a.tag, a.name) animal_info,
				i.breeding_date,
				i.bull_id,
				b.name bull_name,
				CASE
					WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					WHEN NOT EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
						  AND t.test_date > i.breeding_date
						  AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
					) AND age(i.breeding_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN age(i.breeding_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
					ELSE 'FAILED'
				END AS birth_status,
				CASE 
					WHEN c.child_name IS NULL THEN 'Sem Cria'
					ELSE c.child_name
				END AS child_information,
				i.observation,
				i.created_at
			FROM breeding_entries i
				LEFT JOIN animals a ON a.id = i.animal_id
				LEFT JOIN animals b ON b.id = i.bull_id
				LEFT JOIN LATERAL (
					SELECT
					CONCAT_WS(
						' - ', 
						a.tag, 
						COALESCE(a.name, a.sex), 
						TO_CHAR(a.birth_date, 'DD/MM/YYYY')
					) AS child_name
					FROM animals a
					WHERE a.mother_id = i.animal_id
						AND a.birth_date > i.breeding_date
						AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					ORDER BY a.birth_date
					LIMIT 1
				) c ON TRUE
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
		)
		SELECT * FROM cte i
	`, util.MinGestantionDays, util.MaxGestationDays)
	orderExpression := " ORDER BY "

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "i", 2)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := util.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExpression, cursorExpression)

	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	orderExpression += sortExpression
	query += whereExpression + orderExpression
	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return util.GetPage[BreedingDB](r.DB, query, sort, 100, args...)
}

func (r *BreedingRepository) GetEntriesFoot(
	userId string,
	filter *BreedingEntryFilter,
) (*BreedingFoot, error) {

	statusQuery := fmt.Sprintf(`
		WITH cte AS  (
			SELECT
				i.*,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.breeding_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
								AND t.test_date BETWEEN i.breeding_date AND a.birth_date
								AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM breeding_entries i
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
		)
		SELECT pregnancy_status, birth_status
		FROM cte i
	`, util.MinGestantionDays, util.MaxGestationDays)

	filterExpression, _, err := util.GetFilterExpressions(filter, "i", 2)
	if err != nil {
		return nil, err
	}

	whereExpression := ""
	if filterExpression != "" {
		whereExpression = " WHERE " + filterExpression
	}

	statusQuery += whereExpression

	query := fmt.Sprintf(`
		WITH status AS (%s),
		totals AS (
			SELECT 
				COUNT(*) totals,
				COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') birth_success,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success
			FROM status
		)
        SELECT 
            totals,
            COALESCE(birth_success::float / NULLIF(totals, 0), 0) * 100 average_birth_rate,
            COALESCE(pregnancy_success::float / NULLIF(totals, 0), 0) * 100 average_pregnancy_rate
		FROM totals
    `, statusQuery)

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return util.GetOne[BreedingFoot](r.DB, query, args...)
}

func (r *BreedingRepository) FindEntriesByGroup(userId string, date time.Time) (*[]BreedingDB, error) {

	query := fmt.Sprintf(`
        SELECT 
            i.id,
			CONCAT_WS(' - ', b.tag, b.name) AS bull_name,
            CONCAT_WS(' - ', a.tag, a.name) animal_info,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.breeding_date
						AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.breeding_date
						AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'SUCCESS'
				) THEN 'SUCCESS'
				WHEN NOT EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
					  AND t.test_date > i.breeding_date
					  AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
				) AND age(i.breeding_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS pregnancy_status,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.breeding_date
						AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN age(i.breeding_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,
			CASE
				WHEN c.child_name IS NULL THEN 'Sem Cria'
				ELSE child_name
			END AS child_information,
            i.observation
        FROM breeding_entries i
            LEFT JOIN animals a ON a.id = i.animal_id
            LEFT JOIN animals b ON b.id = i.bull_id
			LEFT JOIN LATERAL (
				SELECT
				CONCAT_WS(
					' - ', 
					a.tag, 
					COALESCE(a.name, a.sex), 
					TO_CHAR(a.birth_date, 'DD/MM/YYYY')
				) AS child_name
				FROM animals a
				WHERE a.mother_id = i.animal_id
					AND  a.birth_date > i.breeding_date
					AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				ORDER BY a.birth_date
				LIMIT 1
			) c ON TRUE
		WHERE i.user_id = $1 AND i.deleted_at IS NULL AND i.breeding_date = $2
        ORDER BY COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0)
	`, util.MinGestantionDays, util.MaxGestationDays)

	return util.GetList[BreedingDB](r.DB, query, userId, date)
}

func (r *BreedingRepository) GetEntriesByGroupFoot(userId string, date time.Time) (*BreedingFoot, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.breeding_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
								AND t.test_date BETWEEN i.breeding_date AND a.birth_date
								AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM breeding_entries i
			WHERE i.user_id = $1 
				AND i.breeding_date = $2
				AND i.deleted_at IS NULL
		),
        counting AS (
            SELECT 
                COUNT(*) totals,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') birth_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status
        )
        SELECT 
            totals,
            COALESCE(birth_success::float / NULLIF(totals, 0), 0) * 100 average_birth_rate,
            COALESCE(pregnancy_success::float / NULLIF(totals, 0), 0) * 100 average_pregnancy_rate
        FROM counting
    `, util.MinGestantionDays, util.MaxGestationDays)

	return util.GetOne[BreedingFoot](r.DB, query, userId, date)
}

func (r *BreedingRepository) FindGroups(userId string) (*[]BreedingGroup, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				i.breeding_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.breeding_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
								AND t.test_date BETWEEN i.breeding_date AND a.birth_date
								AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM breeding_entries i
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
		),
        totals AS (
            SELECT 
                breeding_date,
				COUNT(*) cow_number,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') birth_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status i
            GROUP BY breeding_date
        ),
        rates AS (
            SELECT
                breeding_date,
                cow_number,
                (birth_success::float / cow_number::float)*100 birth_rate,
                (pregnancy_success::float / cow_number::float)*100 pregnancy_rate
            FROM totals
        )
        SELECT 
            s.breeding_date,
            s.cow_number,
            s.birth_rate,
            s.pregnancy_rate,
            COALESCE(
				(s.birth_rate / NULLIF(LAG(s.birth_rate) OVER win, 0)) - 1, 0
			) * 100 AS birth_comparison_rate,
            COALESCE(
				(s.pregnancy_rate / NULLIF(LAG(s.pregnancy_rate) OVER win, 0)) - 1, 0
			) * 100 AS pregnancy_comparison_rate
        FROM rates s
		WINDOW win AS (ORDER BY s.breeding_date)
        ORDER BY s.breeding_date DESC
    `, util.MinGestantionDays, util.MaxGestationDays)

	return util.GetList[BreedingGroup](r.DB, query, userId)
}

func (r *BreedingRepository) AddBreedingBull(id string, userId string) *log.APIError {

	query := `
		UPDATE animals
		SET is_breeding_bull = TRUE
		WHERE user_id = $1 AND id = $2
    `

	err := util.Exec(r.DB, query, userId, id)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil

}

func (r *BreedingRepository) GetBreedingById(id string, userId string) (*BreedingEntrySave, error) {
	query := `
		SELECT id, animal_id, bull_id, breeding_date, observation, user_id
		FROM breeding_entries
		WHERE id = $1 AND user_id = $2
	`
	return util.GetOne[BreedingEntrySave](r.DB, query, id, userId)
}

func (r *BreedingRepository) Delete(id string, changeFather bool, userId string) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if !changeFather {

		fatherQuery := fmt.Sprintf(`
			WITH breeding_target AS (
				SELECT 
					animal_id, 
					breeding_date
				FROM breeding_entries
				WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
			),

			children AS (
				SELECT 
					a.id AS child_id,
					a.mother_id,
					(a.birth_date::timestamptz - INTERVAL '%[3]d days') AS conception_date
				FROM animals a
				JOIN breeding_target b ON a.mother_id = b.animal_id
				WHERE a.birth_date > b.breeding_date
				  AND age(a.birth_date, b.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				  AND a.user_id = $2 AND a.deleted_at IS NULL
			),

			potential_fathers AS (
				SELECT 
					c.child_id,
					bull.id AS father_id,
					(COALESCE(bull.exit_date, NOW()) - bull.entry_date) AS duration
				FROM children c
				JOIN pasture_entries mother_past 
					ON mother_past.animal_id = c.mother_id
					AND c.conception_date BETWEEN mother_past.entry_date AND COALESCE(mother_past.exit_date, NOW())
					AND mother_past.deleted_at IS NULL
					AND mother_past.user_id = $2
				JOIN pasture_entries bull_past 
					ON bull_past.pasture_id = mother_past.pasture_id
					AND c.conception_date BETWEEN bull_past.entry_date AND COALESCE(bull_past.exit_date, NOW())
					AND bull_past.deleted_at IS NULL
					AND bull_past.user_id = $2
				JOIN animals bull 
					ON bull.id = bull_past.animal_id 
					AND bull.sex = 'M' 
					AND bull.name IS NOT NULL
			),

			ranked_fathers AS (
				SELECT 
					child_id, 
					father_id,
					ROW_NUMBER() OVER (PARTITION BY child_id ORDER BY duration DESC) AS rn
				FROM potential_fathers
			)

			UPDATE animals a
			SET father_id = rf.father_id
			FROM ranked_fathers rf
			WHERE a.id = rf.child_id 
				AND rf.rn = 1 
				AND a.user_id = $2
		`, util.MinGestantionDays, util.MaxGestationDays, util.AverageGestationDays)

		err = util.ExecTx(tx, fatherQuery, id, userId)
		if err != nil {
			return err
		}

	}

	query := `
		UPDATE breeding_entries
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2
    `

	err = util.ExecTx(tx, query, id, userId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *BreedingRepository) AddBreeding(entry *BreedingEntrySave) *log.APIError {

	var query string
	if entry.Overwrite {
		query = `
			UPDATE breeding_entries
			SET bull_id = :bull_id,
				observation = :observation
			WHERE animal_id = :animal_id
				AND breeding_date = :breeding_date
				AND user_id = :user_id
		`
	} else {
		query = `
			INSERT INTO breeding_entries (animal_id, bull_id, breeding_date, observation, user_id)
			VALUES (:animal_id, :bull_id, :breeding_date, :observation, :user_id)
		`
	}

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *BreedingRepository) Update(newEntry *BreedingEntrySave) (*BreedingDB, *log.APIError) {

	query := `
		UPDATE breeding_entries
		SET bull_id = :bull_id, 
	 		breeding_date = :breeding_date, 
			observation = :observation
		WHERE id = :id AND user_id = :user_id
	`

	err := util.NamedExec(r.DB, query, newEntry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	selectQuery := fmt.Sprintf(`
		SELECT 
			i.id,
			i.animal_id,
			CONCAT_WS(' - ', a.tag, a.name) animal_info,
			i.breeding_date,
			i.bull_id,
			b.name AS bull_name,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.breeding_date
						AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.breeding_date
						AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'SUCCESS'
				) THEN 'SUCCESS'
				WHEN NOT EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
					  AND t.test_date > i.breeding_date
					  AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
				) AND age(i.breeding_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS pregnancy_status,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.breeding_date
						AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN age(i.breeding_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,
			CASE 
				WHEN c.child_name IS NULL THEN 'Sem Cria'
				ELSE c.child_name
			END AS child_information,
			i.observation
		FROM breeding_entries i
			LEFT JOIN animals a ON a.id = i.animal_id
			LEFT JOIN animals b ON b.id = i.bull_id
			LEFT JOIN LATERAL (
				SELECT
				CONCAT_WS(
					' - ', 
					a.tag, 
					COALESCE(a.name, a.sex), 
					TO_CHAR(a.birth_date, 'DD/MM/YYYY')
				) AS child_name
				FROM animals a
				WHERE a.mother_id = i.animal_id
					AND a.birth_date > i.breeding_date
					AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					AND NOT EXISTS (
						SELECT 1
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.pregnancy_status = 'FAILED'
							AND t.test_date BETWEEN a.birth_date AND i.breeding_date
					)
				ORDER BY a.birth_date
				LIMIT 1
			) c ON TRUE
		WHERE i.id = $1
			AND i.user_id = $2
			AND i.deleted_at IS NULL
	`, util.MinGestantionDays, util.MaxGestationDays)

	res, err := util.GetOne[BreedingDB](r.DB, selectQuery, newEntry.Id, newEntry.UserId)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return res, nil
}

func (r *BreedingRepository) UpdateGroup(date time.Time, group *BreedingGroup) (*BreedingGroup, *log.APIError) {
	mainQuery := `
		UPDATE breeding_entries
		SET breeding_date = $1
		WHERE breeding_date = $2 AND user_id = $3
	`

	err := util.Exec(r.DB, mainQuery, group.BreedingDate, date, group.UserId)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	responseQuery := fmt.Sprintf(`
		WITH status AS (
			SELECT
				i.breeding_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.breeding_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.breeding_date
							AND age(t.test_date, i.breeding_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.breeding_date
							AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
								AND t.test_date BETWEEN i.breeding_date AND a.birth_date
								AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM breeding_entries i
			WHERE i.breeding_date = $1 
				AND i.user_id = $2 
				AND i.deleted_at IS NULL
		),
        totals AS (
            SELECT 
                breeding_date,
				COUNT(*) cow_number,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') birth_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status i
            GROUP BY breeding_date
        ),
        rates AS (
            SELECT
                breeding_date,
                cow_number,
                (birth_success::float / cow_number::float)*100 birth_rate,
                (pregnancy_success::float / cow_number::float)*100 pregnancy_rate
            FROM totals
        )
        SELECT 
            s.breeding_date,
            s.cow_number,
            s.birth_rate,
            s.pregnancy_rate,
            COALESCE(
				(s.birth_rate / NULLIF(LAG(s.birth_rate) OVER win, 0)) - 1, 0
			) * 100 AS birth_comparison_rate,
            COALESCE(
				(s.pregnancy_rate / NULLIF(LAG(s.pregnancy_rate) OVER win, 0)) - 1, 0
			) * 100 AS pregnancy_comparison_rate
        FROM rates s
		WINDOW win AS (ORDER BY s.breeding_date)
    `, util.MinGestantionDays, util.MaxGestationDays)

	response, err := util.GetOne[BreedingGroup](r.DB, responseQuery, group.BreedingDate, group.UserId)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return response, nil
}

func (r *BreedingRepository) DeleteGroup(date time.Time, userId string) *log.APIError {

	tx, err := r.DB.Beginx()
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	selectQuery := fmt.Sprintf(`
		SELECT 
			a.id, 
			a.birth_date,
			a.mother_id,
			user_id
		FROM breeding_entries i
			JOIN animals a ON a.mother_id = i.animal_id
				AND a.birth_date > i.breeding_date
				AND age(a.birth_date, i.breeding_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				AND NOT EXISTS (
					SELECT 1
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.pregnancy_status = 'FAILED'
						AND t.test_date BETWEEN a.birth_date AND i.breeding_date
				)
		WHERE i.user_id = $1 
			AND i.deleted_at IS NULL
			AND i.breeding_date = $2
	`, util.MinGestantionDays, util.MaxGestationDays)

	children, err := util.GetListTx(tx, selectQuery, animals.AnimalSave{}, userId, date)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	findFatherQuery := fmt.Sprintf(`
		WITH entries_query AS (
			SELECT p.pasture_id
			FROM pasture_entries p
			WHERE p.deleted_at IS NULL
				AND p.animal_id = $1
				AND ($2::timestamptz - INTERVAL '%d days') BETWEEN p.entry_date AND COALESCE(p.exit_date, NOW())
				AND p.user_id = $3
		)
		SELECT pe.animal_id
		FROM pasture_entries pe
			CROSS JOIN entries_query e
			JOIN animals a ON a.id = pe.animal_id
				AND a.animal_type = 'REPRODUCTION_ANIMAL'
				AND a.sex = 'M'
		WHERE pe.deleted_at IS NULL
			AND pe.user_id = $3
			AND ($2::timestamptz - INTERVAL '308 days') BETWEEN pe.entry_date AND COALESCE(pe.exit_date, NOW())
			AND pe.pasture_id = e.pasture_id
		ORDER BY (COALESCE(pe.exit_date, NOW()) - pe.entry_date) DESC
		LIMIT 1
	`, util.AverageGestationDays)

	for _, child := range *children {
		var fatherId sql.NullString
		err := util.GetPrimitiveTx(tx, findFatherQuery, &fatherId, child.MotherId, child.BirthDate, child.UserId)
		if err != nil {
			return log.InternalServerAPIError(err)
		}

		child.FatherId = nil
		if fatherId.Valid {
			child.FatherId = &fatherId.String
		}
	}

	updateFathersQuery := `
		UPDATE animals
		SET father_id = :father_id
		WHERE id = :id AND user_id = :user_id
	`

	err = util.NamedExecTx(tx, updateFathersQuery, children)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	mainQuery := `
		UPDATE breeding_entries
		SET deleted_at = NOW()
		WHERE breeding_date = $1 AND user_id = $2
	`

	err = util.ExecTx(tx, mainQuery, date, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	err = tx.Commit()
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}
