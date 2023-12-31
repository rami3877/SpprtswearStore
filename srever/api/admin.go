package api

import (
	"db"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"structs"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type admin struct {
	adminGroup *gin.RouterGroup
	adminUsers map[string]string
}

func (admin *admin) GetUsers() {

	f, _ := os.OpenFile("adminuser.json", os.O_RDONLY|os.O_CREATE, 0660)
	data, _ := io.ReadAll(f)

	if err := json.Unmarshal(data, &admin.adminUsers); err != nil {
		log.Fatal(err)
	}
	f.Close()
}

func (admin *admin) setOrderApi() {

	admin.adminGroup.GET("/orders", func(ctx *gin.Context) {
		orders := db.MainDB.Orders.Get()
		if orders == nil {
			ctx.String(http.StatusNoContent, "")
		}
		ctx.JSON(http.StatusOK, orders)
	})
}

func (admin *admin) setProductApi() {

	admin.adminGroup.POST("/product/model/updataImage", func(ctx *gin.Context) {

		type LinkImage struct {
			Id        int      `json:"id"`
			Container string   `json:"container"`
			Kind      string   `json:"kind"`
			Image     []string `json:"image"`
		}
		linkImage := LinkImage{}

		if err := ctx.ShouldBindJSON(&linkImage); err != nil {
			ctx.String(http.StatusBadRequest, "check json")
			return
		}
		if err := db.MainDB.Stock.UpateLinkesImage(linkImage.Id, linkImage.Container, linkImage.Kind, linkImage.Image); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return

		}
		ctx.String(http.StatusAccepted, "UPDATA")

	})

	admin.adminGroup.POST("/product/model/sizes", func(ctx *gin.Context) {
		type Size struct {
			Id        int                       `json:"id"`
			SizeName  string                    `json:"sizeName"`
			Container string                    `json:"container"`
			Kind      string                    `json:"kind"`
			C         map[string]map[string]int `json:"Size"`
		}
		newSize := Size{}
		if err := db.MainDB.Stock.UpdataSizeFromModel(newSize.Id, newSize.Container, newSize.Kind, newSize.SizeName, newSize.C[newSize.SizeName]); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusAccepted, "done")
	})

	admin.adminGroup.GET("/product", func(ctx *gin.Context) {
		AllContainer := db.MainDB.Stock.GetAllContainer()
		if len(AllContainer) != 0 {

			ctx.HTML(http.StatusOK, "product.html", gin.H{
				"Container": AllContainer,
			})
		} else {
			ctx.HTML(http.StatusOK, "product.html", gin.H{
				"Container": "",
			})
		}

	})
	admin.adminGroup.POST("/product/container", func(ctx *gin.Context) {
		newContainer := ctx.PostForm("Container")
		if len(newContainer) == 0 {
			ctx.String(http.StatusNotAcceptable, "name container is empty")
		} else if regexp.MustCompile(`\s`).Match([]byte(newContainer)) {
			ctx.String(http.StatusNotAcceptable, "has space")
		} else if err := db.MainDB.Stock.AddNewContainer(newContainer); err != nil {
			ctx.String(http.StatusNotAcceptable, err.Error())
		} else {
			ctx.String(http.StatusCreated, "Created")
		}

	})

	admin.adminGroup.GET("/product/container", func(ctx *gin.Context) {
		Containers := db.MainDB.Stock.GetAllContainerAndKind()
		ctx.JSON(http.StatusFound, &Containers)

	})

	admin.adminGroup.POST("/product/kind", func(ctx *gin.Context) {
		type kind struct {
			Container string `json:"container"`
			Kindname  string `json:"kind"`
		}
		var newKind kind
		if err := ctx.ShouldBindJSON(&newKind); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		if err := db.MainDB.Stock.NewKindtoContainer(newKind.Container, newKind.Kindname); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusCreated, "Created")

	})

	admin.adminGroup.DELETE("/product/kind", func(ctx *gin.Context) {
		type kind struct {
			Container string `json:"container"`
			Kindname  string `json:"kind"`
		}
		var deleteKind kind
		if err := ctx.ShouldBindJSON(&deleteKind); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		if err := db.MainDB.Stock.DeleteKind(deleteKind.Container, deleteKind.Kindname); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusAccepted, "delete")

	})
	admin.adminGroup.DELETE("/product/container", func(ctx *gin.Context) {
		Container, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.String(http.StatusNotAcceptable, "try agein")
		}
		newContainer := string(Container)
		if len(newContainer) == 0 {
			ctx.String(http.StatusNotAcceptable, "name container is empty")
		} else if regexp.MustCompile(`\s`).Match([]byte(newContainer)) {
			ctx.String(http.StatusNotAcceptable, "has space")
		} else if err := db.MainDB.Stock.DeleteContainer(newContainer); err != nil {
			ctx.String(http.StatusNotAcceptable, err.Error())
		} else {
			ctx.String(http.StatusCreated, "delete")
		}

	})

	admin.adminGroup.GET("/product/model", func(ctx *gin.Context) {
		type kind struct {
			Container string `json:"container"`
			Kindname  string `json:"kind"`
		}
		var ki kind
		if err := ctx.ShouldBindJSON(&ki); err != nil {
			if io.EOF == err {
				ctx.String(http.StatusBadRequest, "empty Body")
			} else {
				ctx.String(http.StatusBadRequest, err.Error())
			}
			return
		}
		fg, err := db.MainDB.Stock.GetAllModelsInKind(ki.Container, ki.Kindname)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		if fg == nil {
			ctx.String(http.StatusNoContent, "no data ")
		}
		ctx.JSON(http.StatusFound, &fg)
	})
	admin.adminGroup.POST("/product/model", func(ctx *gin.Context) {
		type kind struct {
			Container string        `json:"container"`
			Kindname  string        `json:"kind"`
			Model     structs.Model `json:"model"`
		}
		var ki kind
		err := ctx.ShouldBindJSON(&ki)
		if err != nil {
			ctx.String(http.StatusOK, "check json ")
			return
		}
		if err := db.MainDB.Stock.AddModelToKind(ki.Container, ki.Kindname, &ki.Model); err != nil {
			ctx.String(http.StatusOK, err.Error())
			return
		}
		ctx.String(http.StatusOK, "create")

	})
	admin.adminGroup.DELETE("/product/model", func(ctx *gin.Context) {
		type kind struct {
			Container string `json:"container"`
			Kindname  string `json:"kind"`
			Id        int    `josn:"id"`
		}
		deleteModel := kind{}
		if err := ctx.ShouldBindJSON(&deleteModel); err != nil {
			ctx.String(http.StatusBadRequest, "check json")
			return
		}
		if err := db.MainDB.Stock.DeleteModelFromKind(deleteModel.Container, deleteModel.Kindname, deleteModel.Id); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusAccepted, "delete %d", deleteModel.Id)
	})

}

