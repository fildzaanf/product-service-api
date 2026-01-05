package port

import (
	"context"
)

type UserQueryClientInterface interface {
	GetUserByID(ctx context.Context, id string) error
}
