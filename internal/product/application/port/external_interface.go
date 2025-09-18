package port

import (
	"context"
	"product-service-api/internal/product/adapter/client/pb"
)

type UserQueryClientInterface interface {
	GetUserByID(ctx context.Context, id string) (*pb.UserResponse, error)
}
