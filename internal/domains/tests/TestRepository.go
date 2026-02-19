package tests

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type TestEntryRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *TestEntryRepository {
	return &TestEntryRepository{db}
}

func (r *TestEntryRepository) GetPregnancyRate(userId string) (*[]util.GraphData, error) {
	query := `
        WITH cte AS (
            SELECT 
                test_date,
                COUNT(*) AS totals,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') AS pregnancies
            FROM pregnancy_tests
            WHERE deleted_at IS NULL AND user_id = $1 
			GROUP BY 1
			ORDER BY test_date DESC
            LIMIT 10
        )
        SELECT 
            test_date AS date,
            (pregnancies::float / NULLIF(totals, 0)) * 100 AS value
        FROM cte
        ORDER BY test_date
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *TestEntryRepository) GetAnimalsNumber(userId string) (*[]util.GraphData, error) {
	query := `
        WITH cte AS (
            SELECT 
                test_date AS date,
                COUNT(*) AS value
            FROM pregnancy_tests
            WHERE deleted_at IS NULL AND user_id = $1 
			GROUP BY 1
			ORDER BY test_date DESC
            LIMIT 10
        )
        SELECT date, value
        FROM cte
        ORDER BY date
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *TestEntryRepository) GetBirthRate(userId string) (*[]util.GraphData, error) {
	query := `
        WITH cte AS (
            SELECT
                test_date,
                COUNT(*) AS totals,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS' 
					AND EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = t.animal_id
							AND a.birth_date > t.test_date
							AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
					)
				) AS births
            FROM pregnancy_tests t
            WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY test_date
			ORDER BY test_date DESC
            LIMIT 10
        )
        SELECT 
            test_date AS date,
            (births::float / NULLIF(totals, 0)) * 100 AS value
        FROM cte
        ORDER BY test_date
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *TestEntryRepository) GetPregnancyTestHist(userId string) (*[]PregnancyTestHist, error) {
	query := `
        WITH cte AS (
            SELECT 
                test_date,
                COUNT(*) totals,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS' 
					AND EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = t.animal_id
							AND a.birth_date > t.test_date
							AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
					)
				) AS births,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancies
            FROM pregnancy_tests t
            WHERE user_id = $1 AND deleted_at IS NULL 
			GROUP BY 1
			ORDER BY test_date DESC
            LIMIT 36
        )
        SELECT 
            test_date,
            totals,
			pregnancies,
			births
        FROM cte
        ORDER BY test_date
    `
	return util.GetList[PregnancyTestHist](r.DB, query, userId)
}

func (r *TestEntryRepository) GetLastEntries(userId string) (*[]TestDB, error) {
	query := fmt.Sprintf(`
		WITH max_cte AS (
			SELECT MAX(test_date) AS max_date
			FROM pregnancy_tests 
			WHERE user_id = $1 AND deleted_at IS NULL
		)

        SELECT
			t.id,
            t.test_date,
			test_date::date + pregnancy_time AS birth_forecast,
            t.pregnancy_status,
            t.observation,
			CASE
				WHEN pregnancy_status = 'FAILED' THEN 'FAILED'
				WHEN EXISTS (
					SELECT 1 
					FROM animals a
					WHERE a.mother_id = t.animal_id
						AND a.birth_date > t.test_date
						AND age(a.birth_date, t.test_date) <= INTERVAL '%[1]d days'
				) THEN 'SUCCESS'
				WHEN age(t.test_date) < INTERVAL '%[1]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,

			t.animal_id,
			a.tag AS animal_tag,
			a.name AS animal_name
        FROM pregnancy_tests t
			CROSS JOIN max_cte m
            LEFT JOIN animals a ON a.id = t.animal_id
        WHERE t.test_date = m.max_date
			AND t.user_id = $1 
			AND t.deleted_at IS NULL
        ORDER BY COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0)
    `, util.MaxGestationDays)
	return util.GetList[TestDB](r.DB, query, userId)
}

func (r *TestEntryRepository) GetLastGroups(userId string) (*[]TestGroups, error) {
	query := `
        WITH totals AS (
            SELECT 
                test_date,
                COUNT(*) animals_number,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS' 
					AND EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = t.animal_id
							AND a.birth_date > t.test_date
							AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
					)
				) AS birth_success
            FROM pregnancy_tests t
            WHERE deleted_at IS NULL AND user_id = $1 
            GROUP BY 1
            LIMIT 6
        ),
        rates AS (
            SELECT
                test_date,
                animals_number,
                (pregnancy_success::float / NULLIF(animals_number, 0)) * 100 pregnancy_rate,
                (birth_success::float / NULLIF(animals_number, 0)) * 100 birth_rate
            FROM totals
        )
        SELECT
            test_date,
            animals_number,
            pregnancy_rate,
            birth_rate,
            COALESCE(
				(pregnancy_rate / LAG(pregnancy_rate) OVER win) - 1, 0
			) * 100 AS pregnancy_comparison,
            COALESCE(
				(birth_rate / LAG(birth_rate) OVER win) - 1, 0
			) * 100 AS birth_comparison
        FROM rates
		WINDOW win AS (ORDER BY test_date)
        ORDER BY test_date DESC
    `
	return util.GetList[TestGroups](r.DB, query, userId)
}

func (r *TestEntryRepository) GetNextBirths(userId string) (*[]util.GraphData, error) {
	query := `
        SELECT 
			DATE_TRUNC('month', t.test_date::date + t.pregnancy_time) AS date,
            COUNT(*) AS value
        FROM pregnancy_tests t
        WHERE 
            deleted_at IS NULL 
            AND user_id = $1
            AND t.test_date + (t.pregnancy_time * INTERVAL '1 day') > NOW()
            AND pregnancy_status = 'SUCCESS'
            AND age(t.test_date) < INTERVAL '340 days'
			AND NOT EXISTS (
				SELECT 1
				FROM animals a
				WHERE a.mother_id = t.animal_id
					AND a.birth_date > t.test_date
					AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
			)
        GROUP BY 1
        ORDER BY 1
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *TestEntryRepository) GetBestResults(userId string) (*[]TestAnimal, error) {
	query := `
		WITH birth_test_enriched AS (
			SELECT 
				bt.animal_id,
				bt.test_date,
				bt.pregnancy_status,
				EXISTS (
					SELECT 1
					FROM animals a
					WHERE a.mother_id = bt.animal_id
					  AND a.birth_date > bt.test_date
					  AND age(a.birth_date, bt.test_date) <= INTERVAL '340 days'
				) AS has_valid_birth
			FROM pregnancy_tests bt
			WHERE bt.deleted_at IS NULL
			  AND bt.user_id = $1
		),
		totals AS (
			SELECT 
				animal_id,
				COUNT(*) AS totals,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') AS pregnancy_success,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS' AND has_valid_birth) AS birth_success
			FROM birth_test_enriched
			GROUP BY animal_id
			HAVING COUNT(*) >= 5
		),
		rates AS (
			SELECT
				animal_id,
				totals,
				(pregnancy_success::float / totals) * 100 AS pregnancy_rate,
				(birth_success::float / totals) * 100 AS birth_rate
			FROM totals
		),
		general_totals AS (
			SELECT 
				COUNT(*) AS totals,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') AS pregnancy_success,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS' AND has_valid_birth) AS birth_success
			FROM birth_test_enriched
		),
		general_rates AS (
			SELECT
				(pregnancy_success::float / NULLIF(totals, 0)) * 100 AS total_pregnancy_rate,
				(birth_success::float / NULLIF(totals, 0)) * 100 AS total_birth_rate
			FROM general_totals
		),
		scores AS (
			SELECT
				animal_id,
				totals,
				pregnancy_rate,
				birth_rate,
				(birth_rate - AVG(birth_rate) OVER ()) / NULLIF(STDDEV(birth_rate) OVER (), 0) AS birth_score,
				(pregnancy_rate - AVG(pregnancy_rate) OVER ()) / NULLIF(STDDEV(pregnancy_rate) OVER (), 0) AS pregnancy_score
			FROM rates 
		)
		SELECT
			CONCAT_WS(' - ', a.tag, a.name) AS animal_name,
			s.totals,
			s.pregnancy_rate,
			s.birth_rate,
			COALESCE((s.pregnancy_rate / NULLIF(gr.total_pregnancy_rate, 0)) - 1, 0) * 100 AS pregnancy_comparison,
			COALESCE((s.birth_rate / NULLIF(gr.total_birth_rate, 0)) - 1, 0) * 100 AS birth_comparison
		FROM scores s
		CROSS JOIN general_rates gr
		JOIN animals a ON a.id = s.animal_id
		WHERE (s.birth_score + s.pregnancy_score) > 0
		ORDER BY (s.birth_score * 0.7 + s.pregnancy_score * 0.3) DESC
		LIMIT 10;
    `
	return util.GetList[TestAnimal](r.DB, query, userId)
}

