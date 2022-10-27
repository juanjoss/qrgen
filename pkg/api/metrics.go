package api

import (
	"context"
	"log"

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
	s.meter = meterProvider.Meter("github.com/juanjoss/qrgen/pkg/api")

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
}
