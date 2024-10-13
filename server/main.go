package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "portfolio/portfolio"
	"portfolio/server/model"
	"portfolio/server/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type grpcServer struct {
	pb.UnimplementedPortfolioServer
	uc usecase.StreamPortUC
}

func NewGrpcServer() pb.PortfolioServer {
	uc := usecase.NewStreamPortUC()
	return &grpcServer{
		uc: uc,
	}
}

func (s *grpcServer) GetPerformance(req *pb.GetPerformanceRequest, stream grpc.ServerStreamingServer[pb.GetPerformanceResponse]) error {
	log.Printf("Received: %v", req)
	port := model.NewPortfolio(1000000, []model.Position{
		{
			Symbol: "SSI",
			Qty:    100,
		},
		{
			Symbol: "VND",
			Qty:    200,
		},
	})
	navCh, err := s.uc.GetNAV(stream.Context(), port)
	if err != nil {
		log.Printf("Failed to get NAV: %v", err)
		return err
	}

	for nav := range navCh {
		if err := stream.Send(&pb.GetPerformanceResponse{
			Nav: nav,
			Pnl: 0,
		}); err != nil {
			log.Printf("Failed to send: %v", err)
			return err
		}
	}
	return nil
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
