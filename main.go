package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var flagAddr = flag.String("addr", ":3030", "specify the address and port to listen on")

func main() {
	flag.Parse()

	store := NewCache()
	a := AppServer{
		Store: store,
		Addr:  *flagAddr,
	}

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-done
		log.Printf("received signal: %v", s)
		log.Println("shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := a.Shutdown(ctx); err != nil {
			log.Fatalf("error shutting down server: %v", err)
		}
		log.Println("done shutting down")
	}()
	log.Printf("starting server on address %s", *flagAddr)
	if err := a.Run(); err != nil {
		log.Fatalf("error running http server: %v", err)
	}
}
