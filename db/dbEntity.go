package db

import "fmt"

type DatabaseConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DbName   string
}

func (db *DatabaseConfig) ReturnDatabaseInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", db.Host, db.Port, db.User, db.Password, db.DbName)
}
