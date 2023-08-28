package main

import (
	"avitoTech/config"
	"avitoTech/internal/delimer/handlers"
	"avitoTech/internal/delimer/repository"
	"avitoTech/internal/delimer/usecase"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	conf := &config.Config{
		Postgres: struct {
			User     string `json:"user"`
			Password string `json:"password"`
			Host     string `json:"host"`
			Port     string `json:"port"`
			DbName   string `json:"dbName"`
		}(struct {
			User     string
			Password string
			Host     string
			Port     string
			DbName   string
		}{User: "almaz", Password: "almaz", Host: "127.0.0.1", Port: "5432", DbName: "postgres"}),
		ServerPort:   ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	rep, err := repository.NewPostgres(context.Background(), conf)
	if err != nil {
		log.Panic(err)
	}
	ser := usecase.NewService(rep)
	f := handlers.NewFiberServer(conf, ser)
	err = f.StartServer(conf)
	if err != nil {
		os.Exit(2)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	fmt.Println("Bye")
	os.Exit(0)
}
