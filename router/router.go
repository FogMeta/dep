package router

import (
	"net/http"

	"github.com/FogMeta/libra-os/api/result"
	"github.com/FogMeta/libra-os/service"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
	Router.Use(cors())
	v1 := Router.Group("v1")
	{
		// spaces
		v1.GET("/spaces")

		// providers
		v1.GET("/providers")
	}

}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

var jwtService service.JWTService

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			c.Abort()
			c.JSON(http.StatusOK, result.Result{
				Code: result.UserTokenExpired,
				Msg:  "token expired, need login again",
			})
			return
		}
		uid, newToken, err := jwtService.Validate(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, result.Result{
				Code: result.UserTokenInvalid,
				Msg:  "token invalid, need login again",
			})
			return
		}
		if newToken != "" {
			c.Header("new-token", newToken)
			c.Request.Header.Set("Authorization", newToken)
		}
		c.Set("uid", uid)
		c.Next()
	}
}
