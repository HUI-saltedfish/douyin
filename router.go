package main

import (
	"net/http"
	"simple-demo/controller"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", TokenAuth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/", controller.MessageChat)
	apiRouter.POST("/message/action/", controller.MessageAction)
}

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		// get user by token
		// Parse the JWT string and store the result in `claims`.
		claims := &controller.Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return controller.JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusOK, controller.Response{
					StatusCode: 1,
					StatusMsg:  "Unauthorized access",
				})
				return
			} else if err == jwt.ErrTokenExpired {
				c.JSON(http.StatusOK, controller.Response{
					StatusCode: 1,
					StatusMsg:  "Token expired",
				})
			} else if err == jwt.ErrInvalidKey {
				c.JSON(http.StatusOK, controller.Response{
					StatusCode: 1,
					StatusMsg:  "Invalid key",
				})
			}
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		if !tkn.Valid {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  "Unauthorized access",
			})
			return
		}
		c.Next()
	}
}
