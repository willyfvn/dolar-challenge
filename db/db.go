package db

import (
	"database/sql"
	"fmt"
)

func StartDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./mybanco.db")
	if err != nil {
		panic(err)
	}
	fmt.Println("Banco de dados conectado com sucesso")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY, bid TEXT)")
	if err != nil {
		panic(err)
	}

	return db
}

func InsertCotacao(db *sql.DB, cotacao string) error {
	_, err := db.Exec("INSERT INTO cotacao (bid) VALUES (?)", cotacao)
	if err != nil {
		return err
	}
	fmt.Println("Cotação inserida com sucesso")
	return nil
}