func (r *TestEntryRepository) GetWorstResults(userId string) (*[]TestAnimal, error) {
	query := `
		WITH birth_test_enriched AS (
			SELECT 
				bt.animal_id,
				bt.test_date,
				bt.pregnancy_status,
				EXISTS (
					SELECT 1
					FROM animals a
					WHERE a.mother_id = bt.animal_id
					  AND a.birth_date > bt.test_date
					  AND age(a.birth_date, bt.test_date) <= INTERVAL '340 days'
				) AS has_valid_birth
			FROM pregnancy_tests bt
			WHERE bt.deleted_at IS NULL
			  AND bt.user_id = $1
		),
		totals AS (
			SELECT 
				animal_id,
				COUNT(*) AS totals,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') AS pregnancy_success,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS' AND has_valid_birth) AS birth_success
			FROM birth_test_enriched
			GROUP BY animal_id
			HAVING COUNT(*) >= 5
		),
		rates AS (
			SELECT
				animal_id,
				totals,
				(pregnancy_success::float / totals) * 100 AS pregnancy_rate,
				(birth_success::float / totals) * 100 AS birth_rate
			FROM totals
		),
		general_totals AS (
			SELECT 
				COUNT(*) AS totals,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') AS pregnancy_success,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS' AND has_valid_birth) AS birth_success
			FROM birth_test_enriched
		),
		general_rates AS (
			SELECT
				(pregnancy_success::float / NULLIF(totals, 0)) * 100 AS total_pregnancy_rate,
				(birth_success::float / NULLIF(totals, 0)) * 100 AS total_birth_rate
			FROM general_totals
		),
		scores AS (
			SELECT
				animal_id,
				totals,
				pregnancy_rate,
				birth_rate,
				(birth_rate - AVG(birth_rate) OVER ()) / NULLIF(STDDEV(birth_rate) OVER (), 0) AS birth_score,
				(pregnancy_rate - AVG(pregnancy_rate) OVER ()) / NULLIF(STDDEV(pregnancy_rate) OVER (), 0) AS pregnancy_score
			FROM rates 
		)
		SELECT
			CONCAT_WS(' - ', a.tag, a.name) AS animal_name,
			totals,
			pregnancy_rate,
			birth_rate,
			COALESCE((pregnancy_rate / NULLIF(total_pregnancy_rate, 0)) - 1, 0) * 100 AS pregnancy_comparison,
			COALESCE((birth_rate / NULLIF(total_birth_rate, 0)) - 1, 0) * 100 AS birth_comparison
		FROM scores s
		CROSS JOIN general_rates gr
		JOIN animals a ON a.id = s.animal_id
		WHERE (birth_score + pregnancy_score) < 0
		ORDER BY (-birth_score * 0.7 - pregnancy_score * 0.3) DESC
		LIMIT 10;
    `
	return util.GetList[TestAnimal](r.DB, query, userId)
}

