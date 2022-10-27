package grpc

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"log"
	"net"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	pb "github.com/juanjoss/qrgen/pkg/grpc/qrgen"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedQrGeneratorServer
}

func (s *server) GenerateQR(ctx context.Context, req *pb.QrGenRequest) (*pb.QrImage, error) {
	log.Printf("GRPC request: generating QR code for source %v", req.Source)

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

func ListenAndServe(port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("unable to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterQrGeneratorServer(s, &server{})

	log.Printf("GRPC server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("unable to serve: %v", err)
	}
}
