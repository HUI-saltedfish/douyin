package controller

import (
	"context"
	"net/http"
	"simple-demo/foundationPb"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

// var userIdSequence = int64(1)

// type UserLoginResponse struct {
// 	Response
// 	UserId int64  `json:"user_id,omitempty"`
// 	Token  string `json:"token"`
// }

// type UserResponse struct {
// 	Response
// 	User User `json:"user"`
// }

// func Register(c *gin.Context) {
// 	username := c.Query("username")
// 	password := c.Query("password")

// 	token := username + password

// 	if _, exist := usersLoginInfo[token]; exist {
// 		c.JSON(http.StatusOK, UserLoginResponse{
// 			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
// 		})
// 	} else {
// 		atomic.AddInt64(&userIdSequence, 1)
// 		newUser := User{
// 			Id:   userIdSequence,
// 			Name: username,
// 		}
// 		usersLoginInfo[token] = newUser
// 		c.JSON(http.StatusOK, UserLoginResponse{
// 			Response: Response{StatusCode: 0},
// 			UserId:   userIdSequence,
// 			Token:    username + password,
// 		})
// 	}
// }

// func Login(c *gin.Context) {
// 	username := c.Query("username")
// 	password := c.Query("password")

// 	token := username + password

// 	if user, exist := usersLoginInfo[token]; exist {
// 		c.JSON(http.StatusOK, UserLoginResponse{
// 			Response: Response{StatusCode: 0},
// 			UserId:   user.Id,
// 			Token:    token,
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, UserLoginResponse{
// 			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
// 		})
// 	}
// }

// func UserInfo(c *gin.Context) {
// 	token := c.Query("token")

// 	if user, exist := usersLoginInfo[token]; exist {
// 		c.JSON(http.StatusOK, UserResponse{
// 			Response: Response{StatusCode: 0},
// 			User:     user,
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, UserResponse{
// 			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
// 		})
// 	}
// }

var JwtKey = []byte("douyin")

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	registerRequest := foundationPb.DouyinUserRegisterRequest{
		Username: username,
		Password: password,
	}

	registerResponse, err := foundationClient.UserRegister(context.Background(), &registerRequest)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	if registerResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: registerResponse.StatusMsg},
		})
		return
	}

	var creds = Credentials{
		Username: username,
		Password: password,
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(1 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   registerResponse.UserId,
		Token:    tokenString,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	usersLoginInRequest := foundationPb.DouyinUserLoginRequest{
		Username: username,
		Password: password,
	}

	usersLoginInResponse, err := foundationClient.UserLogin(context.Background(), &usersLoginInRequest)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	if usersLoginInResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: usersLoginInResponse.StatusMsg},
		})
		return
	}

	var creds = Credentials{
		Username: username,
		Password: password,
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(1 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   usersLoginInResponse.UserId,
		Token:    tokenString,
	})
}

func UserInfo(c *gin.Context) {
	userId := c.Query("user_id")

	// get user by user_id
	uID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user_id is not valid",
		})
	}

	userInfoRequest := foundationPb.DouyinUserRequest{
		UserId:   uID,
		Username: c.GetString("username"),
	}

	userInfoResponse, err := foundationClient.User(context.Background(), &userInfoRequest)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	if userInfoResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  userInfoResponse.StatusMsg,
		})
	}

	var userRes = UserResponse{
		Response: Response{
			StatusCode: userInfoResponse.StatusCode,
			StatusMsg:  userInfoResponse.StatusMsg,
		},
		User: User{
			Id:            userInfoResponse.User.Id,
			Name:          userInfoResponse.User.Name,
			FollowCount:   userInfoResponse.User.FollowCount,
			FollowerCount: userInfoResponse.User.FollowerCount,
			IsFollow:      userInfoResponse.User.IsFollow,
		},
	}

	c.JSON(http.StatusOK, userRes)

}
