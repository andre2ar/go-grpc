package main

import (
	"context"
	"fmt"
	"github.com/andre2ar/gRPC-go/pb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User {
		Id: "0",
		Name: "André",
		Email: "andre2ar@outlook.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not add user: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User {
		Id: "0",
		Name: "André",
		Email: "andre2ar@outlook.com",
	}

	resStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not add user: %v", err)
	}

	for {
		stream, err := resStream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Stream error: %v", err)
		}

		fmt.Println(stream)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User {
		&pb.User {
			Id: "0",
			Name: "André",
			Email: "andre2ar@outlook.com",
		},
		&pb.User {
			Id: "0",
			Name: "André2",
			Email: "andre2ar2@outlook.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient)  {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User {
		&pb.User {
			Id: "0",
			Name: "André",
			Email: "andre2ar@outlook.com",
		},
		&pb.User {
			Id: "0",
			Name: "André2",
			Email: "andre2ar2@outlook.com",
		},
	}

	wait := make(chan int)
	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Println(res)
		}
		close(wait)
	}()

	<-wait
}