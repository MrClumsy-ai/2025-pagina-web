package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DbRepository struct {
	Db *sql.DB
}

func (r *DbRepository) Connect() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS mascotas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		edad INTEGER,
		altura_cm INTEGER,
		foto64 BLOB,
		descripcion TEXT)`)
	if err != nil {
		return err
	}
	_, err = r.Db.Exec(`CREATE TABLE IF NOT EXISTS mensajes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		email TEXT NOT NULL,
		mensaje TEXT NOT NULL)`)
	if err != nil {
		return err
	}
	return nil
}
