package auth

import (
	"NULL/casbin/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		roles := session.Get("role")
		// 请求的path
		p := c.Request.URL.Path
		// 请求的方法
		m := c.Request.Method
		// 认证检查权限
		var res bool
		for _, role := range roles.([]string) {
			res = models.Enforcer.Enforce(role, p, m)
		}
		if !res {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"msg": "Unauthorized",
				"data":    "",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
