package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

const PORT = ":8000"
const URL = "http://localhost" + PORT

type Templates struct {
	templates *template.Template
}

type Mascota struct {
	Id          int    `form:"id" json:"id"`
	Nombre      string `form:"nombre" json:"nombre"`
	Edad        int    `form:"edad" json:"edad"`
	Altura_cm   int    `form:"altura" json:"altura"`
	Foto64      string `form:"foto64" json:"foto64"`
	Descripcion string `form:"descripcion" json:"descripcion"`
}

type Mensaje struct {
	Id      int    `form:"id" json:"id"`
	Nombre  string `form:"nombre" json:"nombre"`
	Email   string `form:"email" json:"email"`
	Mensaje string `form:"mensaje" json:"mensaje"`
}

type InformacionGeneral struct {
	Nombre    string `form:"nombre" json:"nombre"`
	Edad      int    `form:"edad" json:"edad"`
	Email     string `form:"email" json:"email"`
	Direccion string `form:"direccion" json:"direccion"`
	Telefono  int    `form:"telefono" json:"telefono"`
}

type Mascotas []Mascota

func (t *Templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	tmpl := template.Must(template.ParseGlob("views/*.tmpl"))
	tmpl = template.Must(tmpl.ParseGlob("views/partials/*.tmpl"))
	return &Templates{
		templates: tmpl,
	}
}

