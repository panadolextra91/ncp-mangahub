package grpc

import (
	"context"
	"strings"

	"github.com/user/mangahub/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtSecret string
}

func NewAuthInterceptor(secret string) *AuthInterceptor {
	return &AuthInterceptor{jwtSecret: secret}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		claims, err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		// Add claims to context for service use
		newCtx := context.WithValue(ctx, "claims", claims)
		return handler(newCtx, req)
	}
}

func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		_, err := i.authorize(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) (*auth.Claims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is missing")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header is missing")
	}

	// Support "Bearer <token>"
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := auth.ValidateToken(token, i.jwtSecret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// RBAC: AdminService requires Admin role
	if strings.Contains(method, "AdminService") {
		if claims.Role != "admin" {
			return nil, status.Error(codes.PermissionDenied, "admin role required")
		}
	}

	return claims, nil
}
