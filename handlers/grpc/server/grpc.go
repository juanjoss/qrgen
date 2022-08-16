package grpc

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"log"
	"net"

	pb "github.com/juanjoss/qrgen/handlers/grpc/qrgen"
	"github.com/juanjoss/qrgen/qr"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedQrGeneratorServer
}

func (s *server) GenerateQR(ctx context.Context, req *pb.QrGenRequest) (*pb.QrImage, error) {
	log.Printf("GRPC request: generating QR code for source %v", req.Source)

	qr, err := qr.Generate(req.Source)
	if err != nil {
		return &pb.QrImage{}, err
	}

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, *qr)
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
