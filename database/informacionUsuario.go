package database

import (
	"fmt"
	"strconv"
)

type InformacionGeneral struct {
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
	telefono_p, err := strconv.Atoi(telefono)
	if err != nil {
		return InformacionGeneral{}, err
	}
	if telefono_p < 1_000_000_000 || telefono_p > 9_999_999_999 {
		return InformacionGeneral{}, fmt.Errorf("Telefono debe tener 10 digitos")
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
