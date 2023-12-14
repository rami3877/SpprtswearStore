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



func InitApi ()*Api{
	 return&Api{}
}

func (api *Api) Setup(serverEngin *gin.Engine) {
	api.setUserApi(serverEngin)
	api.setGuestApi(serverEngin)
	api.setAdminApi(serverEngin)
}

func (api *Api) setUserApi(serverEngin *gin.Engine) {
	api.user.userGroup = serverEngin.Group("/user")
	api.user.setLoginApi()
	api.user.setLogoutApi()
	api.user.setMiddleware()
	api.user.setCommintApi()
	api.user.setRegisterApi()
	api.user.setCheckoutApi()
	api.user.setInformationApi()
}

func (api *Api) setAdminApi(serverEngin *gin.Engine) {
	api.admin.adminGroup = serverEngin.Group("/admin")
	api.admin.setLoginApi()
	api.admin.setOrderApi()
	api.admin.setLogoutApi()
	api.admin.setMiddleware()
	api.admin.setProductApi()
	api.admin.GetUsers()
	// api.admin.setAdminPage()
}
func (api *Api) setGuestApi(serverEngin *gin.Engine) {

	serverEngin.GET("/product", func(ctx *gin.Context) {
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

	serverEngin.GET("/AllContainerAndKind", func(ctx *gin.Context) {
		fg := db.MainDB.Stock.GetAllContainerAndKind()
		if len(fg) == 0 {
			ctx.String(http.StatusNoContent, "")
		} else {
			ctx.JSON(http.StatusOK, &fg)
		}

	})
}
