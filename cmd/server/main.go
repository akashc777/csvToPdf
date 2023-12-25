package main

import (
	"fmt"
	"github.com/akashc777/csvToPdf/helpers"
	"github.com/akashc777/csvToPdf/initializers/env"
	"github.com/akashc777/csvToPdf/initializers/postgresInit"
	"github.com/akashc777/csvToPdf/router"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {

	// load .env
	env.LoadEnvVariables()

}

type Config struct {
	Port string
}

type Application struct {
	Config Config
}

func (app *Application) Serve() error {
	port := app.Config.Port

	helpers.MessageLogs.InfoLog.Printf("API is listening on port : %+v", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router.Routes(),
	}

	return srv.ListenAndServe()

}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := Config{
		Port: os.Getenv("PORT"),
	}

	// connect to postgres
	dsn := os.Getenv("DSN")
	err = postgresInit.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to db")
	}

	defer func() {
		postgresInit.DBConn.SqlDB.Close()
	}()
	//Initialise OAuth
	//err = oauth.InitOAuth()
	//if err != nil {
	//	log.Fatal("Failed to configure oauth")
	//}

	app := &Application{
		Config: cfg,
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
