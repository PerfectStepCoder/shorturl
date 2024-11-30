package main

import (
	"context"
	"strings"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
	"google.golang.org/grpc"
)

type contextKey string

// UserKeyUID - идентификатор пользователя который передается в контексте.
const UserKeyUID contextKey = "userUID"

// UnaryInterceptorAuth - перехватчик для авторизации
func UnaryInterceptorAuth(ignoredMethods []string) grpc.UnaryServerInterceptor {
	return func (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		
		// Проверяем, если метод в списке игнорируемых
		for _, method := range ignoredMethods {
			if strings.HasSuffix(info.FullMethod, method) {
				// Пропускаем проверку
				return handler(ctx, req)
			}
		}

		var token string
		
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("token")
			if len(values) > 0 {
				token = values[0]
			}
		}
	
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}
	
		if userUID, err := ExtractUserUID(token); err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		} else {
			ctx = context.WithValue(ctx, UserKeyUID, userUID)
		}
		
		return handler(ctx, req)
	}
}

