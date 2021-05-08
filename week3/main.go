package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)
	server := HTTPServer()
	debugServer := HTTPDebug()

	g.Go(func() error {
		log.Println("start http server")
		return server.ListenAndServe()
	})

	g.Go(func() error {
		log.Println("start http debug server")
		return debugServer.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()
		log.Print("other goroutine done")
		g.Go(func() error {
			log.Println("server shutting down")
			return server.Shutdown(ctx)
		})
		g.Go(func() error {
			log.Println("debug server shutting down")
			return debugServer.Shutdown(ctx)
		})
		return nil
	})

	g.Go(func() error {
		return OsSignal(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Println("received err, all goroutine quits: err = ", err)
	}
}

func HTTPServer() *http.Server {
	index := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/", index)
	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	return &server
}

func HTTPDebug() *http.Server {
	server := http.Server{
		Addr:    ":8001",
		Handler: http.DefaultServeMux,
	}
	return &server
}

func OsSignal(ctx context.Context) error {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.Println("continuing receive signal")
	select {
	case <-ctx.Done():
		log.Println("other goroutine quit, context has done (OsSignal)")
		return ctx.Err()
	case <-sc:
		log.Println("received signal from chan")
		return errors.New("received os signal")
	}
}
