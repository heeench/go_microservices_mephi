package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"go-microservice/handlers"
	"go-microservice/metrics"
	"go-microservice/services"
	"go-microservice/utils"
)

func main() {
	userService := services.NewUserService()
	logger := utils.NewAuditLogger()
	notifier := utils.NewNotifier()

	limiter := utils.NewRateLimiter(1000, 5000)

	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return metrics.MetricsMiddleware(next)
	})
	router.Use(func(next http.Handler) http.Handler {
		return utils.RateLimitMiddleware(limiter, next)
	})

	userHandler := &handlers.UserHandler{
		Service: userService,
		Logger:  logger,
		Notify:  notifier,
	}
	userHandler.RegisterRoutes(router)

	router.Path("/metrics").Handler(metrics.Handler())

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
