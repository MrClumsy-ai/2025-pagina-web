package database

import "database/sql"

type MascotaRepository struct {
	Db *sql.DB
}

type Mascota struct {
	Id          int
	Nombre      string
	Edad        int
	Altura_cm   int
	Foto        []byte
	Foto64      string
	Descripcion string
}

func (r *MascotaRepository) CreateTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS mascotas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT,
		edad INTEGER,
		altura_cm INT,
		foto BLOB,
		descripcion TEXT
	)`)
	return err
}

func (r *MascotaRepository) Insert(mascota Mascota) error {
	_, err := r.Db.Exec("INSERT INTO mascotas (nombre, edad, altura, foto, descripcion) VALUES (?, ?, ?, ?, ?)",
		mascota.Nombre, mascota.Edad, mascota.Altura_cm, mascota.Foto, mascota.Descripcion)
	return err
}

func (r *MascotaRepository) GetAll() ([]Mascota, error) {
	rows, err := r.Db.Query("select * from mascotas")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var mascotas []Mascota
	for rows.Next() {
		var mascota Mascota
		err := rows.Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto, &mascota.Descripcion)
		if err != nil {
			return nil, err
		}
		mascotas = append(mascotas, mascota)
	}
	return mascotas, nil
}

func (r *MascotaRepository) GetById(id int) (Mascota, error) {
	var mascota Mascota
	err := r.Db.QueryRow("SELECT * FROM mascotas WHERE ID = ?", id).Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto, &mascota.Descripcion)
	if err != nil {
		return Mascota{}, err
	}
	return mascota, nil
}

func (r *MascotaRepository) Update(mascota Mascota) error {
	_, err := r.Db.Exec("UPDATE mascotas SET nombre = ?, edad = ?, altura_cm = ?, foto = ?, descripcion = ? WHERE id = ?",
		&mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto, &mascota.Descripcion, &mascota.Id)
	return err
}

func (r *MascotaRepository) Delete(id int) error {
	_, err := r.Db.Exec("DELETE FROM mascotas WHERE id = ?", id)
	return err
}
