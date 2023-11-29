package srever

import (
	"db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"structs"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Srever int

var app *gin.Engine

type String string

func (s *String) byte() []byte {
	return []byte(*s)
}
func (s *String) string() string {
	return string(*s)
}

var adminUsers map[string]String = make(map[string]String)

func (*Srever) Run() {
	adminUsers[os.Getenv("USER")] = "12345"
	if os.Getenv("USER") == "rami" {
		app = gin.New()

		app.Use(gin.Logger())
		app.LoadHTMLGlob("template/*.html")
		httpServer := &http.Server{
			Addr:    ":8080",
			Handler: app,
		}

		app.GET("/", mainPage)
		//--------------- $admin$ ---------------------
		admin := app.Group("/admin")
		admin.Use(func(ctx *gin.Context) {
			if ctx.FullPath() == "/admin/login" {
				ctx.Next()
				return
			}

			if ctx.FullPath() == "/admin/logout" {
				ctx.Next()
				return
			}
			v, err := ctx.Cookie("session")
			if err != nil {
				ctx.Redirect(http.StatusMovedPermanently, "/admin/login")
				ctx.Abort()
				return
			} else if err = bcrypt.CompareHashAndPassword([]byte(v), []byte("12345")); err != nil {
				ctx.String(http.StatusOK, err.Error())
				ctx.Abort()
				return
			}

			ctx.Next()

		})

		admin.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "adminMainPage.html", nil)

		})

		admin.GET("/login", func(ctx *gin.Context) {
			v, err := ctx.Cookie("session")
			if err == nil {
				ctx.Redirect(http.StatusMovedPermanently, "/admin/")
				return
			} else if err = bcrypt.CompareHashAndPassword([]byte(v), []byte("12345")); err == nil {
				ctx.Redirect(http.StatusMovedPermanently, "/admin/")
				return
			}

			ctx.HTML(http.StatusOK, "adminLogin.html", nil)
		})

		admin.POST("/login", func(ctx *gin.Context) {
			password, ok := adminUsers[ctx.PostForm("username")]
			if !ok || password.string() != ctx.PostForm("password") {
				ctx.HTML(http.StatusBadRequest, "adminLogin.html", gin.H{
					"err": "check your password and your username",
				})
			}
			hashpassword, err := bcrypt.GenerateFromPassword(password.byte(), 12)
			if err != nil {
				fmt.Println(err)
				return
			}
			ctx.SetCookie("session", string(hashpassword), 0, "/admin", "", false, true)
			ctx.Redirect(http.StatusMovedPermanently, "/admin")
		})

		admin.GET("/logout", func(ctx *gin.Context) {
			ctx.SetCookie("session", "", -1, "/admin", "", false, true)
			ctx.HTML(http.StatusOK, "adminLogin.html", nil)
		})
		admin.GET("/product", func(ctx *gin.Context) {
			AllContainer := db.MainDB.InStock.GetAllContainer()
			ctx.HTML(http.StatusOK, "product.html", gin.H{
				"Container": AllContainer,
			})
		})
		admin.POST("/product", func(ctx *gin.Context) {
			allPost, _ := io.ReadAll(ctx.Request.Body)
			fmt.Println(string(allPost))
		})
		admin.POST("/product/container", func(ctx *gin.Context) {
			newContainer := ctx.PostForm("Container")
			if len(newContainer) == 0 {
				ctx.String(http.StatusNotAcceptable, "name container is empty")
			} else if regexp.MustCompile(`\s`).Match([]byte(newContainer)) {
				ctx.String(http.StatusNotAcceptable, "has space")
			} else if err := db.MainDB.InStock.AddNewContainer(newContainer); err != nil {
				ctx.String(http.StatusNotAcceptable, err.Error())
			} else {
				ctx.String(http.StatusCreated, "Created")
			}

		})

		admin.GET("/product/container", func(ctx *gin.Context) {
			Containers := db.MainDB.InStock.GetAllContainerAndKind()
			ctx.JSON(http.StatusFound, &Containers)

		})

		admin.POST("/product/kind", func(ctx *gin.Context) {
			type kind struct {
				Container string `json:"container"`
				Kindname  string `json:"kind"`
			}
			var newKind kind
			if err := ctx.ShouldBindJSON(&newKind); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			if err := db.MainDB.InStock.NewKindtoContainer(newKind.Container, newKind.Kindname); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			ctx.String(http.StatusCreated, "Created")

		})

		admin.DELETE("/product/kind", func(ctx *gin.Context) {
			type kind struct {
				Container string `json:"container"`
				Kindname  string `json:"kind"`
			}
			var deleteKind kind
			if err := ctx.ShouldBindJSON(&deleteKind); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			if err := db.MainDB.InStock.DeleteKind(deleteKind.Container, deleteKind.Kindname); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			ctx.String(http.StatusAccepted, "delete")

		})
		admin.DELETE("/product/container", func(ctx *gin.Context) {
			Container, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				ctx.String(http.StatusNotAcceptable, "try agein")
			}
			newContainer := string(Container)
			if len(newContainer) == 0 {
				ctx.String(http.StatusNotAcceptable, "name container is empty")
			} else if regexp.MustCompile(`\s`).Match([]byte(newContainer)) {
				ctx.String(http.StatusNotAcceptable, "has space")
			} else if err := db.MainDB.InStock.DeleteContainer(newContainer); err != nil {
				ctx.String(http.StatusNotAcceptable, err.Error())
			} else {
				ctx.String(http.StatusCreated, "delete")
			}

		})

		admin.GET("/product/model", func(ctx *gin.Context) {
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
			fg, err := db.MainDB.InStock.GetAllModelsInKind(ki.Container, ki.Kindname)
			if err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			if fg == nil {
				ctx.String(http.StatusNoContent, "no data ")
			}
			ctx.JSON(http.StatusFound, &fg)
		})
		admin.POST("/product/model", func(ctx *gin.Context) {
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
			if err := db.MainDB.InStock.AddModelToKind(ki.Container, ki.Kindname, &ki.Model, true); err != nil {
				ctx.String(http.StatusOK, err.Error())
				return
			}
			ctx.String(http.StatusOK, "create")

		})
		admin.DELETE("/product/model", func(ctx *gin.Context) {
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
			if err := db.MainDB.InStock.DeleteModelFromKind(deleteModel.Container, deleteModel.Kindname, deleteModel.Id); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			ctx.String(http.StatusAccepted, "delete %d", deleteModel.Id)
		})

		admin.GET("/orders", func(ctx *gin.Context) {
			orders := db.MainDB.Orders.Get()
			if orders == nil {
				ctx.String(http.StatusNoContent, "")
			}
			ctx.JSON(http.StatusOK, orders)
		})

		admin.POST("/product/model/updataImage", func(ctx *gin.Context) {

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
			if err := db.MainDB.InStock.UpateLinkesImage(linkImage.Id, linkImage.Container, linkImage.Kind, linkImage.Image); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return

			}
			ctx.String(http.StatusAccepted, "UPDATA")

		})

		admin.PUT("/product/model/sizes", func(ctx *gin.Context) {
			type Size struct {
				Id        int          `json:"id"`
				SizeName  string       `json:"sizeName"`
				Container string       `json:"container"`
				Kind      string       `json:"kind"`
				C         structs.Size `json:"Size"`
			}
			newSize := Size{}
			if err := db.MainDB.InStock.UpdataSizeFromModel(newSize.Id, newSize.Container, newSize.Kind, newSize.SizeName, newSize.C); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			ctx.String(http.StatusAccepted, "done")
		})


		admin.GET("/product/model/sizes", func(ctx *gin.Context) {
			type Size struct {
				Id        int    `json:"id"`
				SizeName  string `json:"sizeName"`
				Container string `json:"container"`
				Kind      string `json:"kind"`
			}
			size := Size{}
			var data *map[string][]structs.Size
			data, err := db.MainDB.InStock.GetALLSize(size.Id, size.Container, size.Kind, size.SizeName)
			if err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			ctx.JSON(http.StatusAccepted, *data)
		})

		admin.GET("/shutdown", func(ctx *gin.Context) {
			httpServer.Shutdown(nil)
		})

		httpServer.ListenAndServe()
	} else {
		app = gin.Default()
		app.NoRoute(gin.WrapH(http.FileServer(http.Dir("public"))))
		app.LoadHTMLGlob("./public/*/*.html")
		app.POST("/login", loginPost)
		app.GET("/login", loginGet)
		app.Run(":8080")
	}

}

func adminMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "adminMainPage.html", nil)
}

func mainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func InitSever() *Srever {
	return new(Srever)
}

func m(ctx *gin.Context) {
	_, err := ctx.Cookie("d")
	if err != nil {
		ctx.AsciiJSON(http.StatusOK, `{"err":"dsadasd"}`)
		ctx.Abort()
		return
	}
	ctx.Next()

}
