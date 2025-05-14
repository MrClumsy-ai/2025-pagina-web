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
	"golang.org/x/net/context/ctxhttp"
)

const PORT = "8000"

type Templates struct {
	templates *template.Template
}

type Mascota struct {
	Id          int    `form:"id" json:"id"`
	Nombre      string `form:"nombre" json:"nombre"`
	Edad        int    `form:"edad" json:"edad"`
	Altura_cm   int    `form:"altura" json:"altura"`
	Foto        []byte `form:"foto" json:"foto"`
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
		foto BLOB,
		descripcion TEXT)`)
	if err != nil {
		e.Logger.Fatal("error creating mascota table", err)
	}

	e.GET("/", func(c echo.Context) error {
		response := map[string]any{
			"CurrentRoute": "/",
			"Mascotas":     nil,
		}
		return c.Render(http.StatusOK, "inicio", response)
	})

	e.GET("/adopcion", func(c echo.Context) error {
		response := map[string]any{
			"CurrentRoute": "/adopcion",
		}
		return c.Render(http.StatusOK, "adopcion", response)
	})

	e.GET("/mascotas", func(c echo.Context) error {
		response := map[string]any{
			"CurrentRoute": "/mascotas",
			"Mascotas":     nil,
		}
		return c.Render(http.StatusOK, "mascotas", response)
	})

	e.GET("/contacto", func(c echo.Context) error {
		response := map[string]any{
			"CurrentRoute": "/contacto",
		}
		return c.Render(http.StatusOK, "contacto", response)
	})

	e.GET("/registrar", func(c echo.Context) error {
		return c.Render(http.StatusOK, "registrar", nil)
	})

	// API

	e.GET("/api/mascotas", func(c echo.Context) error {
		rows, err := dbConnection.Query("SELECT * FROM mascotas")
		if err != nil {
			e.Logger.Fatal(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		defer rows.Close()
		var mascotas Mascotas
		for rows.Next() {
			var mascota Mascota
			err := rows.Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto, &mascota.Descripcion)
			if err != nil {
				e.Logger.Fatal(err)
				return c.JSON(http.StatusInternalServerError, nil)
			}
			mascotas = append(mascotas, mascota)
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
			e.Logger.Fatal(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		err = dbConnection.QueryRow("SELECT * FROM mascotas WHERE ID = ?", pId).Scan(&mascota.Id, &mascota.Nombre, &mascota.Edad, &mascota.Altura_cm, &mascota.Foto, &mascota.Descripcion)
		if err != nil {
			return c.JSON(http.StatusNotFound, "mascota not found")
		}
		return c.JSON(http.StatusFound, mascota)
	})

	e.POST("/api/mascotas", func(c echo.Context) error {
		// c.Request().ParseMultipartForm(10 << 20)
		nombre := c.FormValue("nombre")
		edad, err := strconv.Atoi(c.FormValue("edad"))
		altura, err := strconv.Atoi(c.FormValue("altura"))
		descripcion := c.FormValue("descripcion")
		foto, err := c.FormFile("foto")
		if err != nil {
			e.Logger.Fatal(err)
		}
		fotoContenido, err := foto.Open()
		if err != nil {
			e.Logger.Fatal(err)
		}
		defer fotoContenido.Close()
		fotoBytes, err := io.ReadAll(fotoContenido)
		if err != nil {
			e.Logger.Fatal(err)
		}
		reqBody := Mascota{
			Nombre:      nombre,
			Edad:        edad,
			Altura_cm:   altura,
			Descripcion: descripcion,
			Foto:        fotoBytes,
		}
		e.Logger.Printf("request body: %v", reqBody)
		if err := c.Bind(&reqBody); err != nil {
			return c.JSON(http.StatusBadRequest, "Bad request")
		}
		e.Logger.Printf("request body: %v", reqBody)
		if reqBody.Nombre == "" {
			return c.JSON(http.StatusBadRequest, "Ningun nombre dado")
		}
		_, err = dbConnection.Exec("INSERT INTO mascotas (nombre, edad, altura_cm, foto, descripcion) VALUES (?, ?, ?, ?, ?)",
			reqBody.Nombre, reqBody.Edad, reqBody.Altura_cm, reqBody.Foto, reqBody.Descripcion)
		if err != nil {
			e.Logger.Fatal(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		e.Logger.Printf("inserted into db: %v", reqBody)
		return c.JSON(http.StatusCreated, reqBody)
	})

	/* foto, err := os.ReadFile("assets/img/splash.jpg")
	if err != nil {
		e.Logger.Fatal("failed to load file")
	}
	base64str := base64.StdEncoding.EncodeToString(foto) */

	e.Logger.Fatal(e.Start(":" + PORT))
}
