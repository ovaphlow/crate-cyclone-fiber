package main

import (
	"context"
	"log"
	"net"
	pb_schema "ovaphlow/cratecyclone/schema"
	"ovaphlow/cratecyclone/utilities"

	"google.golang.org/grpc"
)

type server struct {
	pb_schema.UnimplementedSchemaServer
}

func (s *server) RetrieveSchema(ctx context.Context, in *pb_schema.Empty) (*pb_schema.RetrieveSchemaReply, error) {
	return &pb_schema.RetrieveSchemaReply{Name: []string{"test", "test1"}}, nil
}

func GRPCServe(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb_schema.RegisterSchemaServer(s, &server{})
	go s.Serve(lis)
	utilities.Slogger.Info("GRPC 服务运行于端口 50051")
}
