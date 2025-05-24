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
