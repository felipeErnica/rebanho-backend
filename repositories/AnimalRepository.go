package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type AnimalRepository struct {
    Impl PageRepositoryImpl[entity.Animal]
    SelectQuery *util.QueryConstructor
}

func (r *AnimalRepository) Init() {
    dateFields:= []string{
        "birth_date",
        "death_date",
        "created_at",
        "deleted_at",
    }

    selectQueryBody:=new(util.QueryConstructor).Select("animals", "id", "name", "identification_number", "birth_date", "sex", "death_date",
        "weaning_date", "status", "average_prod", "average_birth_interval", "max_peak", "isr", "children_quantity", "observation")
        selectQueryBody.AndSelect("mother", "id", "name", "identification_number")
        selectQueryBody.AndSelect("father", "id", "name", "identification_number")
        selectQueryBody.AndSelect("pastures", "id", "name")
        selectQueryBody.From("animals", "")
        selectQueryBody.LeftJoin("animals", "father").On("father.id", "animals.father_id")
        selectQueryBody.LeftJoin("animals", "mother").On("mother.id", "animals.mother_id")
        selectQueryBody.LeftJoin("pastures", "").On("pastures.id", "animals.pasture_id")
        selectQueryBody.Where("animals.deleted_at is null")
    r.SelectQuery = selectQueryBody

    insertQuery:=new(util.QueryConstructor).Insert("animals", "id", "name", "identification_number", "father_id", "mother_id",
        "birth_date", "death_date", "pasture_id", "weaning_date", "status", "average_prod", 
        "average_birth_interval", "max_peak", "isr",
        "children_quantity", "observation", "created_at")
    updateQuery:=new(util.QueryConstructor).Update("animals", "id", "name", "identification_number", "father_id", "mother_id",
        "birth_date", "death_date", "pasture_id", "weaning_date", "status", "average_prod", 
        "average_birth_interval", "max_peak", "isr",
        "children_quantity", "observation", "created_at")

    baseRepo:= &RepositoryImpl[entity.Animal]{
        Repository: r,
        TableName: "animals",
        SelectQueryBody: selectQueryBody.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
    }
    r.Impl = PageRepositoryImpl[entity.Animal]{
        Base: baseRepo,
        PageRepository: r,
        DateFields: dateFields,
    }
}

func (r *AnimalRepository) setNewEntity(model *entity.Animal, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
}

func (r *AnimalRepository) buildListEntity(sqlRows *sql.Rows) (list *[]entity.Animal, err error) {
    var animals []entity.Animal
    for sqlRows.Next() {
        var animal entity.Animal
        err = sqlRows.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.BirthDate, &animal.Sex,
            &animal.DeathDate, &animal.WeaningDate, &animal.Status, &animal.AverageProd, 
            &animal.AverageBirthInterval, &animal.MaxPeak, &animal.Isr, &animal.ChildrenQuantity, &animal.Observation,
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
        &animal.AverageBirthInterval, &animal.MaxPeak, &animal.Isr, &animal.ChildrenQuantity, &animal.Observation,
        &animal.Mother.Id, &animal.Mother.Name, &animal.Mother.IdentificationNumber, 
        &animal.Father.Id, &animal.Father.Name, &animal.Father.IdentificationNumber,
        &animal.Pasture.Id, &animal.Pasture.Name)
    return &animal, err
}

func (r *AnimalRepository) saveOrUpdateScan(query string, animal *entity.Animal) error {
    return execQuery(query, animal.Id, animal.Name, animal.IdentificationNumber, animal.Father.Id, animal.Mother.Id,
        animal.BirthDate, animal.DeathDate, animal.Pasture.Id, animal.WeaningDate, animal.Status, animal.AverageProd, 
        animal.AverageBirthInterval, animal.MaxPeak, animal.Isr, animal.ChildrenQuantity, animal.Observation, animal.CreatedAt)
}

func (r *AnimalRepository) getFields(sort string) (firstField string, secondField string) {
    switch (sort) {
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
    case "max_peak":
        return "animals.max_peak", "animals.id"
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
    switch (sort) {
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
    case "max_peak":
        key = fmt.Sprintf("%f,%s", lastEntry.MaxPeak, lastEntry.Id) 
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

func (r *AnimalRepository) FindPage(sort string, direction string, cursor string) (page *entity.Page[entity.Animal], err error) {
    return r.Impl.FindRandomQueryPage(r.SelectQuery, sort, direction, cursor)
}

func (r *AnimalRepository) FindById(id string) (*entity.Animal, error) {
    return r.Impl.FindById(id)
}

func (r *AnimalRepository) FindByFatherId(fatherId string) (*[]entity.Animal, error) {
    query:=r.SelectQuery.And("mother.id = $1").Order("animals.birth_date asc")
    return r.Impl.FindListByQuery(query, fatherId)
}

func (r *AnimalRepository) FindByPastureId(sort string, direction string, 
    cursor string, pastureId string) (page *entity.Page[entity.Animal], err error) {
    query:= new(util.QueryConstructor).FromQuery(r.Impl.Base.SelectQueryBody)
    query.Where("animals.pastureId = $1")
    return r.Impl.FindRandomQueryPage(query, sort, direction, cursor, pastureId)
}

func (r *AnimalRepository) FindByName(name string) (*[]entity.Animal, error) {
    query:=r.SelectQuery.And("animals.name = $1").And("animals.user_id = $2")
    return r.Impl.FindListByQuery(query, name, GetUserId())
}

func (r *AnimalRepository) FindByIdentificationNumber(number string) (*[]entity.Animal, error) {
    query:=r.SelectQuery.And("animals.name = $1").And("animals.user_id = $2")
    return r.Impl.FindListByQuery(query, number, GetUserId())
}

func (r *AnimalRepository) Add(create entity.Animal) (*entity.Animal, error) {
    return r.Impl.Add(create)
}

func (r *AnimalRepository) Save(animal *entity.Animal) error {
    return r.Impl.Save(animal)
}

func (r *AnimalRepository) Delete(id string) error {
    return r.Impl.SoftDelete(id)
}
