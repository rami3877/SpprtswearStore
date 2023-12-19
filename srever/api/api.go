package api

import (
	"db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Api struct {
	admin
	user
}

func InitApi() *Api {
	return &Api{}
}

func (api *Api) Setup(server *gin.Engine) {
	api.setUserApi(server)
	api.setGuestApi(server)
	api.setAdminApi(server)
}

func (api *Api) setUserApi(server *gin.Engine) {
	api.user.userGroup = server.Group("/user")
	api.user.setLoginApi()
	api.user.setLogoutApi()
	api.user.setMiddleware()
	api.user.setCommintApi()
	api.user.setRegisterApi()
	api.user.setCheckoutApi()
	api.user.setInformationApi()
}

func (api *Api) setAdminApi(server *gin.Engine) {
	api.admin.adminGroup = server.Group("/admin")
	api.admin.setLoginApi()
	api.admin.setOrderApi()
	api.admin.setLogoutApi()
	api.admin.setMiddleware()
	api.admin.setProductApi()
	api.admin.GetUsers()
	// api.admin.setAdminPage()
}

func (api *Api) setGuestApi(server *gin.Engine) {

	server.GET("/product", func(ctx *gin.Context) {
		kind := ctx.Query("kind")
		Containter := ctx.Query("container")
		id, err := strconv.Atoi(ctx.Query("id"))
		if err != nil {
			ctx.String(http.StatusBadRequest, "ERROR")
			return
		}

		data, err := db.MainDB.Stock.GetModelsInKind(id, 10, Containter, kind)
		if err != nil {
			ctx.String(http.StatusBadRequest, "ERROR")
			return
		}
		ctx.JSON(http.StatusAccepted, data)

	})

	server.GET("/AllContainerAndKind", func(ctx *gin.Context) {
		fg := db.MainDB.Stock.GetAllContainerAndKind()
		if len(fg) == 0 {
			ctx.String(http.StatusNoContent, "")
		} else {
			ctx.JSON(http.StatusOK, &fg)
		}

	})

	server.GET("/", func(ctx *gin.Context) {
		 isGust := false
		 if _, err := ctx.Cookie("session") ; err != nil {
			  isGust = true
		 }
		 ctx.HTML(http.StatusOK, "index.html", gin.H{"guest":isGust})
	})



}
