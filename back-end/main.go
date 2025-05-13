package main

import (
	"database/sql"
	"html/template"
	"io"

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

type Counter struct {
	Count int
}

type Mascota struct {
	Foto        []byte
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

	mascotas := []Mascota{
		Mascota{Nombre: "lorem", Edad: 1, Altura: 23, Descripcion: "lorem ipsum lalalala"},
		Mascota{Nombre: "ipsum", Edad: 2, Altura: 12, Descripcion: "lorem ipsum lalalala"},
		Mascota{Nombre: "something", Edad: 3, Altura: 34, Descripcion: "lorem ipsum lalalala"},
		Mascota{Nombre: "else", Edad: 4, Altura: 31, Descripcion: "lorem ipsum lalalala"},
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "inicio", mascotas)
	})
	e.Logger.Fatal(e.Start(":8000"))
}
