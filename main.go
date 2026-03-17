package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rzkhosroshahi/velox/config"
	"github.com/rzkhosroshahi/velox/internal/user"
	"github.com/rzkhosroshahi/velox/pkg/api"
	"github.com/rzkhosroshahi/velox/pkg/db"
	"github.com/rzkhosroshahi/velox/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger
}

func main() {
	conf, err := config.Setup()
	if err != nil {
		log.Panic("error loading config!")
	}

	logger.Init(conf.App.Env)

	db, err := db.New(&conf.DataBase)
	if err != nil {
		log.Fatalln(err)
		log.Panic("can not connect to the database!")
	}
	logger.Log.Info("connected to the database!")

	userStore := user.NewUserStore(db)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)

	r := api.NewRouter(userHandler)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.App.Port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("app is running on port %d", conf.App.Port)
	err = server.ListenAndServe()
	if err != nil {
		log.Panic("setup sever failed!")
	}
}
