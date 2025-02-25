package repositories

import (
	"database/sql"

	"github.com/felipeErnica/rebanho-backend/util"
)

type Repositoty interface {
    InitRepository(db *sql.DB)
}


func LogInitRepository(name string) {
	util.LogInfo("O Repositório " + name + " foi iniciado com sucesso!")
}
