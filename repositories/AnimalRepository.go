package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type AnimalRepository struct {
	Impl        PageRepositoryImpl[entity.Animal]
	SelectQuery util.SelectConstructor
}

func (r *AnimalRepository) Init() {
	dateFields := []string{
		"birth_date",
		"death_date",
		"created_at",
		"deleted_at",
	}

	r.SelectQuery = *util.NewSelectQuery(util.SELECT, 
		*util.NewNamedGroup("animals", "id", "name", "identification_number", "birth_date", "sex", "death_date",
		"weaning_date", "status", "average_prod", "average_birth_interval", 
		"average_peak", "isr", "children_quantity", "observation"),
		*util.NewNamedGroup("mother", "id", "name", "identification_number"), 
		*util.NewNamedGroup("father", "id", "name", "identification_number"),
		*util.NewNamedGroup("pastures", "id", "name")).
		From("animals").
		Joins(
			"left join animals as father on father.id = animals.father_id", 
			"left join animals as mother on mother.id = animals.mother_id", 
			"left join pastures on pastures.id = animals.pasture_id")

	insertQuery := util.NewInsertQuery("animals", "id", "name", "identification_number", "father_id", "mother_id",
		"birth_date", "death_date", "pasture_id", "weaning_date", "status", "average_prod",
		"average_birth_interval", "average_peak", "isr",
		"children_quantity", "observation", "created_at", "user_id")
	updateQuery := util.NewUpdateQuery("animals", "id", "name", "identification_number", "father_id", "mother_id",
		"birth_date", "death_date", "pasture_id", "weaning_date", "status", "average_prod",
		"average_birth_interval", "average_peak", "isr",
		"children_quantity", "observation", "created_at", "user_id")

	baseRepo := &RepositoryImpl[entity.Animal]{
		Repository:      r,
		TableName:       "animals",
		SelectQueryBody: r.SelectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
	}
	r.Impl = PageRepositoryImpl[entity.Animal]{
		Base:           baseRepo,
		PageRepository: r,
		DateFields:     dateFields,
	}
}

func (r *AnimalRepository) setNewEntity(model *entity.Animal, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
	model.UserId = GetUserId()
}

func (r *AnimalRepository) buildListEntity(sqlRows *sql.Rows) (list *[]entity.Animal, err error) {
	var animals []entity.Animal
	for sqlRows.Next() {
		var animal entity.Animal
		err = sqlRows.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.BirthDate, &animal.Sex,
			&animal.DeathDate, &animal.WeaningDate, &animal.Status, &animal.AverageProd,
			&animal.AverageBirthInterval, &animal.AveragePeak, &animal.Isr, &animal.ChildrenQuantity, &animal.Observation,
			&animal.Mother.Id, &animal.Mother.Name, &animal.Mother.IdentificationNumber,
			&animal.Father.Id, &animal.Father.Name, &animal.Father.IdentificationNumber,
			&animal.Pasture.Id, &animal.Pasture.Name)
		if err != nil {
			return nil, err
		}
		animals = append(animals, animal)
	}
	return &animals, err
}

func (r *AnimalRepository) buildEntity(sqlStatement *sql.Row) (model *entity.Animal, err error) {
	var animal entity.Animal
	err = sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.BirthDate, &animal.Sex,
		&animal.DeathDate, &animal.WeaningDate, &animal.Status, &animal.AverageProd,
		&animal.AverageBirthInterval, &animal.AveragePeak, &animal.Isr, &animal.ChildrenQuantity, &animal.Observation,
		&animal.Mother.Id, &animal.Mother.Name, &animal.Mother.IdentificationNumber,
		&animal.Father.Id, &animal.Father.Name, &animal.Father.IdentificationNumber,
		&animal.Pasture.Id, &animal.Pasture.Name)
	return &animal, err
}

func (r *AnimalRepository) saveOrUpdateScan(query string, animal *entity.Animal) error {
	return execQuery(query, animal.Id, animal.Name, animal.IdentificationNumber, animal.Father.Id, animal.Mother.Id,
		animal.BirthDate, animal.DeathDate, animal.Pasture.Id, animal.WeaningDate, animal.Status, animal.AverageProd,
		animal.AverageBirthInterval, animal.AveragePeak, animal.Isr, animal.ChildrenQuantity, animal.Observation, animal.CreatedAt)
}

