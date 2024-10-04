package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "portfolio/portfolio"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type grpcServer struct {
	pb.UnimplementedPortfolioServer
}

func (s *grpcServer) GetPerformance(req *pb.GetPerformanceRequest, stream grpc.ServerStreamingServer[pb.GetPerformanceResponse]) error {
	log.Printf("Received: %v", req)
	for {
		if err := stream.Send(&pb.GetPerformanceResponse{
			Nav: 0,
			Pnl: 0,
		}); err != nil {
			log.Printf("Failed to send: %v", err)
			return err
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPortfolioServer(s, &grpcServer{})
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
