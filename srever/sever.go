package srever

import (
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


func InitSever() *Srever {
	return new(Srever)
}
