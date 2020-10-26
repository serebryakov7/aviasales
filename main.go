package main

import (
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

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	h := handler.NewHandler(logger, anagram.NewAnagram())

	r := http.NewServeMux()

	r.Handle("/anagrams", h)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", *portPtr),
	}

	go func() {
		log.Printf("starting http server at port: %d", *portPtr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Panicf("http server serve error: %v", err)
		}
	}()

	<-osSignal
}
