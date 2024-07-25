package main

import (
	"flag"
	"justpay/internal/config"
	"justpay/internal/handler"
	"log"
	"net/http"
)

var configPath = flag.String("config", "./configs/local.yml", "Path to config file")

func main() {
	flag.Parse()

	log.Printf("init application config '%s'", *configPath)
	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatalf("can't init config: %v", err)
	}

	log.Printf("init application handler")
	appHandler := handler.New(cfg)

	log.Printf("start listening http server on addr '%s'", cfg.HTTPServerListenAddr)
	err = http.ListenAndServe(cfg.HTTPServerListenAddr, appHandler.Init())
	if err != nil {
		log.Fatal(err)
	}
}
