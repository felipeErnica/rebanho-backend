package db

func ConnectPostgres() *DatabaseConfig {
    return &DatabaseConfig{
        Host: "localhost",
        Port: 5431,
        Password: "W4d7d1b7",
        DbName: "rebanho",
        User: "postgres",
    }
}
