package embryoTransfer

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/entity"
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type TransferRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *TransferRepository {
	return &TransferRepository{db}
}

func (r *TransferRepository) GetBirthRateStats(userId string) (*CardEntry, error) {
	query := fmt.Sprintf(`
        WITH totals AS (
            SELECT 
                transfer_date,
                COUNT(*) AS total,
                COUNT(*) FILTER (WHERE EXISTS (
					SELECT 1 
					FROM animals a
					WHERE a.mother_id = et.donor_id
					  AND a.birth_date > et.transfer_date
					  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				)) AS birth_success
            FROM embryo_transfer et
			WHERE et.user_id = $1 AND et.deleted_at IS NULL
			GROUP BY 1
            ORDER BY 1 DESC
            LIMIT 10
        )
        SELECT 
            transfer_date,
			COALESCE(birth_success::float / NULLIF(total, 0), 0) * 100 AS birth_rate
        FROM totals
		ORDER BY 1
    `, util.MinGestantionDays, util.MaxGestationDays)
	result, err := util.GetList[BirthRateEntry](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	birthRates := *result
	var currentRate, previousRate, trend float64

	switch lenght := len(birthRates); lenght {
	case 0:
		currentRate = 0
		previousRate = 0
		trend = 0
	case 1:
		currentRate = birthRates[lenght-1].BirthRate
		previousRate = 0
		trend = 0
	default:
		currentRate = birthRates[lenght-1].BirthRate
		previousRate = birthRates[lenght-2].BirthRate
		trend = util.CalculatePercentageTrend(currentRate, previousRate)
	}

	stats := &CardEntry{
		Hist:    birthRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *TransferRepository) GetPregnancyRateStats(userId string) (*CardEntry, error) {
	query := fmt.Sprintf(`
		WITH insemination_status AS (
			SELECT
				et.transfer_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status
			FROM embryo_transfer et
			WHERE et.user_id = $1 AND et.deleted_at IS NULL 
		),
		cte AS (
			SELECT
				t.transfer_date,
				COUNT(t.*) AS total,
				COUNT(t.*) FILTER (WHERE t.pregnancy_status = 'SUCCESS') AS pregnancy_success
			FROM insemination_status t
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT 
			transfer_date,
			COALESCE(pregnancy_success::float / NULLIF(total, 0), 0) * 100 AS pregnancy_rate
		FROM cte
		ORDER BY 1
    `, util.MinGestantionDays, util.MaxGestationDays)

	result, err := util.GetList[PregnancyRateEntry](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	pregnancyRates := *result
	var currentRate, previousRate, trend float64

	switch lenght := len(pregnancyRates); lenght {
	case 0:
		currentRate = 0
		previousRate = 0
		trend = 0
	case 1:
		currentRate = pregnancyRates[lenght-1].PregnancyRate
		previousRate = 0
		trend = 0
	default:
		currentRate = pregnancyRates[lenght-1].PregnancyRate
		previousRate = pregnancyRates[lenght-2].PregnancyRate
		trend = util.CalculatePercentageTrend(currentRate, previousRate)
	}

	stats := &CardEntry{
		Hist:    pregnancyRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *TransferRepository) GetAnimalsNumber(userId string) (*CardEntry, error) {
	query := `
		WITH cte AS (
			SELECT
				transfer_date,
				COUNT(*) AS animals_number
			FROM embryo_transfer
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT *
		FROM cte
		ORDER BY 1
    `

	result, err := util.GetList[AnimalsNumberEntry](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	pregnancyRates := *result
	var currentRate, previousRate, trend float64

	switch lenght := len(pregnancyRates); lenght {
	case 0:
		currentRate = 0
		previousRate = 0
		trend = 0
	case 1:
		currentRate = pregnancyRates[lenght-1].AnimalsNumber
		previousRate = 0
		trend = 0
	default:
		currentRate = pregnancyRates[lenght-1].AnimalsNumber
		previousRate = pregnancyRates[lenght-2].AnimalsNumber
		trend = util.CalculatePercentageTrend(currentRate, previousRate)
	}

	stats := &CardEntry{
		Hist:    pregnancyRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *TransferRepository) GetTransferHist(userId string) (*[]TransferHist, error) {
	query := fmt.Sprintf(`
        WITH cte AS (
			SELECT 
				et.transfer_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM embryo_transfer et
			WHERE et.user_id = $1 AND et.deleted_at IS NULL
		),
        totals AS (
            SELECT
                transfer_date,
                COUNT(*) animals_number,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') births_number,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancies_number
            FROM cte
            GROUP BY 1
            ORDER BY 1 DESC
            LIMIT 30
        )
        SELECT * FROM totals ORDER BY transfer_date
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[TransferHist](r.DB, query, userId)
}

func (r *TransferRepository) GetFutureBirths(userId string) (*[]FutureBirths, error) {
	query := fmt.Sprintf(`
		WITH upcoming_births AS (
			SELECT (t.test_date + %[3]d - t.pregnancy_time * INTERVAL '1 day') AS birth_forecast
			FROM embryo_transfer et
			JOIN pregnancy_tests t
				ON t.animal_id = et.receiver_id
				AND t.test_date > et.transfer_date
				AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
				AND t.pregnancy_status = 'SUCCESS'
			WHERE et.user_id = $1
			  AND et.deleted_at IS NULL
			  AND NOT EXISTS (
				  SELECT 1
				  FROM animals a
				  WHERE a.mother_id = et.donor_id
					AND a.birth_date > et.transfer_date
					AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					AND NOT EXISTS (
						SELECT 1
						FROM pregnancy_tests t
						WHERE t.animal_id = a.mother_id
							AND t.test_date BETWEEN et.transfer_date AND a.birth_date
							AND t.pregnancy_status = 'FAILED'
					)
			  )
			  AND t.test_date + (%[3]d - t.pregnancy_time * INTERVAL '1 day') >= NOW()  
		)
		SELECT
			DATE_TRUNC('month', birth_forecast) AS birth_forecast,
			COUNT(*) AS births_number
		FROM upcoming_births
		GROUP BY 1
		ORDER BY 1;
	`, util.MinGestantionDays, util.MaxGestationDays, util.AverageGestationDays)
	return util.GetList[FutureBirths](r.DB, query, userId)
}

func (r *TransferRepository) GetBestBull(userId string) (*[]BestAnimals, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				CONCAT_WS(' - ', a.ring_number, a.name) animal_name,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM embryo_transfer et
                LEFT JOIN animals a ON et.bull_id = a.id 
			WHERE et.user_id = $1 AND et.deleted_at IS NULL
		),
        totals AS (
            SELECT
                s.animal_name,
                COUNT(s.*) total,
                COUNT(s.*) FILTER (WHERE s.birth_status = 'SUCCESS') birth_success,
                COUNT(s.*) FILTER (WHERE s.pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status s
            GROUP BY 1
        ),
        rates AS (
            SELECT 
                animal_name,
                total,
                COALESCE(birth_success::float / NULLIF(total, 0), 0) * 100 birth_rate,
                COALESCE(pregnancy_success::float / NULLIF(total, 0), 0) * 100 pregnancy_rate
            FROM totals
        )
        SELECT
			animal_name,
			total,
			birth_rate,
			pregnancy_rate,
			COALESCE(birth_rate / NULLIF(AVG(birth_rate) OVER (), 0) - 1, 0) * 100 AS birth_comparison_rate,
			COALESCE(pregnancy_rate / NULLIF(AVG(pregnancy_rate) OVER (), 0) - 1, 0) * 100 AS pregnancy_comparison_rate
		FROM rates
		ORDER BY birth_rate DESC;
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[BestAnimals](r.DB, query, userId)
}

func (r *TransferRepository) GetBestDonors(userId string) (*[]BestAnimals, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				CONCAT_WS(' - ', a.ring_number, a.name) animal_name,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM embryo_transfer et
                LEFT JOIN animals a ON a.id = et.donor_id
			WHERE et.user_id = $1 AND et.deleted_at IS NULL
		),
        totals AS (
            SELECT
                s.animal_name,
                COUNT(s.*) total,
                COUNT(s.*) FILTER (WHERE s.birth_status = 'SUCCESS') birth_success,
                COUNT(s.*) FILTER (WHERE s.pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status s
            GROUP BY 1
        ),
        rates AS (
            SELECT 
                animal_name,
                total,
                COALESCE(birth_success::float / NULLIF(total, 0), 0) * 100 birth_rate,
                COALESCE(pregnancy_success::float / NULLIF(total, 0), 0) * 100 pregnancy_rate
            FROM totals
        )
        SELECT
			animal_name,
			total,
			birth_rate,
			pregnancy_rate,
			COALESCE(birth_rate / NULLIF(AVG(birth_rate) OVER (), 0) - 1, 0) * 100 AS birth_comparison_rate,
			COALESCE(pregnancy_rate / NULLIF(AVG(pregnancy_rate) OVER (), 0) - 1, 0) * 100 AS pregnancy_comparison_rate
		FROM rates
		ORDER BY birth_rate DESC;
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[BestAnimals](r.DB, query, userId)
}

func (r *TransferRepository) GetBestReceivers(userId string) (*[]BestAnimals, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				CONCAT_WS(' - ', a.ring_number, a.name) animal_name,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM embryo_transfer et
                LEFT JOIN animals a ON a.id = et.receiver_id
			WHERE et.user_id = $1 AND et.deleted_at IS NULL
		),
        totals AS (
            SELECT
                s.animal_name,
                COUNT(s.*) total,
                COUNT(s.*) FILTER (WHERE s.birth_status = 'SUCCESS') birth_success,
                COUNT(s.*) FILTER (WHERE s.pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status s
            GROUP BY 1
        ),
        rates AS (
            SELECT 
                animal_name,
                total,
                COALESCE(birth_success::float / NULLIF(total, 0), 0) * 100 birth_rate,
                COALESCE(pregnancy_success::float / NULLIF(total, 0), 0) * 100 pregnancy_rate
            FROM totals
        )
        SELECT
			animal_name,
			total,
			birth_rate,
			pregnancy_rate,
			COALESCE(birth_rate / NULLIF(AVG(birth_rate) OVER (), 0) - 1, 0) * 100 AS birth_comparison_rate,
			COALESCE(pregnancy_rate / NULLIF(AVG(pregnancy_rate) OVER (), 0) - 1, 0) * 100 AS pregnancy_comparison_rate
		FROM rates
		ORDER BY birth_rate DESC;
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[BestAnimals](r.DB, query, userId)
}

func (r *TransferRepository) GetLastGroups(userId string) (*[]TransferGroup, error) {
	query := fmt.Sprintf(`
		WITH insemination_data AS (
			SELECT
				et.transfer_date,
				CASE
					WHEN EXISTS (
						SELECT 1 
						FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1 	
						FROM animals a
						WHERE a.mother_id = et.donor_id
						  AND a.birth_date > et.transfer_date
						  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM embryo_transfer et
			WHERE et.user_id = $1 AND et.deleted_at IS NULL
		),
		daily_stats AS (
			SELECT
				transfer_date,
				COUNT(*) AS cow_number,
				COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') AS birth_success,
				COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') AS pregnancy_success
			FROM insemination_data
			GROUP BY transfer_date
		),
		rates AS (
			SELECT
				transfer_date,
				cow_number,
				COALESCE(birth_success::float * 100 / NULLIF(cow_number, 0), 0) AS birth_rate,
				COALESCE(pregnancy_success::float * 100 / NULLIF(cow_number, 0), 0) AS pregnancy_rate
			FROM daily_stats
		)
		SELECT
			transfer_date,
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
		WINDOW win AS (ORDER BY transfer_date)
		ORDER BY transfer_date DESC
		LIMIT 5;
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[TransferGroup](r.DB, query, userId)
}

func (r *TransferRepository) GetLastEntries(userId string) (*LastEntry, error) {

	lastDateQuery := `
		SELECT MAX(transfer_date) max_date
		FROM embryo_transfer 
		WHERE deleted_at IS NULL AND user_id = $1
	`

	var lastDate time.Time
	err := util.GetPrimitive(r.DB, lastDateQuery, &lastDate, userId)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT 
			et.id,
			et.transfer_date,
			et.bull_id,
			CONCAT_WS(' - ', r.ring_number, r.name) AS receiver_info,
			CONCAT_WS(' - ', d.ring_number, d.name) AS donor_info,
			b.name AS bull_name,
			CASE
				WHEN EXISTS (
					SELECT 1 
					FROM animals a
					WHERE a.mother_id = et.donor_id
					  AND a.birth_date > et.transfer_date
					  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				) THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = et.receiver_id
					  AND t.test_date > et.transfer_date
					  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
					  AND t.pregnancy_status = 'SUCCESS'
				) THEN 'SUCCESS'
				WHEN NOT EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = et.receiver_id
					  AND t.test_date > et.transfer_date
					  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
				) AND age(et.transfer_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS pregnancy_status,
			CASE
				WHEN EXISTS (
					SELECT 1 
					FROM animals a
					WHERE a.mother_id = et.donor_id
					  AND a.birth_date > et.transfer_date
					  AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				) THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = et.receiver_id
					  AND t.test_date > et.transfer_date
					  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
					  AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN age(et.transfer_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status
		FROM embryo_transfer et
			LEFT JOIN animals r ON r.id = et.receiver_id
			LEFT JOIN animals d ON d.id = et.donor_id
			LEFT JOIN animals b ON b.id = et.bull_id
		WHERE et.user_id = $1 
			AND et.transfer_date = $2
			AND et.deleted_at IS NULL
		ORDER BY COALESCE(REGEXP_REPLACE(r.ring_number, '[^0-9]', '', 'g')::int, 0);
    `, util.MinGestantionDays, util.MaxGestationDays)
	result, err := util.GetList[EmbryoTransfer](r.DB, query, userId, lastDate)
	if err != nil {
		return nil, err
	}

	lastEntry := &LastEntry{
		TransferDate: lastDate,
		Entries:      *result,
	}

	return lastEntry, nil
}

func (r *TransferRepository) FindEntriesPage(
	userId string,
	filter *TransferEntryFilter,
	sort string,
	order string,
	cursor string,
) (*entity.Page[EmbryoTransfer], error) {

	sort = util.AddCommonFields(sort)
	sortMap := map[string]util.SortField{
		"receiver_order": {Field: "et.receiver_order", Order: "asc"},
		"donor_order":    {Field: "et.donor_order", Order: "asc"},
		"receiver_name":  {Field: "et.receiver_name", Order: "asc"},
		"donor_name":     {Field: "et.donor_name", Order: "asc"},
		"transfer_date":  {Field: "coalesce(et.transfer_date, '-infinity')", Order: "asc"},
		"id":             {Field: "et.id", Order: "asc"},
		"created_at":     {Field: "et.created_at", Order: "asc"},
	}

	query := fmt.Sprintf(`
        WITH cte AS (
			SELECT 
				et.id,
				et.receiver_id,
				et.donor_id,
				COALESCE(REGEXP_REPLACE(r.ring_number, '[^0-9]', '', 'g')::int, 0) AS receiver_order,
				COALESCE(REGEXP_REPLACE(d.ring_number, '[^0-9]', '', 'g')::int, 0) AS donor_order,
				r.name AS receiver_name,
				d.name AS donor_name,
				CONCAT_WS(' - ', r.ring_number, r.name) AS receiver_info,
				CONCAT_WS(' - ', d.ring_number, d.name) AS donor_info,
				et.transfer_date,
				et.bull_id,
				b.name AS bull_name,
				CASE
					WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					WHEN NOT EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
					) AND age(et.transfer_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'FAILED'
					) THEN 'FAILED'
					WHEN age(et.transfer_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
					ELSE 'FAILED'
				END AS birth_status,
				CASE 
					WHEN c.child_name IS NULL THEN 'Sem Cria'
					ELSE c.child_name
				END AS child_information,
				et.observation,
				et.created_at
			FROM embryo_transfer et
				LEFT JOIN animals r ON r.id = et.receiver_id
				LEFT JOIN animals d ON d.id = et.donor_id
				LEFT JOIN animals b ON b.id = et.bull_id
				LEFT JOIN LATERAL (
					SELECT
					CONCAT_WS(
						' - ', 
						a.ring_number, 
						COALESCE(a.name, a.sex), 
						TO_CHAR(a.birth_date, 'DD/MM/YYYY')
					) AS child_name
					FROM animals a
					WHERE a.mother_id = et.donor_id
						AND a.birth_date > et.transfer_date
						AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					ORDER BY a.birth_date
					LIMIT 1
				) c ON TRUE
			WHERE et.user_id = $1 AND et.deleted_at IS NULL
		)
		SELECT * FROM cte et
	`, util.MinGestantionDays, util.MaxGestationDays)
	orderExpression := " ORDER BY "

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "et", 2)
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
	return util.GetPage[EmbryoTransfer](r.DB, query, sort, 100, args...)
}

func (r *TransferRepository) GetEntriesFoot(
	userId string,
	filter *TransferEntryFilter,
) (*TransferFoot, error) {

	statusQuery := fmt.Sprintf(`
		WITH cte AS  (
			SELECT
				et.animal_id,
				et.bull_id,
				et.transfer_date,
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = et.donor_id
							AND a.birth_date > et.transfer_date
							AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = et.donor_id
							AND a.birth_date > et.transfer_date
							AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM embryo_transfer et
			WHERE et.user_id = $1 AND et.deleted_at IS NULL
		)
		SELECT pregnancy_status, birth_status
		FROM cte et
	`, util.MinGestantionDays, util.MaxGestationDays)

	filterExpression, _, err := util.GetFilterExpressions(filter, "et", 2)
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

	return util.GetOne[TransferFoot](r.DB, query, args...)
}

func (r *TransferRepository) FindEntriesByGroup(userId string, date time.Time) (*[]EmbryoTransfer, error) {

	query := fmt.Sprintf(`
        SELECT 
            et.id,
			et.bull_id,
			et.receiver_id,
			et.donor_id,
			CONCAT_WS(' - ', b.ring_number, b.name) AS bull_name,
            CONCAT_WS(' - ', r.ring_number, r.name) receiver_info,
            CONCAT_WS(' - ', d.ring_number, d.name) donor_info,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = et.receiver_id
					  AND t.test_date > et.transfer_date
					  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
					  AND t.pregnancy_status = 'SUCCESS'
				) THEN 'SUCCESS'
				WHEN NOT EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = et.receiver_id
					  AND t.test_date > et.transfer_date
					  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
				) AND age(et.transfer_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS pregnancy_status,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = et.receiver_id
					  AND t.test_date > et.transfer_date
					  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
					  AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN age(et.transfer_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,
			CASE
				WHEN c.child_name IS NULL THEN 'Sem Cria'
				ELSE child_name
			END AS child_information,
            et.observation
        FROM embryo_transfer et
            LEFT JOIN animals r ON r.id = et.receiver_id
            LEFT JOIN animals d ON d.id = et.donor_id
            LEFT JOIN animals b ON b.id = et.bull_id
			LEFT JOIN LATERAL (
				SELECT
				CONCAT_WS(
					' - ', 
					a.ring_number, 
					COALESCE(a.name, a.sex), 
					TO_CHAR(a.birth_date, 'DD/MM/YYYY')
				) AS child_name
				FROM animals a
				WHERE a.mother_id = et.donor_id
					AND  a.birth_date > et.transfer_date
					AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				ORDER BY a.birth_date
				LIMIT 1
			) c ON TRUE
		WHERE et.user_id = $1 AND et.deleted_at IS NULL AND et.transfer_date = $2
        ORDER BY COALESCE(REGEXP_REPLACE(r.ring_number, '[^0-9]', '', 'g')::int, 0)
	`, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[EmbryoTransfer](r.DB, query, userId, date)
}

func (r *TransferRepository) GetEntriesByGroupFoot(userId string, date time.Time) (*TransferFoot, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = et.donor_id
							AND a.birth_date > et.transfer_date
							AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = et.donor_id
							AND a.birth_date > et.transfer_date
							AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM embryo_transfer et
			WHERE et.user_id = $1 
				AND et.transfer_date = $2
				AND et.deleted_at IS NULL
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
            COALESCE(birth_success::float / NULLIF(totals, 0), 0) * 100 average_birth_rate,
            COALESCE(pregnancy_success::float / NULLIF(totals, 0), 0) * 100 average_pregnancy_rate
        FROM COUNTING
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetOne[TransferFoot](r.DB, query, userId, date)
}

func (r *TransferRepository) FindGroups(userId string) (*[]TransferGroup, error) {
	query := fmt.Sprintf(`
		WITH status AS (
			SELECT
				et.transfer_date,
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = et.donor_id
							AND a.birth_date > et.transfer_date
							AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = et.donor_id
							AND a.birth_date > et.transfer_date
							AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM embryo_transfer et
			WHERE et.user_id = $1 AND et.deleted_at IS NULL
		),
        totals AS (
            SELECT 
                transfer_date,
				COUNT(*) cow_number,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') birth_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status s
            GROUP BY transfer_date
        ),
        rates AS (
            SELECT
                transfer_date,
                cow_number,
                COALESCE(birth_success::float / NULLIF(cow_number, 0), 0) * 100 birth_rate,
                COALESCE(pregnancy_success::float / NULLIF(cow_number, 0), 0) * 100 pregnancy_rate
            FROM totals
        )
        SELECT 
            s.transfer_date,
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
		WINDOW win AS (ORDER BY s.transfer_date)
        ORDER BY s.transfer_date DESC
    `, util.MinGestantionDays, util.MaxGestationDays)
	return util.GetList[TransferGroup](r.DB, query, userId)
}

func (r *TransferRepository) AddTransfer(entry *EmbryoTransferSave) *log.APIError {

	validateErr := validateAdd(r.DB, entry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		INSERT INTO embryo_transfer (
			receiver_id, 
			donor_id, 
			bull_id, 
			transfer_date, 
			observation, 
			user_id
		)
		VALUES (
			:receiver_id, 
			:donor_id, 
			:bull_id, 
			:transfer_date, 
			:observation, 
			:user_id
		)
    `

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *TransferRepository) Replace(entry *EmbryoTransferSave) *log.APIError {

	query := `
		UPDATE embryo_transfer 
		SET donor_id = :donor_id, 
			bull_id = :bull_id, 
			observation = :observation, 
		WHERE receiver_id = :receiver_id 
			AND transfer_date = :transfer_date
			AND user_id = :user_id
	`

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *TransferRepository) Delete(id string, userId string) *log.APIError {

	oldQuery := `
		SELECT
			id,
			donor_id,
			receiver_id,
			bull_id,
			transfer_date,
			observation,
			user_id
		FROM embryo_transfer
		WHERE id = $1 AND user_id = $2
	`

	oldEntry, err := util.GetOne[EmbryoTransferSave](r.DB, oldQuery, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	validateErr := validateDelete(r.DB, oldEntry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		UPDATE embryo_transfer
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	err = util.Exec(r.DB, query, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *TransferRepository) Update(entry *EmbryoTransferSave) (*EmbryoTransfer, *log.APIError) {

	validateErr := validateUpdate(r.DB, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	query := `
		UPDATE embryo_transfer
		SET donor_id = :donor_id,
			receiver_id = :receiver_id,
			bull_id = :bull_id,
			transfer_date = :transfer_date,
			observation = :observation
		WHERE id = :id AND user_id = :user_id
	`

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	selectQuery := fmt.Sprintf(`
		SELECT 
			et.id,
			et.receiver_id,
			et.donor_id,
			CONCAT_WS(' - ', r.ring_number, r.name) AS receiver_info,
			CONCAT_WS(' - ', d.ring_number, d.name) AS donor_info,
			et.transfer_date,
			et.bull_id,
			b.name AS bull_name,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = et.receiver_id
					  AND t.test_date > et.transfer_date
					  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
					  AND t.pregnancy_status = 'SUCCESS'
				) THEN 'SUCCESS'
				WHEN NOT EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = et.receiver_id
					  AND t.test_date > et.transfer_date
					  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
				) AND age(et.transfer_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS pregnancy_status,
			CASE
				WHEN c.child_name IS NOT NULL THEN 'SUCCESS'
				WHEN EXISTS (
					SELECT 1 
					FROM pregnancy_tests t
					WHERE t.animal_id = et.receiver_id
					  AND t.test_date > et.transfer_date
					  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
					  AND t.pregnancy_status = 'FAILED'
				) THEN 'FAILED'
				WHEN age(et.transfer_date) < INTERVAL '%[2]d days' THEN 'STAND_BY'
				ELSE 'FAILED'
			END AS birth_status,
			CASE 
				WHEN c.child_name IS NULL THEN 'Sem Cria'
				ELSE c.child_name
			END AS child_information,
			et.observation
		FROM embryo_transfer et
			LEFT JOIN animals r ON r.id = et.receiver_id
			LEFT JOIN animals d ON d.id = et.donor_id
			LEFT JOIN animals b ON b.id = et.bull_id
			LEFT JOIN LATERAL (
				SELECT
				CONCAT_WS(
					' - ', 
					a.ring_number, 
					COALESCE(a.name, a.sex), 
					TO_CHAR(a.birth_date, 'DD/MM/YYYY')
				) AS child_name
				FROM animals a
				WHERE a.mother_id = et.donor_id
					AND a.birth_date > et.transfer_date
					AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
				ORDER BY a.birth_date
				LIMIT 1
			) c ON TRUE
		WHERE et.id = $1 AND et.user_id = $2
	`, util.MinGestantionDays, util.MaxGestationDays)

	response, err := util.GetOne[EmbryoTransfer](r.DB, selectQuery, entry.Id, entry.UserId)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return response, nil

}

func (r *TransferRepository) UpdateGroup(entry *TransferGroupSave) (*TransferGroup, *log.APIError) {

	validateErr := validateUpdateGroups(r.DB, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	query := `
		UPDATE embryo_transfer
		SET transfer_date = :transfer_date
		WHERE transfer_date = :old_transfer_date 
			AND user_id = :user_id 
			AND deleted_at IS NULL
	`

	err = util.NamedExecTx(tx, query, entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	returnQuery := fmt.Sprintf(`
		WITH status AS (
			SELECT
				et.transfer_date,
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = et.donor_id
							AND a.birth_date > et.transfer_date
							AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					WHEN EXISTS (
						SELECT 1 
						FROM pregnancy_tests t
						WHERE t.animal_id = et.receiver_id
						  AND t.test_date > et.transfer_date
						  AND age(t.test_date, et.transfer_date) <= INTERVAL '%[2]d days'
						  AND t.pregnancy_status = 'SUCCESS'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS pregnancy_status,
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM animals a
						WHERE a.mother_id = et.donor_id
							AND a.birth_date > et.transfer_date
							AND age(a.birth_date, et.transfer_date) BETWEEN INTERVAL '%[1]d days' AND INTERVAL '%[2]d days'
					) THEN 'SUCCESS'
					ELSE 'FAILED'
				END AS birth_status
			FROM embryo_transfer et
			WHERE et.user_id = :user_id
				AND et.transfer_date = :transfer_date
				AND et.deleted_at IS NULL
		),
        totals AS (
            SELECT 
                transfer_date,
				COUNT(*) cow_number,
                COUNT(*) FILTER (WHERE birth_status = 'SUCCESS') birth_success,
                COUNT(*) FILTER (WHERE pregnancy_status = 'SUCCESS') pregnancy_success
            FROM status s
            GROUP BY transfer_date
        ),
        rates AS (
            SELECT
                transfer_date,
                cow_number,
                COALESCE(birth_success::float / NULLIF(cow_number, 0), 0) * 100 birth_rate,
                COALESCE(pregnancy_success::float / NULLIF(cow_number, 0), 0) * 100 pregnancy_rate
            FROM totals
        )
        SELECT 
            s.transfer_date,
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
		WINDOW win AS (ORDER BY s.transfer_date)
    `, util.MinGestantionDays, util.MaxGestationDays)

	response, err := util.NamedGetTx(tx, returnQuery, TransferGroup{}, entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return response, nil
}

func (r *TransferRepository) DeleteGroup(transferDate time.Time, userId string) *log.APIError {

	query := `
		DELETE FROM embryo_transfer
		WHERE transfer_date = $1 AND user_id = $2
	`

	err := util.Exec(r.DB, query, transferDate, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}
