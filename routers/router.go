package routers

import (
	"NULL/casbin/middleware/auth"
	"NULL/casbin/middleware/cors"
	"NULL/casbin/middleware/jwt"
	"NULL/casbin/pkg/export"
	"NULL/casbin/pkg/qrcode"
	"NULL/casbin/pkg/upload"
	"NULL/casbin/routers/api"
	v1 "NULL/casbin/routers/api/v1"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	//创建基于cookie的存储引擎，secret 参数是用于加密的密钥
	store := cookie.NewStore([]byte("secret"))
	//设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	//store是前面创建的存储引擎，也可以替换成其他存储引擎
	r.Use(sessions.Sessions("mysession", store))

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.CORSMiddleware())

	//API doc
	r.GET("/", api.Home)

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))
	
	//加载 RESTful API doc所需资源
	r.Static("/css", "runtime/static/css")
	r.Static("/img", "runtime/static/img")
	r.Static("/js", "runtime/static/js")

	//上传文件
	r.POST("/file/upload", api.UploadFile)

	//登陆
	r.POST("/login", v1.Login)

	apiv1 := r.Group("/api/v1")
	//token认证
	apiv1.Use(jwt.JWT())
	//权限控制
	apiv1.Use(auth.AuthCheckRole())
	{
		//增加部门
		apiv1.POST("/dept/add", v1.AddDept)
		//删除部门
		apiv1.GET("/dept/del", v1.DeleteDept)
		//修改部门
		apiv1.POST("/dept/upd", v1.UpdateDept)
		//获取部门树
		apiv1.GET("/dept/tree", v1.GetDeptTree)

		//增加人员
		apiv1.POST("/user/adduser", v1.AddUser)
		//删除人员
		apiv1.GET("/user/del", v1.DeleteUser)
		//修改人员
		apiv1.POST("/user/upduser", v1.UpdateUser)
		//获取部门人员列表
		apiv1.GET("/user/deptuser", v1.GetDeptUserList)

		//为用户添加角色
		apiv1.POST("/user/addrole", v1.AddRoleForUser)
		//删除用户的某一角色
		apiv1.POST("/user/delrole", v1.DeleteRoleForUser)
		//删除用户的所有角色
		apiv1.GET("/user/delroles", v1.DeleteRolesForUser)
		//获取用户具有的角色
		apiv1.GET("/user/roles", v1.GetRolesForUser)

		//创建角色
		apiv1.POST("/role/addrole", v1.AddRole)
		//删除一个角色
		apiv1.GET("/role/delrole", v1.DeleteRole)
		//为角色赋于权限
		apiv1.POST("/role/addpolicy", v1.AddRolePolicy)
		//删除角色的某一权限
		apiv1.POST("/role/delpolicy", v1.RemoveRolePolicy)
		//获取具有角色的所有用户
		apiv1.GET("/role/users", v1.GetUsersForRole)
		//获取所有角色列表
		apiv1.GET("/role/all", v1.GetAllRoles)
	}
	return r
}
