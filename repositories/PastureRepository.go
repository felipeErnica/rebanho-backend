package repositories

import (
	"github.com/felipeErnica/rebanho-backend/entity"
)

type PastureRepository struct{}

func (p *PastureRepository) GetAll() (*[]entity.Pasture, error) {
	query := "SELECT * FROM pastures"
	sqlStatement, err := selectQueryList(query)
    defer sqlStatement.Close()
	if err != nil {
		return nil, err
	}
    
    var pastures []entity.Pasture

    for sqlStatement.Next() {
        var pasture entity.Pasture
        err:= sqlStatement.Scan(&pasture.Id, &pasture.Name, &pasture.BullId)
        if err != nil {
            return nil, err
        }
        pastures = append(pastures, pasture)
    }
    
    return &pastures, err
}

func (p *PastureRepository) GetById(id string) (*entity.Pasture, error) {
	query := "SELECT * FROM pastures WHERE id = $1"
    var pasture entity.Pasture
	sqlStatement:= selectQueryOne(query, id)
    err:= sqlStatement.Scan(&pasture.Id, &pasture.Name, &pasture.BullId)
    return &pasture, err
}

func (p *PastureRepository) Add(newPasture *entity.CreatePasture) (*entity.Pasture, error){
    query:="INSERT INTO pastures(id, name, bull_id) VALUES ($1, $2, $3)"
    pasture:= new(entity.Pasture).NewPasture(newPasture)
    err:= execQuery(query, pasture.Id, pasture.Name, pasture.BullId)
    return pasture, err
}

func (p *PastureRepository) Save(pasture *entity.Pasture) (*entity.Pasture, error){
    query:="UPDATE pastures SET name = $1, bull_id = $2 WHERE id = $3"
    err:= execQuery(query, pasture.Name, pasture.BullId, pasture.Id)
    return pasture, err
}
