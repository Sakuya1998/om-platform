package data

import (
	"om-platform/internal/biz"
	"om-platform/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewGreeterRepo,
	NewUserRepo,
	NewRoleRepo,
	NewPermissionRepo,
	NewUserRoleRepo,
	NewRolePermissionRepo,
)

// Data .
type Data struct {
	// TODO wrapped database client
	userRepo           biz.UserRepo
	roleRepo           biz.RoleRepo
	permissionRepo     biz.PermissionRepo
	userRoleRepo       biz.UserRoleRepo
	rolePermissionRepo biz.RolePermissionRepo
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	logger = log.With(logger, "module", "data")
	logHelper := log.NewHelper(logger)

	// 创建Data实例
	d := &Data{}

	// 创建各个仓储实例
	d.userRepo = NewUserRepo(d, logger)
	d.roleRepo = NewRoleRepo(d, logger)
	d.permissionRepo = NewPermissionRepo(d, logger)
	d.userRoleRepo = NewUserRoleRepo(d, logger)
	d.rolePermissionRepo = NewRolePermissionRepo(d, logger)

	cleanup := func() {
		logHelper.Info("closing the data resources")
	}

	return d, cleanup, nil
}
