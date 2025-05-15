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
	dbConnection, err := sql.Open("sqlite3", "pagina.db")
	if err != nil {
		e.Logger.Fatal("error connecting to db: ", err)
	}
	defer dbConnection.Close()
	_, err = dbConnection.Exec(`CREATE TABLE IF NOT EXISTS mascotas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		edad INTEGER,
		altura_cm INTEGER,
		foto64 BLOB,
		descripcion TEXT)`)
	if err != nil {
		e.Logger.Fatal("error creating mascota table", err)
	}

	getMascotas := func() (Mascotas, error) {
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

	// routes
	e.GET("/", func(c echo.Context) error {
		rows, err := dbConnection.Query("SELECT * FROM mascotas LIMIT 3")
		if err != nil {
			e.Logger.Error(err)
			return c.HTML(http.StatusInternalServerError, "<h1>500: Internal server error</h1>")
		}
		defer rows.Close()
		var mascotas Mascotas
		for rows.Next() {
			var mascota Mascota
			err := rows.Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto64, &mascota.Descripcion)
			if err != nil {
				e.Logger.Error(err)
				return c.HTML(http.StatusInternalServerError, "<h1>500: Internal server error</h1>")
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
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/adopcion",
		}
		return c.Render(http.StatusOK, "adopcion", response)
	})

	e.GET("/mascotas", func(c echo.Context) error {
		mascotas, err := getMascotas()
		if err != nil {
			e.Logger.Error(err)
			c.HTML(http.StatusInternalServerError, "<h1>Internal server error</h1>")
		}
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/mascotas",
			"Mascotas":     mascotas,
		}
		return c.Render(http.StatusOK, "mascotas", response)
	})

	e.GET("/contacto", func(c echo.Context) error {
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/contacto",
		}
		return c.Render(http.StatusOK, "contacto", response)
	})

	e.GET("/registrar", func(c echo.Context) error {
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/registrar",
		}
		return c.Render(http.StatusOK, "registrar", response)
	})

	// API
	e.GET("/api/mascotas", func(c echo.Context) error {
		mascotas, err := getMascotas()
		if err != nil {
			e.Logger.Fatal(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		if len(mascotas) == 0 {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusFound, mascotas)
	})

	e.GET("/api/mascotas/:id", func(c echo.Context) error {
		var mascota Mascota
		pId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		err = dbConnection.QueryRow("SELECT * FROM mascotas WHERE ID = ?", pId).Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto64, &mascota.Descripcion)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusNotFound, "mascota not found")
		}
		return c.JSON(http.StatusFound, mascota)
	})

	e.POST("/api/mascotas", func(c echo.Context) error {
		nombre := c.FormValue("nombre")
		edad, err := strconv.Atoi(c.FormValue("edad"))
		altura, err := strconv.Atoi(c.FormValue("altura"))
		descripcion := c.FormValue("descripcion")
		foto, err := c.FormFile("foto")
		if err != nil {
			return c.HTML(http.StatusBadRequest, "<h1>400: Bad request</h1>")
		}
		src, err := foto.Open()
		if err != nil {
			e.Logger.Error(err)
			return c.HTML(http.StatusInternalServerError, "<h1>500: Internal server error</h1>")
		}
		defer src.Close()
		fotoBytes, err := io.ReadAll(src)
		if err != nil {
			e.Logger.Error(err)
			return c.HTML(http.StatusInternalServerError, "<h1>500: Internal server error</h1>")
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
			return c.HTML(http.StatusBadRequest, "<h1>Bad request</h1>")
		}
		if reqBody.Nombre == "" {
			return c.HTML(http.StatusBadRequest, "<h1>Ningun nombre dado</h1>")
		}
		_, err = dbConnection.Exec("INSERT INTO mascotas (nombre, edad, altura_cm, foto64, descripcion) VALUES (?, ?, ?, ?, ?)",
			reqBody.Nombre, reqBody.Edad, reqBody.Altura_cm, reqBody.Foto64, reqBody.Descripcion)
		if err != nil {
			e.Logger.Error(err)
			return c.HTML(http.StatusInternalServerError, "<h1>500: Internal server error</h1>")
		}
		e.Logger.Printf("inserted into db: %v", reqBody.Nombre)
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/registrar",
			"ReqBody":      reqBody,
		}
		return c.Render(http.StatusCreated, "registrado", response)
	})

	e.Logger.Fatal(e.Start(PORT))
}
