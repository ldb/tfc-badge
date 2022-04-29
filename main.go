package main

import (
	"context"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var flagAddr = flag.String("addr", ":3030", "specify the address and port to listen on")
var metricsAddr = flag.String("metrics", ":9080", "specify the address and port to listen on for metrics")
var debug = flag.Bool("debug", false, "enable debug log output")

func main() {
	flag.Parse()

	store := NewCache()
	a := AppServer{
		Store: store,
		Debug: *debug,
		Addr:  *flagAddr,
	}

	mc := &MetricsCollector{
		Store: store,
	}
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(mc)

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
	go func() {
		if err := a.Run(); err != nil {
			log.Fatalf("error running http server: %v", err)
		}
	}()
	log.Printf("starting server on address %s", *metricsAddr)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(*metricsAddr, promhttp.HandlerFor(reg, promhttp.HandlerOpts{})); err != nil {
		log.Fatalf("error running metrics server: %v", err)
	}
}
