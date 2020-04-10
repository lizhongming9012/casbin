package v1

import (
	"NULL/casbin/models"
	"NULL/casbin/pkg/app"
	"NULL/casbin/pkg/e"
	"NULL/casbin/pkg/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type LoginForm struct {
	Username string `json:"username" valid:"Required"`
	Password string `json:"password" valid:"Required"`
}

type AddUserForm struct {
	Username string `json:"username" valid:"Required"`
	Password string `json:"password" valid:"Required"`
	Name     string `json:"name" valid:"Required"`
	Mobile   string `json:"mobile" valid:"Required"`
	DeptID   uint   `json:"dept_id" valid:"Required"`
}

type UpdateUserForm struct {
	ID       string `json:"userid" valid:"Required"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	DeptID   uint   `json:"dept_id"`
}

func Login(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		session = sessions.Default(c)
		form    LoginForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		log.Printf("BIND: %v", errCode)
		appG.Response(httpCode, errCode, nil)
		return
	}
	user, err := models.Login(form.Username, form.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	//生成token
	token, err := util.GenerateToken(form.Username, form.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	//将用户角色存储到session
	roles := make([]string, 0)
	for _, role := range user.UserRoles {
		roles = append(roles, role.RoleKey)
	}
	session.Set("role", roles)
	if err := session.Save(); err != nil {
		log.Printf("session.Save() err:%v", err.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"token": token,
		"data":  user,
	})

}

//增加人员
func AddUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddUserForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	user := models.User{
		ID:       form.Username,
		Username: form.Username,
		Password: form.Password,
		Name:     form.Name,
		Mobile:   form.Mobile,
		DeptID:   form.DeptID,
	}
	if err := models.AddUser(&user); err != nil {
		if strings.Contains(err.Error(), "for key 'PRIMARY'") {
			appG.Response(http.StatusOK, e.ERROR, "the username has exist")
			return
		}
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_USER_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改人员
func UpdateUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form UpdateUserForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	user := models.User{
		ID:       form.ID,
		Username: form.Username,
		Password: form.Password,
		Name:     form.Name,
		Mobile:   form.Mobile,
		DeptID:   form.DeptID,
	}
	if err := models.UpdateUser(&user); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取部门人员列表
func GetDeptUserList(c *gin.Context) {
	appG := app.Gin{C: c}
	deptId, err := strconv.Atoi(c.Query("dept_id"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	dus, err := models.GetUserList(uint(deptId))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(dus) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, dus)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

//删除人员
func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}
	if err := models.DeleteUser(c.Query("userid")); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