func (r *AnimalRepository) getFields(sort string) (firstField string, secondField string) {
	switch sort {
	case "name":
		return "animals.name", "animals.id"
	case "identification_number":
		return "animals.animal_order", "animals.id"
	case "birth_date":
		return "animals.birth_date", "animals.id"
	case "death_date":
		return "animals.death_date", "animals.id"
	case "average_prod":
		return "animals.average_prod", "animals.id"
	case "average_birth_interval":
		return "animals.average_birth_interval", "animals.id"
	case "average_peak":
		return "animals.average_peak", "animals.id"
	case "children_quantity":
		return "animals.children_quantity", "animals.id"
	case "isr":
		return "animals.isr", "animals.id"
	case "deleted_at":
		return "animals.deleted_at", "animals.id"
	default:
		return "animals.created_at", "animals.id"
	}
}

func (r *AnimalRepository) createKey(sort string, lastEntry *entity.Animal) string {
	var key string
	switch sort {
	case "name":
		key = fmt.Sprintf("%s,%s", "null", lastEntry.Id)
		if lastEntry.Name != nil {
			key = fmt.Sprintf("%s,%s", *lastEntry.Name, lastEntry.Id)
		}
	case "identification_number":
		key = fmt.Sprintf("%s,%s", "null", lastEntry.Id)
		if lastEntry.IdentificationNumber != nil {
			key = fmt.Sprintf("%s,%s", *lastEntry.IdentificationNumber, lastEntry.Id)
		}
	case "birth_date":
		key = fmt.Sprintf("%s,%s", "null", lastEntry.Id)
		if lastEntry.BirthDate != nil {
			key = fmt.Sprintf("%s,%s", lastEntry.BirthDate.Format(time.RFC3339Nano), lastEntry.Id)
		}
	case "death_date":
		key = fmt.Sprintf("%s,%s", "null", lastEntry.Id)
		if lastEntry.DeathDate != nil {
			key = fmt.Sprintf("%s,%s", lastEntry.DeathDate.Format(time.RFC3339Nano), lastEntry.Id)
		}
	case "average_prod":
		key = fmt.Sprintf("%f,%s", lastEntry.AverageProd, lastEntry.Id)
	case "average_birth_interval":
		key = fmt.Sprintf("%f,%s", lastEntry.AverageBirthInterval, lastEntry.Id)
	case "average_peak":
		key = fmt.Sprintf("%f,%s", lastEntry.AveragePeak, lastEntry.Id)
	case "children_quantity":
		key = fmt.Sprintf("%d,%s", lastEntry.ChildrenQuantity, lastEntry.Id)
	case "isr":
		key = fmt.Sprintf("%f,%s", lastEntry.Isr, lastEntry.Id)
	case "deleted_at":
		key = fmt.Sprintf("%s,%s", "null", lastEntry.Id)
		if lastEntry.DeletedAt != nil {
			key = fmt.Sprintf("%s,%s", lastEntry.DeletedAt.Format(time.RFC3339Nano), lastEntry.Id)
		}
	default:
		key = fmt.Sprintf("%s,%s", lastEntry.CreatedAt.Format(time.RFC3339Nano), lastEntry.Id)
	}
	return key
}

