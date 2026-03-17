package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rzkhosroshahi/velox/api"
	"github.com/rzkhosroshahi/velox/config"
	"github.com/rzkhosroshahi/velox/pkg/logger"
)

func main() {
	conf, err := config.Setup()
	if err != nil {
		fmt.Println("error loading config!")
	}
	logger.Init(conf.App.Env)

	r := api.NewRouter()

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
		log.Fatalf("setup sever failed!")
	}
}
