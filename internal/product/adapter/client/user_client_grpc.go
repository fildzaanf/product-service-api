package client

import (
	"context"
	"product-service-api/internal/product/adapter/client/pb"
	"product-service-api/internal/product/application/port"
	"product-service-api/pkg/middleware"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type userGRPCClient struct {
	client pb.UserQueryServiceClient
}

func NewUserGRPCClient(conn *grpc.ClientConn) port.UserQueryClientInterface {
	return &userGRPCClient{
		client: pb.NewUserQueryServiceClient(conn),
	}
}

func (c *userGRPCClient) GetUserByID(ctx context.Context, userID string) error {
	var token string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if authHeader, exists := md["authorization"]; exists && len(authHeader) > 0 {
			token = authHeader[0]
		}
	}

	if token == "" {
		if t, ok := ctx.Value(middleware.ClaimToken).(string); ok {
			token = t
		}
	}

	if token == "" {
		return status.Error(codes.Unauthenticated, "missing raw token in context")
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	md := metadata.New(map[string]string{
		"authorization": token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	userRequest := &pb.GetUserByIDRequest{
		Id: userID,
	}

	_, err := c.client.GetUserByID(ctx, userRequest)
	if err != nil {
		return err
	}

	return nil
}
