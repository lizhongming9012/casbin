package v1

import (
	"NULL/casbin/models"
	"NULL/casbin/pkg/app"
	"NULL/casbin/pkg/e"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type RoleForm struct {
	RoleKey  string `json:"role_key" valid:"Required"`  //角色代码
	RoleName string `json:"role_name" valid:"Required"` //角色名称
}
type PolicyForm struct {
	RoleKey string `json:"role_key" valid:"Required"` //角色代码
	Path    string `json:"path" valid:"Required"`     //请求的path
	Method  string `json:"method" valid:"Required"`   //请求的方法
}

type RoleForUserForm struct {
	Userid  string `json:"userid" valid:"Required"`   //用户ID
	RoleKey string `json:"role_key" valid:"Required"` //角色代码
}

//创建角色
func AddRole(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form RoleForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//判断角色是否存在
	role, err := models.GetRole(form.RoleKey)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if role != nil {
		appG.Response(http.StatusOK, e.ERROR, "the role has exist")
		return
	}
	if err = models.AddRole(&form); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//为角色赋于权限
func AddRolePolicy(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form PolicyForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	res := models.AddRolePolicy(form.RoleKey, form.Path, form.Method)
	if !res {
		appG.Response(http.StatusOK, e.ERROR, "rule already exists")
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除角色的某一权限
func RemoveRolePolicy(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form PolicyForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	res := models.RemoveRolePolicy(form.RoleKey, form.Path, form.Method)
	if !res {
		appG.Response(http.StatusOK, e.ERROR, "the role does not have the policy")
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除一个角色
func DeleteRole(c *gin.Context) {
	appG := app.Gin{C: c}
	roleKey := c.Query("roleKey")
	if err := models.DeleteRole(roleKey); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//为用户添加角色
func AddRoleForUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form RoleForUserForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	res := models.AddRoleForUser(form.Userid, form.RoleKey)
	if !res {
		log.Println("the user already has the role")
	}
	//在user_role表中添加用户角色
	err := models.AddUserRole(&form)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除用户的某一角色
func DeleteRoleForUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form RoleForUserForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	res := models.DeleteRoleForUser(form.Userid, form.RoleKey)
	if !res {
		log.Println("the user does not have the role")
	}
	//在user_role表中删除用户角色
	err := models.DeleteUserRole(form.Userid, form.RoleKey)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除用户的所有角色
func DeleteRolesForUser(c *gin.Context) {
	appG := app.Gin{C: c}
	userid := c.Query("userid")
	res := models.DeleteRolesForUser(userid)
	if !res {
		log.Println("the user does not have any roles")
	}
	//在user_role表中删除用户所有角色
	err := models.DeleteUserRoles(userid)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取用户具有的角色
func GetRolesForUser(c *gin.Context) {
	appG := app.Gin{C: c}
	userid := c.Query("userid")
	res, err := models.GetRolesForUser(userid)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(res) > 0 {
		su := make([]*models.Role, 0)
		fa := make([]string, 0)
		//遍历roleKeys,取Role表对象结果集
		for _, roleKey := range res {
			role, err := models.GetRole(roleKey)
			if err != nil {
				log.Println(err)
				fa = append(fa, roleKey)
			}
			if role != nil {
				su = append(su, role)
			}
		}
		if len(fa) > 0 {
			appG.Response(http.StatusOK, e.ERROR,
				map[string]interface{}{"success": su, "fail": fa})
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, su)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取具有角色的所有用户
func GetUsersForRole(c *gin.Context) {
	appG := app.Gin{C: c}
	roleKey := c.Query("roleKey")
	res, err := models.GetUsersForRole(roleKey)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(res) > 0 {
		su := make([]*models.User, 0)
		fa := make([]string, 0)
		//遍历userids,取User表对象结果集
		for _, userid := range res {
			user, err := models.GetUser(userid)
			if err != nil {
				log.Println(err)
				fa = append(fa, userid)
			}
			if user != nil {
				su = append(su, user)
			}
		}
		if len(fa) > 0 {
			appG.Response(http.StatusOK, e.ERROR,
				map[string]interface{}{"success": su, "fail": fa})
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, su)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取所有角色列表
func GetAllRoles(c *gin.Context) {
	appG := app.Gin{C: c}
	res := models.GetAllRoles()
	if len(res) > 0 {
		su := make([]*models.Role, 0)
		fa := make([]string, 0)
		//遍历roleKeys,取Role表对象结果集
		for _, roleKey := range res {
			role, err := models.GetRole(roleKey)
			if err != nil {
				log.Println(err)
				fa = append(fa, roleKey)
			}
			if role != nil {
				su = append(su, role)
			}
		}
		if len(fa) > 0 {
			appG.Response(http.StatusOK, e.ERROR,
				map[string]interface{}{"success": su, "fail": fa})
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, su)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
