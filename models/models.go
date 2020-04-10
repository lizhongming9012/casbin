package models

import (
	"NULL/casbin/pkg/setting"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var db *gorm.DB

// Setup initializes the database instance
func Setup() {
	var err error
	// 获取数据库连接
	db, err = gorm.Open(setting.DatabaseSetting.Type,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			setting.DatabaseSetting.User,
			setting.DatabaseSetting.Password,
			setting.DatabaseSetting.Host,
			setting.DatabaseSetting.Name))

	// 将数据库连接同步给插件，插件用来操作数据库
	//注意：该插件会根据CasbinRule表是否存在,自行创建CasbinRule表
	Enforcer, err = Casbin()

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	//db.LogMode(true)
	db.SingularTable(true)
	CheckTable()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CheckTable() {
	if !db.HasTable("dept") {
		db.CreateTable(Dept{})
	} else {
		db.AutoMigrate(Dept{})
	}
	if !db.HasTable("user") {
		db.CreateTable(User{})
	} else {
		db.AutoMigrate(User{})
	}
	if !db.HasTable("user_role") {
		db.CreateTable(UserRole{})
	} else {
		db.AutoMigrate(UserRole{})
	}
	if !db.HasTable("role") {
		db.CreateTable(Role{})
	} else {
		db.AutoMigrate(Role{})
	}
}
