package services

import (
	"context"
	"fmt"
	"github.com/andre2ar/gRPC-go/pb"
	"io"
	"log"
	"time"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	fmt.Println("New user added")

	return &pb.User{
		Id: "123",
		Name: req.GetName(),
		Email: req.GetEmail(),
	}, nil
}

func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResultStream{
		Status: "Init",
		User: &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Inserting",
		User: &pb.User{
			Id: "0",
		},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Concluded",
		User: &pb.User{
			Id: "0",
			Name: req.Name,
			Email: req.Email,
		},
	})

	return nil
}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		users = append(users, &pb.User{
			Id: req.GetId(),
			Name: req.GetName(),
			Email: req.GetEmail(),
		})

		fmt.Println("Adding user", req.GetEmail())
	}

	return nil
}

func (*UserService) AddUserStreamBoth(stream pb.UserService_AddUserStreamBothServer) error  {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error receiving stream from client: %v", err)
		}

		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User: req,
		})
		if err != nil {
			log.Fatalf("Error sending stream for client: %v", err)
		}
	}
}