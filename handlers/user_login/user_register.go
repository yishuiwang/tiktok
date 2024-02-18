package user_login

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/models"
	"tiktok/service/user_login"
)

type UserRegisterResponse struct {
	models.CommonResponse
	*user_login.LoginResponse
}

func UserRegisterHandler(c *gin.Context) {
	username := c.Query("username")
	rawVal, _ := c.Get("password")
	password, ok := rawVal.(string)
	if !ok {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResponse: models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "parse password error",
			},
		})
		return
	}
	registerResponse, err := user_login.PostUserLogin(username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResponse: models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		CommonResponse: models.CommonResponse{StatusCode: 0},
		LoginResponse:  registerResponse,
	})
}
