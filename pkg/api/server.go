package api

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type server struct {
	router  *mux.Router
	port    string
	fs      http.FileSystem
	tracer  trace.Tracer
	meter   metric.Meter
	metrics *metrics
}

func NewServer(port string, embedFS embed.FS) *server {
	// build FS
	buildFS, err := fs.Sub(embedFS, "view")
	if err != nil {
		log.Fatal(err)
	}

	return &server{
		router: mux.NewRouter(),
		port:   port,
		fs:     http.FS(buildFS),
	}
}

func (s *server) ListenAndServe() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// init tracer
	shutdownTracer, err := s.initTracer()
	if err != nil {
		log.Fatalf("unable to init tracer: %v", err)
	}
	defer func() {
		if err := shutdownTracer(ctx); err != nil {
			log.Fatalf("unable to shutdown tracer provider: %v", err)
		}
	}()

	// init metrics
	shutdownMeter, err := s.initMeter()
	if err != nil {
		log.Fatalf("unable to init meter: %v", err)
	}
	defer func() {
		if err := shutdownMeter(ctx); err != nil {
			log.Fatalf("unable to shutdown metrics provider: %v", err)
		}
	}()

	// register routes
	s.router.HandleFunc("/", http.FileServer(s.fs).ServeHTTP).Methods(http.MethodGet)
	s.router.HandleFunc("/qr", s.generateQR).Methods(http.MethodPost)
	s.router.Handle("/metrics", promhttp.Handler())

	// create logging and recovery middlewares
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.router)
	recoveryRouter := handlers.RecoveryHandler()(loggedRouter)

	// run the http server
	log.Printf("qr service running on port %s", s.port)
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%s", s.port),
			recoveryRouter,
		),
	)
}
