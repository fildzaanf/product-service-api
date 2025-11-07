package client

import (
	"context"
	"fmt"
	"time"

	"product-service-api/internal/product/adapter/client/pb"
	"product-service-api/internal/product/application/port"
	"product-service-api/pkg/middleware"

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

func (c *userGRPCClient) GetUserByID(ctx context.Context, userID string) (*pb.UserResponse, error) {
	fmt.Println("========== [DEBUG][userGRPCClient.GetUserByID] START ==========")
	fmt.Printf("[DEBUG] Target userID: %s\n", userID)

	// Deklarasi token di awal
	var token string

	// Ambil token dari metadata incoming context (jika ada)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if authHeader, exists := md["authorization"]; exists && len(authHeader) > 0 {
			token = authHeader[0]
			fmt.Printf("[DEBUG] Token ditemukan di incoming metadata: %s\n", token)
		}
	}

	// Fallback: ambil dari ctx.Value jika metadata kosong
	if token == "" {
		if t, ok := ctx.Value(middleware.ClaimToken).(string); ok {
			token = t
			fmt.Printf("[DEBUG] Token ditemukan di ctx.Value: %s\n", token)
		}
	}

	// Jika tetap tidak ada token
	if token == "" {
		fmt.Println("[ERROR] Missing raw token in context or metadata!")
		return nil, status.Error(codes.Unauthenticated, "missing raw token in context")
	}

	fmt.Printf("[DEBUG] Extracted token from context: %s\n", token)

	// Buat timeout
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Tambahkan token ke metadata outgoing
	md := metadata.New(map[string]string{
		"authorization": token, // rawToken sudah termasuk "Bearer ..."
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Debug metadata outgoing
	if mdOut, ok := metadata.FromOutgoingContext(ctx); ok {
		fmt.Printf("[DEBUG] Outgoing metadata from context: %+v\n", mdOut)
	} else {
		fmt.Println("[ERROR] No outgoing metadata found in context")
	}

	// Buat request gRPC
	userRequest := &pb.GetUserByIDRequest{Id: userID}
	fmt.Printf("[DEBUG] Sending gRPC request to user-service: %+v\n", userRequest)

	// Panggil user-service
	userResp, err := c.client.GetUserByID(ctx, userRequest)
	if err != nil {
		fmt.Printf("[ERROR] gRPC call failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("[DEBUG] gRPC call success. Response: %+v\n", userResp)
	fmt.Println("========== [DEBUG][userGRPCClient.GetUserByID] END ==========")
	return userResp, nil
}
