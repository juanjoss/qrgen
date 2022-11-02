package grpc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/sdk/metric"
)

type metrics struct {
	qrgenRequestsCounter syncfloat64.Counter
}

func (s *server) initMeter() (func(context.Context) error, error) {
	metricExporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(metricExporter))
	s.meter = meterProvider.Meter("qrgen-meter")

	s.initMetrics()

	return meterProvider.Shutdown, nil
}

func (s *server) initMetrics() {
	qrgenRC, err := s.meter.SyncFloat64().Counter("qrgen-requests-counter")
	if err != nil {
		log.Fatalf("unable to create qrgen requests counter: %v", err)
	}

	s.metrics = &metrics{
		qrgenRequestsCounter: qrgenRC,
	}

	go s.serveMetrics()
}

func (s *server) serveMetrics() {
	log.Printf("serving metrics at port %s /metrics", os.Getenv("HTTP_METRICS_PORT"))

	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("HTTP_METRICS_PORT")), nil)
	if err != nil {
		log.Printf("unable to serve metrics: %v", err)
	}
}
