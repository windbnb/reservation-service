package main

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/windbnb/reservation-service/tracer"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/windbnb/reservation-service/handler"
	"github.com/windbnb/reservation-service/repository"
	"github.com/windbnb/reservation-service/router"
	"github.com/windbnb/reservation-service/service"
	"github.com/windbnb/reservation-service/util"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	db := util.ConnectToDatabase()

	tracer, closer := tracer.Init("reservation-service")
	opentracing.SetGlobalTracer(tracer)
	router := router.ConfigureRouter(&handler.Handler{
		Tracer:  tracer,
		Closer:  closer,
		Service: &service.ReservationRequestService{Repo: &repository.Repository{Db: db}}})

	servicePath, servicePathFound := os.LookupEnv("SERVICE_PATH")
	if !servicePathFound {
		servicePath = "localhost:8083"
	}

	srv := &http.Server{Addr: servicePath, Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
