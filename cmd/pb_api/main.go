package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/jcgfreitas/pb_api/internal/handlers"
	"github.com/jcgfreitas/pb_api/internal/repository"
	"github.com/jcgfreitas/pb_api/internal/service"
	"github.com/sirupsen/logrus"

	"github.com/jcgfreitas/pb_api/pkg/gormdb/postgres"
)

func main() {
	// flag management
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish")
	host := flag.String("host", "localhost", "postgres db hostname")
	user := flag.String("user", "postgres", "postgres db username")
	dbName := flag.String("name", "postgres", "postgres db name")
	password := flag.String("password", "password1", "postgres db password")
	port := flag.String("port", "5432", "postgres db port number")
	debug := flag.Bool("debug", true, "debug logger level")
	drop := flag.Bool("dropTable", false, "drop coupons table rows")
	flag.Parse()

	// start logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	if *debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	// start db connection
	logger.Info("starting db connection")
	db, err := postgres.Open(*host, *port, *user, *dbName, *password)
	if err != nil {
		logger.WithError(err).Fatal("failed to connect to database")
	}
	defer db.Close()

	// create handler and its chained dependencies
	repo := repository.New(db)
	if *drop {
		repo.Reset()
	}
	s := service.NewService(repo, logger)
	h := handlers.NewHandlers(s, logger)

	// router creation and assignment of handlers
	r := mux.NewRouter()
	r.HandleFunc(h.CreateCouponPath(), h.CreateCouponHandler).Methods("POST")
	r.HandleFunc(h.GetCouponsPath(), h.GetCouponsHandler).Methods("GET")
	r.HandleFunc(h.GetCouponPath(), h.GetCouponHandler).Methods("GET")
	r.HandleFunc(h.DeleteCouponPath(), h.DeleteCouponHandler).Methods("DELETE")
	r.HandleFunc(h.UpdateCouponPath(), h.UpdateCouponHandler).Methods("POST")

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	logger.Info("starting server")
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.WithError(err).Fatal("failed to start server")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Info("shutting down")
	os.Exit(0)

}
