package storage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"test_effective_mobile_task/internal/config"
	_ "github.com/lib/pq"
)

func InitDatabase(cfg config.DatabaseConfig) (*sql.DB,error) {
	
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			cfg.User,
		 	cfg.Password,
		  	cfg.Addr,
		   	cfg.DB,	
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