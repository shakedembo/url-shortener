package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"url-shortener/config"
	"url-shortener/internal"
	"url-shortener/internal/utils"
	"url-shortener/pkg"
	"url-shortener/server"
)

var logger = log.New(os.Stdout, "url-shortener ", log.LstdFlags|log.Lshortfile)

func main() {
	logger.Print("Initializing...")
	defer logger.Print("Bye Bye :)")

	timeout := 30 * time.Second

	convertFunc := func(in string) []byte { return []byte(in) }
	hash := utils.NewFnv32aHashProvider(convertFunc, logger)
	dbConfig := config.NewDbConfig()
	dao := internal.NewCassandraDAO(dbConfig, logger)
	service := internal.NewUrlShorteningService(hash, dao, logger)
	controller := internal.NewUrlShortenerController(service)

	port := "8080" // os.Getenv("SERVER_PORT")

	server.Init(logger)
	server.AddHandler[pkg.CreateRequest, pkg.CreateResponse]("/urls/", controller.Create, timeout, server.POST)
	server.AddHandler[pkg.ReadRequest, pkg.ReadResponse]("/urls/", controller.Read, timeout, server.GET)
	server.AddHandler[pkg.UpdateRequest, pkg.UpdateResponse]("/urls/", controller.Replace, timeout, server.PUT)
	server.AddHandler[pkg.DeleteRequest, pkg.DeleteResponse]("/urls/", controller.Delete, timeout, server.DELETE)

	go server.Listen(port)

	logger.Print("--- Initialization Completed Successfully ---")

	gracefullyQuit()
}

func gracefullyQuit() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
}
