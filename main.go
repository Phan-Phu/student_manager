package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"student_manager/routes"
	"student_manager/services"
	"time"
)

func main() {
	// load config
	services.LoadConfig()

	// init mongo
	services.InitMongoDB()

	//init data
	services.InitializeRepository()

	// config API
	routes.InitGin()
	router := routes.ConfigRoute()

	// use redis

	// config server
	server := &http.Server{
		Addr:         services.Config.ServerAddr + ":" + services.Config.ServerPort,
		WriteTimeout: time.Second * 300,
		ReadTimeout:  time.Second * 300,
		IdleTimeout:  time.Second * 300,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 15 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