func (r *TestEntryRepository) FindEntriesPage(
	filter *TestFilter,
	sort string,
	order string,
	cursor string,
	limit int,
	userId string,
) (*[]TestDB, error) {

	sortMap := map[string]util.SortField{
		"animal_order":   {Field: "cte.animal_order", Order: "asc"},
		"test_date":      {Field: "cte.test_date", Order: "desc"},
		"birth_forecast": {Field: "coalesce(cte.birth_forecast, '-infinity')", Order: "desc"},
		"animal_name":    {Field: "cte.animal_name", Order: "asc"},
		"id":             {Field: "cte.id", Order: "asc"},
		"created_at":     {Field: "cte.created_at", Order: "asc"},
	}

	query := fmt.Sprintf(`
        WITH cte AS (
			SELECT
				t.id,
				t.test_date,
				t.test_date::date + t.pregnancy_time AS birth_forecast,
				t.pregnancy_status,
				t.observation,
				CASE
					WHEN c.id IS NOT NULL THEN 'SUCCESS'
					WHEN age(t.test_date) < INTERVAL '%[1]d days' THEN 'STAND_BY'
					ELSE 'FAILED'
				END AS birth_status,

				t.animal_id,
				a.name AS animal_name,
				a.tag AS animal_tag,
				COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0) animal_order,

				c.id AS calf_id,
				c.tag AS calf_tag,
				c.name AS calf_name,
				c.sex AS calf_sex,
				c.birth_date AS calf_birth_date,
				c.death_date AS calf_death_date,

				t.created_at
			FROM pregnancy_tests t 
				LEFT JOIN animals a ON a.id = t.animal_id
				LEFT JOIN animals c ON t.pregnancy_status = 'SUCCESS'
					AND c.mother_id = t.animal_id
					AND c.birth_date > t.test_date
					AND age(a.birth_date, t.test_date) <= INTERVAL '%[1]d days'
					AND NOT EXISTS (
						SELECT 1
						FROM pregnancy_tests t1
						WHERE t.pregnancy_status = 'FAILED'
							AND t1.animal_id = c.mother_id
							AND t1.test_date < c.birth_date
							AND t1.test_date BETWEEN t.test_date AND c.birth_date
					)
			WHERE t.user_id = $1 AND t.deleted_at IS NULL
		)

		SELECT * FROM cte
    `, util.MaxGestationDays)

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "cte", 2)
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
	sortExpression = " ORDER BY " + sortExpression 

	query += whereExpression + sortExpression + fmt.Sprintf(" LIMIT %d", limit)
	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)

	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return util.GetList[TestDB](r.DB, query, args...)
}

