package api

import (
	"fmt"
	"net/http"
	"proweb-backend/database"
	"strconv"

	"github.com/labstack/echo/v4"
)

func getRoot(c echo.Context) error {
	fmt.Printf("GET /\n")
	mascotas, err := repository.Get3Mascotas()
	if err != nil {
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusInternalServerError,
			"Message": err,
		}
		return c.Render(http.StatusInternalServerError, "error", response)
	}
	response := map[string]any{
		"URL":          _URL,
		"CurrentRoute": "/",
		"Mascotas":     mascotas,
	}
	return c.Render(http.StatusOK, "inicio", response)
}

func getAdopcion(c echo.Context) error {
	fmt.Printf("GET /adopcion\n")
	response := map[string]any{
		"URL":          _URL,
		"CurrentRoute": "/adopcion",
	}
	return c.Render(http.StatusOK, "adopcion", response)
}

func getAdopcionById(c echo.Context) error {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	fmt.Printf("GET /adopcion/%v\n", pId)
	mascota, err := repository.GetMascotaById(pId)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusNotFound,
			"Message": "Mascota no encontrada",
		}
		return c.Render(http.StatusNotFound, "error", response)
	}
	response := map[string]any{
		"URL":          _URL,
		"CurrentRoute": "/adopcion",
		"Mascota":      mascota,
	}
	return c.Render(http.StatusOK, "solicitudAdopcion", response)
}

func getContacto(c echo.Context) error {
	fmt.Printf("GET /contacto\n")
	response := map[string]any{
		"URL":          _URL,
		"CurrentRoute": "/contacto",
	}
	return c.Render(http.StatusOK, "contacto", response)
}

func getMascotas(c echo.Context) error {
	fmt.Printf("GET /mascotas\n")
	mascotas, err := repository.GetAllMascotas()
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusInternalServerError,
			"Message": "Error getting all mascotas",
		}
		return c.Render(http.StatusInternalServerError, "error", response)
	}
	response := map[string]any{
		"URL":          _URL,
		"CurrentRoute": "/mascotas",
		"Mascotas":     mascotas,
	}
	return c.Render(http.StatusOK, "mascotas", response)
}

// TODO: modificar esta ruta para otra tmpl renderizada
func getMascotaById(c echo.Context) error {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusInternalServerError,
			"Message": "Param id unprocessable",
		}
		return c.Render(http.StatusInternalServerError, "error", response)
	}
	fmt.Printf("GET /mascotas/%v\n", pId)
	mascota, err := repository.GetMascotaById(pId)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusNotFound,
			"Message": "Not found",
		}
		return c.Render(http.StatusNotFound, "error", response)
	}
	mascotas := database.Mascotas{mascota}
	response := map[string]any{
		"URL":          _URL,
		"CurrentRoute": "/mascotas",
		"Mascotas":     mascotas,
	}
	return c.Render(http.StatusOK, "mascotas", response)
}

func getRegistrar(c echo.Context) error {
	fmt.Printf("GET /registrar\n")
	response := map[string]any{
		"URL":          _URL,
		"CurrentRoute": "/registrar",
	}
	return c.Render(http.StatusOK, "registrar", response)
}

func postAdopcion(c echo.Context) error {
	pId, err := strconv.Atoi(c.Param("id"))
	fmt.Printf("POST /adopcion/%v\n", pId)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusInternalServerError,
			"Message": "Param id unprocessable",
		}
		return c.Render(http.StatusInternalServerError, "error", response)
	}
	mascota, err := repository.GetMascotaById(pId)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusNotFound,
			"Message": "Mascota no encontrada",
		}
		return c.Render(http.StatusNotFound, "error", response)
	}
	// info general
	nombre := c.FormValue("nombre")
	edad := c.FormValue("edad")
	email := c.FormValue("email")
	direccion := c.FormValue("direccion")
	telefono := c.FormValue("telefono")
	// info hogar
	tipoDeVivienda := c.FormValue("tipo-de-vivienda")
	esPropiedadVivienda := c.FormValue("is-propiedad-propia")
	personasEnHogar := c.FormValue("personas-en-hogar")
	tienePatio := c.FormValue("tiene-patio")
	// experiencia-mascotas
	mascotasAnteriormente := c.FormValue("mascotas-anteriormente")
	tieneMascotas := c.FormValue("tiene-mascotas")
	tieneVeterinario := c.FormValue("tiene-veterinario")
	// compromisos
	cuidadosNecesarios := c.FormValue("cuidados-necesarios")
	visitasSeguimiento := c.FormValue("visitas-seguimiento")
	responsabilidad := c.FormValue("responsabilidad")
	// confirmacion
	firma := c.FormValue("firma")
	// error checking
	infoGeneral, err := database.ValidarInfoGeneral(nombre, edad, email, direccion, telefono)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	infoHogar, err := database.ValidarInfoHogar(tipoDeVivienda, esPropiedadVivienda, tienePatio, personasEnHogar)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	experienciaMascotas, err := database.ValidarExperienciaMascotas(mascotasAnteriormente, tieneMascotas, tieneVeterinario)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	compromisos, err := database.ValidarCompromisos(cuidadosNecesarios, visitasSeguimiento, responsabilidad)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	confirmacion, err := database.ValidarConfirmacion(firma)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	solicitud, err := repository.InsertarSolicitud(pId, infoGeneral, infoHogar, experienciaMascotas, compromisos, confirmacion)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	response := map[string]any{
		"URL":       _URL,
		"Mascota":   mascota,
		"Solicitud": solicitud,
	}
	return c.Render(http.StatusOK, "solicitudEnviada", response)
}

