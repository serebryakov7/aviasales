package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/serebryakov7/aviasales/anagram"
	"github.com/serebryakov7/aviasales/handler"
)

var (
	portPtr = flag.Int("port", 8080, "http server port")
)

func main() {
	flag.Parse()

	var logger = log.New(os.Stdout, "", log.LstdFlags)
	logger.Printf("successfully initialized logger\n")

	var (
		done     = make(chan struct{})
		osSignal = make(chan os.Signal, 1)
	)

	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	var (
		r = http.NewServeMux()
		h = handler.NewHandler(logger, anagram.NewAnagram())
	)

	r.Handle("/anagrams", h)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", *portPtr),
	}

	go func() {
		sig := <-osSignal

		log.Printf("os signal received: %s\n", sig.String())
		log.Printf("perform graceful shutdown\n")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("http server Shutdown: %v\n", err)
		}

		close(done)
	}()

	log.Printf("starting http server at port: %d\n", *portPtr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Panicf("http server serve error: %v\n", err)
	}

	<-done
}
