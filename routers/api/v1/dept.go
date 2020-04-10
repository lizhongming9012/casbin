package v1

import (
	"NULL/casbin/models"
	"NULL/casbin/pkg/app"
	"NULL/casbin/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AddDeptForm struct {
	Parentid uint   `json:"parentid" valid:"Required"`
	Name     string `json:"name" valid:"Required"`
}

type UpdateDeptForm struct {
	ID       uint   `json:"id" valid:"Required"`
	Parentid uint   `json:"parentid"`
	Name     string `json:"name"`
}

//增加部门
func AddDept(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddDeptForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if err := models.AddDept(&form); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改部门
func UpdateDept(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form UpdateDeptForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	dept := models.Dept{
		ID:       form.ID,
		Parentid: form.Parentid,
		Name:     form.Name,
	}
	if err := models.UpdateDept(&dept); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取部门树
func GetDeptTree(c *gin.Context) {
	var appG = app.Gin{C: c}
	deptId, err := strconv.Atoi(c.Query("dept_id"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	tree, err := models.GetDeptTree(uint(deptId))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, tree)
}

//删除部门
func DeleteDept(c *gin.Context) {
	appG := app.Gin{C: c}
	deptId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if models.IsParent(uint(deptId)) {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_DEVDETP_IS_PARENT, nil)
		return
	}
	if models.IsDeptUserExist(uint(deptId)) {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_DEVDETP_NOT_NULL, nil)
		return
	}
	if err := models.DeleteDept(uint(deptId)); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