func postContacto(c echo.Context) error {
	fmt.Printf("POST /contacto\n")
	nombre := c.FormValue("nombre")
	email := c.FormValue("email")
	mensaje := c.FormValue("mensaje")
	reqBody, err := database.ValidarMensaje(nombre, email, mensaje)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	g_mensaje, err := repository.InsertMensaje(reqBody)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusInternalServerError,
			"Message": err,
		}
		return c.Render(http.StatusInternalServerError, "error", response)
	}
	response := map[string]any{
		"URL":          _URL,
		"CurrentRoute": "/registrar",
		"Mensaje":      g_mensaje,
	}
	return c.Render(http.StatusCreated, "contacto", response)
}

func postMascota(c echo.Context) error {
	fmt.Printf("POST /mascotas\n")
	nombre := c.FormValue("nombre")
	edad := c.FormValue("edad")
	altura := c.FormValue("altura")
	descripcion := c.FormValue("descripcion")
	foto, err := c.FormFile("foto")
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	mascota, err := database.ValidarMascota(nombre, edad, altura, foto, descripcion)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	mascota_g, err := repository.InsertMascota(mascota)
	if err != nil {
		e.Logger.Error(err)
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusBadRequest,
			"Message": err,
		}
		return c.Render(http.StatusBadRequest, "error", response)
	}
	response := map[string]any{
		"URL":          _URL,
		"CurrentRoute": "/registrar",
		"Mascota":      mascota_g,
	}
	return c.Render(http.StatusCreated, "registrado", response)
}

// ######################### JSON #########################
// ------------------------- GET -------------------------
func jsonGetContacto(c echo.Context) error {
	fmt.Printf("GET /json/contacto\n")
	mensajes, err := repository.GetAllMensajes()
	if err != nil {
		e.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
	}
	if len(mensajes) == 0 {
		return c.JSON(http.StatusNotFound, "Mensajes not found")
	}
	return c.JSON(http.StatusFound, mensajes)
}

func jsonGetMascotas(c echo.Context) error {
	fmt.Printf("GET /json/mascotas\n")
	mascotas, err := repository.GetAllMascotas()
	if err != nil {
		e.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	if len(mascotas) == 0 {
		return c.JSON(http.StatusNotFound, "Mascotas not found")
	}
	return c.JSON(http.StatusFound, mascotas)
}

func jsonGetMascotaById(c echo.Context) error {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		e.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	e.Logger.Printf("GET /json/mascotas/%v", pId)
	mascota, err := repository.GetMascotaById(pId)
	if err != nil {
		e.Logger.Error(err)
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusFound, mascota)
}

func jsonDelContacto(c echo.Context) error {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		e.Logger.Error(err)
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	fmt.Printf("DELETE /json/contacto/%v\n", pId)
	mensaje, err := repository.DelMensaje(pId)
	if err != nil {
		e.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, mensaje)
}

func jsonDelMascota(c echo.Context) error {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		e.Logger.Error(err)
		return c.JSON(http.StatusUnprocessableEntity, "param id no procesable")
	}
	fmt.Printf("DELETE /json/mascotas/%v\n", pId)
	mascota, err := repository.DelMascota(pId)
	return c.JSON(http.StatusOK, mascota)
}
