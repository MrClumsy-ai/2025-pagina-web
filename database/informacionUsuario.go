package database

import (
	"fmt"
	"strconv"
)

type InformacionGeneral struct {
	Id        int    `form:"id" json:"id"`
	IdMascota int    `form:"idMascota" json:"idMascota"`
	Nombre    string `form:"nombre" json:"nombre"`
	Edad      int    `form:"edad" json:"edad"`
	Email     string `form:"email" json:"email"`
	Direccion string `form:"direccion" json:"direccion"`
	Telefono  int    `form:"telefono" json:"telefono"`
}

type InformacionHogar struct {
	TipoDeVivienda      string `form:"tipo-de-vivienda" json:"tipo-de-vivienda"`
	EsPropiedadVivienda bool   `form:"is-propiedad-vivienda" json:"is-propiedad-vivienda"`
	TienePatio          bool   `form:"tiene-patio" json:"tiene-patio"`
	PersonasEnHogar     int    `form:"personas-en-hogar" json:"personas-en-hogar"`
}

type ExperienciaMascotas struct {
	MascotasAnteriormente bool `form:"mascotas-anteriormente" json:"mascotas-anteriormente"`
	TieneMascotas         bool `form:"tiene-mascotas" json:"tiene-mascotas"`
	TieneVeterinario      bool `form:"tiene-veterinario" json:"tiene-veterinario"`
}

type Compromisos struct {
	CuidadosNecesarios bool `form:"cuidados-necesarios" json:"cuidados-necesarios"`
	VisitasSeguimiento bool `form:"visitas-seguimiento" json:"visitas-seguimiento"`
	Responsabilidad    bool `form:"responsabilidad" json:"responsabilidad"`
}

type Confirmacion struct {
	Firma string `form:"firma" json:"firma"`
}

type Solicitud struct {
	InformacionGeneral
	InformacionHogar
	ExperienciaMascotas
	Compromisos
	Confirmacion
}

type Solicitudes []Solicitud

func ValidarInfoGeneral(nombre, edad, email, direccion, telefono string) (InformacionGeneral, error) {
	if nombre == "" {
		return InformacionGeneral{}, fmt.Errorf("Nombre no solicitado")
	}
	edad_p, err := strconv.Atoi(edad)
	if err != nil {
		return InformacionGeneral{}, fmt.Errorf("Edad no procesable")
	}
	if edad_p < 18 || edad_p > 100 {
		return InformacionGeneral{}, fmt.Errorf("Edad no es valida")
	}
	if email == "" {
		return InformacionGeneral{}, fmt.Errorf("Email no solicitado")
	}
	if direccion == "" {
		return InformacionGeneral{}, fmt.Errorf("Direccion no solicitada")
	}
	if len(telefono) != 10 {
		return InformacionGeneral{}, fmt.Errorf("Telefono debe tener 10 digitos")
	}
	telefono_p, err := strconv.Atoi(telefono)
	if err != nil {
		return InformacionGeneral{}, fmt.Errorf("Telefono no procesable")
	}
	return InformacionGeneral{
		Nombre:    nombre,
		Edad:      edad_p,
		Email:     email,
		Direccion: direccion,
		Telefono:  telefono_p,
	}, nil
}

func ValidarInfoHogar(tipoDeVivienda, esPropiedadVivienda, tienePatio, personasEnHogar string) (InformacionHogar, error) {
	personasEnHogar_d, err := strconv.Atoi(personasEnHogar)
	if err != nil {
		return InformacionHogar{}, fmt.Errorf("Personas en hogar no procesable")
	}
	// error handling
	if tipoDeVivienda == "" {
		return InformacionHogar{}, fmt.Errorf("Tipo de vivienda sin texto")
	}
	var esPropiedadVivienda_d bool
	switch esPropiedadVivienda {
	case "Si":
		esPropiedadVivienda_d = true
	case "No":
		esPropiedadVivienda_d = false
	default:
		return InformacionHogar{}, fmt.Errorf("Es propiedad vivienda sin respuesta")
	}
	var tienePatio_d bool
	switch tienePatio {
	case "Si":
		tienePatio_d = true
	case "No":
		tienePatio_d = false
	default:
		return InformacionHogar{}, fmt.Errorf("Tiene patio respuesta no procesable")
	}
	if personasEnHogar_d < 1 {
		return InformacionHogar{}, fmt.Errorf("Personas en hogar no puede ser menor a 1")
	}
	return InformacionHogar{
		TipoDeVivienda:      tipoDeVivienda,
		EsPropiedadVivienda: esPropiedadVivienda_d,
		TienePatio:          tienePatio_d,
		PersonasEnHogar:     personasEnHogar_d,
	}, nil
}

