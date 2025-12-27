package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"

	"go-microservice/handlers"
	"go-microservice/metrics"
	"go-microservice/services"
	"go-microservice/utils"
)

func main() {
	// Оптимизация для максимальной производительности
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	userService := services.NewUserService()
	logger := utils.NewAuditLogger()
	notifier := utils.NewNotifier()

	// Увеличиваем rate limiter для достижения >1000 RPS
	limiter := utils.NewRateLimiter(5000, 10000)

	router := httprouter.New()
	
	// Middleware обертка для httprouter
	wrapMiddleware := func(handler httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Rate limiting
			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			
			// Metrics
			start := time.Now()
			metrics.TotalRequests.WithLabelValues(r.Method, r.URL.Path).Inc()
			
			// Вызываем handler
			handler(w, r, ps)
			
			metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
		}
	}

	userHandler := &handlers.UserHandler{
		Service: userService,
		Logger:  logger,
		Notify:  notifier,
	}
	userHandler.RegisterRoutes(router, wrapMiddleware)

	router.Handler(http.MethodGet, "/metrics", metrics.Handler())

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
		// Оптимизация для высокой нагрузки
		ReadHeaderTimeout: 2 * time.Second,
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
