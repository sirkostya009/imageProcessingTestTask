package web

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"imageProcessingTestTask/db"
)

type Handlers struct {
	db *db.Queries
}

func New(pool *pgxpool.Pool) Handlers {
	return Handlers{db.New(pool)}
}

func (h *Handlers) AddHandlers(app *echo.Echo) {
	app.POST("/register", h.Register)
	app.POST("/login", h.Login)

	imgGroup := app.Group("/images")
	imgGroup.Use(h.AuthMiddleware)
	imgGroup.GET("", h.GetImages)
	imgGroup.POST("", h.UploadImage)
}
