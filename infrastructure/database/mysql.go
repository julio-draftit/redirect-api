package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func NewMysql() (*sql.DB, error) {
	con := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"u165437626_admin",
		"Jacare123@",
		"redirectbussines.com",
		"u165437626_redirect",
	)

	db, err := sql.Open("mysql", con)
	if err != nil {
		db.Close()
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
