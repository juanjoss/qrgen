package grpc

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	pb "github.com/juanjoss/qrgen/pkg/grpc/qrgen"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedQrGeneratorServer
	port    string
	tracer  trace.Tracer
	meter   metric.Meter
	metrics *metrics
}

func NewServer(port string) *server {
	return &server{
		port: port,
	}
}

func (s *server) ListenAndServe() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		log.Fatalf("unable to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterQrGeneratorServer(grpcServer, s)

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

	log.Printf("GRPC server listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("unable to serve: %v", err)
	}
}

func (s *server) GenerateQR(ctx context.Context, req *pb.QrGenRequest) (*pb.QrImage, error) {
	_, span := s.tracer.Start(ctx, "generateQR")
	defer span.End()

	s.metrics.qrgenRequestsCounter.Add(ctx, 1)

	encodedData, err := qr.Encode(req.Source, qr.L, qr.Auto)
	if err != nil {
		log.Printf("unable to encode source: %v", err)
		return &pb.QrImage{}, err
	}

	barcode, err := barcode.Scale(encodedData, 256, 256)
	if err != nil {
		log.Printf("unable to scale encoded barcode: %v", err)
		return &pb.QrImage{}, err
	}

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, barcode)
	if err != nil {
		return &pb.QrImage{}, err
	}

	return &pb.QrImage{Barcode: buffer.Bytes()}, nil
}
