package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rzkhosroshahi/velox/api"
	"github.com/rzkhosroshahi/velox/config"
	"github.com/rzkhosroshahi/velox/internal/user"
	"github.com/rzkhosroshahi/velox/pkg/db"
	"github.com/rzkhosroshahi/velox/pkg/logger"
)

type Application struct {
	db     *sqlx.DB
	Router chi.Router
	Config *config.Config
}

func NewApplication() *Application {
	conf, err := config.Setup()

	if err != nil {
		panic("error loading config!")
	}
	logger.Init(conf.App.Env)

	db, err := db.New(&conf.DataBase)
	if err != nil {
		logger.Log.Panic("can not connect to the database!")
	}
	logger.Log.Info("connected to the database!")

	// user service
	userStore := user.NewUserStore(db)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)

	r := api.NewRouter(userHandler)

	app := &Application{
		db:     db,
		Router: r,
		Config: conf,
	}
	return app
}
