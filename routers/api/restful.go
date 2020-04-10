package api

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
)

func Home(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static/apidoc.html")
		if err := t.Execute(c.Writer, nil); err != nil {
			log.Printf("err:-->%v", err)
		}
	}
}
