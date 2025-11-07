package middleware

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)
type claimKey string

const (
	ClaimUserID claimKey = "user_id"
	ClaimRole   claimKey = "role"
	ClaimToken  claimKey = "token"
)

func JWTUnaryInterceptor() grpc.UnaryServerInterceptor {
	publicMethodsEnv := os.Getenv("PUBLIC_METHODS")
	publicMethods := strings.Split(publicMethodsEnv, ",")

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		for _, method := range publicMethods {
			if info.FullMethod == strings.TrimSpace(method) {
				return handler(ctx, req)
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		rawToken := authHeader[0]

		tokenString := strings.TrimSpace(strings.TrimPrefix(rawToken, "Bearer "))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}

		id, _ := claims["id"].(string)
		role, _ := claims["role"].(string)

		if id == "" || role == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}

		ctx = context.WithValue(ctx, ClaimUserID, id)
		ctx = context.WithValue(ctx, ClaimRole, role)
		ctx = context.WithValue(ctx, ClaimToken, tokenString)

		return handler(ctx, req)
	}
}

func ExtractTokenFromContext(ctx context.Context) (string, string, string, error) {
	id, ok := ctx.Value(ClaimUserID).(string)
	if !ok {
		return "", "", "", errors.New("user id not found in context")
	}

	role, ok := ctx.Value(ClaimRole).(string)
	if !ok {
		return "", "", "", errors.New("role not found in context")
	}

	token, ok := ctx.Value(ClaimToken).(string)
	if !ok || token == "" {
		return "", "", "", errors.New("raw token not found in context")
	}

	return id, role, token, nil
}
