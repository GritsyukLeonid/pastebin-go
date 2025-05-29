package main

import (
	"log"
	"net"

	"github.com/GritsyukLeonid/pastebin-go/internal/grpc/grpcimpl"
	"github.com/GritsyukLeonid/pastebin-go/internal/pb"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	repository.LoadData()

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	reflection.Register(s)

	srv := &grpcimpl.Server{}

	pb.RegisterUserServiceServer(s, srv)
	pb.RegisterPasteServiceServer(s, srv)
	pb.RegisterStatsServiceServer(s, srv)
	pb.RegisterShortURLServiceServer(s, srv)

	log.Println("gRPC server listening on :9090")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
