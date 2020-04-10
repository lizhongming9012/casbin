package models

type UserRole struct {
	ID      uint   `gorm:"primary_key"`
	UserID  string `json:"userid" gorm:"column:userid;COMMENT:'用户ID'"`
	RoleKey string `json:"role_key" gorm:"column:role_key;COMMENT:'角色代码'"`
}

func AddUserRole(data interface{}) error {
	if err := db.Table("user_role").Create(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUserRole(userid, roleKey string) error {
	if err := db.Where("userid=? and role_key=?", userid, roleKey).
		Delete(UserRole{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUserRoles(userid string) error {
	if err := db.Where("userid=? ", userid).Delete(UserRole{}).Error; err != nil {
		return err
	}
	return nil
}
