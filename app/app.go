package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rzkhosroshahi/velox/api"
	"github.com/rzkhosroshahi/velox/config"
	"github.com/rzkhosroshahi/velox/internal/token"
	"github.com/rzkhosroshahi/velox/internal/user"
	"github.com/rzkhosroshahi/velox/pkg/db"
	"github.com/rzkhosroshahi/velox/pkg/logger"
	"github.com/rzkhosroshahi/velox/pkg/redis"
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

	redisClient, err := redis.New(&conf.Redis)
	if err != nil {
		logger.Log.Panic("can not connect to the redis!")
	}
	// token service
	tokenService := token.NewService(redisClient, conf.App.JWTSecretKey)
	tokenHandler := token.NewHandler(tokenService)
	// user service
	userStore := user.NewUserStore(db)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService, tokenService)

	r := api.NewRouter(userHandler, tokenHandler)

	app := &Application{
		db:     db,
		Router: r,
		Config: conf,
	}
	return app
}
