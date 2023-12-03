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
	mdl := []interface{}{
		models.KeyboardEvent{},
		models.KeyboardFlow{},
		models.User{},
		models.Password{},
		models.AccessToken{},
	}
	if err := database.MakeMigrations(mdl); err != nil {
		return err
	}
	app.db = db

	// serve front-end
	app.fiber.Static("/", "./web")

	// wire
	userService := services.NewUserService(app.db, app.config)
	keyboardService := services.NewKeyboardService(app.db, app.config, userService)

	inputCon := controllers.NewInputController(keyboardService, userService)
	userCon := controllers.NewUserController(userService)

	// handle requests
	app.fiber.Post("/auth/register", userCon.Register)
	app.fiber.Post("/auth/login", userCon.Login)
	app.fiber.Post("/user/has-secret", userCon.UserHasSecret)
	app.fiber.Post("/user/get-secret", userCon.GetSecret)
	app.fiber.Post("/user/set-secret", userCon.SetSecret)

	app.fiber.Post("/get-passwords", inputCon.GetPasswords)
	app.fiber.Post("/process", inputCon.Process)
	app.fiber.Post("/history", inputCon.History)
	app.fiber.Post("/history/:id", inputCon.History)
	app.fiber.Post("/get-token", inputCon.GetToken)

	return app.fiber.Listen(app.config.AppUrl)
}
