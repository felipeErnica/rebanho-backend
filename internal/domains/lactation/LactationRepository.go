package lactation

import (
	"fmt"

	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type LactationRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *LactationRepository {
	return &LactationRepository{db}
}

type SaveValidation struct {
	LactationExists  bool `db:"lac_exists"`
	InvalidNew       bool `db:"invalid_new"`
	InvalidStart     bool `db:"invalid_start"`
	InvalidEnd       bool `db:"invalid_end"`
	InvalidCalf      bool `db:"invalid_calf"`
	InvalidEmptyEnd  bool `db:"invalid_empty_end"`
	DifferentPasture bool `db:"different_pasture"`
}

func (r *LactationRepository) CheckLactationConflicts(lac LactationHistSave) (*SaveValidation, error) {
	query := `
		SELECT 
			EXISTS (
				SELECT 1
				FROM lactations 
				WHERE animal_id = :animal_id
					AND id IS DISTINCT FROM :id
					AND start_date = :start_date
					AND user_id = :user_id
					AND deleted_at IS NULL
			) AS lac_exists,
			EXISTS (
				SELECT 1 
				FROM lactations
				WHERE :end_date IS NULL
					AND id IS DISTINCT FROM :id
					AND animal_id = :animal_id
					AND end_date IS NULL
					AND  user_id = :user_id
					AND deleted_at IS NULL
			) AS invalid_new,
			EXISTS (
				SELECT 1 
				FROM lactations l
				WHERE l.animal_id = :animal_id
					AND l.id IS DISTINCT FROM :id
					AND l.start_date < :start_date
					AND l.end_date >= :start_date
					AND l.deleted_at IS NULL
					AND l.user_id = :user_id
			) AS invalid_start,
			EXISTS (
				SELECT 1 
				FROM lactations l
				WHERE :end_date IS NULL
					AND l.animal_id = :animal_id
					AND l.id IS DISTINCT FROM :id
					AND l.start_date > :start_date
					AND l.user_id = :user_id
					AND l.deleted_at IS NULL
			) AS invalid_empty_end,
			EXISTS (
				SELECT 1 
				FROM lactations l
				WHERE l.animal_id = :animal_id
					AND l.id IS DISTINCT FROM :id
					AND l.start_date > :start_date
					AND l.start_date <= :end_date
					AND l.deleted_at IS NULL
					AND l.user_id = :user_id
			) AS invalid_end,
			EXISTS (
				SELECT 1
				FROM lactations 	
				WHERE id IS DISTINCT FROM :id
					AND calf_id = :calf_id
					AND user_id = :user_id
					AND deleted_at IS NULL
			) AS invalid_calf,
			(
				SELECT COALESCE(pasture_id <> :pasture_id, FALSE)
				FROM pasture_entries
				WHERE animal_id = :animal_id
					AND exit_date IS NULL
					AND user_id = :user_id
					AND deleted_at IS NULL
			) AS different_pasture
	`
	return util.NamedGet(r.DB, query, SaveValidation{}, lac)
}

func (r *LactationRepository) GetLastLactatingEntries(userId string) (*[]util.GraphData, error) {
	query := `
        WITH cte AS (
            SELECT
                entry_date AS date,
                COUNT(*) value
            FROM milk_entries
            WHERE user_id = $1 AND deleted_at IS NULL
            GROUP BY 1
            ORDER BY 1 DESC
            LIMIT 10
        )
        SELECT * FROM cte ORDER BY date
    `

	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *LactationRepository) GetLastDryEntries(userId string) (*[]util.GraphData, error) {
	query := `
		WITH cte AS (
            SELECT
                entry_date AS date,
                COUNT(*) AS value
            FROM milk_entries m
				JOIN lactations l ON l.animal_id = m.animal_id
					AND m.entry_date = l.end_date
					AND l.deleted_at IS NULL
            WHERE m.user_id = $1 AND m.deleted_at IS NULL
            GROUP BY 1
            ORDER BY 1 DESC
            LIMIT 10
        )
        SELECT * FROM cte ORDER BY date
    `

	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *LactationRepository) GetDairyTypes(userId string) (*DairyTypes, error) {

	query := `
		WITH cte AS (
			SELECT
				a.id,
				EXISTS (
					SELECT 1
					FROM lactations l
					WHERE l.user_id = $1
						AND l.deleted_at IS NULL
						AND l.animal_id = a.id
						AND l.end_date IS NULL
				) AS is_lactating
			FROM animals a
			WHERE a.user_id = $1 
				AND a.deleted_at IS NULL
				AND a.animal_type = 'DAIRY_ANIMAL'
				AND a.death_date IS NULL
		)
		SELECT
			COUNT(*) FILTER (WHERE is_lactating = FALSE) AS dry,
			COUNT(*) FILTER (WHERE is_lactating = TRUE) AS lactating
		FROM cte
	`

	result, err := util.GetOne[DairyTypes](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *LactationRepository) GetBestAnimals(userId string) (*[]AnimalsRating, error) {
	query := `
        WITH lac_prod AS (
            SELECT
                l.id,
                AVG(m.quantity) avg_prod
            FROM lactations l
                JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),

        lac_tbl AS (
            SELECT
                l.animal_id,
                s.avg_prod,
				l.end_date,
                EXTRACT(days FROM l.end_date - l.start_date) + 1 period,
                (EXTRACT(days FROM l.end_date - l.start_date) + 1)*s.avg_prod total,
                EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) lac_interval
            FROM lactations l
                JOIN lac_prod s USING (id)
            WHERE EXISTS (
				SELECT 1
				FROM animals a
				WHERE a.id = l.animal_id AND a.death_date IS NULL
			) 
        ),
        lac_stats AS (
            SELECT
                CONCAT_WS(' - ', a.tag, a.name) animal_name,
                COUNT(l.*) lac_num,
                AVG(l.period) FILTER (WHERE l.end_date IS NOT NULL) avg_period,
                AVG(l.avg_prod) FILTER (WHERE l.end_date IS NOT NULL) avg_prod,
                AVG(l.total) FILTER (WHERE l.end_date IS NOT NULL) avg_total,
                AVG(l.lac_interval) FILTER (WHERE l.lac_interval IS NOT NULL) avg_interval
            FROM lac_tbl l 
				JOIN animals a ON a.id = l.animal_id
            GROUP BY 1
			HAVING COUNT(l.*) >= 3
        ),
        gn_stats AS (
            SELECT
                AVG(avg_period) gn_avg_period,
                AVG(avg_prod) gn_avg_prod,
                AVG(avg_total) gn_avg_total,
                AVG(avg_interval) gn_avg_interval,
				STDDEV(avg_total) dev_total,
				STDDEV(avg_interval) dev_interval
            FROM lac_stats
        ),
		cte AS (
			SELECT
				l.*,
				((l.avg_period / NULLIF(s.gn_avg_period, 0)) - 1) AS period_rate,
				((l.avg_prod / NULLIF(s.gn_avg_prod, 0)) - 1) AS prod_rate,
				((l.avg_total / NULLIF(s.gn_avg_total, 0)) - 1) AS total_rate,
				((l.avg_interval / NULLIF(s.gn_avg_interval, 0)) - 1) AS interval_rate,
				(l.avg_total - gn_avg_total ) / NULLIF(s.dev_total, 0) AS total_score,
				(l.avg_interval - gn_avg_interval) / NULLIF(s.dev_interval, 0) AS interval_score
			FROM lac_stats l
				CROSS JOIN gn_stats s
		)
		SELECT cte.*
		FROM cte
		WHERE (total_score * 0.6 - interval_score * 0.4) > 0
		ORDER BY (total_score * 0.6 - interval_score * 0.4) DESC
		LIMIT 10
    `
	return util.GetList[AnimalsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetWorstAnimals(userId string) (*[]AnimalsRating, error) {
	query := `
        WITH lac_prod AS (
            SELECT
                l.id,
                AVG(m.quantity) avg_prod
            FROM lactations l
                JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),
        lac_tbl AS (
            SELECT
                l.animal_id,
                s.avg_prod,
				l.end_date,
                EXTRACT(days FROM l.end_date - l.start_date) + 1 period,
                (EXTRACT(days FROM l.end_date - l.start_date) + 1)*s.avg_prod total,
                EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) lac_interval
            FROM lactations l
                JOIN lac_prod s USING (id)
            WHERE EXISTS (
				SELECT 1
				FROM animals a
				WHERE a.id = l.animal_id AND a.death_date IS NULL
			)
        ),
        lac_stats AS (
            SELECT
                CONCAT_WS(' - ', a.tag, a.name) animal_name,
                COUNT(l.*) lac_num,
                AVG(l.period) FILTER (WHERE l.end_date IS NOT NULL) avg_period,
                AVG(l.avg_prod) FILTER (WHERE l.end_date IS NOT NULL)avg_prod,
                AVG(l.total) FILTER (WHERE l.end_date IS NOT NULL) avg_total,
                AVG(l.lac_interval) FILTER (WHERE l.lac_interval IS NOT NULL) avg_interval
            FROM lac_tbl l 
				JOIN animals a ON a.id = l.animal_id
            GROUP BY 1
			HAVING COUNT(l.*) >= 3
        ),
        gn_stats AS (
            SELECT
                AVG(avg_period) gn_avg_period,
                AVG(avg_prod) gn_avg_prod,
                AVG(avg_total) gn_avg_total,
                AVG(avg_interval) gn_avg_interval,
				STDDEV(avg_total) dev_total,
				STDDEV(avg_interval) dev_interval
            FROM lac_stats
        ),
		cte AS (
			SELECT
				l.*,
				((l.avg_period / NULLIF(s.gn_avg_period, 0)) - 1) AS period_rate,
				((l.avg_prod / NULLIF(s.gn_avg_prod, 0)) - 1) AS prod_rate,
				((l.avg_total / NULLIF(s.gn_avg_total, 0)) - 1) AS total_rate,
				((l.avg_interval / NULLIF(s.gn_avg_interval, 0)) - 1) AS interval_rate,
				(l.avg_total - gn_avg_total ) / NULLIF(s.dev_total, 0) AS total_score,
				(l.avg_interval - gn_avg_interval) / NULLIF(s.dev_interval, 0) AS interval_score
			FROM lac_stats l
				CROSS JOIN gn_stats s
		)
		SELECT cte.*
		FROM cte
		WHERE (-total_score * 0.6 + interval_score * 0.4) > 0
		ORDER BY (-total_score * 0.6 + interval_score * 0.4) DESC
		LIMIT 10
    `
	return util.GetList[AnimalsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetBestMothers(userId string) (*[]ParentsRating, error) {
	query := `
        WITH lac_stats AS (
            SELECT
                l.id,
                AVG(m.quantity) avg_prod
            FROM lactations l
                JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),
        lac_tbl AS (
            SELECT
                l.animal_id,
                s.avg_prod,
				l.end_date,
                EXTRACT(days FROM l.end_date - l.start_date) + 1 period,
                (EXTRACT(days FROM l.end_date - l.start_date) + 1)*s.avg_prod total,
                EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) lac_interval
            FROM lactations l
                JOIN lac_stats s USING (id)
        ),
        cte_animals AS (
            SELECT
               	l.animal_id,
                COUNT(l.*) lac_num,
                AVG(l.period) FILTER (WHERE l.end_date IS NOT NULL) avg_period,
                AVG(l.avg_prod) FILTER (WHERE l.end_date IS NOT NULL) avg_prod,
                AVG(l.total) FILTER (WHERE l.end_date IS NOT NULL) avg_total,
                AVG(l.lac_interval) FILTER (WHERE l.lac_interval IS NOT NULL) avg_interval
            FROM lac_tbl l 
            GROUP BY 1
        ),
        mother_stats AS (
            SELECT
                CONCAT_WS(' - ', f.tag, f.name) parent_name,
				COUNT(cte.*) children_number,
                AVG(avg_period) avg_period,
                AVG(avg_prod) avg_prod,
                AVG(avg_total) avg_total,
                AVG(avg_interval) avg_interval
            FROM cte_animals cte
				JOIN animals a ON a.id = cte.animal_id
				JOIN animals f ON 
					f.id = a.mother_id
					AND f.death_date IS NULL
            GROUP BY 1
			HAVING COUNT(cte.*) >= 3
        ),
        gn_stats AS (
            SELECT
                AVG(avg_period) avg_period,
                AVG(avg_prod) avg_prod,
                AVG(avg_total) avg_total,
                AVG(avg_interval) avg_interval,
                STDDEV(avg_interval) dev_interval,
                STDDEV(avg_total) dev_total
            FROM cte_animals
        ),
		cte AS (
			SELECT
				m.*,
				((m.avg_period / NULLIF(s.avg_period, 0)) - 1) period_rate,
				((m.avg_prod / NULLIF(s.avg_prod, 0)) - 1) prod_rate,
				((m.avg_total / NULLIF(s.avg_total, 0)) - 1) total_rate,
				((m.avg_interval / NULLIF(s.avg_interval, 0)) - 1) interval_rate,
				(m.avg_total - s.avg_total) / NULLIF(s.dev_total, 0) AS total_score,
				(m.avg_interval - s.avg_interval) / NULLIF(s.dev_interval, 0) AS interval_score
			FROM mother_stats m
				CROSS JOIN gn_stats s
		)
		SELECT *
		FROM cte
		WHERE (total_score * 0.6 - interval_score * 0.4) > 0
		ORDER BY (total_score * 0.6 - interval_score * 0.4) DESC
		LIMIT 10
    `
	return util.GetList[ParentsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetWorstMothers(userId string) (*[]ParentsRating, error) {
	query := `
        WITH lac_stats AS (
            SELECT
                l.id,
                AVG(m.quantity) avg_prod
            FROM lactations l
                JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),
        lac_tbl AS (
            SELECT
                l.animal_id,
                s.avg_prod,
				l.end_date,
                EXTRACT(days FROM l.end_date - l.start_date) + 1 period,
                (EXTRACT(days FROM l.end_date - l.start_date) + 1)*s.avg_prod total,
                EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) lac_interval
            FROM lactations l
                JOIN lac_stats s USING (id)
        ),
        cte_animals AS (
            SELECT
               	l.animal_id,
                COUNT(l.*) lac_num,
                AVG(l.period) FILTER (WHERE l.end_date IS NOT NULL) avg_period,
                AVG(l.avg_prod) FILTER (WHERE l.end_date IS NOT NULL) avg_prod,
                AVG(l.total) FILTER (WHERE l.end_date IS NOT NULL) avg_total,
                AVG(l.lac_interval) FILTER (WHERE l.lac_interval IS NOT NULL) avg_interval
            FROM lac_tbl l 
            GROUP BY 1
        ),
        mother_stats AS (
            SELECT
                CONCAT_WS(' - ', f.tag, f.name) parent_name,
				COUNT(cte.*) children_number,
                AVG(avg_period) avg_period,
                AVG(avg_prod) avg_prod,
                AVG(avg_total) avg_total,
                AVG(avg_interval) avg_interval
            FROM cte_animals cte
				JOIN animals a ON a.id = cte.animal_id
				JOIN animals f ON 
					f.id = a.mother_id
					AND f.death_date IS NULL
            GROUP BY 1
			HAVING COUNT(cte.*) >= 3
        ),
        gn_stats AS (
            SELECT
                AVG(avg_period) avg_period,
                AVG(avg_prod) avg_prod,
                AVG(avg_total) avg_total,
                AVG(avg_interval) avg_interval,
                STDDEV(avg_interval) dev_interval,
                STDDEV(avg_total) dev_total
            FROM cte_animals
        ),
		cte AS (
			SELECT
				m.*,
				((m.avg_period / NULLIF(s.avg_period, 0)) - 1) period_rate,
				((m.avg_prod / NULLIF(s.avg_prod, 0)) - 1) prod_rate,
				((m.avg_total / NULLIF(s.avg_total, 0)) - 1) total_rate,
				((m.avg_interval / NULLIF(s.avg_interval, 0)) - 1) interval_rate,
				(m.avg_total - s.avg_total) / NULLIF(s.dev_total, 0) AS total_score,
				(m.avg_interval - s.avg_interval) / NULLIF(s.dev_interval, 0) AS interval_score
			FROM mother_stats m
				CROSS JOIN gn_stats s
		)
		SELECT *
		FROM cte
		WHERE (-total_score * 0.6 + interval_score * 0.4) > 0
		ORDER BY (-total_score * 0.6 + interval_score * 0.4) DESC
		LIMIT 10
    `
	return util.GetList[ParentsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetBestFathers(userId string) (*[]ParentsRating, error) {
	query := `
        WITH lac_stats AS (
            SELECT
                l.id,
                AVG(m.quantity) avg_prod
            FROM lactations l
                JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),
        lac_tbl AS (
            SELECT
                l.animal_id,
                s.avg_prod,
				l.end_date,
                EXTRACT(days FROM l.end_date - l.start_date) + 1 period,
                (EXTRACT(days FROM l.end_date - l.start_date) + 1)*s.avg_prod total,
                EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) lac_interval
            FROM lactations l
                JOIN lac_stats s USING (id)
        ),
        cte_animals AS (
            SELECT
               	l.animal_id,
                COUNT(l.*) lac_num,
                AVG(l.period) FILTER (WHERE l.end_date IS NOT NULL) avg_period,
                AVG(l.avg_prod) FILTER (WHERE l.end_date IS NOT NULL) avg_prod,
                AVG(l.total) FILTER (WHERE l.end_date IS NOT NULL) avg_total,
                AVG(l.lac_interval) FILTER (WHERE l.lac_interval IS NOT NULL) avg_interval
            FROM lac_tbl l 
            GROUP BY 1
        ),
        mother_stats AS (
            SELECT
                CONCAT_WS(' - ', f.tag, f.name) parent_name,
				COUNT(cte.*) children_number,
                AVG(avg_period) avg_period,
                AVG(avg_prod) avg_prod,
                AVG(avg_total) avg_total,
                AVG(avg_interval) avg_interval
            FROM cte_animals cte
				JOIN animals a ON a.id = cte.animal_id
				JOIN animals f ON 
					f.id = a.father_id
					AND f.death_date IS NULL
            GROUP BY 1
			HAVING COUNT(cte.*) >= 5
        ),
        gn_stats AS (
            SELECT
                AVG(avg_period) avg_period,
                AVG(avg_prod) avg_prod,
                AVG(avg_total) avg_total,
                AVG(avg_interval) avg_interval,
                STDDEV(avg_interval) dev_interval,
                STDDEV(avg_total) dev_total
            FROM cte_animals
        ),
		cte AS (
			SELECT
				m.*,
				((m.avg_period / NULLIF(s.avg_period, 0)) - 1)  period_rate,
				((m.avg_prod / NULLIF(s.avg_prod, 0)) - 1)  prod_rate,
				((m.avg_total / NULLIF(s.avg_total, 0)) - 1)  total_rate,
				((m.avg_interval / NULLIF(s.avg_interval, 0)) - 1)  interval_rate,
				(m.avg_total - s.avg_total) / NULLIF(s.dev_total, 0) AS total_score,
				(m.avg_interval - s.avg_interval) / NULLIF(s.dev_interval, 0) AS interval_score
			FROM mother_stats m
				CROSS JOIN gn_stats s
		)
		SELECT *
		FROM cte
		WHERE (total_score * 0.6 - interval_score * 0.4) > 0
		ORDER BY (total_score * 0.6 - interval_score * 0.4) DESC
		LIMIT 10
    `
	return util.GetList[ParentsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetWorstFathers(userId string) (*[]ParentsRating, error) {
	query := `
        WITH lac_stats AS (
            SELECT
                l.id,
                AVG(m.quantity) avg_prod
            FROM lactations l
                JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),
        lac_tbl AS (
            SELECT
                l.animal_id,
                s.avg_prod,
				l.end_date,
                EXTRACT(days FROM l.end_date - l.start_date) + 1 period,
                (EXTRACT(days FROM l.end_date - l.start_date) + 1) * s.avg_prod total,
                EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) lac_interval
            FROM lactations l
                JOIN lac_stats s USING (id)
        ),
        cte_animals AS (
            SELECT
               	l.animal_id,
                COUNT(l.*) lac_num,
                AVG(l.period) FILTER (WHERE l.end_date IS NOT NULL) avg_period,
                AVG(l.avg_prod) FILTER (WHERE l.end_date IS NOT NULL) avg_prod,
                AVG(l.total) FILTER (WHERE l.end_date IS NOT NULL) avg_total,
                AVG(l.lac_interval) FILTER (WHERE l.lac_interval IS NOT NULL) avg_interval
            FROM lac_tbl l 
            GROUP BY 1
        ),
        mother_stats AS (
            SELECT
                CONCAT_WS(' - ', f.tag, f.name) parent_name,
				COUNT(cte.*) children_number,
                AVG(avg_period) avg_period,
                AVG(avg_prod) avg_prod,
                AVG(avg_total) avg_total,
                AVG(avg_interval) avg_interval
            FROM cte_animals cte
				JOIN animals a ON a.id = cte.animal_id
				JOIN animals f ON 
					f.id = a.father_id
					AND f.death_date IS NULL
            GROUP BY 1
			HAVING COUNT(cte.*) >= 5
        ),
        gn_stats AS (
            SELECT
                AVG(avg_period) avg_period,
                AVG(avg_prod) avg_prod,
                AVG(avg_total) avg_total,
                AVG(avg_interval) avg_interval,
                STDDEV(avg_interval) dev_interval,
                STDDEV(avg_total) dev_total
            FROM cte_animals
        ),
		cte AS (
			SELECT
				m.*,
				((m.avg_period / NULLIF(s.avg_period, 0)) - 1)  period_rate,
				((m.avg_prod / NULLIF(s.avg_prod, 0)) - 1)  prod_rate,
				((m.avg_total / NULLIF(s.avg_total, 0)) - 1)  total_rate,
				((m.avg_interval / NULLIF(s.avg_interval, 0)) - 1)  interval_rate,
				(m.avg_total - s.avg_total) / NULLIF(s.dev_total, 0) AS total_score,
				(m.avg_interval - s.avg_interval) / NULLIF(s.dev_interval, 0) AS interval_score
			FROM mother_stats m
				CROSS JOIN gn_stats s
			WHERE m.avg_interval > 0 
		)
		SELECT *
		FROM cte
		WHERE (-total_score * 0.6 + interval_score * 0.4) > 0
		ORDER BY (-total_score * 0.6 + interval_score * 0.4) DESC
		LIMIT 10
    `
	return util.GetList[ParentsRating](r.DB, query, userId)
}

func (r *LactationRepository) FindLactationPage(
	filter *LactationHistFilter,
	sort string,
	order string,
	cursor string,
	limit int,
	userId string,
) (*[]LactationDB, error) {

	sortMap := map[string]util.SortField{
		"animal_order":     {Field: "cte.animal_order", Order: "asc"},
		"animal_name":      {Field: "cte.animal_name", Order: "asc"},
		"start_date":       {Field: "cte.start_date", Order: "asc"},
		"end_date":         {Field: "coalesce(cte.end_date, '-infinity')", Order: "asc"},
		"calf_birth_date":  {Field: "coalesce(cte.calf_birth_date, -infinity)", Order: "asc"},
		"avg_production":   {Field: "coalesce(cte.avg_production, 0)", Order: "asc"},
		"lac_period":       {Field: "cte.lac_period", Order: "asc"},
		"total_production": {Field: "coalesce(cte.total_production, 0)", Order: "asc"},
		"lac_interval":     {Field: "coalesce(cte.lac_interval, 0)", Order: "asc"},
		"id":               {Field: "cte.id", Order: "asc"},
		"created_at":       {Field: "cte.created_at", Order: "asc"},
	}

	query := `
        WITH lac_stats AS (
            SELECT
                l.id,
                AVG(COALESCE(m.quantity, 0)) avg_prod,
				MAX(entry_date) max_date,
				MAX(COALESCE(m.quantity, 0)) peak
            FROM lactations l
                LEFT JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),

		cte AS (
			SELECT
				l.id,
				l.animal_id,
				a.tag AS animal_tag,
				a.name AS animal_name,
				COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0) AS animal_order,

				l.calf_id,
				c.tag AS calf_tag,
				c.name AS calf_name,
				c.sex AS calf_sex,
				c.birth_date AS calf_birth_date,
				c.death_date AS calf_death_date,

				l.start_date,
				l.end_date,
				s.avg_prod AS avg_production,
				COALESCE(EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1, 0) AS lac_period,
				COALESCE(EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1, 0) * s.avg_prod AS total_production,
				EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) AS lac_interval,
				s.peak,
				l.observation,
				l.created_at
			FROM lactations l
				JOIN lac_stats s USING (id)
				JOIN animals a ON a.id = l.animal_id
				LEFT JOIN animals c ON c.id = l.calf_id
				LEFT JOIN animals cm ON cm.id = c.mother_id
			WHERE l.user_id = $1 AND l.deleted_at IS NULL
		)
		SELECT * FROM cte
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

	orderBy := " ORDER BY " + sortExpression
	query += whereExpression + orderBy + fmt.Sprintf(" LIMIT %d", limit)
	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	return util.GetList[LactationDB](r.DB, query, args...)
}

func (r *LactationRepository) GetLactationPageFoot(filter *LactationHistFilter, userId string) (*LactationHistFoot, error) {

	lacQuery := "SELECT * FROM cte"

	filterExpression, _, err := util.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		lacQuery += " WHERE " + filterExpression
	}

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	mainQuery := fmt.Sprintf(`
        WITH lac_stats AS (
            SELECT
                l.id,
				MAX(m.entry_date) max_date,
                AVG(m.quantity) avg_prod,
				MAX(m.quantity) peak
            FROM lactations l
                JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),
		cte AS (
			SELECT
				l.animal_id,
				l.calf_id,
				c.birth_date calf_birth_date,
				l.start_date,
				l.end_date,
				s.avg_prod avg_production,
				EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1 lac_period,
				(EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1)*s.avg_prod total_production,
				EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) lac_interval,
				s.peak,
				l.observation
			FROM lactations l
				JOIN lac_stats s USING (id)
				LEFT JOIN animals c ON c.id = l.calf_id
			WHERE l.user_id = $1 AND l.deleted_at IS NULL
		),
		lac AS (%s)
		SELECT 
			COUNT(*) AS total_lacs,
			AVG(lac_interval) FILTER (WHERE lac_interval IS NOT NULL) avg_lac_interval,
			AVG(lac_period) FILTER (WHERE end_date IS NOT NULL) avg_lac_period,
			AVG(total_production) FILTER (WHERE end_date IS NOT NULL) avg_total_production,
			AVG(avg_production) FILTER (WHERE end_date IS NOT NULL) avg_production,
			AVG(peak) FILTER (WHERE end_date IS NOT NULL) avg_peak
		FROM lac
	`, lacQuery)

	return util.GetOne[LactationHistFoot](r.DB, mainQuery, args...)
}

func (r *LactationRepository) FindAnimalsPage(
	filter *AnimalFilter,
	sort string,
	order string,
	cursor string,
	limit int,
	userId string,
) (*[]AnimalDB, error) {
	sortMap := map[string]util.SortField{
		"tag_order":       {Field: "cte.tag_order", Order: "asc"},
		"name":            {Field: "cte.name", Order: "asc"},
		"lac_start":       {Field: "cte.lac_start", Order: "asc"},
		"lac_end":         {Field: "coalesce(cte.lac_end, '-infinity')", Order: "asc"},
		"calf_birth_date": {Field: "coalesce(cte.calf_birth_date, -infinity)", Order: "asc"},
		"lac_average":     {Field: "coalesce(cte.lac_average, 0)", Order: "asc"},
		"lac_period":      {Field: "cte.lac_period", Order: "asc"},
		"lac_total":       {Field: "coalesce(cte.lac_total, 0)", Order: "asc"},
		"lac_interval":    {Field: "coalesce(cte.lac_interval, 0)", Order: "asc"},
		"id":              {Field: "cte.id", Order: "asc"},
		"created_at":      {Field: "cte.created_at", Order: "asc"},
	}

	query := `
        WITH lac_stats AS (
            SELECT
                l.id,
				AVG(COALESCE(m.quantity, 0)) lac_average,
				MAX(entry_date) max_date,
				MAX(COALESCE(m.quantity, 0)) lac_peak
            FROM lactations l
                LEFT JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),

		lac_animals AS (
			SELECT
				a.id,
				a.name,
				a.tag,
				a.birth_date,
				EXISTS (
					SELECT 1
					FROM lactations l
					WHERE l.animal_id = a.id
						AND l.end_date IS NULL
						AND l.user_id = $1
						AND l.deleted_at IS NULL
				) AS is_lactating,
				a.created_at
			FROM animals a
			WHERE a.animal_type = 'DAIRY_ANIMAL'
				AND a.death_date IS NULL
				AND a.user_id = $1
				AND a.deleted_at IS NULL
		),

		cte AS (
			SELECT DISTINCT ON (a.id)
				a.id,
				a.tag,
				a.name,
				a.birth_date,
				COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0) AS tag_order,
				a.is_lactating,
				
				l.calf_id,
				c.name AS calf_name,
				c.tag AS calf_tag,
				c.sex AS calf_sex,
				c.birth_date AS calf_birth_date,
				c.death_date AS calf_death_date,

				l.id AS lac_id,
				l.start_date AS lac_start,
				l.end_date AS lac_end,
				s.lac_average,
				EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1 AS lac_period,
				(EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1) * s.lac_average AS lac_total,
				EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) AS lac_interval,
				s.lac_peak,
				l.observation AS lac_observation,

				a.created_at
			FROM lac_animals a
				LEFT JOIN lactations l ON a.id = l.animal_id
				LEFT JOIN lac_stats s ON l.id = s.id
				LEFT JOIN animals c ON c.id = l.calf_id
			ORDER BY a.id, l.start_date DESC
		)
		SELECT * FROM cte
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

	orderBy := " ORDER BY " + sortExpression + fmt.Sprintf(" LIMIT %d", limit)
	query += whereExpression + orderBy
	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	return util.GetList[AnimalDB](r.DB, query, args...)
}

func (r *LactationRepository) GetAnimalsPageFoot(filter *AnimalFilter, userId string) (*LactationHistFoot, error) {

	lacQuery := "SELECT * FROM cte"

	filterExpression, _, err := util.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		lacQuery += " WHERE " + filterExpression
	}

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	mainQuery := fmt.Sprintf(`
        WITH lac_stats AS (
            SELECT
                l.id,
				MAX(m.entry_date) max_date,
                AVG(m.quantity) avg_prod,
				MAX(m.quantity) peak
            FROM lactations l
                JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),
		lac_animals AS (
			SELECT
				a.id,
				EXISTS(
					SELECT 1
					FROM lactations l
					WHERE l.animal_id = a.id
						AND l.end_date IS NULL
						AND l.user_id = $1
						AND l.deleted_at IS NULL
				) AS is_lactating
			FROM animals a
			WHERE a.animal_type = 'DAIRY_ANIMAL'
				AND a.death_date IS NULL
				AND a.user_id = $1
				AND a.deleted_at IS NULL
		),
		cte AS (
			SELECT DISTINCT ON (a.id)
				a.is_lactating,
				l.animal_id,
				l.calf_id,
				c.birth_date AS calf_birth_date,
				l.start_date,
				l.end_date,
				s.avg_prod avg_production,
				EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1 AS lac_period,
				(EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1) * s.avg_prod AS total_production,
				EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) AS lac_interval,
				s.peak,
				l.observation
			FROM lac_animals a
				LEFT JOIN lactations l ON a.id = l.animal_id
				LEFT JOIN lac_stats s ON l.id = s.id
				LEFT JOIN animals c ON c.id = l.calf_id
			ORDER BY a.id, l.start_date DESC
		),
		lac AS (%s)
		SELECT 
			COUNT(*) AS total_lacs,
			AVG(lac_interval) FILTER (WHERE lac_interval IS NOT NULL) AS avg_lac_interval,
			AVG(lac_period) FILTER (WHERE end_date IS NOT NULL) AS avg_lac_period,
			AVG(total_production) FILTER (WHERE end_date IS NOT NULL) AS avg_total_production,
			AVG(avg_production) FILTER (WHERE end_date IS NOT NULL) AS avg_production,
			AVG(peak) FILTER (WHERE end_date IS NOT NULL) AS avg_peak
		FROM lac
	`, lacQuery)

	return util.GetOne[LactationHistFoot](r.DB, mainQuery, args...)
}
func (r *LactationRepository) FindById(id string, userId string) (*LactationDB, error) {

	query := `
        WITH lac_stats AS (
            SELECT
                l.id,
                AVG(COALESCE(m.quantity, 0)) avg_prod,
				MAX(entry_date) max_date,
				MAX(COALESCE(m.quantity, 0)) peak
            FROM lactations l
                LEFT JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
            GROUP BY 1
        ),

		cte AS (
			SELECT
				l.id,

				l.animal_id,
				a.tag AS animal_tag,
				a.name AS animal_name,
				COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0) AS animal_order,

				l.calf_id,
				c.tag AS calf_tag,
				c.name AS calf_name,
				c.sex AS calf_sex,
				c.birth_date AS calf_birth_date,
				c.death_date AS calf_death_date,

				l.start_date,
				l.end_date,
				s.avg_prod avg_production,
				COALESCE(EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1, 0) lac_period,
				COALESCE(EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1, 0) * s.avg_prod total_production,
				EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) AS lac_interval,
				s.peak,
				l.observation,
				l.created_at
			FROM lactations l
				JOIN lac_stats s USING (id)
				JOIN animals a ON a.id = l.animal_id
				LEFT JOIN animals c ON c.id = l.calf_id
				LEFT JOIN animals cm ON cm.id = c.mother_id
			WHERE l.id = $1
				AND l.user_id = $2
				AND l.deleted_at IS NULL
		)
		SELECT * FROM cte
    `
	return util.GetOne[LactationDB](r.DB, query, id, userId)
}

func (r *LactationRepository) AddLactation(entry *LactationHistSave) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if entry.TransferPasture {
		updateQuery := `
			UPDATE pasture_entries
			SET exit_date = :start_date
			WHERE animal_id = :animal_id
				AND exit_date IS NULL
				AND user_id = :user_id
				AND deleted_at IS NULL
		`
		err = util.NamedExecTx(tx, updateQuery, entry)
		if err != nil {
			return err
		}

		entryQuery := `
			INSERT INTO pasture_entries (animal_id, entry_date, pasture_id, user_id)
			VALUES (:animal_id, :start_date, :pasture_id, :user_id)
		`
		err = util.NamedExecTx(tx, entryQuery, entry)
		if err != nil {
			return err
		}
	}

	if entry.Overwrite {
		updateQuery := `
			UPDATE lactations
			SET end_date = :end_date,
				calf_id = :calf_id,
				observation = :observation
			WHERE animal_id = :animal_id
				AND start_date = :start_date
				AND user_id = :user_id
				AND deleted_at IS NULL
		`
		err = util.NamedExecTx(tx, updateQuery, entry)
		if err != nil {
			return err
		}
	} else {
		insertQuery := `
			INSERT INTO lactations (animal_id, calf_id, start_date, end_date, observation, user_id)
			VALUES (:animal_id, :calf_id, :start_date, :end_date, :observation, :user_id)
		`

		err = util.NamedExecTx(tx, insertQuery, entry)
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

func (r *LactationRepository) UpdateLactation(entry *LactationHistSave) (*LactationDB, error) {

	insertQuery := `
		UPDATE lactations
		SET calf_id = :calf_id,
			start_date = :start_date,
			end_date = :end_date,
			observation = :observation
		WHERE id = :id AND user_id = :user_id
	`

	err := util.NamedExec(r.DB, insertQuery, entry)
	if err != nil {
		return nil, err
	}

	selectQuery := `
		WITH lac_stats AS (
            SELECT
                AVG(COALESCE(m.quantity, 0)) AS avg_prod,
				MAX(entry_date) AS max_date,
				MAX(COALESCE(m.quantity, 0)) AS peak
            FROM lactations l
                LEFT JOIN milk_entries m ON 
                    l.animal_id = m.animal_id
                    AND l.start_date <= m.entry_date
                    AND COALESCE(l.end_date, NOW()) >= m.entry_date
                    AND m.deleted_at IS NULL
                    AND m.user_id = $1
			WHERE l.id = $1
        )

		SELECT
			l.id,

			l.animal_id,
			a.tag AS animal_tag,
			a.name AS animal_name,

			l.calf_id,
			c.name AS calf_name,
			c.tag AS calf_tag,
			c.sex AS calf_sex,
			c.birth_date calf_birth_date,
			c.death_date calf_death_date,

			l.start_date,
			l.end_date,
			s.avg_prod AS avg_production,
			COALESCE(EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1, 0) AS lac_period,
			COALESCE(EXTRACT(days FROM COALESCE(l.end_date, s.max_date) - l.start_date) + 1, 0) * s.avg_prod AS total_production,
			EXTRACT(days FROM l.start_date - LAG(l.end_date) OVER (PARTITION BY l.animal_id ORDER BY l.start_date)) AS lac_interval,
			s.peak,
			l.observation
		FROM lactations l
			CROSS JOIN lac_stats s
			JOIN animals a ON a.id = l.animal_id
			LEFT JOIN animals c ON c.id = l.calf_id
			LEFT JOIN animals cm ON cm.id = c.mother_id
		WHERE l.id = $1
	`

	response, err := util.GetOne[LactationDB](r.DB, selectQuery, entry.Id)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *LactationRepository) DeleteLactation(id string, userId string) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `
		UPDATE lactations
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2
	`
	err = util.ExecTx(tx, query, id, userId)
	if err != nil {
		return err
	}

	entriesQuery := `
		UPDATE milk_entries m
		SET deleted_at = NOW()
		FROM lactations l
		WHERE l.id = $1
			AND m.animal_id = l.animal_id
			AND m.entry_date BETWEEN l.start_date AND COALESCE(l.end_date, NOW());
			AND m.user_id = $2
	`
	err = util.ExecTx(tx, entriesQuery, id, userId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