func ValidarExperienciaMascotas(mascotasAnteriormente, tieneMascotas, tieneVeterinario string) (ExperienciaMascotas, error) {
	var mascotasAnteriormente_d bool
	switch mascotasAnteriormente {
	case "Si":
		mascotasAnteriormente_d = true
	case "No":
		mascotasAnteriormente_d = false
	default:
		return ExperienciaMascotas{}, fmt.Errorf("mascotas anteriormente respuesta no procesable")
	}
	var tieneMascotas_d bool
	switch tieneMascotas {
	case "Si":
		tieneMascotas_d = true
	case "No":
		tieneMascotas_d = false
	default:
		return ExperienciaMascotas{}, fmt.Errorf("tiene mascotas respuesta no procesable")
	}
	var tieneVeterinario_d bool
	switch tieneVeterinario {
	case "Si":
		tieneVeterinario_d = true
	case "No":
		tieneVeterinario_d = false
	default:
		return ExperienciaMascotas{}, fmt.Errorf("tiene veterinario respuesta no procesable")
	}
	return ExperienciaMascotas{MascotasAnteriormente: mascotasAnteriormente_d, TieneMascotas: tieneMascotas_d, TieneVeterinario: tieneVeterinario_d}, nil
}

func ValidarCompromisos(cuidadosNecesarios, visitasSeguimiento, responsabilidad string) (Compromisos, error) {
	var cuidadosNecesarios_d bool
	switch cuidadosNecesarios {
	case "Si":
		cuidadosNecesarios_d = true
	case "No":
		cuidadosNecesarios_d = false
	default:
		return Compromisos{}, fmt.Errorf("tiene veterinario respuesta no procesable")
	}
	var visitasSeguimiento_d bool
	switch visitasSeguimiento {
	case "Si":
		visitasSeguimiento_d = true
	case "No":
		visitasSeguimiento_d = false
	default:
		return Compromisos{}, fmt.Errorf("tiene veterinario respuesta no procesable")
	}
	var responsabilidad_d bool
	switch responsabilidad {
	case "Si":
		responsabilidad_d = true
	case "No":
		responsabilidad_d = false
	default:
		return Compromisos{}, fmt.Errorf("tiene veterinario respuesta no procesable")
	}
	return Compromisos{CuidadosNecesarios: cuidadosNecesarios_d, VisitasSeguimiento: visitasSeguimiento_d, Responsabilidad: responsabilidad_d}, nil
}

func ValidarConfirmacion(firma string) (Confirmacion, error) {
	if firma == "" {
		return Confirmacion{}, fmt.Errorf("Firma debe ser llenada")
	}
	return Confirmacion{Firma: firma}, nil
}