func (r *TestEntryRepository) GetEntriesFoot(filter *TestFilter, userId string) (*TestFoot, error) {

	countQuery := `
		WITH cte AS (
			SELECT 
				t.*,
				CASE
					WHEN pregnancy_status = 'FAILED' THEN 'FAILED'
					WHEN EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = t.animal_id
							AND a.birth_date > t.test_date
							AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
					) THEN 'SUCCESS'
					WHEN age(t.test_date) < INTERVAL '340 days' THEN 'STAND_BY'
					ELSE 'FAILED'
				END AS birth_status
			FROM pregnancy_tests t
			WHERE t.user_id = $1 AND t.deleted_at IS NULL
		)
		SELECT 
			COUNT(*) totals,
			COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success,
			COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') AS birth_success
		FROM cte t
    `

	filterExpression, _, err := util.GetFilterExpressions(filter, "t", 2)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExpression)
	countQuery += whereExpression

	query := fmt.Sprintf(`
        WITH count_query AS (%s)
        SELECT 
            totals,
            COALESCE(birth_success::float / NULLIF(totals, 0), 0) * 100 birth_rate,
            COALESCE(pregnancy_success::float / NULLIF(totals, 0), 0) * 100 pregnancy_rate
        FROM count_query
    `, countQuery)

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	return util.GetOne[TestFoot](r.DB, query, args...)
}

func (r *TestEntryRepository) FindGroups(userId string) (*[]TestGroups, error) {
	query := `
        WITH totals AS (
            SELECT 
                test_date,
                COUNT(*) animals_number,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS' 
					AND EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = t.animal_id
							AND a.birth_date > t.test_date
							AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
					)
				) AS birth_success
            FROM pregnancy_tests t
            WHERE deleted_at IS NULL AND user_id = $1 
            GROUP BY 1
        ),
        rates AS (
            SELECT
                g.test_date,
                g.animals_number,
                (g.pregnancy_success::float / g.animals_number::float)*100 pregnancy_rate,
                (g.birth_success::float / g.animals_number::float)*100 birth_rate
            FROM totals g
        )
        SELECT
            test_date,
            animals_number,
            pregnancy_rate,
            birth_rate,
            COALESCE((pregnancy_rate / LAG(pregnancy_rate) OVER win) - 1, 0) *100 pregnancy_comparison,
            COALESCE((birth_rate / LAG(birth_rate) OVER win) - 1, 0) * 100 birth_comparison
        FROM rates
		WINDOW win AS (ORDER BY test_date)
        ORDER BY test_date DESC
    `
	return util.GetList[TestGroups](r.DB, query, userId)
}

