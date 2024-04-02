package main

import (
	"context"
	"log"
	"net"
	pb_schema "ovaphlow/cratecyclone/schema"
	"ovaphlow/cratecyclone/utility"

	"google.golang.org/grpc"
)

type server struct {
	pb_schema.UnimplementedSchemaServer
}

func (s *server) RetrieveSchema(ctx context.Context, in *pb_schema.Empty) (*pb_schema.RetrieveSchemaReply, error) {
	return &pb_schema.RetrieveSchemaReply{Schema: []string{"test", "test1"}}, nil
}

func (s *server) RetrieveTable(ctx context.Context, in *pb_schema.RetrieveTableRequest) (*pb_schema.RetrieveTableReply, error) {
	return &pb_schema.RetrieveTableReply{Table: []string{in.Schema, "test table", "test1 table"}}, nil
}

func GRPCServe(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb_schema.RegisterSchemaServer(s, &server{})
	go s.Serve(lis)
	utility.Slogger.Info("GRPC 服务运行于端口 50051")
}
