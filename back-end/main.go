package main

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	e := echo.New()
	dbConnection, err := sql.Open("sqlite3", "pagina.db")
	if err != nil {
		e.Logger.Fatal("error connecting to db: ", err)
	}
	defer dbConnection.Close()

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `<h1>Welcome</h1>`)
	})
	e.Logger.Fatal(e.Start(":8000"))
}
