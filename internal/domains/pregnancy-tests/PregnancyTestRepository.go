package pregnancyTests

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/entity"
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

func (r *TestEntryRepository) GetPregnancyRate(userId string) (*[]PregnancyHist, error) {
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
            test_date,
            (pregnancies::float / NULLIF(totals, 0)) * 100 pregnancy_rate
        FROM cte
        ORDER BY test_date
    `
	return util.GetList[PregnancyHist](r.DB, query, userId)
}

func (r *TestEntryRepository) GetAnimalsNumber(userId string) (*[]AnimalsNumberHist, error) {
	query := `
        WITH cte AS (
            SELECT 
                test_date,
                COUNT(*) AS totals
            FROM pregnancy_tests
            WHERE deleted_at IS NULL AND user_id = $1 
			GROUP BY 1
			ORDER BY test_date DESC
            LIMIT 10
        )
        SELECT test_date, totals
        FROM cte
        ORDER BY test_date
    `
	return util.GetList[AnimalsNumberHist](r.DB, query, userId)
}

func (r *TestEntryRepository) GetBirthRate(userId string) (*[]BirthHist, error) {
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
            test_date,
            (births::float / NULLIF(totals, 0)) * 100 AS birth_rate
        FROM cte
        ORDER BY test_date
    `
	return util.GetList[BirthHist](r.DB, query, userId)
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

func (r *TestEntryRepository) GetLastEntries(userId string) (*LastEntries, error) {
	dateQuery := `
		SELECT MAX(test_date)
		FROM pregnancy_tests 
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	var lastDate time.Time
	err := util.GetPrimitive(r.DB, dateQuery, &lastDate, userId)
	if err != nil {
		return nil, err
	}

	query := `
        SELECT
			t.id,
			t.animal_id,
            CONCAT_WS(' - ', a.ring_number, a.name) animal_info,
            t.test_date,
            test_date + (pregnancy_time*INTERVAL '1 day') AS birth_forecast,
            t.pregnancy_status,
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
			END AS birth_status,
            t.observation
        FROM pregnancy_tests t
            LEFT JOIN animals a ON a.id = t.animal_id
        WHERE t.user_id = $1 
			AND t.test_date = $2
			AND t.deleted_at IS NULL
        ORDER BY COALESCE(REGEXP_REPLACE(a.ring_number, '[^0-9]', '', 'g')::int, 0)
    `
	result, err := util.GetList[TestEntry](r.DB, query, userId, lastDate)
	if err != nil {
		return nil, err
	}

	lastEntry := &LastEntries{
		TestDate: lastDate,
		Entries:  *result,
	}

	return lastEntry, nil
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

func (r *TestEntryRepository) GetNextBirths(userId string) (*[]NextBirths, error) {
	query := `
        SELECT 
            DATE_TRUNC('month', t.test_date + (t.pregnancy_time * INTERVAL '1 day')) AS birth_forecast,
            COUNT(*) birth_numbers
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
	return util.GetList[NextBirths](r.DB, query, userId)
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
			CONCAT_WS(' - ', a.ring_number, a.name) AS animal_name,
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
			FROM general_rates
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
			CONCAT_WS(' - ', a.ring_number, a.name) AS animal_name,
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
	filter *TestEntryFilter,
	sort string,
	order string,
	cursor string,
	userId string,
) (*entity.Page[TestEntry], error) {

	sort = util.AddCommonFields(sort)
	sortMap := map[string]util.SortField{
		"animal_order":   {Field: "cte.animal_order", Order: "asc"},
		"test_date":      {Field: "cte.test_date", Order: "desc"},
		"birth_forecast": {Field: "coalesce(cte.birth_forecast, '-infinity')", Order: "desc"},
		"animal_name":    {Field: "cte.animal_name", Order: "asc"},
		"id":             {Field: "cte.id", Order: "asc"},
		"created_at":     {Field: "cte.created_at", Order: "asc"},
	}

	query := `
        WITH cte AS (
			SELECT
				t.id,
				t.test_date,
				t.animal_id,
				a.name AS animal_name,
				CONCAT_WS(' - ', a.ring_number, a.name) animal_info,
				COALESCE(REGEXP_REPLACE(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
				t.test_date + (t.pregnancy_time*INTERVAL '1 day') AS birth_forecast,
				t.pregnancy_status,
				CASE
					WHEN pregnancy_status = 'FAILED' THEN 'FAILED'
					WHEN child_name IS NOT NULL THEN 'SUCCESS'
					WHEN age(t.test_date) < INTERVAL '340 days' THEN 'STAND_BY'
					ELSE 'FAILED'
				END AS birth_status,
				CASE 
					WHEN pregnancy_status = 'FAILED' THEN 'Sem Cria'
					WHEN child_name IS NOT NULL THEN child_name
					ELSE 'Sem Cria'
				END AS child_information,
				t.observation,
				t.created_at
			FROM pregnancy_tests t 
				LEFT JOIN animals a ON a.id = t.animal_id
				LEFT JOIN LATERAL (
					SELECT CONCAT_WS(
						' - ',
						a.ring_number,
						COALESCE(a.name, a.sex),
						TO_CHAR(a.birth_date, 'DD/MM/YYYY')
					) AS child_name
					FROM animals a
					WHERE a.mother_id = t.animal_id
						AND a.birth_date > t.test_date
						AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
					LIMIT 1
				) c ON TRUE
			WHERE t.user_id = $1 AND t.deleted_at IS NULL
		)
		SELECT *
		FROM cte
    `
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

	query += whereExpression + sortExpression
	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)

	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return util.GetPage[TestEntry](r.DB, query, sort, 100, args...)
}

