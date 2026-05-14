package storage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"github.com/joho/godotenv"
)

func InitDatabase() (*sql.DB,error) {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file:", "error:", err)
		return nil,err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	 os.Getenv("DB_HOST"),
     os.Getenv("DB_PORT"),
     os.Getenv("DB_USER"),
     os.Getenv("DB_PASSWORD"),
     os.Getenv("DB_NAME"),
	 os.Getenv("DB_SSLMODE"),
)
	database,err := sql.Open("postgres",dsn )
	  if err != nil {
   	   return nil, fmt.Errorf("invalid arguments for opening db: %w", err)
}
err = database.Ping()
if err != nil {
    slog.Error("Ошибка подключения к БД:", "error:", err)
}
    return database,nil
}