func (r *TestEntryRepository) FindEntriesByGroup(
	sort string,
	order string,
	testDate time.Time,
	userId string,
) (*[]TestDB, error) {

	sortMap := map[string]util.SortField{
		"animal_order":   {Field: "coalesce(regexp_replace(a.tag, '[^0-9]', '', 'g')::int, 0)", Order: "desc"},
		"birth_forecast": {Field: "coalesce(t.birth_forecast, '-infinity')", Order: "desc"},
		"animal_name":    {Field: "a.name", Order: "asc"},
	}

	query := fmt.Sprintf(`
		SELECT
			t.id,
			t.test_date,
			t.test_date::date + t.pregnancy_time AS birth_forecast,
			t.pregnancy_status,
			t.observation,
			CASE
				WHEN c.id IS NOT NULL THEN 'SUCCESS'
				WHEN age(t.test_date) < INTERVAL '%[1]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,

			t.animal_id,
			a.name AS animal_name,
			a.tag AS animal_tag,
			COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0) animal_order,

			c.id AS calf_id,
			c.tag AS calf_tag,
			c.name AS calf_name,
			c.sex AS calf_sex,
			c.birth_date AS calf_birth_date,
			c.death_date AS calf_death_date,

			t.created_at
		FROM pregnancy_tests t 
			LEFT JOIN animals a ON a.id = t.animal_id
			LEFT JOIN animals c ON t.pregnancy_status = 'SUCCESS'
				AND c.mother_id = t.animal_id
				AND c.birth_date > t.test_date
				AND age(a.birth_date, t.test_date) <= INTERVAL '%[1]d days'
				AND NOT EXISTS (
					SELECT 1
					FROM pregnancy_tests t1
					WHERE t.pregnancy_status = 'FAILED'
						AND t1.animal_id = c.mother_id
						AND t1.test_date < c.birth_date
						AND t1.test_date BETWEEN t.test_date AND c.birth_date
				)
		WHERE t.user_id = $1 
			AND t.test_date = $2
			AND t.deleted_at IS NULL
    `, util.MaxGestationDays)

	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query = query + " ORDER BY " + sortExpression
	return util.GetList[TestDB](r.DB, query, userId, testDate)
}

func (r *TestEntryRepository) GetEntriesByGroupFoot(testDate time.Time, userId string) (*TestFoot, error) {
	query := `
        WITH count_query AS (
            SELECT 
                COUNT(*) totals,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS' 
					AND EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = t.animal_id
							AND a.birth_date > t.test_date
							AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
					)
				) birth_success
            FROM pregnancy_tests t
            WHERE t.test_date = $1 AND t.user_id = $2 AND t.deleted_at IS NULL
        )
        SELECT 
            totals,
            COALESCE(birth_success::float / NULLIF(totals, 0), 0) * 100 birth_rate,
            COALESCE(pregnancy_success::float / NULLIF(totals, 0), 0) * 100 pregnancy_rate
        FROM count_query
    `
	return util.GetOne[TestFoot](r.DB, query, testDate, userId)
}

func (r *TestEntryRepository) CheckEntryExistence(entry *TestSave) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM pregnancy_tests
			WHERE test_date = :test_date
				AND animal_id = :animal_id
				AND user_id = :user_id
				AND id IS DISTINCT FROM :id
				AND deleted_at IS NULL
		)
	`
	var exists bool
	err := util.NamedPrimitive(r.DB, query, &exists, entry)
	return exists, err
}

func (r *TestEntryRepository) CheckGroupExistence(entry *TestGroupSave) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM pregnancy_tests
			WHERE test_date = :test_date
				AND user_id = :user_id
				AND deleted_at IS NULL
		)
	`
	var exists bool
	err := util.NamedPrimitive(r.DB, query, &exists, entry)
	return exists, err
}

