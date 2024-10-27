package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"workflowsvr/api"
)

func main() {

	ctx := context.Background()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/start", api.StartWorkflow)

	srv := &http.Server{
		Addr: ":3000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				logrus.Warn("stopping http server gracefully, please wait...")
				if err := srv.Shutdown(context.Background()); err != nil {
					logrus.Errorf("got error when shutting down server err: %s", err)
				}
				return
			}
		}
	}()

	logrus.Warn("started at port:", 3000)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logrus.Errorf("got error when starting http server err: %s", err)
	}
}
