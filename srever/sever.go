package srever

import (
	"net/http"
	"srever/api"

	"github.com/gin-gonic/gin"
)

type Srever int

func (*Srever) Run() {

	//gin.SetMode(gin.ReleaseMode)
	serverEngin := gin.New()
	serverEngin.Use(gin.Logger())
	serverEngin.Use(gin.Recovery())

	httpServer := &http.Server{
		 Addr:    ":8080",
		Handler: serverEngin,
	}

	serverEngin.NoRoute(gin.WrapH(http.FileServer(http.Dir("public"))))
	serverEngin.LoadHTMLGlob("templates/*.html")
	api.InitApi().Setup(serverEngin)

	httpServer.ListenAndServe()

}

func InitSever() *Srever {
	return new(Srever)
}
