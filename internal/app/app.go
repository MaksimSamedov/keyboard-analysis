package app

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"keyboard-analysis/internal/config"
	"keyboard-analysis/internal/database"
	"keyboard-analysis/internal/models"
	"keyboard-analysis/internal/services"
	"keyboard-analysis/internal/transport/controllers"
	"log"
)

type App struct {
	fiber  *fiber.App
	config config.Config
	logger *log.Logger
	db     *gorm.DB
}

func New(conf config.Config, logger *log.Logger) *App {
	return &App{
		fiber:  fiber.New(conf.Fiber),
		config: conf,
		logger: logger,
	}
}

func (app *App) Run() error {
	// setup database
	db, err := database.NewConnection(app.config)
	if err != nil {
		return err
	}
	// migrate
	mdl := []interface{}{models.KeyboardEvent{}, models.KeyboardFlow{}}
	if err := database.MakeMigrations(mdl); err != nil {
		return err
	}
	app.db = db

	// serve front-end
	app.fiber.Static("/", "./web")

	// handle requests
	keyboardService := services.NewKeyboardService(app.db)
	inputCon := controllers.NewInputController(keyboardService)
	app.fiber.Post("/process", inputCon.Process)
	app.fiber.Get("/history", inputCon.History)
	app.fiber.Get("/history/:id", inputCon.History)

	return app.fiber.Listen(app.config.AppUrl)
}
