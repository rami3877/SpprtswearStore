package srever

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func loginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginPost(c *gin.Context) {
	Newuser := User{}

	if c.GetHeader("Content-Type") == "application/x-www-form-urlencoded" {
		Newuser.Username = c.PostForm("username")
		Newuser.Password = c.PostForm("password")
		Newuser.Email = c.PostForm("email")
	} else if c.GetHeader("Content-Type") == "application/json" {
		if err := c.ShouldBindJSON(&Newuser); err != nil {
			c.AsciiJSON(http.StatusBadRequest, `{"err":"invilde json"}`)
			return
		}
	}

	if Newuser.Username == "" {
		c.AsciiJSON(http.StatusBadRequest, `{"err":"username is empty"}`)
		return
	}
	if Newuser.Password == "" {
		c.AsciiJSON(http.StatusBadRequest, `{"err":"password is empty"}`)
		return
	}
	if Newuser.Email == "" {
		c.AsciiJSON(http.StatusBadRequest, `{"err":"email is empty"}`)
		return
	}
	c.AsciiJSON(http.StatusOK, `{"isOK":true}`)
	fmt.Println(Newuser)
}
