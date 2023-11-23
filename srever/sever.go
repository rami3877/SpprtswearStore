package srever

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Srever int

var sever *gin.Engine

func (*Srever) Run() {

	sever = gin.Default()
	sever.NoRoute(gin.WrapH(http.FileServer(http.Dir("public"))))
	sever.LoadHTMLGlob("./public/*/*.html")

	sever.POST("/login", loginPost)
	sever.GET("/login", loginGet)
	sever.Run(":8080")

}

func loginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginPost(c *gin.Context) {
	_, srt := c.GetPostForm("password")
	fmt.Println(srt)
}

func InitSever() *Srever {
	return new(Srever)
}
