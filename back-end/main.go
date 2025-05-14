package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

const PORT = "8000"

type Templates struct {
	templates *template.Template
}

type Mascota struct {
	Id          int
	Nombre      string
	Edad        int
	Altura_cm   int
	Foto        []byte
	Foto64      string
	Descripcion string
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
		nombre TEXT,
		edad INTEGER,
		altura_cm INTEGER,
		foto BLOB,
		descripcion TEXT)`)
	if err != nil {
		e.Logger.Fatal("error creating mascota table", err)
	}

	e.GET("/", func(c echo.Context) error {
		e.Logger.Printf("retrieved file")
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
		var reqBody Mascota
		if err := c.Bind(&reqBody); err != nil {
			return c.JSON(http.StatusBadRequest, "Bad request")
		}
		e.Logger.Printf("request body: %v", reqBody)
		if reqBody.Nombre == "" {
			return c.JSON(http.StatusBadRequest, "No name given")
		}
		_, err = dbConnection.Exec("INSERT INTO mascotas (name) VALUES (?)", reqBody.Name)
		if err != nil {
			e.Logger.Fatal(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		log.Println("inserted into db")
		return c.JSON(http.StatusCreated, reqBody.Name)
	})

	/* foto, err := os.ReadFile("assets/img/splash.jpg")
	if err != nil {
		e.Logger.Fatal("failed to load file")
	}
	base64str := base64.StdEncoding.EncodeToString(foto) */

	e.Logger.Fatal(e.Start(":" + PORT))
}
