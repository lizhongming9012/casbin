package models

import "github.com/jinzhu/gorm"

// 用户
type User struct {
	ID        string `json:"userid" gorm:"primary_key;column:userid"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	DeptID    uint   `json:"dept_id"`
	Dept      Dept
	UserRoles []UserRole
}

func Login(u, p string) (*User, error) {
	var user User
	err := db.Preload("Dept").Preload("UserRoles").
		Where("username = ? and password = ?", u, p).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &user, nil
}

func AddUser(data interface{}) error {
	if err := db.Table("user").Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUser(u *User) error {
	if err := db.Table("user").Where("userid=?", u.ID).Updates(u).Error; err != nil {
		return err
	}
	return nil
}

func GetUserList(dept_id uint) ([]*User, error) {
	var dus []*User
	if err := db.Preload("Dept").
		Where("dept_id=?", dept_id).Find(&dus).Error; err != nil {
		return nil, err
	}
	return dus, nil
}

func GetUser(userid string) (*User, error) {
	var u User
	err := db.Preload("Dept").Where("userid=?", userid).First(&u).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &u, nil
}

func DeleteUser(userid string) error {
	if err := db.Where("userid=?", userid).Delete(User{}).Error; err != nil {
		return err
	}
	//删除用户,并删除用户的所有角色
	DeleteRolesForUser(userid)
	return nil
}

func IsDeptUserExist(depId uint) bool {
	var du User
	if err := db.Where("dept_id=?", depId).First(&du).Error; err != nil {
		return false
	}
	return true
}
