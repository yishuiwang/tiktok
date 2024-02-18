package user_login

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/models"
	"tiktok/service/user_login"
)

type UserLoginResponse struct {
	models.CommonResponse
	*user_login.LoginResponse
}

func UserLoginHandler(c *gin.Context) {
	username := c.Query("username")
	raw, _ := c.Get("password")
	password, ok := raw.(string)
	if !ok {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "password error",
			},
		})
	}
	userLoginResponse, err := user_login.QueryUserLogin(username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: models.CommonResponse{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: models.CommonResponse{StatusCode: 0},
		LoginResponse:  userLoginResponse,
	})
}
