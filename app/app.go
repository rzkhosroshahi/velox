package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rzkhosroshahi/velox/api"
	"github.com/rzkhosroshahi/velox/config"
	"github.com/rzkhosroshahi/velox/internal/user"
	"github.com/rzkhosroshahi/velox/pkg/db"
	"github.com/rzkhosroshahi/velox/pkg/logger"
	"go.uber.org/zap"
)

type Application struct {
	db     *sqlx.DB
	Router chi.Router
	Logger *zap.Logger
	Config *config.Config
}

func NewApplication() *Application {
	conf, err := config.Setup()

	if err != nil {
		panic("error loading config!")
	}
	logger := logger.Init(conf.App.Env)

	db, err := db.New(&conf.DataBase)
	if err != nil {
		logger.Panic("can not connect to the database!")
	}
	logger.Info("connected to the database!")

	// user service
	userStore := user.NewUserStore(db)
	userService := user.NewService(userStore, logger)
	userHandler := user.NewHandler(userService, logger)

	r := api.NewRouter(userHandler)

	app := &Application{
		db:     db,
		Router: r,
		Config: conf,
		Logger: logger,
	}
	return app
}
