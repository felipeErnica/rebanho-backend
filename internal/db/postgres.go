package db

func ConnectPostgres() *DatabaseConfig {
    return &DatabaseConfig{
        Host: "localhost",
        Port: 5432,
        Password: "W4d7d1b7",
        DbName: "rebanho",
        User: "postgres",
    }
}
