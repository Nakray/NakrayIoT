package main

import (
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"NakrayIoT/internal/redis"
	"NakrayIoT/internal/server"
	tempPb "NakrayIoT/proto"
)

func main() {
	log.SetOutput(os.Stderr)
	redisClient := redis.NewRedisClient(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	defer redisClient.Close()

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	grpcServer := grpc.NewServer()
	tempService := server.NewTemperatureService(redisClient)
	tempPb.RegisterTemperatureServiceServer(grpcServer, tempService)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
