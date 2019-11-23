package main

//go:generate protoc -I pkg/api/v1/ pkg/api/v1/api.proto --go_out=plugins=grpc:pkg/api/v1

import (
	"context"
	pb "github.com/vllry/sensor-collector/pkg/api/v1"
	"github.com/vllry/sensor-collector/pkg/webview"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	grpcPort = ":50051"
)

type server struct {
	pb.PostDataServer
	lastTemperature map[int32]float32
}

func (s *server) PostTemperature(ctx context.Context, in *pb.SensorData) (*pb.DataResponse, error) {
	log.Printf("Received: %v %v", in.GetSensorId(), in.GetValue())
	s.lastTemperature[in.GetSensorId()] = in.GetValue()
	return &pb.DataResponse{Ok: true}, nil
}

func (s *server) PostHumidity(ctx context.Context, in *pb.SensorData) (*pb.DataResponse, error) {
	log.Printf("Received: %v %v", in.GetSensorId(), in.GetValue())
	return &pb.DataResponse{Ok: true}, nil
}

func main() {
	lastTemperatureData := make(map[int32]float32) // TODO use a channel or something.

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPostDataServer(
		s,
		&server{
			lastTemperature: lastTemperatureData,
		},
	)

	go webview.RunWebserver(lastTemperatureData)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
