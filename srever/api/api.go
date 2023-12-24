package api

import (
	"db"
	"net/http"
	"strconv"
	"strings"
	"structs"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	api.admin.setAdminPage()
}

func (api *Api) setGuestApi(server *gin.Engine) {

	server.GET("/product", func(ctx *gin.Context) {
		kind := ctx.Query("kind")
		Containter := ctx.Query("container")
		id, err := strconv.Atoi(ctx.Query("id"))
		if err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		data, err := db.MainDB.Stock.GetModelsInKind(id, 10, Containter, kind)
		if err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, data)

	})

	server.GET("/AllContainerAndKind", func(ctx *gin.Context) {
		fg := db.MainDB.Stock.GetAllContainerAndKind()
		if len(fg) == 0 {
			ctx.JSON(http.StatusOK, "no data")
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

		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.HTML(http.StatusOK, "index.html", gin.H{"guest": true})
			return
		}
		_, m, d := time.Now().Date()

		infoUser := strings.Split(v, ",")
		if len(infoUser) != 2 {
			ctx.SetCookie("session", "", -1, "/", "", false, true)
			ctx.HTML(http.StatusOK, "index.html", gin.H{"guest": true})
			return
		}

		user := structs.User{}
		if err := db.MainDB.Users.GetUser(infoUser[0], &user); err != nil {
			ctx.SetCookie("session", "", -1, "/user", "", false, true)
			ctx.HTML(http.StatusOK, "index.html", gin.H{"guest": true})
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(infoUser[1]), []byte(user.Password+m.String()+strconv.Itoa(d))); err != nil {
			ctx.SetCookie("session", "", -1, "/", "", false, true)
			ctx.HTML(http.StatusOK, "index.html", gin.H{"guest": true})
			return
		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{"guest": false})
	})

}
