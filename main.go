package main

import (
	"context"
	"log"
	"net"
	"os"
	"github.com/hrand1005/training-notebook/set"
	"google.golang.org/grpc"
)

type setService struct {
	logger *log.Logger
	set.UnimplementedSetServer
}

func newSetService() *setService {
	l := log.New(os.Stdout, "SET SERVICE: ", log.Ltime)
	return &setService{
		logger: l,
	}
}

func (s *setService) Create(ctx context.Context, req *set.CreateRequest) (*set.CreateResponse, error) {
	s.logger.Printf("Recieved create request:\n%#v", req)
	return &set.CreateResponse{
		Id: req.Id,
		Name: req.Name,
		Reps: req.Reps,
		Intensity: req.Intensity,
		Volume: req.Volume,
	}, nil
}

func (s *setService) Read(ctx context.Context, req *set.ReadRequest) (*set.ReadResponse, error) {
	s.logger.Printf("Recieved read request:\n%#v", req)

	return &set.ReadResponse{}, nil
}

const setServicePort = ":8080"

func main() {
	// define server options and pass them into NewServer
	srv := grpc.NewServer()
	setService := newSetService()
	set.RegisterSetServer(srv, setService)

	lis, err := net.Listen("tcp", setServicePort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s, err: %v", setServicePort, err)
	}
	srv.Serve(lis)
}