func (r *AnimalRepository) filterQuery(filter *entity.AnimalFilter) (query util.SelectConstructor, args []any) {
	query = r.SelectQuery
	query.Where("animals.user_id = $1 and animals.deleted_at is null")
	args = append(args, GetUserId())
	numParam := 2

	if !filter.IsFiltered {
		return
	}

	if filter.Name != nil {
		query.AppendWhere(fmt.Sprintf("and animals.name like $%d", numParam))
		numParam++
		name := fmt.Sprintf("%%%s%%", *filter.Name)
		args = append(args, name)
	}

	if filter.IdentificationNumber != nil {
		query.AppendWhere(fmt.Sprintf("and animals.identification_number like $%d", numParam))
		numParam++
		name := fmt.Sprintf("%%%s%%", *filter.IdentificationNumber)
		args = append(args, name)
	}

	if filter.Sex != nil {
		query.AppendWhere(fmt.Sprintf("and animals.sex = $%d", numParam))
		numParam++
		args = append(args, *filter.Sex)
	}

	if filter.MinWeaningDate != nil {
		query.AppendWhere(fmt.Sprintf("and animals.weaning_date >= $%d", numParam))
		numParam++
		query.AppendWhere(fmt.Sprintf("and animals.weaning_date <= $%d", numParam))
		numParam++
		args = append(args, *filter.MinWeaningDate)
		args = append(args, *filter.MaxWeaningDate)
	}

	if filter.Fathers != nil {
		arr := *filter.Fathers
		params := fmt.Sprintf("$%d", numParam)
		numParam++
		args = append(args, arr[0])
		for i := 1; i < len(arr); i++ {
			params += fmt.Sprintf(", $%d", numParam)
			args = append(args, arr[i])
			numParam++
		}
		query.AppendWhere(fmt.Sprintf("and animals.father_id IN (%s)", params))
		numParam++
	}

	if filter.Mothers != nil {
		arr := *filter.Mothers
		params := fmt.Sprintf("$%d", numParam)
		numParam++
		args = append(args, arr[0])
		for i := 1; i < len(arr); i++ {
			params += fmt.Sprintf(", $%d", numParam)
			args = append(args, arr[i])
			numParam++
		}
		query.AppendWhere(fmt.Sprintf("and animals.mother_id IN (%s)", params))
		numParam++
	}

	if filter.MinBirthDate != nil {
		query.AppendWhere(fmt.Sprintf("and animals.birth_date >= $%d", numParam))
		numParam++
		query.AppendWhere(fmt.Sprintf("and animals.birth_date <= $%d", numParam))
		numParam++
		args = append(args, *filter.MinBirthDate)
		args = append(args, *filter.MaxBirthDate)
	}

	if filter.MinDeathDate != nil {
		query.AppendWhere(fmt.Sprintf("and animals.death_date >= $%d", numParam))
		numParam++
		query.AppendWhere(fmt.Sprintf("and animals.death_date <= $%d", numParam))
		numParam++
		args = append(args, *filter.MinDeathDate)
		args = append(args, *filter.MaxDeathDate)
	}

	if filter.Pastures != nil {
		arr := *filter.Pastures
		params := fmt.Sprintf("$%d", numParam)
		numParam++
		args = append(args, arr[0])
		for i := 1; i < len(arr); i++ {
			params += fmt.Sprintf(", $%d", numParam)
			args = append(args, arr[i])
			numParam++
		}
		query.AppendWhere(fmt.Sprintf("and animals.pasture_id IN (%s)", params))
		numParam++
	}

	if filter.Status != nil {
		arr := *filter.Status
		params := fmt.Sprintf("$%d", numParam)
		numParam++
		args = append(args, arr[0])
		for i := 1; i < len(arr); i++ {
			params += fmt.Sprintf(", $%d", numParam)
			args = append(args, arr[i])
			numParam++
		}
		query.AppendWhere(fmt.Sprintf("and animals.status IN (%s)", params))
		numParam++
	}

	if filter.MinIsr != nil {
		query.AppendWhere(fmt.Sprintf("and animals.isr >= $%d", numParam))
		numParam++
		query.AppendWhere(fmt.Sprintf("and animals.isr <= $%d", numParam))
		numParam++
		args = append(args, *filter.MinIsr)
		args = append(args, *filter.MaxIsr)
	}

	if filter.MinAverageBirthInterval != nil {
		query.AppendWhere(fmt.Sprintf("and animals.average_birth_interval >= $%d", numParam))
		numParam++
		query.AppendWhere(fmt.Sprintf("and animals.average_birth_interval <= $%d", numParam))
		numParam++
		args = append(args, *filter.MinAverageBirthInterval)
		args = append(args, *filter.MaxAverageBirthInterval)
	}

	if filter.MinAverageProd != nil {
		query.AppendWhere(fmt.Sprintf("and animals.average_prod >= $%d", numParam))
		numParam++
		query.AppendWhere(fmt.Sprintf("and animals.average_prod <= $%d", numParam))
		numParam++
		args = append(args, *filter.MinAverageProd)
		args = append(args, *filter.MaxAverageProd)
	}

	if filter.MinAveragePeak != nil {
		query.AppendWhere(fmt.Sprintf("and animals.average_peak >= $%d", numParam))
		numParam++
		query.AppendWhere(fmt.Sprintf("and animals.average_peak <= $%d", numParam))
		numParam++
		args = append(args, *filter.MinAveragePeak)
		args = append(args, *filter.MaxAveragePeak)
	}

	if filter.MinChildrenQuantity != nil {
		query.AppendWhere(fmt.Sprintf("and animals.children_quantity >= $%d", numParam))
		numParam++
		query.AppendWhere(fmt.Sprintf("and animals.children_quantity <= $%d", numParam))
		numParam++
		args = append(args, *filter.MinChildrenQuantity)
		args = append(args, *filter.MaxChildrenQuantity)
	}

	return query, args
}

