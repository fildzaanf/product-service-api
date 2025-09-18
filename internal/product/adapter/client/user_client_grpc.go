package client

import (
	"context"
	"time"

	"product-service-api/internal/product/application/port"
	"product-service-api/internal/product/adapter/client/pb"

	"google.golang.org/grpc"
)

type userGRPCClient struct {
	client pb.UserQueryServiceClient
}

func NewUserGRPCClient(conn *grpc.ClientConn) port.UserQueryClientInterface {
	return &userGRPCClient{
		client: pb.NewUserQueryServiceClient(conn),
	}
}

func (c *userGRPCClient) GetUserByID(ctx context.Context, userID string) (*pb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	userRequest := &pb.GetUserByIDRequest{
		   Id: userID,
	}
	return c.client.GetUserByID(ctx, userRequest)
}