func (r *DbRepository) InsertarSolicitud(idMascota int, infoGeneral InformacionGeneral, infoHogar InformacionHogar, expMascotas ExperienciaMascotas, compromisos Compromisos, confirmacion Confirmacion) (Solicitud, error) {
	_, err := r.Db.Exec(`
		INSERT INTO solicitudes (
			mascotaId,
			nombre, edad, direccion, email, telefono,
			tipoDeVivienda, esPropiedadVivienda, personasEnHogar, tienePatio,
			mascotasAnteriormente, tieneMascotas, tieneVeterinario,
			cuidadosNecesarios, visitasSeguimiento, responsabilidad,
			firma)
		VALUES (
			?,
			?,?,?,?,?,
			?,?,?,?,
			?,?,?,
			?,?,?,
			?)
		`,
		idMascota,
		infoGeneral.Nombre, infoGeneral.Edad, infoGeneral.Direccion, infoGeneral.Email, infoGeneral.Telefono,
		infoHogar.TipoDeVivienda, infoHogar.EsPropiedadVivienda, infoHogar.PersonasEnHogar, infoHogar.TienePatio,
		expMascotas.MascotasAnteriormente, expMascotas.TieneMascotas, expMascotas.TieneVeterinario,
		compromisos.CuidadosNecesarios, compromisos.VisitasSeguimiento, compromisos.Responsabilidad,
		confirmacion.Firma)
	if err != nil {
		return Solicitud{}, err
	}
	row := r.Db.QueryRow("SELECT * FROM solicitudes ORDER BY ROWID DESC LIMIT 1")
	var solicitud Solicitud
	err = row.Scan(
		&solicitud.Id,
		&solicitud.IdMascota,
		&solicitud.Nombre,
		&solicitud.Edad,
		&solicitud.Direccion,
		&solicitud.Email,
		&solicitud.Telefono,

		&solicitud.TipoDeVivienda,
		&solicitud.EsPropiedadVivienda,
		&solicitud.PersonasEnHogar,
		&solicitud.TienePatio,

		&solicitud.MascotasAnteriormente,
		&solicitud.TieneMascotas,
		&solicitud.TieneVeterinario,

		&solicitud.CuidadosNecesarios,
		&solicitud.VisitasSeguimiento,
		&solicitud.Responsabilidad,

		&solicitud.Firma)
	if err != nil {
		return Solicitud{}, err
	}
	fmt.Printf("inserted id: %v\n", solicitud.Id)
	return solicitud, nil
}

func (r *DbRepository) GetAllSolicitudes() (Solicitudes, error) {
	rows, err := r.Db.Query("SELECT * FROM solicitudes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var solicitudes Solicitudes
	for rows.Next() {
		var solicitud Solicitud
		err := rows.Scan(
			&solicitud.Id,
			&solicitud.IdMascota,
			&solicitud.Nombre,
			&solicitud.Edad,
			&solicitud.Direccion,
			&solicitud.Email,
			&solicitud.Telefono,

			&solicitud.TipoDeVivienda,
			&solicitud.EsPropiedadVivienda,
			&solicitud.PersonasEnHogar,
			&solicitud.TienePatio,

			&solicitud.MascotasAnteriormente,
			&solicitud.TieneMascotas,
			&solicitud.TieneVeterinario,

			&solicitud.CuidadosNecesarios,
			&solicitud.VisitasSeguimiento,
			&solicitud.Responsabilidad,

			&solicitud.Firma)
		if err != nil {
			return nil, err
		}
		solicitudes = append(solicitudes, solicitud)
	}
	if len(solicitudes) == 0 {
		return Solicitudes{}, nil
	}
	return solicitudes, nil
}

func (r *DbRepository) GetSolicitudById(id int) (Solicitud, error) {
	var solicitud Solicitud
	err := r.Db.QueryRow("SELECT * FROM solicitudes WHERE ID = ?", id).Scan(
		&solicitud.Id,
		&solicitud.IdMascota,
		&solicitud.Nombre,
		&solicitud.Edad,
		&solicitud.Direccion,
		&solicitud.Email,
		&solicitud.Telefono,

		&solicitud.TipoDeVivienda,
		&solicitud.EsPropiedadVivienda,
		&solicitud.PersonasEnHogar,
		&solicitud.TienePatio,

		&solicitud.MascotasAnteriormente,
		&solicitud.TieneMascotas,
		&solicitud.TieneVeterinario,

		&solicitud.CuidadosNecesarios,
		&solicitud.VisitasSeguimiento,
		&solicitud.Responsabilidad,

		&solicitud.Firma)
	if err != nil {
		return Solicitud{}, err
	}
	return solicitud, nil
}

func (r *DbRepository) DelSolicitud(id int) (Solicitud, error) {
	solicitud, err := r.GetSolicitudById(id)
	if err != nil {
		return Solicitud{}, err
	}
	fmt.Println(solicitud)
	_, err = r.Db.Exec("DELETE FROM solicitudes WHERE id = ?", solicitud.Id)
	if err != nil {
		return Solicitud{}, err
	}
	fmt.Println(solicitud.Id)
	return solicitud, nil
}
