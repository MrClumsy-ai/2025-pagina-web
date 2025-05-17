package main

import (
	"database/sql"
	"encoding/base64"
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
	e.Logger.Printf("Database connection established: %v", dbLocation)
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
		e.Logger.Printf("db: selecting id: %v", id)
		err = dbConnection.QueryRow("SELECT * FROM mascotas WHERE ID = ?", id).Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto64, &mascota.Descripcion)
		if err != nil {
			return Mascota{}, err
		}
		e.Logger.Printf("mascota: %v", mascota)
		return mascota, nil
	}

	// routes
	e.GET("/", func(c echo.Context) error {
		e.Logger.Print("GET /")
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
		e.Logger.Print("GET /adopcion")
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/adopcion",
		}
		return c.Render(http.StatusOK, "adopcion", response)
	})

	e.GET("/mascotas", func(c echo.Context) error {
		e.Logger.Print("GET /mascotas")
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
		e.Logger.Print("GET /mascotas/id")
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
		e.Logger.Printf("mascotas: %v", mascotas)
		return c.Render(http.StatusOK, "mascotas", response)
	})

	e.GET("/contacto", func(c echo.Context) error {
		e.Logger.Print("GET /contacto")
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/contacto",
		}
		return c.Render(http.StatusOK, "contacto", response)
	})

	e.POST("/contacto", func(c echo.Context) error {
		e.Logger.Printf("POST /contacto")
		nombre := c.FormValue("nombre")
		email := c.FormValue("email")
		mensaje := c.FormValue("mensaje")
		reqBody := Mensaje{
			Nombre:  nombre,
			Email:   email,
			Mensaje: mensaje,
		}
		e.Logger.Printf("%v (%v): %v", reqBody.Nombre, reqBody.Email, reqBody.Mensaje)
		if err := c.Bind(&reqBody); err != nil {
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusBadRequest,
				"Message": "Bad request",
			}
			return c.Render(http.StatusBadRequest, "error", response)
		}
		if reqBody.Nombre == "" {
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusBadRequest,
				"Message": "Ningun nombre dado",
			}
			return c.Render(http.StatusBadRequest, "error", response)
		}
		if reqBody.Email == "" {
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusBadRequest,
				"Message": "Ningun email dado",
			}
			return c.Render(http.StatusBadRequest, "error", response)
		}
		if reqBody.Mensaje == "" {
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
		e.Logger.Printf("inserted into db: %v", reqBody.Nombre)
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
		e.Logger.Printf("inserted id: %v", g_mensaje.Id)
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/registrar",
			"Mensaje":      g_mensaje,
		}
		return c.Render(http.StatusCreated, "contacto", response)
	})

	e.GET("/registrar", func(c echo.Context) error {
		e.Logger.Print("GET /registrar")
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/registrar",
		}
		return c.Render(http.StatusOK, "registrar", response)
	})

	e.POST("/mascotas", func(c echo.Context) error {
		e.Logger.Print("POST /mascotas")
		nombre := c.FormValue("nombre")
		edad, err := strconv.Atoi(c.FormValue("edad"))
		if err != nil {
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusUnprocessableEntity,
				"Message": "Edad no es procesable",
			}
			return c.Render(http.StatusUnprocessableEntity, "error", response)
		}
		altura, err := strconv.Atoi(c.FormValue("altura"))
		if err != nil {
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
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusBadRequest,
				"Message": "Bad request",
			}
			return c.Render(http.StatusBadRequest, "error", response)
		}
		if reqBody.Nombre == "" {
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
		e.Logger.Printf("inserted into db: %v", reqBody.Nombre)
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
		e.Logger.Printf("inserted id: %v", mascota.Id)
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/registrar",
			"Mascota":      mascota,
		}
		return c.Render(http.StatusCreated, "registrado", response)
	})

	// JSON
	e.GET("/json/mascotas", func(c echo.Context) error {
		e.Logger.Print("GET /json/mascotas")
		mascotas, err := getAllMascotas()
		if err != nil {
			e.Logger.Fatal(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		if len(mascotas) == 0 {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusFound, mascotas)
	})

	e.GET("/json/mascotas/:id", func(c echo.Context) error {
		e.Logger.Print("GET /json/mascotas/id")
		pId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		mascota, err := getMascotaById(pId)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusNotFound, "mascota not found")
		}
		return c.JSON(http.StatusFound, mascota)
	})

	e.GET("/json/contacto", func(c echo.Context) error {
		e.Logger.Print("GET /json/contacto")
		rows, err := dbConnection.Query("SELECT * FROM mensajes")
		if err != nil {
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
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusNotFound,
				"Message": "Not found",
			}
			return c.JSON(http.StatusNotFound, response)
		}
		return c.JSON(http.StatusFound, mensajes)
	})

	e.GET("/*", func(c echo.Context) error {
		response := map[string]any{
			"URL":     URL,
			"Code":    http.StatusNotFound,
			"Message": "Not found",
		}
		return c.Render(http.StatusNotFound, "error", response)
	})

	e.Logger.Fatal(e.Start(PORT))
}
