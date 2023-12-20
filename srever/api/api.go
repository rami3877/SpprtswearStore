package api

import (
	"db"
	"net/http"
	"strconv"
	"strings"

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

	server.NoRoute(func(ctx *gin.Context) {
		Path := ctx.Request.URL.Path[1:]
		if db.MainDB.Stock.IsExist(Path) {
			containerAndkind := db.MainDB.Stock.GetAllContainerAndKind()
			ctx.HTML(http.StatusOK, "selection.html", gin.H{
				"container": ctx.Request.URL.Path[1:],
				"kinds":     containerAndkind[Path],
			},
			)
			ctx.Abort()
			return
		}
		pathArray := strings.Split(Path, "/")
		if len(pathArray) < 2 {
			ctx.Next()
			return
		}
		if listModels, err := db.MainDB.Stock.GetModelsInKind(0, 10, pathArray[0], pathArray[1]); err == nil {
			ctx.HTML(http.StatusOK, "products.html", gin.H{

				"container":  pathArray[0],
				"kind":       pathArray[1],
				"listModels": listModels,
			})

			ctx.Abort()
			return
		}

	}, gin.WrapH(http.FileServer(http.Dir("public"))))
	server.GET("/", func(ctx *gin.Context) {
		isGust := false
		if _, err := ctx.Cookie("session"); err != nil {
			isGust = true
		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{"guest": isGust})
	})

}
