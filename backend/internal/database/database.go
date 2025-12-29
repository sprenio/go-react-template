package database

import (
	"database/sql"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

var dsn string

func SetDSN(dataSourceName string) {
	dsn = dataSourceName
}
func GetDB() (*sql.DB, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(25)            // max 25 jednoczesnych połączeń
	conn.SetMaxIdleConns(15)            // trzymaj do 15 połączeń w idle
	conn.SetConnMaxLifetime(30 * time.Minute) // restartuj połączenia co 30 min
	conn.SetConnMaxIdleTime(5 * time.Minute)  // idle połączenia max 5 min
	return conn, nil
}

