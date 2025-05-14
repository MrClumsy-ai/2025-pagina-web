package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	// "encoding/base64"
	"html/template"
	"io"
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
	mascotaRepository := &database.MascotaRepository{Db: dbConnection}
	err = mascotaRepository.CreateTable()
	if err != nil {
		e.Logger.Fatal("error creating mascota table: ", err)
	}

	e.GET("/", func(c echo.Context) error {
		/* foto, err := os.ReadFile("assets/img/splash.jpg")
		if err != nil {
			e.Logger.Fatal("failed to load file")
		}
		base64str := base64.StdEncoding.EncodeToString(foto) */
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

	e.GET("/api/mascotas", func(c echo.Context) error {
		res, err := mascotaRepository.GetAll()
		if err != nil {
			e.Logger.Fatal("error getting all", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		r, err := json.Marshal(res)
		e.Logger.Printf("%v", r)
		return c.JSON(http.StatusOK, r)
	})

	e.Logger.Fatal(e.Start(":8000"))
}
