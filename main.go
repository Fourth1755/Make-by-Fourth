package main

import (
	"Fourth1755/Make-by-Fourth/cloudpockets"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func InitDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error")
	}
	createTb := `
	CREATE TABLE IF NOT EXISTS cloud_pockets ( id SERIAL PRIMARY KEY, name TEXT,category TEXT,currency TEXT,balance FLOAT,account TEXT);
	`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
	fmt.Println("create table success")
	return db
}

func main() {
	LoadEnvVariables()
	db := InitDB()
	fmt.Println("start at port", os.Getenv("PORT"))

	e := echo.New()
	hCloudPockets := cloudpockets.NewApplication(db)

	e.GET("/cloud-pockets", hCloudPockets.GetAllCloudPockets)
	e.POST("/cloud-pockets", hCloudPockets.CreateCloudPockets)

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	fmt.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("bye bye")
}
