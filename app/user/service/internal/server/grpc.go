package server

import (
	userv1 "om-platform/api/user/service/v1"
	helloworldv1 "om-platform/app/user/service/api/helloworld/v1"
	"om-platform/app/user/service/internal/conf"
	"om-platform/app/user/service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService, userSvc *service.UserManagementService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	// 注册Greeter服务
	helloworldv1.RegisterGreeterServer(srv, greeter)
	// 注册用户管理服务
	userv1.RegisterUserManagementServiceServer(srv, userSvc)
	return srv
}
