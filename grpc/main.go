package main

import (
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}
