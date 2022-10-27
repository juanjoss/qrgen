package api

import (
	"fmt"
	"image/png"
	"log"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *server) generateQR(w http.ResponseWriter, r *http.Request) {
	_, span := s.tracer.Start(r.Context(), "generateQR")
	defer span.End()

	s.metrics.qrgenRequestsCounter.Add(r.Context(), 1)

	err := r.ParseForm()
	if err != nil {
		span.SetStatus(codes.Error, "unable to parse form")
		span.RecordError(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to parse form: %v", err.Error())
		return
	}
	source := r.FormValue("source")

	span.SetAttributes(attribute.String("source", source))

	encodedData, err := qr.Encode(source, qr.L, qr.Auto)
	if err != nil {
		span.SetStatus(codes.Error, "unable to encode source")
		span.RecordError(err)

		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable to encode source: %v", err)
		return
	}

	barcode, err := barcode.Scale(encodedData, 256, 256)
	if err != nil {
		span.SetStatus(codes.Error, "unable to scale encoded barcode")
		span.RecordError(err)

		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable to scale encoded barcode: %v", err)
		return
	}

	png.Encode(w, barcode)
}