func main() {
	// init
	e := echo.New()
	e.Static("/assets", "assets")
	e.Renderer = newTemplate()
	const dbLocation = "pagina.db"
	dbConnection, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		e.Logger.Fatal("Error connecting to %v: %v", dbLocation, err)
	}
	fmt.Printf("Database connection established: %v\n", dbLocation)
	defer dbConnection.Close()
	_, err = dbConnection.Exec(`CREATE TABLE IF NOT EXISTS mascotas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		edad INTEGER,
		altura_cm INTEGER,
		foto64 BLOB,
		descripcion TEXT)`)
	if err != nil {
		e.Logger.Fatal("Error creating mascota table", err)
	}
	_, err = dbConnection.Exec(`CREATE TABLE IF NOT EXISTS mensajes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		email TEXT NOT NULL,
		mensaje TEXT NOT NULL)`)
	if err != nil {
		e.Logger.Fatal("Error creating mensajes table", err)
	}

	getAllMascotas := func() (Mascotas, error) {
		rows, err := dbConnection.Query("SELECT * FROM mascotas")
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

	getMascotaById := func(id int) (Mascota, error) {
		var mascota Mascota
		err = dbConnection.QueryRow("SELECT * FROM mascotas WHERE ID = ?", id).Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto64, &mascota.Descripcion)
		if err != nil {
			e.Logger.Error(err)
			return Mascota{}, err
		}
		return mascota, nil
	}

	getMensajeById := func(id int) (Mensaje, error) {
		var mensaje Mensaje
		err = dbConnection.QueryRow("SELECT * FROM mensajes WHERE ID = ?", id).Scan(&mensaje.Id, &mensaje.Nombre, &mensaje.Email, &mensaje.Mensaje)
		if err != nil {
			e.Logger.Error(err)
			return Mensaje{}, err
		}
		return mensaje, nil
	}
	e.GET("/*", func(c echo.Context) error {
		fmt.Printf("GET catchall\n")
		response := map[string]any{
			"URL":     URL,
			"Code":    http.StatusNotFound,
			"Message": "Not found",
		}
		return c.Render(http.StatusNotFound, "error", response)
	})

	// ------------------------- GET -------------------------
	e.GET("/", func(c echo.Context) error {
		fmt.Printf("GET /\n")
		rows, err := dbConnection.Query("SELECT * FROM mascotas LIMIT 3")
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Internal server error",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		defer rows.Close()
		var mascotas Mascotas
		for rows.Next() {
			var mascota Mascota
			err := rows.Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto64, &mascota.Descripcion)
			if err != nil {
				e.Logger.Error(err)
				response := map[string]any{
					"URL":     URL,
					"Code":    http.StatusInternalServerError,
					"Message": "Internal server error",
				}
				return c.Render(http.StatusInternalServerError, "error", response)
			}
			mascotas = append(mascotas, mascota)
		}
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/",
			"Mascotas":     mascotas,
		}
		return c.Render(http.StatusOK, "inicio", response)
	})

	e.GET("/adopcion", func(c echo.Context) error {
		fmt.Printf("GET /adopcion\n")
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/adopcion",
		}
		return c.Render(http.StatusOK, "adopcion", response)
	})

	e.GET("/adopcion/:id", func(c echo.Context) error {
		pId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		fmt.Printf("GET /adopcion/%v\n", pId)
		mascota, err := getMascotaById(pId)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusNotFound,
				"Message": "Mascota no encontrada",
			}
			return c.Render(http.StatusNotFound, "error", response)
		}
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/adopcion",
			"Mascota":      mascota,
		}
		return c.Render(http.StatusOK, "solicitudAdopcion", response)
	})

	e.GET("/contacto", func(c echo.Context) error {
		fmt.Printf("GET /contacto\n")
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/contacto",
		}
		return c.Render(http.StatusOK, "contacto", response)
	})

	e.GET("/mascotas", func(c echo.Context) error {
		fmt.Printf("GET /mascotas\n")
		mascotas, err := getAllMascotas()
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Error getting all mascotas",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/mascotas",
			"Mascotas":     mascotas,
		}
		return c.Render(http.StatusOK, "mascotas", response)
	})

	e.GET("/mascotas/:id", func(c echo.Context) error {
		pId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Param id unprocessable",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		fmt.Printf("GET /mascotas/%v\n", pId)
		mascota, err := getMascotaById(pId)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusNotFound,
				"Message": "Not found",
			}
			return c.Render(http.StatusNotFound, "error", response)
		}
		mascotas := Mascotas{mascota}
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/mascotas",
			"Mascotas":     mascotas,
		}
		return c.Render(http.StatusOK, "mascotas", response)
	})

	e.GET("/registrar", func(c echo.Context) error {
		fmt.Printf("GET /registrar\n")
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/registrar",
		}
		return c.Render(http.StatusOK, "registrar", response)
	})

	// ------------------------- POST -------------------------
	e.POST("/adopcion/info-general/:id", func(c echo.Context) error {
		pId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Param id unprocessable",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		fmt.Printf("POST /adopcion/info-general/%v\n", pId)
		nombre := c.FormValue("nombre")
		if nombre == "" {
			e.Logger.Errorf("nombre: %v", nombre)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnauthorized,
				"Message": "nombre no solicitado",
			}
			return c.Render(http.StatusUnauthorized, "error", response)
		}
		edad, err := strconv.Atoi(c.FormValue("edad"))
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnprocessableEntity,
				"Message": "Edad no procesable",
			}
			return c.Render(http.StatusUnprocessableEntity, "error", response)
		}
		if edad < 18 || edad > 100 {
			e.Logger.Errorf("edad: %v", edad)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnauthorized,
				"Message": "Edad no es valida",
			}
			return c.Render(http.StatusUnauthorized, "error", response)
		}
		email := c.FormValue("email")
		if email == "" {
			e.Logger.Errorf("email: %v", email)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnauthorized,
				"Message": "Email no solicitado",
			}
			return c.Render(http.StatusUnauthorized, "error", response)
		}
		direccion := c.FormValue("direccion")
		if direccion == "" {
			e.Logger.Errorf("direccion: %v", direccion)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnauthorized,
				"Message": "Direccion no solicitada",
			}
			return c.Render(http.StatusUnauthorized, "error", response)
		}
		telefono, err := strconv.Atoi(c.FormValue("telefono"))
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnprocessableEntity,
				"Message": "telefono no procesable",
			}
			return c.Render(http.StatusUnprocessableEntity, "error", response)
		}
		e.Logger.Printf("telefono: %v", telefono)
		if telefono < 1_000_000_000 || telefono > 9_999_999_999 {
			e.Logger.Errorf("telefono: %v", telefono)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnauthorized,
				"Message": "Telefono debe tener 10 digitos",
			}
			return c.Render(http.StatusUnauthorized, "error", response)
		}
		informacionGeneral := InformacionGeneral{
			Nombre:    nombre,
			Edad:      edad,
			Email:     email,
			Direccion: direccion,
			Telefono:  telefono,
		}
		_, err = getMascotaById(pId)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusNotFound,
				"Message": "Mascota no encontrada",
			}
			return c.Render(http.StatusNotFound, "error", response)
		}
		response := map[string]any{
			"URL":                URL,
			"InformacionGeneral": informacionGeneral,
			"MascotaId":          pId,
		}
		return c.Render(http.StatusOK, "informacionHogar", response)
	})

	// TODO: POST para /adopcion/info-hogar/:id (desde informacionHogar.tmpl)
	e.POST("/adopcion/info-hogar/:id", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "/adopcion/info-hogar/:id")
	})

	e.POST("/contacto", func(c echo.Context) error {
		fmt.Printf("POST /contacto\n")
		nombre := c.FormValue("nombre")
		email := c.FormValue("email")
		mensaje := c.FormValue("mensaje")
		reqBody := Mensaje{
			Nombre:  nombre,
			Email:   email,
			Mensaje: mensaje,
		}
		if err := c.Bind(&reqBody); err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusBadRequest,
				"Message": "Bad request",
			}
			return c.Render(http.StatusBadRequest, "error", response)
		}
		if reqBody.Nombre == "" {
			e.Logger.Errorf("reqBody.Nombre: %v", reqBody.Nombre)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusBadRequest,
				"Message": "Ningun nombre dado",
			}
			return c.Render(http.StatusBadRequest, "error", response)
		}
		if reqBody.Email == "" {
			e.Logger.Errorf("reqBody.Email: %v", reqBody.Email)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusBadRequest,
				"Message": "Ningun email dado",
			}
			return c.Render(http.StatusBadRequest, "error", response)
		}
		if reqBody.Mensaje == "" {
			e.Logger.Errorf("reqBody.Mensaje: %v", reqBody.Mensaje)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusBadRequest,
				"Message": "Ningun mensaje dado",
			}
			return c.Render(http.StatusBadRequest, "error", response)
		}
		_, err = dbConnection.Exec("INSERT INTO mensajes (nombre, email, mensaje) VALUES (?, ?, ?)",
			reqBody.Nombre, reqBody.Email, reqBody.Mensaje)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Error inserting into DB",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		row := dbConnection.QueryRow("SELECT * FROM mensajes ORDER BY ROWID DESC LIMIT 1")
		var g_mensaje Mensaje
		err = row.Scan(&g_mensaje.Id, &g_mensaje.Nombre, &g_mensaje.Email, &g_mensaje.Mensaje)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Error retrieving row",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		fmt.Printf("inserted (%v): %v (%v): %v\n", g_mensaje.Id, g_mensaje.Nombre, g_mensaje.Email, g_mensaje.Mensaje)
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/registrar",
			"Mensaje":      g_mensaje,
		}
		return c.Render(http.StatusCreated, "contacto", response)
	})

	e.POST("/mascotas", func(c echo.Context) error {
		fmt.Printf("POST /mascotas\n")
		nombre := c.FormValue("nombre")
		edad, err := strconv.Atoi(c.FormValue("edad"))
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnprocessableEntity,
				"Message": "Edad no es procesable",
			}
			return c.Render(http.StatusUnprocessableEntity, "error", response)
		}
		altura, err := strconv.Atoi(c.FormValue("altura"))
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnprocessableEntity,
				"Message": "Altura no es procesable",
			}
			return c.Render(http.StatusUnprocessableEntity, "error", response)
		}
		descripcion := c.FormValue("descripcion")
		foto, err := c.FormFile("foto")
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnprocessableEntity,
				"Message": "Foto no es procesable",
			}
			return c.Render(http.StatusUnprocessableEntity, "error", response)
		}
		src, err := foto.Open()
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnprocessableEntity,
				"Message": "Foto no pudo ser abierta",
			}
			return c.Render(http.StatusUnprocessableEntity, "error", response)
		}
		defer src.Close()
		fotoBytes, err := io.ReadAll(src)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnprocessableEntity,
				"Message": "Foto no pudo ser convertida a bytes",
			}
			return c.Render(http.StatusUnprocessableEntity, "error", response)
		}
		foto64 := base64.StdEncoding.EncodeToString(fotoBytes)
		reqBody := Mascota{
			Nombre:      nombre,
			Edad:        edad,
			Altura_cm:   altura,
			Descripcion: descripcion,
			Foto64:      foto64,
		}
		if err := c.Bind(&reqBody); err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusBadRequest,
				"Message": "Bad request",
			}
			return c.Render(http.StatusBadRequest, "error", response)
		}
		if reqBody.Nombre == "" {
			e.Logger.Errorf("reqBody.Nombre: %v", reqBody.Nombre)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusBadRequest,
				"Message": "Ningun nombre dado",
			}
			return c.Render(http.StatusBadRequest, "error", response)
		}
		_, err = dbConnection.Exec("INSERT INTO mascotas (nombre, edad, altura_cm, foto64, descripcion) VALUES (?, ?, ?, ?, ?)",
			reqBody.Nombre, reqBody.Edad, reqBody.Altura_cm, reqBody.Foto64, reqBody.Descripcion)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Error inserting into DB",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		fmt.Printf("inserted into db: %v\n", reqBody.Nombre)
		row := dbConnection.QueryRow("SELECT * FROM mascotas ORDER BY ROWID DESC LIMIT 1")
		var mascota Mascota
		err = row.Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto64, &mascota.Descripcion)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Error retrieving row",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		fmt.Printf("inserted id: %v\n", mascota.Id)
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/registrar",
			"Mascota":      mascota,
		}
		return c.Render(http.StatusCreated, "registrado", response)
	})

	// ######################### JSON #########################
	// ------------------------- GET -------------------------
	e.GET("/json/contacto", func(c echo.Context) error {
		fmt.Printf("GET /json/contacto\n")
		rows, err := dbConnection.Query("SELECT * FROM mensajes")
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "error selecting from mensajes",
			}
			return c.JSON(http.StatusInternalServerError, response)
		}
		defer rows.Close()
		var mensajes []Mensaje
		for rows.Next() {
			var mensaje Mensaje
			err := rows.Scan(&mensaje.Id, &mensaje.Nombre, &mensaje.Email, &mensaje.Mensaje)
			if err != nil {
				e.Logger.Error(err)
				response := map[string]any{
					"URL":     URL,
					"Code":    http.StatusInternalServerError,
					"Message": "Error scanning row",
				}
				return c.JSON(http.StatusInternalServerError, response)
			}
			mensajes = append(mensajes, mensaje)
		}
		if len(mensajes) == 0 {
			e.Logger.Errorf("len(mensajes): %v", len(mensajes))
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusNotFound,
				"Message": "Not found",
			}
			return c.JSON(http.StatusNotFound, response)
		}
		return c.JSON(http.StatusFound, mensajes)
	})

	e.GET("/json/mascotas", func(c echo.Context) error {
		fmt.Printf("GET /json/mascotas\n")
		mascotas, err := getAllMascotas()
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		if len(mascotas) == 0 {
			e.Logger.Errorf("len(mascotas): %v", len(mascotas))
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusFound, mascotas)
	})

	e.GET("/json/mascotas/:id", func(c echo.Context) error {
		pId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		e.Logger.Printf("GET /json/mascotas/%v", pId)
		mascota, err := getMascotaById(pId)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusNotFound, "mascota not found")
		}
		return c.JSON(http.StatusFound, mascota)
	})

	// ------------------------- DELETE -------------------------
	e.DELETE("/json/contacto/:id", func(c echo.Context) error {
		pId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusUnprocessableEntity, "param id no procesable")
		}
		fmt.Printf("DELETE /json/contacto/%v\n", pId)
		mensaje, err := getMensajeById(pId)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusNotFound, "mensaje not found")
		}
		_, err = dbConnection.Exec("DELETE FROM mensajes WHERE id = ?", mensaje.Id)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, "mensaje no se pudo borrar")
		}
		return c.JSON(http.StatusOK, mensaje)
	})

	e.DELETE("/json/mascotas/:id", func(c echo.Context) error {
		pId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusUnprocessableEntity, "param id no procesable")
		}
		fmt.Printf("DELETE /json/mascotas/%v\n", pId)
		mascota, err := getMascotaById(pId)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusNotFound, "mascota no encontrada")
		}
		_, err = dbConnection.Exec("DELETE FROM mascotas WHERE id = ?", mascota.Id)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, "mascota no se pudo borrar")
		}
		return c.JSON(http.StatusOK, mascota)
	})

	e.Logger.Fatal(e.Start(PORT))
}
