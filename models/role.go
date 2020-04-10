package models

import "github.com/jinzhu/gorm"

type Role struct {
	ID       uint   `gorm:"primary_key"`
	RoleKey  string `json:"role_key" gorm:"column:role_key;COMMENT:'角色代码'"`
	RoleName string `json:"role_name" gorm:"column:role_name;COMMENT:'角色名称'"`
}

func AddRole(data interface{}) error {
	if err := db.Table("role").Create(data).Error; err != nil {
		return err
	}
	return nil
}

func GetRole(roleKey string) (*Role, error) {
	var r Role
	err := db.Where("role_key=?", roleKey).First(&r).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &r, nil
}

func DeleteRole(roleKey string) error {
	err := db.Where("role_key=?", roleKey).Delete(Role{}).Error
	if err != nil {
		return err
	}
	//删除一个角色,并清空角色策略
	RemoveRole(roleKey)
	return nil
}
