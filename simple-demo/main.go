package main

import (
	"simple-demo/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	// go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	// grpc init

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	defer controller.FoundationPbConn.Close()
	defer controller.InteractionPbConn.Close()
}
