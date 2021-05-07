package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	logger := log.New(os.Stdout, "my app: ", log.LstdFlags)

	server := runServer(logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() (err error) {
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			case <-c:
				server.SetKeepAlivesEnabled(false)
				err = server.Shutdown(ctx)
				return
			}
		}
	})

	logger.Println("http server Running on http://:8080")
	_ = server.ListenAndServe()

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		logger.Fatalf("shutdown server found error: %v\n", err)
	}

	logger.Println("server gracefully shutdown")
}

func runServer(logger *log.Logger) *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/show", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	})

	return &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
