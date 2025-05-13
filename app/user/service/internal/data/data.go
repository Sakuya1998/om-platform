package data

import (
	"om-platform/app/user/service/internal/biz"
	"om-platform/app/user/service/internal/conf"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewUserRepo, NewRoleRepo, NewPermissionRepo)

// Data .
type Data struct {
	// TODO wrapped database client

	// 用户管理相关数据
	users           map[string]*biz.User
	roles           map[string]*biz.Role
	permissions     map[string]*biz.Permission
	userRoles       map[string][]string // userID -> roleIDs
	rolePermissions map[string][]string // roleID -> permissionNames

	mutex sync.RWMutex
}

// NewData 创建数据访问层实例
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	// 初始化数据存储
	return &Data{
		users:           make(map[string]*biz.User),
		roles:           make(map[string]*biz.Role),
		permissions:     make(map[string]*biz.Permission),
		userRoles:       make(map[string][]string),
		rolePermissions: make(map[string][]string),
	}, cleanup, nil
}
