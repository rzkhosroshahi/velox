package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rzkhosroshahi/velox/app"
)

func main() {
	app := app.NewApplication()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.App.Port),
		Handler:      app.Router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Sugar().Infof("app is running on port %d", app.Config.App.Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Panic("setup sever failed!")
	}
}
