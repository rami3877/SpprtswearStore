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

type user struct {
	userGroup *gin.RouterGroup
}

func (user *user) setCheckoutApi() {

	user.userGroup.GET("/checkout", func(ctx *gin.Context) {
		 ctx.HTML(http.StatusOK , "payment.html", gin.H{})
	})
	user.userGroup.POST("/buy", func(ctx *gin.Context) {
		v, _ := ctx.Cookie("session")

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
					ctx.String(http.StatusBadRequest, err_.Error())
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
}

func (user *user) setInformationApi() {

	user.userGroup.POST("/phone", func(ctx *gin.Context) {
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

	user.userGroup.GET("/phone", func(ctx *gin.Context) {
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
	user.userGroup.POST("/visa", func(ctx *gin.Context) {
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

	user.userGroup.GET("/visa", func(ctx *gin.Context) {
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

	user.userGroup.POST("/addr", func(ctx *gin.Context) {
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

	user.userGroup.GET("/addr", func(ctx *gin.Context) {
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
	user.userGroup.GET("/lastName", func(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK , user.LastName)
	})

	user.userGroup.POST("/lastName", func(ctx *gin.Context) {
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
	user.userGroup.GET("/firstName", func(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK , user.FirstName)
	})

	user.userGroup.POST("/firstName", func(ctx *gin.Context) {
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
	user.userGroup.POST("/password", func(ctx *gin.Context) {

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

}

func (user *user) setLogoutApi() {
	user.userGroup.GET("/logout", func(ctx *gin.Context) {
		ctx.SetCookie("session", "", -1, "/", "", false, true)

	})
}

func (user *user) setLoginApi() {

	user.userGroup.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "loginAndRegister.html", gin.H{})
	})
	user.userGroup.POST("/login", func(ctx *gin.Context) {
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
		ctx.SetCookie("session", userlogin.Username+","+string(pasHash), 0, "/", "", false, true)
		ctx.String(http.StatusAccepted, "ok")

	})
}
func (user *user) setRegisterApi() {

	user.userGroup.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "loginAndRegister.html", gin.H{})
	})
	user.userGroup.POST("/register", func(ctx *gin.Context) {
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
		newuser.UserEmail = useradd.Email

		if err := db.MainDB.Users.AddNew(newuser); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		_, m, d := time.Now().Date()
		pasHash, _ := bcrypt.GenerateFromPassword([]byte(newuser.Password+m.String()+strconv.Itoa(d)), 12)
		ctx.SetCookie("session", newuser.Username+","+string(pasHash), 0, "/", "", false, true)
		ctx.String(http.StatusAccepted, "create")

	})
}

func (user *user) setCommintApi() {

	user.userGroup.DELETE("/commint", func(ctx *gin.Context) {

		type UserCommint struct {
			Commint   string `json:"commint"`
			Container string `json:"container"`
			Kind      string `json:"kind"`
			Idmodel   int    `json:"idmodel"`
		}
		commint := UserCommint{}

		if err := ctx.ShouldBindJSON(&commint); err != nil {
			ctx.String(http.StatusNotAcceptable, "check json")
			return
		}

		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		infoUser := strings.Split(v, ",")

		if err := db.MainDB.Stock.DeleteCommint(commint.Idmodel, commint.Container, commint.Kind, infoUser[0], commint.Commint); err != nil {
			ctx.String(http.StatusNotAcceptable, err.Error())
			return
		}
		ctx.String(http.StatusOK, "delete")

	})
	user.userGroup.POST("/commint", func(ctx *gin.Context) {

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

		v, _ := ctx.Cookie("session")
		infoUser := strings.Split(v, ",")
		commint.Username = infoUser[0]

		if err := db.MainDB.AddCommint(commint.Idmodel, commint.Container, commint.Kind, structs.UserCommint{Username: commint.Username, Commint: commint.Commint, Stars: commint.Stars}); err != nil {
			ctx.String(http.StatusNotAcceptable, err.Error())
			return
		}
		ctx.String(http.StatusOK, "add")

	})

}

func (user *user) setMiddleware() {

	user.userGroup.Use(func(ctx *gin.Context) {
		if ctx.FullPath() == "/user/login" {

			ctx.Redirect(http.StatusSeeOther, "/user/login")
			return
		}
		if ctx.FullPath() == "/user/logout" {
			ctx.Redirect(http.StatusSeeOther, "/user/login")
			return
		}

		if ctx.FullPath() == "/user/register" {
			ctx.Redirect(http.StatusSeeOther, "/user/login")
			return
		}
		// check if will work on /user becase you set to be on root path
		v, err := ctx.Cookie("session")
		if err != nil {
			ctx.Redirect(http.StatusSeeOther, "/user/login")
			ctx.Abort()
			return
		}
		_, m, d := time.Now().Date()

		infoUser := strings.Split(v, ",")
		if len(infoUser) != 2 {
			ctx.SetCookie("session", "", -1, "/", "", false, true)
			ctx.Redirect(http.StatusSeeOther, "/user/login")
			ctx.Abort()
			return
		}

		user := structs.User{}
		if err := db.MainDB.Users.GetUser(infoUser[0], &user); err != nil {
			ctx.SetCookie("session", "", -1, "/user", "", false, true)
			ctx.Redirect(http.StatusSeeOther, "/user/login")
			ctx.Abort()
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(infoUser[1]), []byte(user.Password+m.String()+strconv.Itoa(d))); err != nil {
			ctx.SetCookie("session", "", -1, "/", "", false, true)
			ctx.Redirect(http.StatusSeeOther, "/user/login")
			ctx.Abort()
			return
		}
		ctx.Next()
	})

}
