package api

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"proweb-backend/database"

	"github.com/labstack/echo/v4"
)

var (
	e          *echo.Echo
	_PORT      string
	_URL       string
	repository *database.DbRepository
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

func StartServer(url, port string) {
	_PORT = port
	_URL = url
	e = echo.New()
	e.Static("/assets", "assets")
	e.Renderer = newTemplate()
	dbConnection, err := sql.Open("sqlite3", "database/pagina.db")
	if err != nil {
		e.Logger.Fatal("Error connecting to db")
	}
	defer dbConnection.Close()
	repository = &database.DbRepository{Db: dbConnection}
	err = repository.Connect()
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.GET("/*", func(c echo.Context) error {
		fmt.Printf("GET catchall\n")
		response := map[string]any{
			"URL":     _URL,
			"Code":    http.StatusNotFound,
			"Message": "Not found",
		}
		return c.Render(http.StatusNotFound, "error", response)
	})

	e.GET("/", getRoot)
	e.GET("/adopcion", getAdopcion)
	e.GET("/adopcion/:id", getAdopcionById)
	e.GET("/contacto", getContacto)
	e.GET("/mascotas", getMascotas)
	e.GET("/mascotas/:id", getMascotaById)
	e.GET("/registrar", getRegistrar)
	e.POST("/adopcion/:id", postAdopcion)
	e.POST("/contacto", postContacto)
	e.POST("/mascotas", postMascota)

	// json
	e.GET("/json/contacto", jsonGetContacto)
	e.GET("/json/mascotas", jsonGetMascotas)
	e.GET("/json/mascotas/:id", jsonGetMascotaById)
	e.GET("/json/solicitudes", jsonGetSolicitudes)
	e.DELETE("/json/contacto/:id", jsonDelContacto)
	e.DELETE("/json/mascotas/:id", jsonDelMascota)
	e.Logger.Fatal(e.Start(_PORT))
}
