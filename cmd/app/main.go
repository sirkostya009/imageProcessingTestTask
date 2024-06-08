package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"imageProcessingTestTask/web"
	"os"
)

func main() {
	app := echo.New()

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	h := web.New(pool)
	h.AddHandlers(app)

	err = app.Start(":8080")
	if err != nil {
		panic(err)
	}
}
