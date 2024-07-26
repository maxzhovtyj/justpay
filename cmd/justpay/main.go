package main

import (
	"context"
	"flag"
	"justpay/internal/config"
	"justpay/internal/handler"
	"justpay/internal/service"
	"justpay/internal/storage"
	"justpay/pkg/postgresql"
	"log"
	"net/http"
)

var configPath = flag.String("config", "./configs/config.yml", "Path to config file")

func main() {
	flag.Parse()

	log.Printf("init application config '%s'", *configPath)
	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatalf("can't init config: %v", err)
	}

	conn, err := postgresql.NewConn(cfg.DBSourceName)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer func() {
		err = conn.Close(context.Background())
		if err != nil {
			return
		}
	}()

	log.Printf("init application storage")
	appStorage := storage.New(conn)

	log.Printf("init application service")
	appService := service.New(appStorage)

	log.Printf("init application handler")
	appHandler := handler.New(cfg, appService)

	log.Printf("start listening http server on addr '%s'", cfg.HTTPServerListenAddr)
	err = http.ListenAndServe(cfg.HTTPServerListenAddr, appHandler.Init())
	if err != nil {
		log.Fatal(err)
	}
}
