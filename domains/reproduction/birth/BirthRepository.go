package birth

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type BirthRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *BirthRepository {
	return &BirthRepository{db}
}

func (r *BirthRepository) GetBestIntervals(userId string) (*[]IntervalAnimal, error) {
	query := `
        with interval_list as (
			select 
				b.mother_id,
				extract(days from c.birth_date - lag(c.birth_date) over (partition by b.mother_id order by c.birth_date)) birth_interval
			from births b join animals c on c.id = b.calf_id
			where b.user_id = $1 and b.deleted_at is null
		),
		average_list as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                count(l.*) birth_numbers,
                avg(l.birth_interval) interval_average
            from interval_list l
                left join animals a on a.id = l.mother_id
            group by animal_name
        ),
        birth_stats as (
            select 
                max(birth_numbers) max_births,
				avg(interval_average) interval_median
            from average_list
        ),
        scores as (
            select 
                al.*,
                al.birth_numbers / bs.max_births  birth_score,
                (375 / al.interval_average) interval_score
            from average_list al, birth_stats bs
            where al.interval_average is not null
        ),
        best_scored as (
            select 
                animal_name,
                interval_average,
                birth_numbers,
                (interval_score*0.6 + birth_score*0.4) reproductive_score,
                ((interval_average / b.interval_median) - 1) * 100 average_rate
            from scores, birth_stats b
            order by reproductive_score desc
            limit 10
        )
        select 
            animal_name,
            interval_average,
            birth_numbers,
            average_rate
        from best_scored
        order by interval_average
    `
	return repositoriesUtil.GetList[IntervalAnimal](r.DB, query, userId)
}

func (r *BirthRepository) GetWorstIntervals(userId string) (*[]IntervalAnimal, error) {
	query := `
        with interval_list as (
			select 
				b.mother_id,
				extract(days from c.birth_date - lag(c.birth_date) over (partition by b.mother_id order by c.birth_date)) birth_interval
			from births b join animals c on c.id = b.calf_id
			where b.user_id = $1 and b.deleted_at is null
		),
		average_list as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                count(l.*) birth_numbers,
                avg(l.birth_interval) interval_average
            from interval_list l
                left join animals a on a.id = l.mother_id
            group by animal_name
        ),
        birth_stats as (
            select 
                max(birth_numbers) max_births,
				avg(interval_average) interval_median
            from average_list
        ),
        scores as (
            select 
                al.*,
                (1 - (al.birth_numbers / bs.max_births))  birth_score,
                (1 - (375 / al.interval_average)) interval_score
            from average_list al, birth_stats bs
            where al.interval_average is not null
        ),
        best_scored as (
            select 
                animal_name,
                interval_average,
                birth_numbers,
                (interval_score*0.6 + birth_score*0.4) reproductive_score,
                ((interval_average / b.interval_median) - 1) * 100 average_rate
            from scores, birth_stats b
            order by reproductive_score desc
            limit 10
        )
        select 
            animal_name,
            interval_average,
            birth_numbers,
            average_rate
        from best_scored
        order by interval_average
    `
	return repositoriesUtil.GetList[IntervalAnimal](r.DB, query, userId)
}

func (r *BirthRepository) GetBirthIntervalHistory(userId string) (*IntervalStats, error) {
	query := `
		with year_birth_interval as (
			select 
				date_trunc('year', c.birth_date) birth_date,
				extract(days from c.birth_date - lag(c.birth_date) over (partition by b.mother_id order by c.birth_date)) birth_interval
			from births b
				join animals c on c.id = b.calf_id
			where b.user_id = $1 and b.deleted_at is null
		),
		cte as (
			select 
				birth_date,
				avg(birth_interval) interval_average
			from year_birth_interval
			where birth_interval <> 0
			group by 1
			order by birth_date desc
			limit 10
		)
		select * from cte order by birth_date
    `
	results, err := repositoriesUtil.GetList[BirthIntervalHist](r.DB, query, userId)

	if err != nil {
		return nil, err
	}

	intervalHist := *results
	var currentInterval, previousInterval, intervalTrend float64

	switch lenght := len(intervalHist); lenght {
	case 0:
		currentInterval = 0
		previousInterval = 0
		intervalTrend = 0
	case 1:
		currentInterval = intervalHist[lenght-1].IntervalAverage
		previousInterval = 0
		intervalTrend = 0
	default:
		currentInterval = intervalHist[lenght-1].IntervalAverage
		previousInterval = intervalHist[lenght-2].IntervalAverage
		intervalTrend = ((currentInterval / previousInterval) - 1) * 100
	}

	intervalResult := IntervalStats{
		IntervalTrend:     intervalTrend,
		BirthIntervalHist: intervalHist,
		CurrentInterval:   currentInterval,
	}

	return &intervalResult, nil
}

