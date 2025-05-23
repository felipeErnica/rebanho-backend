package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type MilkRepository struct {
	Impl        PageRepositoryImpl[entity.MilkEntry]
	SelectQuery util.SelectConstructor
}

func (r *MilkRepository) Init() {

	dateFields := []string{
		"entry_date",
		"created_at",
	}

    r.SelectQuery = *util.NewSelectQuery(util.SELECT, 
        *util.NewNamedGroup("milk", "id", "entry_date", "milk_quantity", "milk_entries.lactation_id"),
        *util.NewNamedGroup("animal", "id", "identification_number", "animal_order", "name"),
        *util.NewNamedGroup("pasture", "id", "name")).
        From("milk_entries as milk").
        Joins(
            "left join animals as animal on animal.id = milk.animal_id",
            "left join pastures as pasture on pasture.id = milk.pasture_id")

    insertQuery := util.NewInsertQuery("milk_entries", "id", "entry_date", "milk_quantity", "animal_id", 
        "pasture_id", "lactation_id", "created_at", "user_id")
    updateQuery := util.NewUpdateQuery("milk_entries", "id", "entry_date", "milk_quantity", "animal_id", 
        "pasture_id", "lactation_id", "created_at", "user_id")

	mainRepository := &RepositoryImpl[entity.MilkEntry]{
		Repository:      r,
		SelectQueryBody: r.SelectQuery,
		UpdateQuery:     *updateQuery,
		InsertQuery:     *insertQuery,
		TableName:       "milk_entries",
	}

	r.Impl = PageRepositoryImpl[entity.MilkEntry]{
		Base:           mainRepository,
		PageRepository: r,
		DateFields:     dateFields,
	}
}

func (r *MilkRepository) getFields(sort string) (firstField string, secondField string) {
	switch sort {
	case "name":
		return "animal.name", "milk_entries.id"
	case "identification_number":
		return "animal.animal_order", "milk_entries.id"
	case "entry_date":
		return "milk_entries.entry_date", "milk_entries.id"
	default:
		return "milk_entries.created_at", "milk_entries.id"
	}
}

func (r *MilkRepository) createKey(sort string, lastEntry *entity.MilkEntry) (key string) {
	switch sort {
	case "name":
		return fmt.Sprintf("%s,%s", lastEntry.AnimalName, lastEntry.Id)
	case "identification_number":
		return fmt.Sprintf("%d,%s", lastEntry.AnimalOrder, lastEntry.Id)
	case "entry_date":
		return fmt.Sprintf("%s,%s", lastEntry.EntryDate, lastEntry.Id)
	default:
		return fmt.Sprintf("%s,%s", lastEntry.CreatedAt, lastEntry.Id)
	}
}

func (r *MilkRepository) setNewEntity(model *entity.MilkEntry, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
}

func (r *MilkRepository) buildEntity(row *sql.Row) (model *entity.MilkEntry, err error) {
	var milk entity.MilkEntry
	err = row.Scan(&milk.Id, &milk.EntryDate, &milk.MilkQuantity, &milk.LactationId,
		&milk.AnimalId, &milk.AnimalNumber, &milk.AnimalOrder, &milk.AnimalName,
		&milk.PastureId, &milk.PastureName)
	return &milk, err
}

func (r *MilkRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.MilkEntry, err error) {
	var entries []entity.MilkEntry
	for rows.Next() {
		var milk entity.MilkEntry
		err = rows.Scan(&milk.Id, &milk.EntryDate, &milk.MilkQuantity, &milk.LactationId,
			&milk.AnimalId, &milk.AnimalNumber, &milk.AnimalOrder, &milk.AnimalName,
			&milk.PastureId, &milk.PastureName)
		if err != nil {
			return
		}
		entries = append(entries, milk)
	}
	return &entries, err
}

func (r *MilkRepository) saveOrUpdateScan(query string, model *entity.MilkEntry) error {
	return execQuery(query, model.Id, model.EntryDate, model.MilkQuantity,
		model.AnimalId, model.PastureId, model.LactationId, model.CreatedAt)
}

func (r *MilkRepository) FindPage(sort string, direction string, cursor string) (*entity.Page[entity.MilkEntry], error) {
	return r.Impl.FindPage(sort, direction, cursor)
}

func (r *MilkRepository) FindByAnimal(animalId string) (*[]entity.MilkEntry, error) {
	query := r.SelectQuery.Where("milk_entries.deleted_at is null and milk_entries.animal_id = $1")
	return r.Impl.FindListByQuery(query, animalId)
}

func (r *MilkRepository) FindByEntryDate(entryDate time.Time) (*[]entity.MilkEntry, error) {
	query := r.SelectQuery.Where("milk_entries.deleted_at is null and milk_entries.entry_date = $1")
	return r.Impl.FindListByQuery(query, entryDate)
}

func (r *MilkRepository) Add(newMilk *entity.MilkEntry) (*entity.MilkEntry, error) {
	return r.Impl.Add(newMilk)
}

func (r *MilkRepository) Save(milk *entity.MilkEntry) error {
	return r.Impl.Save(milk)
}

func (r *MilkRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
