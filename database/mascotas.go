package database

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
)

type Mascota struct {
	Id          int    `form:"id" json:"id"`
	Nombre      string `form:"nombre" json:"nombre"`
	Edad        int    `form:"edad" json:"edad"`
	Altura_cm   int    `form:"altura" json:"altura"`
	Foto64      string `form:"foto64" json:"foto64"`
	Descripcion string `form:"descripcion" json:"descripcion"`
}

type Mascotas []Mascota

func (r *DbRepository) GetAllMascotas() (Mascotas, error) {
	rows, err := r.Db.Query("SELECT * FROM mascotas")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var mascotas Mascotas
	for rows.Next() {
		var mascota Mascota
		err := rows.Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto64, &mascota.Descripcion)
		if err != nil {
			return nil, err
		}
		mascotas = append(mascotas, mascota)
	}
	if len(mascotas) == 0 {
		return Mascotas{}, nil
	}
	return mascotas, nil
}

func (r *DbRepository) GetMascotaById(id int) (Mascota, error) {
	var mascota Mascota
	err := r.Db.QueryRow("SELECT * FROM mascotas WHERE ID = ?", id).Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto64, &mascota.Descripcion)
	if err != nil {
		return Mascota{}, err
	}
	return mascota, nil
}

func (r *DbRepository) Get3Mascotas() (Mascotas, error) {
	rows, err := r.Db.Query("SELECT * FROM mascotas LIMIT 3")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var mascotas Mascotas
	for rows.Next() {
		var mascota Mascota
		err := rows.Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto64, &mascota.Descripcion)
		if err != nil {
			return nil, err
		}
		mascotas = append(mascotas, mascota)
	}
	if len(mascotas) == 0 {
		return Mascotas{}, nil
	}
	return mascotas, nil
}

func ValidarMascota(nombre, edad, altura_cm string, foto *multipart.FileHeader, descripcion string) (Mascota, error) {
	if nombre == "" {
		return Mascota{}, fmt.Errorf("Nombre no puede estar vacio")
	}
	edad_d, err := strconv.Atoi(edad)
	if err != nil {
		return Mascota{}, err
	}
	altura_cm_d, err := strconv.Atoi(altura_cm)
	if err != nil {
		return Mascota{}, err
	}
	src, err := foto.Open()
	if err != nil {
		return Mascota{}, err
	}
	defer src.Close()
	fotoBytes, err := io.ReadAll(src)
	if err != nil {
		return Mascota{}, err
	}
	if descripcion == "" {
		return Mascota{}, fmt.Errorf("Descripcion no puede estar vacio")
	}
	foto64 := base64.StdEncoding.EncodeToString(fotoBytes)
	return Mascota{Nombre: nombre, Edad: edad_d, Altura_cm: altura_cm_d, Foto64: foto64, Descripcion: descripcion}, nil
}

func (r *DbRepository) InsertMascota(mascota Mascota) (Mascota, error) {
	_, err := r.Db.Exec("INSERT INTO mascotas (nombre, edad, altura_cm, foto64, descripcion) VALUES (?, ?, ?, ?, ?)",
		mascota.Nombre, mascota.Edad, mascota.Altura_cm, mascota.Foto64, mascota.Descripcion)
	if err != nil {
		return Mascota{}, err
	}
	row := r.Db.QueryRow("SELECT * FROM mascotas ORDER BY ROWID DESC LIMIT 1")
	var mascota_g Mascota
	err = row.Scan(&mascota_g.Id, &mascota_g.Nombre, &mascota_g.Edad, &mascota_g.Altura_cm, &mascota_g.Foto64, &mascota_g.Descripcion)
	if err != nil {
		return Mascota{}, err
	}
	fmt.Printf("inserted id: %v\n", mascota_g.Id)
	return mascota_g, nil
}

func (r *DbRepository) DelMascota(id int) (Mascota, error) {
	mascota, err := r.GetMascotaById(id)
	if err != nil {
		return Mascota{}, err
	}
	_, err = r.Db.Exec("DELETE FROM mascotas WHERE id = ?", mascota.Id)
	if err != nil {
		return Mascota{}, err
	}
	return mascota, nil
}
