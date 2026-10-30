package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rcontrisha/backend-eventhub/internal/broker"
	"rcontrisha/backend-eventhub/internal/config"
	"rcontrisha/backend-eventhub/internal/router"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title													API Eventhub
// @version												1.0
// @description										Documentation of API Eventhub

// @host													localhost:8181
// @BasePath											/

// @securityDefinitions.apikey		BearerToken
// @in														header
// @name													Authorization
// @description										Bearer Token used as identity for accessing backend

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file.")
	}

	pdb := config.NewPsqlDb(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	pool, err := pdb.Connect()
	if err != nil {
		log.Println("Cannot Connect to DB\nReason: ", err.Error())
		return
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Println("Database is not ready\nReason: ", err.Error())
		return
	}
	log.Println("Database is ready")

	rdb := config.NewRedisClient(os.Getenv("REDIS_USER"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")).Connect()
	defer rdb.Close()

	str, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Redis is not ready\nReason: ", err.Error())
		return
	}
	log.Println(str)
	log.Println("Redis is ready")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	brk := broker.NewBroker()
	go brk.Run(ctx)

	r := gin.Default()

	router.MainRouter(r, pool, rdb, brk)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("Server error: ", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("Shutdown error: ", err)
	}
}
