package insemination

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type InseminationRepository struct {
	DB *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) *InseminationRepository {
	return &InseminationRepository{db}
}

func (r *InseminationRepository) GetBirthRateStats(userId string) (*BirthRateStats, error) {
	query := `
        with totals as (
            select distinct
                date_trunc('month', insemination_date) date_month,
                count(*) over (order by date_trunc('month', insemination_date)) total,
                count(*) filter (where status = 'SUCCESS') over (order by date_trunc('month', insemination_date)) birth_success
            from insemination_entries 
            where 
                age(insemination_date) > interval '11 months'
                and user_id = $1
                and deleted_at is null
            order by 1
            limit 10
        )
        select 
            date_month,
            (birth_success::float/total::float)*100 birth_rate
        from totals
    `
	result, err := repositoriesUtil.GetList[BirthRateHist](r.DB, query, userId)
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

	stats := &BirthRateStats{
		Hist:    birthRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *InseminationRepository) GetPregnancyRateStats(userId string) (*PregnancyRateStats, error) {
	query := `
        with totals as (
            select distinct
                date_trunc('month', insemination_date) date_month,
                count(*) over (order by date_trunc('month', insemination_date)) total,
                count(*) filter (where pregnancy_status = 'SUCCESS') over (order by date_trunc('month', insemination_date)) pregnancy_success
            from insemination_entries 
            where 
                age(insemination_date) > interval '11 months'
                and user_id = $1
                and deleted_at is null
            order by 1
            limit 10
        )
        select 
            date_month,
            (pregnancy_success::float/total::float)*100 pregnancy_rate
        from totals
    `
	result, err := repositoriesUtil.GetList[PregnancyRateHist](r.DB, query, userId)
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

	stats := &PregnancyRateStats{
		Hist:    pregnancyRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *InseminationRepository) GetInseminationStats(userId string) (*[]InseminationHist, error) {
	query := `
        with totals as (
            select
                date_trunc('month', insemination_date) date_month,
                count(*) total,
                count(*) filter (where status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from insemination_entries 
            where 
                age(insemination_date) > interval '11 months'
                and user_id = $1
                and deleted_at is null
            group by 1
            order by 1
            limit 60
        )
        select 
            date_month,
            total,
            (birth_success::float/total::float)*100 birth_rate,
            (pregnancy_success::float/total::float)*100 pregnancy_rate
        from totals
    `
	return repositoriesUtil.GetList[InseminationHist](r.DB, query, userId)
}

func (r *InseminationRepository) GetPregnantNumbers(userId string) (*PregnantsNumber, error) {
	query := `
        select count(*) filter (where pregnancy_status = 'SUCCESS' and status = 'STAND_BY') pregnants_number
        from insemination_entries 
        where user_id = $1 and  deleted_at is null
    `
	return repositoriesUtil.GetOne[PregnantsNumber](r.DB, query, userId)
}

func (r *InseminationRepository) GetBestBull(userId string) (*[]InseminationBulls, error) {
	query := `
        with totals as (
            select
                bull.name bull_name,
                count(*) total,
                count(*) filter (where status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from insemination_entries i
                left join animals bull on i.bull_id = bull.id 
            where 
                age(i.insemination_date) > interval '11 months'
                and i.user_id = $1 and i.deleted_at is null
            group by 1
        ),
        cte as (
            select 
                bull_name,
                total,
                (birth_success::float/total::float)*100 birth_rate,
                (pregnancy_success::float/total::float)*100 pregnancy_rate
            from totals
        ),
        average as (
            select 
                avg(birth_rate) birth_rate,
                avg(pregnancy_rate) pregnancy_rate
            from cte
        ) 
        select
            cte.*,
            ((cte.birth_rate / avg.birth_rate) - 1)*100 birth_comparison_rate,
            ((cte.birth_rate / avg.pregnancy_rate) - 1)*100 pregnancy_comparison_rate
        from cte, average avg
        order by cte.birth_rate desc
    `
	return repositoriesUtil.GetList[InseminationBulls](r.DB, query, userId)
}

func (r *InseminationRepository) GetLastGroups(userId string) (*[]InseminationGroup, error) {
	query := `
        with totals_entries as (
            select
                count(*) totals,
                count(*) filter (where status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from insemination_entries
            where user_id = $1 and deleted_at is null
        ),
        general_average as (
            select 
                (birth_success::float / totals::float)*100 average_birth_rate,
                (pregnancy_success::float / totals::float)*100 average_pregnancy_rate 
            from totals_entries
        ),
        totals_group as (
            select 
                insemination_date,
                bull_id,
                count(*) cow_number,
                count(*) filter (where status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from insemination_entries i
            where user_id = $1 and deleted_at is null
            group by insemination_date, bull_id
        ),
        stats as (
            select
                insemination_date,
                bull_id,
                cow_number,
                (birth_success::float / cow_number::float)*100 birth_rate,
                (pregnancy_success::float / cow_number::float)*100 pregnancy_rate
            from totals_group
        )
        select 
            b.name bull_name,
            s.bull_id,
            s.insemination_date,
            s.cow_number,
            s.birth_rate,
            s.pregnancy_rate,
            ((s.birth_rate / avg.average_birth_rate) - 1)*100 birth_comparison_rate,
            ((s.pregnancy_rate / avg.average_pregnancy_rate) - 1)*100 pregnancy_comparison_rate
        from general_average avg, stats s
            left join animals b on b.id = s.bull_id
        order by s.insemination_date desc
        limit 10
    `
	return repositoriesUtil.GetList[InseminationGroup](r.DB, query, userId)
}

func (r *InseminationRepository) GetLastEntries(userId string) (*[]InseminationEntry, error) {
	query := `
        with last_date as (
            select max(insemination_date) max_date
            from insemination_entries 
            where deleted_at is null and user_id = $1
        )
        select 
            concat_ws(' - ', a.ring_number, a.name) animal_name,
            i.insemination_date,
            b.name bull_name,
            i.pregnancy_status,
            i.status,
            i.observation
        from last_date l, insemination_entries i
            left join animals a on a.id = i.animal_id
            left join animals b on b.id = i.bull_id
        where i.user_id = $1 and i.deleted_at is null and i.insemination_date = l.max_date
        order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
    `
	return repositoriesUtil.GetList[InseminationEntry](r.DB, query, userId)
}

func (r *InseminationRepository) FindEntriesPage(
	userId string,
	filter InseminationEntryFilter,
	sort string,
	order string,
	cursor string,
) (*entity.Page[InseminationEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"animal_order":      {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"name":              {Field: "a.name", Order: "asc"},
		"insemination_date": {Field: "coalesce(i.insemination_date, '-infinity')", Order: "asc"},
		"id":                {Field: "i.id", Order: "asc"},
		"created_at":        {Field: "i.created_at", Order: "asc"},
	}

	query := `
        select 
            i.id,
            coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
            concat_ws(' - ', a.ring_number, a.name) animal_name,
            a.name,
            i.insemination_date,
            b.name bull_name,
            i.observation,
            i.pregnancy_status,
            i.status,
            i.created_at
        from insemination_entries i
            left join animals a on a.id = i.animal_id
            left join animals b on b.id = i.bull_id
    `
	whereExpression := "where i.user_id = $1 and i.deleted_at is null"
	orderExpression := " order by "

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "i", 2)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	if cursorExpression != "" {
		whereExpression += " and " + cursorExpression
	}

	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	orderExpression += sortExpression
	query += whereExpression + orderExpression
	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return repositoriesUtil.GetPage[InseminationEntry](r.DB, query, sort, 100, args...)
}

func (r *InseminationRepository) FindEntriesByGroup(userId string, bullId string, inseminationDate time.Time) (*[]InseminationEntry, error) {
	query := `
        select 
            i.id,
            concat_ws(' - ', a.ring_number, a.name) animal_name,
            i.observation,
            i.pregnancy_status,
            i.status
        from insemination_entries i
            left join animals a on a.id = i.animal_id
        where i.user_id = $1 and i.deleted_at is null and i.bull_id = $2 and i.insemination_date = $3
        order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
    `
	return repositoriesUtil.GetList[InseminationEntry](r.DB, query, userId, bullId, inseminationDate)
}

func (r *InseminationRepository) GetEntriesByGroupFoot(
	userId string,
	bullId string,
	inseminationDate time.Time,
) (*InseminationFooter, error) {
	query := `
        with counting as (
            select 
                count(*) totals,
                count(*) filter (where status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from insemination_entries
            where user_id = $1 and deleted_at is null and bull_id = $2 and insemination_date = $3
        )
        select 
            totals,
            (birth_success::float / totals::float)*100 average_birth_rate,
            (pregnancy_success::float / totals::float)*100 average_pregnancy_rate
        from counting
    `
	return repositoriesUtil.GetOne[InseminationFooter](r.DB, query, userId, bullId, inseminationDate)
}

func (r *InseminationRepository) FindGroups(userId string) (*[]InseminationGroup, error) {
	query := `
        with totals_entries as (
            select
                count(*) totals,
                count(*) filter (where status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from insemination_entries
            where user_id = $1 and deleted_at is null
        ),
        general_average as (
            select 
                (birth_success::float / totals::float)*100 average_birth_rate,
                (pregnancy_success::float / totals::float)*100 average_pregnancy_rate 
            from totals_entries
        ),
        totals_group as (
            select 
                insemination_date,
                bull_id,
                count(*) cow_number,
                count(*) filter (where status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from insemination_entries i
            where user_id = $1 and deleted_at is null
            group by insemination_date, bull_id
        ),
        stats as (
            select
                insemination_date,
                bull_id,
                cow_number,
                (birth_success::float / cow_number::float)*100 birth_rate,
                (pregnancy_success::float / cow_number::float)*100 pregnancy_rate
            from totals_group
        )
        select 
            b.name bull_name,
            s.bull_id,
            s.insemination_date,
            s.cow_number,
            s.birth_rate,
            s.pregnancy_rate,
            ((s.birth_rate / avg.average_birth_rate) - 1)*100 birth_comparison_rate,
            ((s.pregnancy_rate / avg.average_pregnancy_rate) - 1)*100 pregnancy_comparison_rate
        from general_average avg, stats s
            left join animals b on b.id = s.bull_id
        order by s.insemination_date desc
    `
	return repositoriesUtil.GetList[InseminationGroup](r.DB, query, userId)
}

func (r *InseminationRepository) GetGroupsFoot(userId string) (*InseminationFooter, error) {
	query := `
        with total_stats as (
            select
                count(*) totals,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries
            where user_id = $1 and deleted_at is null
        ),
        general_average as (select (successful::float / totals::float)*100 average_birth_rate from total_stats),
        grouped_insemination as (
            select distinct bull_id, insemination_date 
            from insemination_entries
            where user_id = $1 and deleted_at is null    
        )
        select totals, average_birth_rate from general_average a, (select count(*) totals from grouped_insemination)
    `
	return repositoriesUtil.GetOne[InseminationFooter](r.DB, query, userId)
}

func (r *InseminationRepository) GetEntriesFoot(userId string, filter InseminationEntryFilter) (*InseminationFooter, error) {
	totalQuery := `
        select
            count(*) totals,
            count(*) filter (where status = 'SUCCESS') birth_success,
            count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
        from insemination_entries
        where user_id = $1 and deleted_at is null
    `
	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "insemination_entries", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		totalQuery += " and " + filterExpression
	}

	query := fmt.Sprintf(`
        with total_stats as (%s)
        select 
            totals,
            (birth_success::float / totals::float)*100 average_birth_rate,
            (pregnancy_success::float / totals::float)*100 average_pregnancy_rate
            from total_stats
    `, totalQuery)

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetOne[InseminationFooter](r.DB, query, args...)
}

func (r *InseminationRepository) SearchInseminationBulls(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select distinct a.id, a.name label
        from animals a join insemination_entries i on i.bull_id = a.id
        where i.user_id = $1 and i.deleted_at is null and a.name ilike $2
        order by a.name
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}
