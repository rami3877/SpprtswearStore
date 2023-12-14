package srever

import (
	"db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"structs"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Srever int

var api *gin.Engine

type String string

func (s *String) byte() []byte {
	return []byte(*s)
}
func (s *String) string() string {
	return string(*s)
}

var AdminUsers map[string]String = make(map[string]String)

func (*Srever) Run() {

	f, _ := os.OpenFile("adminuser.json", os.O_RDONLY|os.O_CREATE, 0660)
	data, _ := io.ReadAll(f)
	if err := json.Unmarshal(data, &AdminUsers); err != nil {
		log.Fatal(err)
	}

	api = gin.New()

	api.Use(gin.Logger())

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: api,
	}

	api.NoRoute(gin.WrapH(http.FileServer(http.Dir("public"))))

	//--------------- $admin$ ---------------------
	admin := api.Group("/admin")
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

		ctx.String(http.StatusOK, "ok")
	})

	admin.POST("/login", func(ctx *gin.Context) {
		password, ok := AdminUsers[ctx.PostForm("username")]
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
		ctx.String(http.StatusOK, "ok")
	})

	admin.GET("/logout", func(ctx *gin.Context) {
		ctx.SetCookie("session", "", -1, "/admin", "", false, true)
		ctx.String(http.StatusOK, "ok")
	})
	admin.GET("/product", func(ctx *gin.Context) {
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
	admin.POST("/product/container", func(ctx *gin.Context) {
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

	admin.GET("/product/container", func(ctx *gin.Context) {
		Containers := db.MainDB.Stock.GetAllContainerAndKind()
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
		if err := db.MainDB.Stock.NewKindtoContainer(newKind.Container, newKind.Kindname); err != nil {
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
		if err := db.MainDB.Stock.DeleteKind(deleteKind.Container, deleteKind.Kindname); err != nil {
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
		} else if err := db.MainDB.Stock.DeleteContainer(newContainer); err != nil {
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
		if err := db.MainDB.Stock.AddModelToKind(ki.Container, ki.Kindname, &ki.Model); err != nil {
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
		if err := db.MainDB.Stock.DeleteModelFromKind(deleteModel.Container, deleteModel.Kindname, deleteModel.Id); err != nil {
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
		if err := db.MainDB.Stock.UpateLinkesImage(linkImage.Id, linkImage.Container, linkImage.Kind, linkImage.Image); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return

		}
		ctx.String(http.StatusAccepted, "UPDATA")

	})

	admin.POST("/product/model/sizes", func(ctx *gin.Context) {
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

	admin.GET("/shutdown", func(ctx *gin.Context) {
		httpServer.Shutdown(nil)
	})
	// -------------MainApi-------------------
	api.GET("/product", func(ctx *gin.Context) {
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

	api.GET("/AllContainerAndKind", func(ctx *gin.Context) {
		fg := db.MainDB.Stock.GetAllContainerAndKind()
		if len(fg) == 0 {
			ctx.String(http.StatusNoContent, "")
		} else {
			ctx.JSON(http.StatusOK, &fg)
		}

	})

	//------------USER----------
	user := api.Group("user")

	user.Use(func(ctx *gin.Context) {
		if ctx.FullPath() == "/user/login" {
			ctx.Next()
			return
		}
		if ctx.FullPath() == "/user/logout" {
			ctx.Next()
			return
		}

		if ctx.FullPath() == "/user/register" {
			ctx.Next()
			return
		}
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error login")
			ctx.Abort()
			return
		}
		_, m, d := time.Now().Date()

		infoUser := strings.Split(v, ",")
		if len(infoUser) != 2 {
			ctx.String(http.StatusBadRequest, "error")
			ctx.Abort()
			return
		}

		user := structs.User{}
		if err := db.MainDB.Users.GetUser(infoUser[0], &user); err != nil {
			ctx.String(http.StatusBadRequest, "error")
			ctx.Abort()
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(infoUser[1]), []byte(user.Password+m.String()+strconv.Itoa(d))); err != nil {
			ctx.String(http.StatusBadRequest, "error")
			ctx.Abort()
			return
		}
		ctx.Next()
	})

	user.DELETE("/commint", func(ctx *gin.Context) {

		type UserCommint struct {
			Username  string `josn:"username"`
			Commint   string `json:"commint"`
			Stars     int    `json:"stars"`
			Container string `json:"container"`
			Kind      string `json:"kind"`
			Idmodel   int    `json:"idmodel"`
		}
		commint := UserCommint{}
		if err := ctx.ShouldBindJSON(&commint); err != nil {
			ctx.String(http.StatusNotAcceptable, "check json")
			return
		}
		if commint.Commint == "" {
			ctx.String(http.StatusNotAcceptable, "commint empty")
			return
		}
		if commint.Stars <= 0 || commint.Stars > 6 {
			ctx.String(http.StatusNotAcceptable, "stars invaild")
			return
		}

		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		infoUser := strings.Split(v, ",")
		commint.Username = infoUser[0]

		if err := db.MainDB.AddCommint(commint.Idmodel, commint.Container, commint.Kind, structs.UserCommint{Username: commint.Username, Commint: commint.Commint, Stars: commint.Stars}); err != nil {
			ctx.String(http.StatusNotAcceptable, err.Error())
			return
		}
		ctx.String(http.StatusOK, "add")

	})
	user.POST("/commint", func(ctx *gin.Context) {

		type UserCommint struct {
			Username  string `josn:"username"`
			Commint   string `json:"commint"`
			Stars     int    `json:"stars"`
			Container string `json:"container"`
			Kind      string `json:"kind"`
			Idmodel   int    `json:"idmodel"`
		}
		commint := UserCommint{}
		if err := ctx.ShouldBindJSON(&commint); err != nil {
			ctx.String(http.StatusNotAcceptable, "check json")
			return
		}
		if commint.Commint == "" {
			ctx.String(http.StatusNotAcceptable, "commint empty")
			return
		}
		if commint.Stars <= 0 || commint.Stars > 6 {
			ctx.String(http.StatusNotAcceptable, "stars invaild")
			return
		}

		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		infoUser := strings.Split(v, ",")
		commint.Username = infoUser[0]

		if err := db.MainDB.AddCommint(commint.Idmodel, commint.Container, commint.Kind, structs.UserCommint{Username: commint.Username, Commint: commint.Commint, Stars: commint.Stars}); err != nil {
			ctx.String(http.StatusNotAcceptable, err.Error())
			return
		}
		ctx.String(http.StatusOK, "add")

	})

	user.POST("/register", func(ctx *gin.Context) {
		type User1 struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}
		var useradd User1
		newuser := structs.User{}

		if err := ctx.ShouldBindJSON(&useradd); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		newuser.Username = useradd.Username
		newuser.Password = useradd.Password
		newuser.UserEmail= useradd.Email

		if err := db.MainDB.Users.AddNew(newuser); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusAccepted, "create")

	})
	user.POST("/login", func(ctx *gin.Context) {
		type loginUser struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		userlogin := loginUser{}
		if err := ctx.ShouldBindJSON(&userlogin); err != nil {
			ctx.String(http.StatusBadRequest, "check json")
			return
		}
		user := structs.User{}
		if err := db.MainDB.Users.GetUser(userlogin.Username, &user); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		_, m, d := time.Now().Date()
		pasHash, _ := bcrypt.GenerateFromPassword([]byte(userlogin.Password+m.String()+strconv.Itoa(d)), 12)
		ctx.SetCookie("session", userlogin.Username+","+string(pasHash), 0, "/user", "", false, true)
		ctx.String(http.StatusAccepted, "ok")

	})
	user.POST("/logout", func(ctx *gin.Context) {
		ctx.SetCookie("session", "", -1, "/user", "", false, true)
		ctx.HTML(http.StatusOK, "adminLogin.html", nil)

	})
	user.POST("/phone", func(ctx *gin.Context) {
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		phone := ctx.PostForm("phone")
		if err := db.MainDB.Users.UpdataPhone(username, phone); err != nil {
			ctx.String(http.StatusNotAcceptable, err.Error())
			return
		}
		ctx.String(http.StatusAccepted, "update")
	})

	user.GET("/phone", func(ctx *gin.Context) {
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		user := structs.User{}
		if err := db.MainDB.Users.GetUser(username, &user); err != nil {
			ctx.String(http.StatusNotAcceptable, err.Error())
			return
		}
		if user.Phone == "" {
			ctx.String(http.StatusNoContent, "")
		} else {
			ctx.String(http.StatusFound, user.Phone)
		}

	})
	user.POST("/visa", func(ctx *gin.Context) {
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		visa := structs.Visa{}
		if err := ctx.ShouldBindJSON(&visa); err != nil {
			ctx.String(http.StatusBadRequest, "check the json")
			return
		}
		if err := db.MainDB.Users.AddVisa(username, visa); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusAccepted, "add visa")
	})

	user.GET("/visa", func(ctx *gin.Context) {
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		user := structs.User{}

		if err := db.MainDB.Users.GetUser(username, &user); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		if len(user.UserVisa) == 0 {
			ctx.String(http.StatusNoContent, "")
			return
		}
		visa := []string{}
		for _, v := range user.UserVisa {
			toWebVisa := ""
			toWebVisa += v.Number[:3]
			for i := 0; i < 9; i++ {
				toWebVisa += "*"
			}
			toWebVisa += v.Number[len(v.Number)-3:]
			visa = append(visa, toWebVisa)
		}
		ctx.JSON(http.StatusFound, visa)
	})

	user.POST("/addr", func(ctx *gin.Context) {
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		addr := structs.Addr{}

		if err := ctx.ShouldBindJSON(&addr); err != nil {
			ctx.String(http.StatusBadRequest, "check the json plz")
			return
		}
		if err := db.MainDB.Users.UpdataAddr(username, addr); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusAccepted, "update")

	})

	user.GET("/addr", func(ctx *gin.Context) {
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		user := structs.User{}

		if err := db.MainDB.Users.GetUser(username, &user); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		if user.UserAddr.City == "" {
			ctx.JSON(http.StatusFound, user.UserAddr)
		} else {
			ctx.String(http.StatusNoContent, "")
		}

	})
	user.GET("/lastName", func(ctx *gin.Context) {
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		user := structs.User{}

		if err := db.MainDB.Users.GetUser(username, &user); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		if user.LastName != "" {
			ctx.String(http.StatusOK, user.LastName)
		} else {

			ctx.String(http.StatusNoContent, "")
		}
	})

	user.POST("/lastName", func(ctx *gin.Context) {
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		lastName := ctx.PostForm("lastName")
		if err := db.MainDB.Users.UpdataLastName(username, lastName); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusAccepted, "update")

	})
	user.GET("/firstName", func(ctx *gin.Context) {
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		user := structs.User{}

		if err := db.MainDB.Users.GetUser(username, &user); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		if user.FirstName != "" {
			ctx.String(http.StatusOK, user.FirstName)
		} else {

			ctx.String(http.StatusNoContent, "")
		}
	})

	user.POST("/firstName", func(ctx *gin.Context) {
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		firstName := ctx.PostForm("firstName")
		if err := db.MainDB.Users.UpdataFirstName(username, firstName); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusAccepted, "update")

	})
	user.POST("/password", func(ctx *gin.Context) {

		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		type _password struct {
			NewPassowrd string `json:"newPassowrd"`
			OldPassowrd string `json:"oldPassowrd"`
		}
		password := _password{}
		if err := ctx.ShouldBindJSON(&password); err != nil {
			ctx.String(http.StatusBadRequest, "check json")
			return
		}

		if err := db.MainDB.Users.UpdataPassword(username, password.OldPassowrd, password.NewPassowrd); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		ctx.String(http.StatusAccepted, "update")

	})

	user.POST("/buy", func(ctx *gin.Context) {
		v, errs := ctx.Cookie("session")
		if errs != nil {
			ctx.String(http.StatusBadRequest, "error")
			return
		}
		username := strings.Split(v, ",")[0]
		orders := []db.Orders{}
		type Err struct {
			Id  int    `json:"id"`
			ERR string `json:"err"`
		}
		err := []Err{}
		if err := ctx.ShouldBindJSON(&orders); err != nil {
			ctx.String(http.StatusBadRequest, "check json")
			return
		}
		for _, v := range orders {
			v.Username = username
			if err_ := db.MainDB.Buy(v); err_ != nil {
				if err_ != nil && err_ != db.ErrDataBase {
					ctx.String(http.StatusLocked, err_.Error())
					return
				}
				err = append(err, Err{Id: v.IdModel, ERR: err_.Error()})
			}
		}
		if len(err) == 0 {
			ctx.AsciiJSON(http.StatusOK, `{"allIsOk":true}`)
		} else {
			ctx.JSON(http.StatusNotAcceptable, &err)
		}

	})

	httpServer.ListenAndServe()

}

func InitSever() *Srever {
	return new(Srever)
}
