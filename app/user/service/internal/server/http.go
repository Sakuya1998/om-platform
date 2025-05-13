package server

import (
	userv1 "om-platform/api/user/service/v1"
	helloworldv1 "om-platform/app/user/service/api/helloworld/v1"
	"om-platform/app/user/service/internal/conf"
	"om-platform/app/user/service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, userSvc *service.UserManagementService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	// 注册Greeter服务
	helloworldv1.RegisterGreeterHTTPServer(srv, greeter)
	// 注册用户管理服务
	userv1.RegisterUserManagementServiceHTTPServer(srv, userSvc)
	return srv
}
