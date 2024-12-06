package main

import (
	"context"
	"log"

	"github.com/PerfectStepCoder/shorturl/internal/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

// UserKeyUID - идентификатор пользователя который передается в контексте.
const UserKeyUID contextKey = "userUID"


// UnaryInterceptorAuth - перехватчик для авторизации
func UnaryInterceptorAuth(ignoredMethods map[string]struct{}) grpc.UnaryServerInterceptor {
	return func (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		
        // Проверяем, если метод в списке игнорируемых
		log.Printf("Info: %v", info.FullMethod)
		log.Print(ignoredMethods)

        if _, ok := ignoredMethods[info.FullMethod]; ok {
            // Пропускаем проверку
            return handler(ctx, req)
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


// UnaryInterceptorTrustedSubnet - проверка принадлежности IP адреса к доверенной подсети.
func UnaryInterceptorTrustedSubnet(ignoredMethods map[string]struct{}, trustedSubnet string) grpc.UnaryServerInterceptor {
	return func (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		
        // Проверяем, если метод в списке игнорируемых
        if _, ok := ignoredMethods[info.FullMethod]; ok {
            // Пропускаем проверку
            return handler(ctx, req)
        }

		var realIP string
		
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("X-Real-IP")
			if len(values) > 0 {
				realIP = values[0]
			}
		}
	
		if len(realIP) == 0 {
			return nil, status.Error(codes.InvalidArgument, "missing X-Real-IP")
		}
	
		if allowIP, err := handlers.IsIPInCIDR(realIP, trustedSubnet); err != nil {
			return nil, status.Error(codes.PermissionDenied, "invalid X-Real-IP");
		} else {
			if !allowIP {
				return nil, status.Error(codes.PermissionDenied, "X-Real-IP not allowed");
			}
		}

		return handler(ctx, req)
	}
}