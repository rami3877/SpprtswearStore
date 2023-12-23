package srever

import (
	"srever/api"

	"github.com/gin-gonic/gin"
)

type Srever int

func (*Srever) Run() {

	//gin.SetMode(gin.ReleaseMode)
	serverEngin := gin.New()
	serverEngin.Use(gin.Logger())
	serverEngin.Use(gin.Recovery())
	serverEngin.LoadHTMLGlob("templates/*.html")
	api.InitApi().Setup(serverEngin)

	serverEngin.Run(":8080")

}

func InitSever() *Srever {
	return new(Srever)
}
