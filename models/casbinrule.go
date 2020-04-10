package models

import (
	"fmt"
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
)

var Enforcer *casbin.Enforcer

//从数据库中读取模型文件对应的策略信息。Casbin 根据策略字段判断是否有权限
func Casbin() (*casbin.Enforcer, error) {
	//将数据库连接同步给插件，插件用来操作数据库
	Apter := gormadapter.NewAdapterByDB(db)
	Enforcer := casbin.NewEnforcer("conf/rbac.conf", Apter)
	//开启权限认证日志
	//Enforcer.EnableLog(true)
	//加载数据库中的策略
	if err := Enforcer.LoadPolicy(); err != nil {
		fmt.Printf("casbin rbac model or policy init error,message:%v", err)
		return nil, err

	}
	return Enforcer, nil
}

//数据库操作直接调用Casbin官方API

//创建角色,并赋于权限
func AddRolePolicy(roleKey, path, method string) bool {
	return Enforcer.AddPolicy(roleKey, path, method)
}

//删除角色的某一权限
func RemoveRolePolicy(roleKey, path, method string) bool {
	return Enforcer.RemovePolicy(roleKey, path, method)
}

//删除一个角色
func RemoveRole(roleKey string) {
	Enforcer.DeleteRole(roleKey)
}

//为用户添加角色。 如果用户已经拥有该角色（aka不受影响），则返回false。
func AddRoleForUser(userid, roleKey string) bool {
	return Enforcer.AddRoleForUser(userid, roleKey)
}

//删除用户的某一角色。 如果用户没有该角色（aka不受影响），则返回false
func DeleteRoleForUser(userid, roleKey string) bool {
	return Enforcer.DeleteRoleForUser(userid, roleKey)
}

//删除用户的所有角色。 如果用户没有任何角色（aka不受影响），则返回false。
func DeleteRolesForUser(userid string) bool {
	return Enforcer.DeleteRolesForUser(userid)
}

//获取用户具有的角色
func GetRolesForUser(userid string) ([]string, error) {
	return Enforcer.GetRolesForUser(userid)
}

//获取具有角色的所有用户
func GetUsersForRole(roleKey string) ([]string, error) {
	return Enforcer.GetUsersForRole(roleKey)
}

//获取当前策略中显示的角色列表
func GetAllRoles() []string {
	return Enforcer.GetAllRoles()
}
