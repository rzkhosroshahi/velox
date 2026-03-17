package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rzkhosroshahi/velox/api"
	"github.com/rzkhosroshahi/velox/pkg/logger"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	logger.Init("development")

	r := api.NewRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	fmt.Printf("app is running on port %d", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("setup sever failed!")
	}
}