func (r *AnimalRepository) FindPage(sort string, direction string,
	cursor string, filter *entity.AnimalFilter) (page *entity.Page[entity.Animal], err error) {
	query, args := r.filterQuery(filter)
	return r.Impl.FindRandomQueryPage(&query, sort, direction, cursor, args...)
}

func (r *AnimalRepository) FindMaxValues() (maxValues *entity.AnimalMaxValues, err error) {
	query := util.NewSelectQuery(util.MAX, 
        *util.NewGroup( "weaning_date", "birth_date", "death_date",
		"isr", "average_birth_interval", "average_prod", "average_peak", "children_quantity")).
		From("animals").
		Where("user_id = $1 and deleted_at is null")
	row := selectQueryOne(query.Build(), GetUserId())
	maxValues = new(entity.AnimalMaxValues)

	err = row.Scan(
		&maxValues.MaxWeaningDate,
		&maxValues.MaxBirthDate,
		&maxValues.MaxDeathDate,
		&maxValues.MaxIsr,
		&maxValues.MaxAverageBirthInterval,
		&maxValues.MaxAverageProd,
		&maxValues.MaxAveragePeak,
		&maxValues.MaxChildrenQuantity,
	)
	if err != nil {
		return
	}

	return maxValues, err
}

func (r *AnimalRepository) FindMinValues() (minValues *entity.AnimalMinValues, err error) {
	query := util.NewSelectQuery(util.MIN, *util.NewGroup( "weaning_date", "birth_date", "death_date",)).
		From("animals").Where("user_id = $1 and deleted_at is null")
	row := selectQueryOne(query.Build(), GetUserId())
	minValues = new(entity.AnimalMinValues)

	err = row.Scan(
		&minValues.MinWeaningDate,
		&minValues.MinBirthDate,
		&minValues.MinDeathDate,
	)
	if err != nil {
		return
	}

	return minValues, err
}

func (r *AnimalRepository) FindById(id string) (*entity.Animal, error) {
	return r.Impl.FindById(id)
}

func (r *AnimalRepository) FindByFatherId(fatherId string) (*[]entity.Animal, error) {
	query := r.SelectQuery.Where("animals.father_id = $1 and animals.deleted_at is null").
		OrderBy("animals.birth_date asc")
	return r.Impl.FindListByQuery(query, fatherId)
}

func (r *AnimalRepository) FindByMotherId(motherId string) (*[]entity.Animal, error) {
	query := r.SelectQuery.Where("animals.mother_id = $1 and animals.deleted_at is null").
		OrderBy("animals.birth_date asc")
	return r.Impl.FindListByQuery(query, motherId)
}

func (r *AnimalRepository) FindByPastureId(sort string, direction string,
	cursor string, pastureId string) (page *entity.Page[entity.Animal], err error) {
	query := r.SelectQuery.Where("animals.pastureId = $1")
	return r.Impl.FindRandomQueryPage(query, sort, direction, cursor, pastureId)
}

func (r *AnimalRepository) FindByName(name string) (*[]entity.Animal, error) {
	query := r.SelectQuery. Where("animals.name = $1 and animals.user_id = $2 and animals.deleted_at is null")
	return r.Impl.FindListByQuery(query, name, GetUserId())
}

func (r *AnimalRepository) FindByIdentificationNumber(number string) (*[]entity.Animal, error) {
	query := r.SelectQuery.Where("animals.name = $1 and animals.user_id = $2 and animals.deleted_at is null")
	return r.Impl.FindListByQuery(query, number, GetUserId())
}

func (r *AnimalRepository) Add(create *entity.Animal) (*entity.Animal, error) {
	return r.Impl.Add(create)
}

func (r *AnimalRepository) Save(animal *entity.Animal) error {
	return r.Impl.Save(animal)
}

func (r *AnimalRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