func (admin *admin) setLogoutApi() {

	admin.adminGroup.GET("/logout", func(ctx *gin.Context) {
		ctx.SetCookie("session", "", -1, "/admin", "", false, true)
		ctx.String(http.StatusOK, "ok")
	})
}

func (admin *admin) setLoginApi() {

	//	admin.adminGroup.GET("/login", func(ctx *gin.Context) {
	//		v, err := ctx.Cookie("session")
	//		infoUserAdmin := strings.Split(v, ",")
	//		if err == nil || len(infoUserAdmin) < 2 {
	//			//ctx.Redirect(http.StatusMovedPermanently, "/admin/")
	//			ctx.String(http.StatusBadRequest, "ERROR")
	//			return
	//		} else if err = bcrypt.CompareHashAndPassword([]byte(infoUserAdmin[1]), []byte(admin.adminUsers[infoUserAdmin[0]])); err == nil {
	//			//ctx.Redirect(http.StatusMovedPermanently, "/admin/")
	//			ctx.String(http.StatusBadRequest, "ERROR")
	//			return
	//		}
	//
	//		ctx.String(http.StatusOK, "ok")
	//	})

	admin.adminGroup.POST("/login", func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password, ok := admin.adminUsers[username]
		if !ok || string(password) != ctx.PostForm("password") {
			//	ctx.HTML(http.StatusBadRequest, "adminLogin.html", gin.H{
			//		"err": "check your password and your username",
			//	})
			ctx.String(http.StatusBadRequest, "check your password and your username")
			return
		}
		hashpassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

		ctx.SetCookie("session", username+","+string(hashpassword), 0, "/admin", "", false, true)
		ctx.String(http.StatusOK, "ok")
	})

}

func (admin *admin) setMiddleware() {

	admin.adminGroup.Use(func(ctx *gin.Context) {
		if ctx.FullPath() == "/admin/login" {
			ctx.Next()
			return
		}

		if ctx.FullPath() == "/admin/logout" {
			ctx.Next()
			return
		}
		v, err := ctx.Cookie("session")
		infoAdmin := strings.Split(v, ",")
		if err != nil {
			ctx.Redirect(http.StatusMovedPermanently, "/admin/login")
			ctx.Abort()
			return
		} else if err = bcrypt.CompareHashAndPassword([]byte(infoAdmin[1]), []byte(admin.adminUsers[infoAdmin[0]])); err != nil {
			ctx.String(http.StatusOK, err.Error())
			ctx.Abort()
			return
		}

		ctx.Next()

	})
}

func (admin *admin) setAdminPage() {
	admin.adminGroup.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminMainPage.html", nil)

	})
}