func (r *TestEntryRepository) GetEntriesFoot(filter *TestEntryFilter, userId string) (*TestEntryFoot, error) {

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
	return util.GetOne[TestEntryFoot](r.DB, query, args...)
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
) (*[]TestEntry, error) {

	sortMap := map[string]util.SortField{
		"animal_order":   {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "desc"},
		"birth_forecast": {Field: "coalesce(t.birth_forecast, '-infinity')", Order: "desc"},
		"animal_name":    {Field: "a.name", Order: "asc"},
	}

	query := `
        SELECT
            t.id,
            t.test_date,
            t.animal_id,
            CONCAT_WS(' - ', a.ring_number, a.name) AS animal_info,
            COALESCE(REGEXP_REPLACE(a.ring_number, '[^0-9]', '', 'g')::int, 0) AS animal_order,
            t.test_date + (t.pregnancy_time * INTERVAL '1 day') AS birth_forecast,
            t.pregnancy_status,
			CASE
				WHEN pregnancy_status = 'FAILED' THEN 'FAILED'
				WHEN child_name IS NOT NULL THEN 'SUCCESS'
				WHEN age(t.test_date) < INTERVAL '340 days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,
			CASE 
				WHEN pregnancy_status = 'FAILED' THEN 'Sem Cria'
				WHEN child_name IS NOT NULL THEN child_name
				ELSE 'Sem Cria'
			END AS child_information,
            t.observation
        FROM pregnancy_tests t 
			LEFT JOIN animals a ON a.id = t.animal_id
			LEFT JOIN LATERAL (
				SELECT CONCAT_WS(
					' - ',
					a.ring_number,
					COALESCE(a.name, a.sex),
					TO_CHAR(a.birth_date, 'DD/MM/YYYY')
				) AS child_name
				FROM animals a
				WHERE a.mother_id = t.animal_id
					AND a.birth_date > t.test_date
					AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
				LIMIT 1
			) c ON TRUE
		WHERE t.user_id = $1 AND t.test_date = $2 AND t.deleted_at IS NULL
    `
	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query = query + " ORDER BY " + sortExpression
	return util.GetList[TestEntry](r.DB, query, userId, testDate)
}

func (r *TestEntryRepository) GetEntriesByGroupFoot(testDate time.Time, userId string) (*TestEntryFoot, error) {
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
	return util.GetOne[TestEntryFoot](r.DB, query, testDate, userId)
}

func (r *TestEntryRepository) CheckEntryExistence(entry *TestEntrySave) (bool, error) {
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

func (r *TestEntryRepository) CheckGroupExistence(entry *TestGroups) (bool, error) {
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

func (r *TestEntryRepository) Add(entry *TestEntrySave) *log.APIError {

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

func (r *TestEntryRepository) Update(entry *TestEntrySave) (*TestEntry, *log.APIError) {

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
		return nil, log.InternalServerAPIError(err)
	}

	selectQuery := `
		SELECT
			t.id,
			t.test_date,
			t.animal_id,
			CONCAT_WS(' - ', a.ring_number, a.name) animal_info,
			t.test_date + (INTERVAL '1 day' * t.pregnancy_time) AS t.birth_forecast,
			t.pregnancy_status,
			CASE
				WHEN pregnancy_status = 'FAILED' THEN 'FAILED'
				WHEN child_name IS NOT NULL THEN 'SUCCESS'
				WHEN age(t.test_date) < INTERVAL '340 days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,
			CASE 
				WHEN pregnancy_status = 'FAILED' THEN 'Sem Cria'
				WHEN child_name IS NOT NULL THEN child_name
				ELSE 'Sem Cria'
			END AS child_information,
			t.observation
		FROM pregnancy_tests t 
			LEFT JOIN animals a ON a.id = t.animal_id
			LEFT JOIN LATERAL (
				SELECT CONCAT_WS(
					' - ',
					a.ring_number,
					COALESCE(a.name, a.sex),
					TO_CHAR(a.birth_date, 'DD/MM/YYYY')
				) AS child_name
				FROM animals a
				WHERE a.mother_id = t.animal_id
					AND a.birth_date > t.test_date
					AND age(a.birth_date, t.test_date) <= INTERVAL '340 days'
				LIMIT 1
			) c ON TRUE
		WHERE t.id = $1 AND t.user_id = $2
	`
	result, err := util.GetOne[TestEntry](r.DB, selectQuery, entry.Id, entry.UserId)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
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

func (r *TestEntryRepository) UpdateBatch(group *TestGroups) (*TestGroups, *log.APIError) {

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

func (r *TestEntryRepository) DeleteBatch(testDate time.Time, userId string) *log.APIError {
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
