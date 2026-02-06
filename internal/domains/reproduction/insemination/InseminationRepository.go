package insemination

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/entity"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type InseminationRepository struct {
	DB *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) *InseminationRepository {
	return &InseminationRepository{db}
}

type InseminationValidation struct {
	Repeated    bool `db:"repeated"`
	HasChildren bool `db:"has_children"`
	IsPregnant  bool `db:"is_pregnant"`
}

func (r *InseminationRepository) CheckAddConflicts(entry InseminationEntrySave) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM insemination_entries
			WHERE animal_id = :animal_id
				AND insemination_date = :insemination_date
				AND user_id = :user_id
				AND deleted_at IS NULL
		) AS exists
	`
	var exists bool
	err := util.NamedPrimitive(r.DB, query, &exists, entry)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *InseminationRepository) CheckUpdateConflicts(entry InseminationEntrySave) (*InseminationValidation, error) {
	query := fmt.Sprintf(`
		SELECT 
			EXISTS (
				SELECT 1
				FROM insemination_entries
				WHERE animal_id = :animal_id
					AND insemination_date = :insemination_date
					AND user_id = :user_id
					AND id <> :id
					AND deleted_at IS NULL
			) AS repeated,
			EXISTS (
				SELECT 1
				FROM animals a
				WHERE a.deleted_at IS NULL 
					AND a.user_id = :user_id
					AND a.mother_id = :animal_id
					AND a.birth_date > :insemination_date
					AND age(a.birth_date, :insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					AND a.father_id = :bull_id
			) AS has_children,
			(
				SELECT EXISTS (
					SELECT 1
					FROM pregnancy_tests t
					WHERE t.pregnancy_status = 'SUCCESS'
						AND t.animal_id = :animal_id
						AND t.test_date > :insemination_date
						AND age(t.test_date, :insemination_date) <= INTERVAL '%[2]d days'
						AND t.user_id = :user_id
						AND t.deleted_at IS NULL
				) AND NOT EXISTS (
					SELECT 1
					FROM pregnancy_tests t
					WHERE t.pregnancy_status = 'FAILED'
						AND t.animal_id = :animal_id
						AND t.test_date > :insemination_date
						AND age(t.test_date, :insemination_date) <= INTERVAL '%[2]d days'
						AND t.user_id = :user_id
						AND t.deleted_at IS NULL
				)
			) AS is_pregnant
	`, util.MinGestantionDays, util.MaxGestationDays)
	return util.NamedGet(r.DB, query, InseminationValidation{}, entry)
}

func (r *InseminationRepository) CheckDeleteConflicts(params *InseminationEntryDelete) (*InseminationValidation, error) {
	query := fmt.Sprintf(`
		WITH entry_cte AS (
			SELECT 
				animal_id,
				bull_id,
				insemination_date
			FROM insemination_entries
			WHERE id = :id AND user_id = :user_id
		)

		SELECT
			EXISTS (
				SELECT 1
				FROM animals a
				CROSS JOIN entry_cte e
				WHERE a.deleted_at IS NULL 
					AND a.user_id = :user_id
					AND a.mother_id = e.animal_id
					AND a.birth_date > e.insemination_date
					AND a.father_id = e.bull_id
					AND age(a.birth_date, e.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
			) AS has_children,
			 EXISTS (
				SELECT 1
				FROM pregnancy_tests t
				CROSS JOIN entry_cte e
				WHERE t.pregnancy_status = 'SUCCESS'
					AND t.animal_id = e.animal_id
					AND t.test_date < e.insemination_date
					AND age(t.test_date, e.insemination_date) <= INTERVAL '%[2]d days'
					AND t.user_id = :user_id
					AND t.deleted_at IS NULL
			) AND NOT EXISTS (
				SELECT 1
				FROM pregnancy_tests t
				CROSS JOIN entry_cte e
				WHERE t.pregnancy_status = 'FAILED'
					AND t.animal_id = e.animal_id
					AND t.test_date > e.insemination_date
					AND age(t.test_date, e.insemination_date) <= INTERVAL '%[2]d days'
					AND t.user_id = :user_id
					AND t.deleted_at IS NULL
			) AS is_pregnant
	`, util.MinGestantionDays, util.MaxGestationDays)
	return util.NamedGet(r.DB, query, InseminationValidation{}, params)
}

func (r *InseminationRepository) CheckGroupConflicts(group InseminationGroupSave) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM insemination_entries
			WHERE insemination_date <> :old_insemination_date
				AND insemination_date = :insemination_date
				AND user_id = :user_id
				AND deleted_at IS NULL
		) 
	`
	var exists bool
	err := util.NamedPrimitive(r.DB, query, &exists, group)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *InseminationRepository) GetBirthRateStats(userId string) (*[]BirthRateHist, error) {
	query := fmt.Sprintf(`
		WITH totals AS (
			SELECT 
				i.insemination_date,
				COUNT(i.*) AS total,
				COUNT(DISTINCT i.id) FILTER (WHERE a.birth_date IS NOT NULL) AS birth_success
			FROM insemination_entries i
			LEFT JOIN animals a ON a.mother_id = i.animal_id
				AND a.birth_date > i.insemination_date
				AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				AND NOT EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.pregnancy_status = 'FAILED'
						AND t.test_date BETWEEN i.insemination_date AND a.birth_date
				)
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT 
			insemination_date, 
			(birth_success::float / NULLIF(total, 0)::float) * 100 AS birth_rate
		FROM totals
		ORDER BY 1
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[BirthRateHist](r.DB, query, userId)
}

func (r *InseminationRepository) GetPregnancyRateStats(userId string) (*[]PregnancyRateHist, error) {
	query := fmt.Sprintf(`
		WITH insemination_status AS (
			SELECT
				i.insemination_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
						  AND t.test_date > i.insemination_date
						  AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
						  AND t.test_date > i.insemination_date
						  AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status
			FROM insemination_entries i
			WHERE i.user_id = $1 AND i.deleted_at IS NULL 
		),
		cte AS (
			SELECT
				t.insemination_date,
				COUNT(t.*) AS total,
				COUNT(t.*) FILTER (WHERE t.pregnancy_status = 'SUCCESS') AS pregnancy_success
			FROM insemination_status t
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT 
			insemination_date,
			(pregnancy_success::float / NULLIF(total, 0)) * 100 AS pregnancy_rate
		FROM cte
		ORDER BY 1
    `, util.MinGestantionDays, util.MaxGestationDays)

	return util.GetList[PregnancyRateHist](r.DB, query, userId)
}

func (r *InseminationRepository) GetAnimalsNumber(userId string) (*[]AnimalsHist, error) {
	query := `
		WITH cte AS (
			SELECT
				insemination_date,
				COUNT(*) AS animals_number
			FROM insemination_entries
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT *
		FROM cte
		ORDER BY 1
    `
	return util.GetList[AnimalsHist](r.DB, query, userId)
}

func (r *InseminationRepository) GetInseminationStats(userId string) (*[]InseminationHist, error) {
	query := fmt.Sprintf(`
        WITH cte AS (
			SELECT 
				i.insemination_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
						  AND t.test_date > i.insemination_date
						  AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
						  AND t.test_date > i.insemination_date
						  AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM insemination_entries i
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
		),
        totals AS (
            SELECT
                insemination_date,
                COUNT(*) total,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') birth_numbers,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_numbers
            FROM cte
            GROUP BY 1
            ORDER BY 1 DESC
            LIMIT 30
        )
        SELECT * FROM totals ORDER BY insemination_date
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[InseminationHist](r.DB, query, userId)
}

func (r *InseminationRepository) GetFutureBirths(userId string) (*[]FutureBirths, error) {
	query := fmt.Sprintf(`
		WITH upcoming_births AS (
			SELECT 
				i.id,
				t.test_date::date + (%[3]d - t.pregnancy_time) AS birth_forecast
			FROM insemination_entries i
				JOIN pregnancy_tests t ON t.animal_id = i.animal_id
					AND t.test_date > i.insemination_date
					AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
					AND t.pregnancy_status = 'SUCCESS'
			WHERE i.user_id = $1
				AND i.deleted_at IS NULL
				AND t.test_date::date + (%[3]d - pregnancy_time) >= NOW()  
				AND NOT EXISTS (
					SELECT 1
					FROM animals a 
					WHERE a.mother_id = i.animal_id
						AND a.birth_date > i.insemination_date
						AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
						AND NOT EXISTS (
							SELECT 1
							FROM pregnancy_tests f
							WHERE f.animal_id = i.animal_id
								AND f.test_date BETWEEN i.insemination_date AND a.birth_date
								AND f.pregnancy_status = 'FAILED'
						)

				)
		)
		SELECT
			DATE_TRUNC('month', birth_forecast) AS birth_forecast,
			COUNT(DISTINCT id) AS births_number
		FROM upcoming_births
		GROUP BY 1
		ORDER BY 1;
	`, util.MinGestantionDays, util.MaxGestationDays, util.MinGestantionDays)
	return util.GetList[FutureBirths](r.DB, query, userId)
}

func (r *InseminationRepository) GetBestBull(userId string) (*[]InseminationBulls, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				CONCAT_WS(' - ', b.ring_number, b.name) bull_name,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
								AND t.test_date BETWEEN i.insemination_date AND a.birth_date
								AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM insemination_entries i
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
	return util.GetList[InseminationBulls](r.DB, query, userId)
}

func (r *InseminationRepository) GetLastGroups(userId string) (*[]InseminationGroup, error) {
	query := fmt.Sprintf(`
		WITH insemination_data AS (
			SELECT
				i.insemination_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM insemination_entries i
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
		),
		daily_stats AS (
			SELECT
				insemination_date,
				COUNT(*) AS cow_number,
				COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') AS birth_success,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') AS pregnancy_success
			FROM insemination_data
			GROUP BY insemination_date
		),
		rates AS (
			SELECT
				insemination_date,
				cow_number,
				(birth_success::float * 100 / NULLIF(cow_number, 0)) AS birth_rate,
				(pregnancy_success::float * 100 / NULLIF(cow_number, 0)) AS pregnancy_rate
			FROM daily_stats
		)
		SELECT
			insemination_date,
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
		WINDOW win AS (ORDER BY insemination_date)
		ORDER BY insemination_date DESC
		LIMIT 5;
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[InseminationGroup](r.DB, query, userId)
}

func (r *InseminationRepository) GetLastEntries(userId string) (*[]InseminationEntry, error) {
	query := fmt.Sprintf(`

		WITH last_date_cte AS (
			SELECT MAX(insemination_date) last_date
			FROM insemination_entries
			WHERE user_id = $1 AND deleted_at IS NULL
		)

		SELECT 
			i.id,
			i.insemination_date,
			i.bull_id,
			CONCAT_WS(' - ', a.ring_number, a.name) AS animal_info,
			b.name AS bull_name,
			CASE
				WHEN EXISTS (
					SELECT 1 
					FROM animals a
					WHERE a.mother_id = i.animal_id
						AND a.birth_date > i.insemination_date
						AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
						AND NOT EXISTS (
							SELECT 1
							FROM pregnancy_tests t
							WHERE t.animal_id = a.mother_id
							AND t.test_date BETWEEN i.insemination_date AND a.birth_date
							AND t.pregnancy_status = 'FAILED'
						)
				) THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.insemination_date
						AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.insemination_date
						AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'SUCCESS'
				) THEN 'SUCCESS'
				WHEN NOT EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.insemination_date
						AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
				) AND age(i.insemination_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS pregnancy_status,
			CASE
				WHEN EXISTS (
					SELECT 1 
					FROM animals a
					WHERE a.mother_id = i.animal_id
						AND a.birth_date > i.insemination_date
						AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
						AND NOT EXISTS (
							SELECT 1
							FROM pregnancy_tests t
							WHERE t.animal_id = a.mother_id
							AND t.test_date BETWEEN i.insemination_date AND a.birth_date
							AND t.pregnancy_status = 'FAILED'
						)
				) THEN 'SUCCESS'
				WHEN age(i.insemination_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status
		FROM insemination_entries i
		CROSS JOIN last_date_cte l
		LEFT JOIN animals a ON a.id = i.animal_id
		LEFT JOIN animals b ON b.id = i.bull_id
		WHERE i.user_id = $1 
			AND i.insemination_date = l.last_date
			AND i.deleted_at IS NULL
		ORDER BY COALESCE(REGEXP_REPLACE(a.ring_number, '[^0-9]', '', 'g')::int, 0);
    `, util.MinGestantionDays, util.MaxGestationDays)

	return util.GetList[InseminationEntry](r.DB, query, userId)
}

func (r *InseminationRepository) FindEntriesPage(
	userId string,
	filter *InseminationEntryFilter,
	sort string,
	order string,
	cursor string,
) (*entity.Page[InseminationEntry], error) {

	sort = util.AddCommonFields(sort)
	sortMap := map[string]util.SortField{
		"animal_order":      {Field: "cte.animal_order", Order: "asc"},
		"animal_name":       {Field: "cte.animal_name", Order: "asc"},
		"insemination_date": {Field: "coalesce(cte.insemination_date, '-infinity')", Order: "asc"},
		"id":                {Field: "cte.id", Order: "asc"},
		"created_at":        {Field: "cte.created_at", Order: "asc"},
	}

	query := fmt.Sprintf(`
        WITH cte AS (
			SELECT 
				i.id,
				i.animal_id,
				a.name AS animal_name,
				COALESCE(REGEXP_REPLACE(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
				CONCAT_WS(' - ', a.ring_number, a.name) animal_info,
				i.insemination_date,
				i.bull_id,
				b.name AS bull_name,
				CASE
					WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					WHEN NOT EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
						  AND t.test_date > i.insemination_date
						  AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
					) AND age(i.insemination_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
						  AND t.test_date > i.insemination_date
						  AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN age(i.insemination_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
					ELSE 'FAILED'
				END AS birth_status,
				CASE 
					WHEN c.child_name IS NULL THEN 'Sem Cria'
					ELSE c.child_name
				END AS child_information,
				i.observation,
				i.created_at
			FROM insemination_entries i
			LEFT JOIN animals a ON a.id = i.animal_id
			LEFT JOIN animals b ON b.id = i.bull_id
			LEFT JOIN LATERAL (
				SELECT
				CONCAT_WS(
					' - ', 
					a.ring_number, 
					COALESCE(a.name, a.sex), 
					TO_CHAR(a.birth_date, 'DD/MM/YYYY')
				) AS child_name
				FROM animals a
				WHERE a.mother_id = i.animal_id
					AND a.birth_date > i.insemination_date
					AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					AND NOT EXISTS (
						SELECT 1
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.pregnancy_status = 'FAILED'
							AND t.test_date BETWEEN a.birth_date AND i.insemination_date
					)
				ORDER BY a.birth_date
				LIMIT 1
			) c ON TRUE
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
		)
		SELECT * FROM cte
	`, util.MinGestantionDays, util.MaxGestationDays)

	orderExpression := " ORDER BY "
	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "cte", 2)
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
	return util.GetPage[InseminationEntry](r.DB, query, sort, 100, args...)
}

func (r *InseminationRepository) GetEntriesFoot(
	userId string,
	filter *InseminationEntryFilter,
) (*InseminationFoot, error) {

	statusQuery := fmt.Sprintf(`
		WITH cte AS  (
			SELECT
				i.*,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM insemination_entries i
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

	return util.GetOne[InseminationFoot](r.DB, query, args...)
}

func (r *InseminationRepository) FindEntriesByGroup(userId string, date time.Time) (*[]InseminationEntry, error) {

	query := fmt.Sprintf(`
        SELECT 
            i.id,
			i.animal_id,
			i.bull_id,
			CONCAT_WS(' - ', b.ring_number, b.name) AS bull_name,
            CONCAT_WS(' - ', a.ring_number, a.name) AS animal_info,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.insemination_date
						AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.insemination_date
						AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'SUCCESS'
				) THEN 'SUCCESS'
				WHEN NOT EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
					  AND t.test_date > i.insemination_date
					  AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
				) AND age(i.insemination_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS pregnancy_status,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.insemination_date
						AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN age(i.insemination_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,
			CASE
				WHEN c.child_name IS NULL THEN 'Sem Cria'
				ELSE child_name
			END AS child_information,
            i.observation
        FROM insemination_entries i
            LEFT JOIN animals a ON a.id = i.animal_id
            LEFT JOIN animals b ON b.id = i.bull_id
			LEFT JOIN LATERAL (
				SELECT
				CONCAT_WS(
					' - ', 
					a.ring_number, 
					COALESCE(a.name, a.sex), 
					TO_CHAR(a.birth_date, 'DD/MM/YYYY')
				) AS child_name
				FROM animals a
				WHERE a.mother_id = i.animal_id
					AND a.birth_date > i.insemination_date
					AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					AND NOT EXISTS (
						SELECT 1
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.pregnancy_status = 'FAILED'
							AND t.test_date BETWEEN a.birth_date AND i.insemination_date
					)
				ORDER BY a.birth_date
				LIMIT 1
			) c ON TRUE
		WHERE i.user_id = $1 AND i.deleted_at IS NULL AND i.insemination_date = $2
        ORDER BY COALESCE(REGEXP_REPLACE(a.ring_number, '[^0-9]', '', 'g')::int, 0)
	`, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[InseminationEntry](r.DB, query, userId, date)
}

func (r *InseminationRepository) GetEntriesByGroupFoot(userId string, date time.Time) (*InseminationFoot, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM insemination_entries i
			WHERE i.user_id = $1 
				AND i.insemination_date = $2
				AND i.deleted_at IS NULL
		),
        COUNTING AS (
            SELECT 
                COUNT(*) totals,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') birth_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status
        )
        SELECT 
            totals,
            (birth_success::float / NULLIF(totals, 0)) * 100 average_birth_rate,
            (pregnancy_success::float / NULLIF(totals, 0)) * 100 average_pregnancy_rate
        FROM COUNTING
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetOne[InseminationFoot](r.DB, query, userId, date)
}

func (r *InseminationRepository) FindGroups(userId string) (*[]InseminationGroup, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				i.insemination_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM insemination_entries i
			WHERE i.user_id = $1 AND i.deleted_at IS NULL
		),
        totals AS (
            SELECT 
                insemination_date,
				COUNT(*) cow_number,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') birth_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status i
            GROUP BY insemination_date
        ),
        rates AS (
            SELECT
                insemination_date,
                cow_number,
                (birth_success::float / cow_number::float)*100 birth_rate,
                (pregnancy_success::float / cow_number::float)*100 pregnancy_rate
            FROM totals
        )
        SELECT 
            s.insemination_date,
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
		WINDOW win AS (ORDER BY s.insemination_date)
        ORDER BY s.insemination_date DESC
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[InseminationGroup](r.DB, query, userId)
}

func (r *InseminationRepository) Add(entry *InseminationEntrySave) error {
	var query string
	if entry.Overwrite {
		query = `
			INSERT INTO insemination_entries (animal_id, bull_id, insemination_date, observation, user_id)
			VALUES (:animal_id, :bull_id, :insemination_date, :observation, :user_id)
		`
	} else {
		query = `
			UPDATE insemination_entries
			SET bull_id = :bull_id,
				observation = :observation
			WHERE animal_id = :animal_id
				AND insemination_date = :insemination_date
				AND user_id = :user_id
		`
	}

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return err
	}

	return nil
}

func (r *InseminationRepository) FindById(id string, userId string) (*InseminationEntrySave, error) {
	query := `
		SELECT id, animal_id, bull_id, insemination_date, user_id
		FROM insemination_entries
		WHERE id = $1 AND user_id = $2
	`
	return util.GetOne[InseminationEntrySave](r.DB, query, id, userId)
}

func (r *InseminationRepository) Delete(params *InseminationEntryDelete) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if !params.ChangeFather {
		fatherQuery := fmt.Sprintf(`
			WITH insemination_cte AS (
				SELECT 
					animal_id, 
					insemination_date
				FROM insemination_entries
				WHERE id = :id 
					AND user_id = :user_id
					AND deleted_at IS NULL
			),

			children AS (
				SELECT 
					a.id AS child_id,
					a.mother_id,
					(a.birth_date::timestamptz - INTERVAL '%[3]d days') AS conception_date
				FROM animals a
				JOIN insemination_cte i ON a.mother_id = i.animal_id
				WHERE a.birth_date > i.insemination_date
				  AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				  AND a.user_id = :user_id AND a.deleted_at IS NULL
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
					AND mother_past.user_id = :user_id
				JOIN pasture_entries bull_past 
					ON bull_past.pasture_id = mother_past.pasture_id
					AND c.conception_date BETWEEN bull_past.entry_date AND COALESCE(bull_past.exit_date, NOW())
					AND bull_past.deleted_at IS NULL
					AND bull_past.user_id = :user_id
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
				AND a.user_id = :user_id
				AND a.deleted_at IS NULL
		`, util.MinGestantionDays, util.MaxGestationDays, util.AverageGestationDays)

		err = util.NamedExecTx(tx, fatherQuery, params)
		if err != nil {
			return err
		}

	}

	query := `
		UPDATE insemination_entries
		SET deleted_at = NOW()
		WHERE id = :id AND user_id = :user_id
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

func (r *InseminationRepository) Update(newEntry *InseminationEntrySave) (*InseminationEntry, error) {
	query := `
		UPDATE insemination_entries
		SET bull_id = :bull_id, 
	 		insemination_date = :insemination_date, 
			observation = :observation
		WHERE id = :id AND user_id = :user_id
	`

	err := util.NamedExec(r.DB, query, newEntry)
	if err != nil {
		return nil, err
	}

	selectQuery := fmt.Sprintf(`
		SELECT 
			i.id,
			i.animal_id,
			CONCAT_WS(' - ', a.ring_number, a.name) animal_info,
			i.insemination_date,
			i.bull_id,
			b.name AS bull_name,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.insemination_date
						AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.insemination_date
						AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'SUCCESS'
				) THEN 'SUCCESS'
				WHEN NOT EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
					  AND t.test_date > i.insemination_date
					  AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
				) AND age(i.insemination_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS pregnancy_status,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = i.animal_id
						AND t.test_date > i.insemination_date
						AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
						AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN age(i.insemination_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,
			CASE 
				WHEN c.child_name IS NULL THEN 'Sem Cria'
				ELSE c.child_name
			END AS child_information,
			i.observation
		FROM insemination_entries i
			LEFT JOIN animals a ON a.id = i.animal_id
			LEFT JOIN animals b ON b.id = i.bull_id
			LEFT JOIN LATERAL (
				SELECT
				CONCAT_WS(
					' - ', 
					a.ring_number, 
					COALESCE(a.name, a.sex), 
					TO_CHAR(a.birth_date, 'DD/MM/YYYY')
				) AS child_name
				FROM animals a
				WHERE a.mother_id = i.animal_id
					AND a.birth_date > i.insemination_date
					AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					AND NOT EXISTS (
						SELECT 1
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.pregnancy_status = 'FAILED'
							AND t.test_date BETWEEN a.birth_date AND i.insemination_date
					)
				ORDER BY a.birth_date
				LIMIT 1
			) c ON TRUE
		WHERE i.id = $1
			AND i.user_id = $2
			AND i.deleted_at IS NULL
	`, util.MinGestantionDays, util.MaxGestationDays)
	res, err := util.GetOne[InseminationEntry](r.DB, selectQuery, newEntry.Id, newEntry.UserId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *InseminationRepository) UpdateGroup(group *InseminationGroupSave) (*InseminationGroup, error) {

	mainQuery := `
		UPDATE insemination_entries
		SET insemination_date = :insemination_date
		WHERE insemination_date = :old_insemination_date 
			AND user_id = :user_id
			AND deleted_at IS NULL
	`

	err := util.NamedExec(r.DB, mainQuery, group)
	if err != nil {
		return nil, err
	}

	responseQuery := fmt.Sprintf(`
		WITH status AS (
			SELECT
				i.insemination_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = i.animal_id
							AND t.test_date > i.insemination_date
							AND age(t.test_date, i.insemination_date) <= INTERVAL '%[2]d days'
							AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = i.animal_id
							AND a.birth_date > i.insemination_date
							AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
							AND NOT EXISTS (
								SELECT 1
								FROM pregnancy_tests t
								WHERE t.animal_id = a.mother_id
									AND t.test_date BETWEEN i.insemination_date AND a.birth_date
									AND t.pregnancy_status = 'FAILED'
							)
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM insemination_entries i
			WHERE i.insemination_date = $1 
				AND i.user_id = $2 
				AND i.deleted_at IS NULL
		),

        totals AS (
            SELECT 
                insemination_date,
				COUNT(*) cow_number,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') birth_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status i
            GROUP BY insemination_date
        ),

        rates AS (
            SELECT
                insemination_date,
                cow_number,
                (birth_success::float / cow_number::float)*100 birth_rate,
                (pregnancy_success::float / cow_number::float)*100 pregnancy_rate
            FROM totals
        )

        SELECT 
            s.insemination_date,
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
		WINDOW win AS (ORDER BY s.insemination_date)
    `, util.MinGestantionDays, util.MaxGestationDays)

	response, err := util.NamedGet(r.DB, responseQuery, InseminationGroup{}, group)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *InseminationRepository) DeleteGroup(params *InseminationGroupDelete) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if !params.ChangeFather {
		fatherQuery := fmt.Sprintf(`
			WITH insemination_cte AS (
				SELECT 
					animal_id, 
					insemination_date
				FROM insemination_entries
				WHERE insemination_date = :insemination_date
					AND user_id = :user_id
					AND deleted_at IS NULL
			),

			children AS (
				SELECT 
					a.id AS child_id,
					a.mother_id,
					(a.birth_date::timestamptz - INTERVAL '%[3]d days') AS conception_date
				FROM animals a
				JOIN insemination_cte i ON a.mother_id = i.animal_id
				WHERE a.birth_date > i.insemination_date
				  AND age(a.birth_date, i.insemination_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				  AND a.user_id = :user_id 
				  AND a.deleted_at IS NULL
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
					AND mother_past.user_id = :user_id
				JOIN pasture_entries bull_past 
					ON bull_past.pasture_id = mother_past.pasture_id
					AND c.conception_date BETWEEN bull_past.entry_date AND COALESCE(bull_past.exit_date, NOW())
					AND bull_past.deleted_at IS NULL
					AND bull_past.user_id = :user_id
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
				AND a.user_id = :user_id
				AND a.deleted_at IS NULL
	`, util.MinGestantionDays, util.MaxGestationDays, util.AverageGestationDays)

		err = util.NamedExecTx(tx, fatherQuery, params)
		if err != nil {
			return err
		}
	}

	mainQuery := `
		UPDATE insemination_entries
		SET deleted_at = NOW()
		WHERE insemination_date = :insemination_date AND user_id = :user_id
	`

	err = util.NamedExecTx(tx, mainQuery, params)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
