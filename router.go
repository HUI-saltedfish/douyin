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
	apiRouter.POST("/publish/action/", TokenAuth(), controller.Publish)
	apiRouter.GET("/publish/list/", TokenAuth(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", TokenAuth(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", TokenAuth(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", TokenAuth(), controller.CommentAction)
	apiRouter.GET("/comment/list/", TokenAuth(), controller.CommentList)

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
		var token string
		token1, ok1 := c.GetQuery("token")
		token2, ok2 := c.GetPostForm("token")
		if !ok1 && !ok2 {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  "Token is inValid",
			})
			// c.Redirect(http.StatusFound, "/douyin/user/login/")
			c.Abort()
			return
		}
		if ok1 {
			token = token1
		} else {
			token = token2
		}
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
			} else {
				c.JSON(http.StatusOK, controller.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				})
			}
			// c.Redirect(http.StatusFound, "/douyin/user/login/")
			c.Abort()
			return
		}
		if !tkn.Valid {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  "Unauthorized access",
			})
			// c.Redirect(http.StatusFound, "/douyin/user/login/")
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}
