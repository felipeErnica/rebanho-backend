package insemination

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type InseminationRepository struct {
	Db *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) *InseminationRepository {
	return &InseminationRepository{db}
}

func (r *InseminationRepository) GetBirthRateStats(userId string) (*BirthRateStats, error) {
	query := `
        with totals as (
            select distinct
                date_trunc('month', g.insemination_date) date_month,
                count(*) over (order by date_trunc('month', g.insemination_date)) total,
                count(*) filter (where status = 'SUCCESS') over (order by date_trunc('month', g.insemination_date)) successful
            from insemination_entries i
                left join insemination_groups g on g.id = i.group_id 
            where 
                age(g.insemination_date) > interval '11 months'
                and i.user_id = $1
                and i.deleted_at is null
            order by 1
            limit 10
        )
        select 
            date_month,
            (successful::float/total::float)*100 birth_rate
        from totals
    `
	result, err := repositoriesUtil.GetList[BirthRateHist](r.Db, query, userId)
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

func (r *InseminationRepository) GetInseminationStats(userId string) (*[]InseminationHist, error) {
	query := `
        with totals as (
            select
                date_trunc('month', g.insemination_date) date_month,
                count(*) total,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries i
                left join insemination_groups g on g.id = i.group_id 
            where 
                age(g.insemination_date) > interval '11 months'
                and i.user_id = $1
                and i.deleted_at is null
            group by 1
            order by 1
            limit 60
        )
        select 
            date_month,
            total,
            (successful::float/total::float)*100 birth_rate
        from totals
    `
	return repositoriesUtil.GetList[InseminationHist](r.Db, query, userId)
}

func (r *InseminationRepository) GetPregnantNumbers(userId string) (*PregnantsNumber, error) {
	query := `
        select count(*) filter (where status = 'PREGNANT') pregnants_number
        from insemination_entries 
        where user_id = $1 and  deleted_at is null
    `
	return repositoriesUtil.GetOne[PregnantsNumber](r.Db, query, userId)
}

func (r *InseminationRepository) GetBestBull(userId string) (*[]InseminationBulls, error) {
	query := `
        with totals as (
            select
                bull.name bull_name,
                count(*) total,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries i
                left join insemination_groups g on g.id = i.group_id 
                left join animals bull on g.bull_id = bull.id 
            where 
                age(g.insemination_date) > interval '11 months'
                and i.user_id = $1 and i.deleted_at is null
            group by 1
        ),
        cte as (
            select 
                bull_name,
                total,
                (successful::float/total::float)*100 birth_rate
            from totals
        ),
        average as (select avg(birth_rate) average_rate from cte) 
        select
            cte.*,
            ((cte.birth_rate / avg.average_rate) - 1)*100 comparison_rate
        from cte, average avg
        order by cte.birth_rate desc
    `
	return repositoriesUtil.GetList[InseminationBulls](r.Db, query, userId)
}

func (r *InseminationRepository) GetLastGroups(userId string) (*[]InseminationGroup, error) {
	query := `
        with totals_entries as (
            select
                count(*) totals,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries
            where user_id = $1 and deleted_at is null
        ),
        general_average as (select (successful::float / totals::float)*100 average_birth_rate from totals_entries),
        totals_group as (
            select 
                group_id,
                count(*) cow_number,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries i
            where user_id = $1 and deleted_at is null
            group by 1
        ),
        stats as (
            select
                group_id, 
                cow_number,
                (successful::float / cow_number::float)*100 birth_rate
            from totals_group
        )
        select 
            g.id,
            g.bull_id,
            b.name bull_name,
            g.insemination_date,
            s.cow_number,
            s.birth_rate,
            ((s.birth_rate / avg.average_birth_rate) - 1)*100 comparison_rate
        from general_average avg, insemination_groups g
            left join animals b on b.id = g.bull_id
            left join stats s on s.group_id = g.id
        where g.user_id = $1 and g.deleted_at is null
        order by g.insemination_date desc
        limit 10
    `
	return repositoriesUtil.GetList[InseminationGroup](r.Db, query, userId)
}

func (r *InseminationRepository) GetLastEntries(userId string) (*[]InseminationEntry, error) {
	query := `
        with last_group as (
            select id group_id, bull_id, insemination_date 
            from insemination_groups 
            where 
                deleted_at is null
                and user_id = $1
            order by insemination_date desc
            limit 1
        )
        select 
            concat_ws(' - ', a.ring_number, a.name) animal_name,
            g.insemination_date,
            b.name bull_name,
            i.status,
            i.observation
        from insemination_entries i
            join last_group g using (group_id)
            left join animals a on a.id = i.animal_id
            left join animals b on b.id = g.bull_id
        order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
    `
	return repositoriesUtil.GetList[InseminationEntry](r.Db, query, userId)
}

func (r *InseminationRepository) FindEntriesPage(
	userId string,
	filter InseminationEntryFilter,
	sort string,
	order string,
	cursor string,
) (*entity.Page[InseminationEntry], error) {

	sortMap := map[string]string{
		"animal_order":      "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)",
		"name":              "a.name",
		"insemination_date": "coalesce(g.insemination_date, '-infinity')",
	}

	query := `
        select 
            i.id,
            coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
            concat_ws(' - ', a.ring_number, a.name) animal_name,
            a.name,
            i.group_id,
            g.insemination_date,
            b.name bull_name,
            i.observation,
            i.status,
            i.created_at
        from insemination_entries i
            left join insemination_groups g on g.id = i.group_id
            left join animals a on a.id = i.animal_id
            left join animals b on b.id = g.bull_id
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

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, "i", cursorArgs, nextParam)
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
	return repositoriesUtil.GetPage[InseminationEntry](r.Db, query, sort, 100, args...)
}

func (r *InseminationRepository) FindEntriesByGroup(userId string, groupId string) (*[]InseminationEntry, error) {
	query := `
        select 
            i.id,
            concat_ws(' - ', a.ring_number, a.name) animal_name,
            i.observation,
            i.status
        from insemination_entries i
            left join animals a on a.id = i.animal_id
        where i.user_id = $1 and i.group_id = $2 and i.deleted_at is null
    `
	return repositoriesUtil.GetList[InseminationEntry](r.Db, query, userId, groupId)
}

func (r *InseminationRepository) GetEntriesByGroupFoot(userId string, groupId string) (*EntryFooter, error) {
	query := `
        with counting as (
            select 
                count(*) totals,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries
            where user_id = $1 and group_id = $2 and deleted_at is null
        )
        select 
            totals,
            (successful::float / totals::float)*100 birth_rate
        from counting
    `
	return repositoriesUtil.GetOne[EntryFooter](r.Db, query, userId, groupId)
}

func (r *InseminationRepository) FindGroups(userId string) (*[]InseminationGroup, error) {
	query := `
        with totals_entries as (
            select
                count(*) totals,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries
            where user_id = $1 and deleted_at is null
        ),
        general_average as (select (successful::float / totals::float)*100 average_birth_rate from totals_entries),
        totals_group as (
            select 
                group_id,
                count(*) cow_number,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries i
            where user_id = $1 and deleted_at is null
            group by 1
        ),
        stats as (
            select
                group_id, 
                cow_number,
                (successful::float / cow_number::float)*100 birth_rate
            from totals_group
        )
        select 
            g.id,
            g.bull_id,
            b.name bull_name,
            g.insemination_date,
            s.cow_number,
            s.birth_rate,
            ((s.birth_rate / avg.average_birth_rate) - 1)*100 comparison_rate
        from general_average avg, insemination_groups g
            left join animals b on b.id = g.bull_id
            left join stats s on s.group_id = g.id
        where g.user_id = $1 and g.deleted_at is null
        order by g.insemination_date desc
    `
	return repositoriesUtil.GetList[InseminationGroup](r.Db, query, userId)
}

func (r *InseminationRepository) GetGroupsFoot(userId string) (*GroupFooter, error) {
	query := `
        with total_stats as (
            select
                count(*) totals,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries
            where user_id = $1 and deleted_at is null
        ),
        general_average as (select (successful::float / totals::float)*100 average_birth_rate from total_stats),
        totals_group as (select count(*) totals from insemination_groups where user_id = $1 and deleted_at is null)
        select totals, average_birth_rate from totals_group, general_average a
    `
    return repositoriesUtil.GetOne[GroupFooter](r.Db, query, userId)
}
