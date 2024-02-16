package main

import (
	"database/sql"
	"fmt"
	"github.com/andre2ar/go-grpc/internal/database"
	"github.com/andre2ar/go-grpc/internal/pb"
	"github.com/andre2ar/go-grpc/internal/service"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryRepository := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryRepository)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	fmt.Println("gRPC server running on http://localhost:50051")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
