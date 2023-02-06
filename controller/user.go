package controller

import (
	"errors"
	"net/http"
	"simple-demo/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
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

var jwtKey = []byte("douyin")

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

	var creds = Credentials{
		Username: username,
		Password: password,
	}

	if _, err := model.GetUserByName(creds.Username); !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := model.User{
			Name:     creds.Username,
			Password: creds.Password,
		}
		model.CreateUser(&newUser)
		u, _ := model.GetUserByName(creds.Username)
		newUser.ID = (*u).ID

		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Minute)
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
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Internal server error"},
			})
			return
		}

		// Finally, we set the client cookie for "token" as the JWT we just generated
		// we also set an expiry time which is the same as the token itself
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(newUser.ID),
			Token:    tokenString,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var creds = Credentials{
		Username: username,
		Password: password,
	}

	user, err := model.GetUserByNameAndPassword(creds.Username, creds.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
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
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Internal server error"},
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   int64(user.ID),
		Token:    tokenString,
	})
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")

	// get user by user_id
	uID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user_id is not valid",
		})
	}

	user, err := model.GetUserById(int(uID))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	// get user by token
	// Parse the JWT string and store the result in `claims`.
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: 1,
				StatusMsg:  "Unauthorized access",
			})
			return
		}
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "Bad request",
		})
		return
	}
	if !tkn.Valid {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "Unauthorized access",
		})
		return
	}

	if claims.Username != user.Name {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "Unauthorized access",
		})
		return
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(user.ID),
			Token:    token,
		})
	}
}
