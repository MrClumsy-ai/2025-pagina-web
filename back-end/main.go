package main

import (
	"database/sql"
	"encoding/base64"
	"html/template"
	"io"
	"os"
	"proweb-backend/database"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type Templates struct {
	templates *template.Template
}

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

type Mascota struct {
	Foto        []byte
	Foto64      string
	Nombre      string
	Edad        int
	Altura      int
	Descripcion string
}

func main() {
	e := echo.New()
	e.Static("/assets", "assets")
	e.Renderer = newTemplate()
	dbConnection, err := sql.Open("sqlite3", "pagina.db")
	if err != nil {
		e.Logger.Fatal("error connecting to db: ", err)
	}
	defer dbConnection.Close()
	mascotaRepository := &database.MascotaRepository{Db: dbConnection}
	err = mascotaRepository.CreateTable()
	if err != nil {
		e.Logger.Fatal("error creating mascota table: ", err)
	}

	var mascotas []Mascota

	e.GET("/", func(c echo.Context) error {
		foto, err := os.ReadFile("assets/img/splash.jpg")
		if err != nil {
			e.Logger.Fatal("failed to load file")
		}
		base64str := base64.StdEncoding.EncodeToString(foto)
		e.Logger.Printf("retrieved file")
		// test
		mascotas = []Mascota{
			{Foto64: base64str, Nombre: "lorem", Edad: 1, Altura: 23, Descripcion: "lorem ipsum lalalala"},
			{Foto64: base64str, Nombre: "ipsum", Edad: 2, Altura: 12, Descripcion: "lorem ipsum lalalala"},
			{Foto64: base64str, Nombre: "something", Edad: 3, Altura: 34, Descripcion: "lorem ipsum lalalala"},
		}
		response := map[string]any{
			"CurrentRoute": "/",
			"Mascotas":     mascotas,
		}
		return c.Render(200, "inicio", response)
	})

	e.GET("/adopcion", func(c echo.Context) error {
		response := map[string]any{
			"CurrentRoute": "/adopcion",
		}
		return c.Render(200, "adopcion", response)
	})

	e.GET("/mascotas", func(c echo.Context) error {
		response := map[string]any{
			"CurrentRoute": "/mascotas",
			"Mascotas":     mascotas,
		}
		return c.Render(200, "mascotas", response)
	})

	e.GET("/contacto", func(c echo.Context) error {
		response := map[string]any{
			"CurrentRoute": "/contacto",
		}
		return c.Render(200, "contacto", response)
	})

	e.GET("/api/mascotas", func(c echo.Context) error {
		return nil
	})
	e.Logger.Fatal(e.Start(":8000"))
}
