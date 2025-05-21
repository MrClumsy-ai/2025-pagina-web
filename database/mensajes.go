package database

import "fmt"

type Mensaje struct {
	Id      int    `form:"id" json:"id"`
	Nombre  string `form:"nombre" json:"nombre"`
	Email   string `form:"email" json:"email"`
	Mensaje string `form:"mensaje" json:"mensaje"`
}

type Mensajes []Mensaje

func (r *DbRepository) GetMensajeById(id int) (Mensaje, error) {
	var mensaje Mensaje
	err := r.Db.QueryRow("SELECT * FROM mensajes WHERE ID = ?", id).Scan(&mensaje.Id, &mensaje.Nombre, &mensaje.Email, &mensaje.Mensaje)
	if err != nil {
		return Mensaje{}, err
	}
	return mensaje, nil
}

func ValidarMensaje(nombre, email, mensaje string) (Mensaje, error) {
	if nombre == "" {
		return Mensaje{}, fmt.Errorf("Nombre no puede estar vacio")
	}
	if email == "" {
		return Mensaje{}, fmt.Errorf("Email no puede estar vacio")
	}
	if mensaje == "" {
		return Mensaje{}, fmt.Errorf("Mensaje no puede estar vacio")
	}
	return Mensaje{Nombre: nombre, Email: email, Mensaje: mensaje}, nil
}

func (r *DbRepository) InsertMensaje(mensaje Mensaje) (Mensaje, error) {
	_, err := r.Db.Exec("INSERT INTO mensajes (nombre, email, mensaje) VALUES (?, ?, ?)",
		mensaje.Nombre, mensaje.Email, mensaje.Mensaje)
	if err != nil {
		return Mensaje{}, err
	}
	row := r.Db.QueryRow("SELECT * FROM mensajes ORDER BY ROWID DESC LIMIT 1")
	var g_mensaje Mensaje
	err = row.Scan(&g_mensaje.Id, &g_mensaje.Nombre, &g_mensaje.Email, &g_mensaje.Mensaje)
	if err != nil {
		return Mensaje{}, err
	}
	fmt.Printf("inserted (%v): %v (%v): %v\n", g_mensaje.Id, g_mensaje.Nombre, g_mensaje.Email, g_mensaje.Mensaje)
	return g_mensaje, nil
}

func (r *DbRepository) GetAllMensajes() (Mensajes, error) {
	rows, err := r.Db.Query("SELECT * FROM mensajes")
	if err != nil {
		return Mensajes{}, err
	}
	defer rows.Close()
	var mensajes Mensajes
	for rows.Next() {
		var mensaje Mensaje
		err := rows.Scan(&mensaje.Id, &mensaje.Nombre, &mensaje.Email, &mensaje.Mensaje)
		if err != nil {
			return Mensajes{}, err
		}
		mensajes = append(mensajes, mensaje)
	}
	return mensajes, nil
}

func (r *DbRepository) DelMensaje(id int) (Mensaje, error) {
	mensaje, err := r.GetMensajeById(id)
	if err != nil {
		return Mensaje{}, err
	}
	_, err = r.Db.Exec("DELETE FROM mensajes WHERE id = ?", mensaje.Id)
	if err != nil {
		return Mensaje{}, err
	}
	return mensaje, nil
}