func (r *TestEntryRepository) Add(entry *TestSave) *log.APIError {

	var query string
	if entry.Overwrite {
		query = `
			INSERT INTO pregnancy_tests (test_date, animal_id, pregnancy_status, pregnancy_time, observation, user_id)
			VALUES (:test_date, :animal_id, :pregnancy_status, :pregnancy_time, :observation, :user_id)
		`
	} else {
		query = `
			UPDATE pregnancy_tests 
			SET pregnancy_status = :pregnancy_status, 
				pregnancy_time = :pregnancy_time, 
				observation = :observation
			WHERE test_date = :test_date
				AND animal_id = :animal_id
				AND user_id = :user_id
		`
	}

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *TestEntryRepository) Update(entry *TestSave) (*TestDB, error) {

	query := `
		UPDATE pregnancy_tests 
		SET test_date = :test_date,
			pregnancy_status = :pregnancy_status, 
			birth_forecast = :pregnancy_time, 
			observation = :observation
		WHERE id = :id AND user_id = :user_id
	`

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return nil, err
	}

	selectQuery := fmt.Sprintf(`
		SELECT
			t.id,
			t.test_date,
			t.test_date::date + t.pregnancy_time AS birth_forecast,
			t.pregnancy_status,
			t.observation,
			CASE
				WHEN c.id IS NOT NULL THEN 'SUCCESS'
				WHEN age(t.test_date) < INTERVAL '%[1]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,

			t.animal_id,
			a.name AS animal_name,
			a.tag AS animal_tag,
			COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0) animal_order,

			c.id AS calf_id,
			c.tag AS calf_tag,
			c.name AS calf_name,
			c.sex AS calf_sex,
			c.birth_date AS calf_birth_date,
			c.death_date AS calf_death_date,

			t.created_at
		FROM pregnancy_tests t 
			LEFT JOIN animals a ON a.id = t.animal_id
			LEFT JOIN animals c ON t.pregnancy_status = 'SUCCESS'
				AND c.mother_id = t.animal_id
				AND c.birth_date > t.test_date
				AND age(a.birth_date, t.test_date) <= INTERVAL '%[1]d days'
				AND NOT EXISTS (
					SELECT 1
					FROM pregnancy_tests t1
					WHERE t.pregnancy_status = 'FAILED'
						AND t1.animal_id = c.mother_id
						AND t1.test_date < c.birth_date
						AND t1.test_date BETWEEN t.test_date AND c.birth_date
				)
				WHERE t.id = :id
					AND t.user_id = :user_id 
					AND t.deleted_at IS NULL
	`, util.MaxGestationDays)
	result, err := util.NamedGet(r.DB, selectQuery, TestDB{}, entry)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TestEntryRepository) Delete(id string, userId string) *log.APIError {
	query := `
		UPDATE pregnancy_tests
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	err := util.Exec(r.DB, query, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *TestEntryRepository) UpdateGroup(group *TestGroupSave) (*TestGroups, *log.APIError) {

	query := `
		UPDATE pregnancy_tests
		SET test_date = :test_date
		WHERE test_date = :old_test_date AND user_id = :user_id
	`
	err := util.NamedExec(r.DB, query, group)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	returnQuery := `
        WITH totals AS (
            SELECT 
                test_date,
                COUNT(*) animals_number,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS' 
					AND EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = t.animal_id
							AND a.birth_date > t.test_date
							AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
					)
				) AS birth_success
            FROM pregnancy_tests t
            WHERE deleted_at IS NULL 
				AND user_id = $1 
				AND test_date = $2
            GROUP BY 1
        ),
        rates AS (
            SELECT
                g.test_date,
                g.animals_number,
                (g.pregnancy_success::float / g.animals_number::float)*100 pregnancy_rate,
                (g.birth_success::float / g.animals_number::float)*100 birth_rate
            FROM totals g
        )
        SELECT
            test_date,
            animals_number,
            pregnancy_rate,
            birth_rate,
            COALESCE((pregnancy_rate / LAG(pregnancy_rate) OVER win) - 1, 0) * 100 pregnancy_comparison,
            COALESCE((birth_rate / LAG(birth_rate) OVER win) - 1, 0) * 100 birth_comparison
        FROM rates
		WINDOW win AS (ORDER BY test_date)
    `
	response, err := util.NamedGet(r.DB, returnQuery, TestGroups{}, group)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return response, nil
}

func (r *TestEntryRepository) DeleteGroup(testDate time.Time, userId string) *log.APIError {
	query := `
		UPDATE pregnancy_tests
		SET deleted_at = NOW()
		WHERE test_date = $1 AND user_id = $2
	`

	err := util.Exec(r.DB, query, testDate, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}