func (r *BirthRepository) GetDeathIndex(userId string) (*DeathStats, error) {
	query := `
        with death_tbl as (
            select 
                date_trunc('year', a.death_date) date,
                count(a.*) deaths
            from animals a
            where
                a.user_id = $1
                and a.death_date is not null
                and age(a.death_date, a.birth_date) < interval '1 year'
                and a.deleted_at is null
            group by 1
        ),
        birth_tbl as (
            select
                date_trunc('year', a.birth_date) date,
                count(b.*) births
            from births b join animals a on a.id = b.calf_id
            where b.user_id = $1 and b.deleted_at is null
            group by 1
        ),
        cte as (
			select
				date date_month,
				coalesce((deaths::float / nullif(births, 0)::float)*100, 0) death_index
			from birth_tbl full join death_tbl using(date)
			order by 1 desc
			limit 10
		)
		select * from cte order by date_month
    `
	results, err := repositoriesUtil.GetList[DeathIndexHist](r.DB, query, userId)

	if err != nil {
		return nil, err
	}

	indexHist := *results
	var currentIndex, previousIndex, indexTrend float64

	switch lenght := len(indexHist); lenght {
	case 0:
		currentIndex = 0
		previousIndex = 0
		indexTrend = 0
	case 1:
		currentIndex = indexHist[lenght-1].DeathIndex
		previousIndex = 0
		indexTrend = 0
	default:
		currentIndex = indexHist[lenght-1].DeathIndex
		previousIndex = indexHist[lenght-2].DeathIndex
		indexTrend = ((currentIndex / previousIndex) - 1) * 100
	}

	deathStats := DeathStats{
		DeathIndexHist:  indexHist,
		DeathIndexTrend: indexTrend,
		DeathIndex:      currentIndex,
	}

	return &deathStats, nil
}

func (r *BirthRepository) GetBirthHistory(userId string) (*[]BirthsByDate, error) {
	query := `
        with death_data as ( 
            select 
                date_trunc('month', death_date) date,
                count(*) death_total 
			from animals  
            where 
                deleted_at is null 
                and death_date is not null
                and age(death_date, birth_date) < interval '1 year'
                and user_id = $1
            group by 1
        ), 
        birth_data as (
            select 
                date_trunc('month', calf.birth_date) date,
                count(births.*) as birth_total
            from births 
                left join animals as calf on births.calf_id = calf.id
            where 
                births.user_id = $1
                and births.deleted_at is null 
            group by 1
        )
        select
            date,
            coalesce(birth_data.birth_total,0) birth_total,
            coalesce(death_data.death_total, 0) death_total
        from birth_data full join death_data using(date)
        where date >= now() - interval '5 years'
        order by date
    `
	return repositoriesUtil.GetList[BirthsByDate](r.DB, query, userId)
}

func (r *BirthRepository) TotalBySex(userId string) (*[]TotalBirthsBySex, error) {
	query := `
        select 
            date_trunc('month', calf.birth_date) birth_month,
            count(births.id) filter (where calf.sex = 'M') males,
            count(births.id) filter (where calf.sex = 'F') females
        from births
            left join animals calf on births.calf_id = calf.id
        where 
            births.user_id = $1 
            and calf.birth_date >= now() - interval '5 years'
            and births.deleted_at is null
        group by birth_month
        order by birth_month
    `
	return repositoriesUtil.GetList[TotalBirthsBySex](r.DB, query, userId)
}

func (r *BirthRepository) FindPage(
	userId string,
	sort string,
	order string,
	filter BirthEntryFilter,
	cursor string,
) (*entity.Page[BirthEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"calf_birth_date": {Field: "c.birth_date", Order: "asc"},
		"mother_order":    {Field: "coalesce(regexp_replace(m.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"mother_name":     {Field: "m.name", Order: "asc"},
		"birth_interval":  {Field: "coalesce(i.birth_interval, 0)", Order: "asc"},
		"id":              {Field: "m.id", Order: "asc"},
		"created_at":      {Field: "m.created_at", Order: "asc"},
	}

	query := `
		with interval_cte as (
			select
				b.id,
				extract(days from c.birth_date - lag(c.birth_date) over (partition by b.mother_id order by c.birth_date)) birth_interval
			from births b join animals c on c.id = b.calf_id
		)
        select 
            b.id,
            b.mother_id,
            b.calf_id,
            concat_ws(' - ', m.ring_number, m.name) mother_name,
            coalesce(regexp_replace(m.ring_number, '[^0-9]', '', 'g')::int, 0) mother_order,
            c.birth_date calf_birth_date,
            c.sex calf_sex,
            case 
                when c.name is null then ''
                else concat_ws(' - ', c.ring_number, c.name)
                end as calf_name,
            c.father_id calf_father_id,
            concat_ws(' - ', f.ring_number, f.name) calf_father,
			i.birth_interval,
            b.observation
        from births b
			join interval_cte i on i.id = b.id
            left join animals m on m.id = b.mother_id
            left join animals c on c.id = b.calf_id
            left join animals f on f.id = c.father_id
    `
	whereExpression := "where b.user_id = $1 and b.deleted_at is null "
	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "b", 2)
	cursorExpression, nextParam, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	if cursorExpression != "" {
		whereExpression += " and " + cursorExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	orderExpression := " order by " + sortExpression
	query += whereExpression + orderExpression
	return repositoriesUtil.GetPage[BirthEntry](r.DB, query, sort, 100, args...)
}

func (r *BirthRepository) FindPageFooter(userId string, filter BirthEntryFilter) (*BirthFooter, error) {
	query := `
        with interval_cte as (
			select 
				b.id,
				extract(days from c.birth_date - lag(c.birth_date) over (partition by b.mother_id order by c.birth_date)) birth_interval
			from births b join animals c on c.id = b.calf_id
		)
		select
			count (*) total,
			avg(i.birth_interval) interval_average
		from births b
			join interval_cte i on i.id = b.id
            left join animals m on m.id = b.mother_id
            left join animals c on c.id = b.calf_id
            left join animals f on f.id = c.father_id
    `
	whereExpression := "where b.user_id = $1 and b.deleted_at is null "
	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "b", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	query += whereExpression
	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	return repositoriesUtil.GetOne[BirthFooter](r.DB, query, args...)
}
