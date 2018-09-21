package main

import (
	"github.com/golang/glog"
	//1
	_ "feedback/init"
	//2
	"feedback/api"
	"feedback/middleware"
	"feedback/room"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.New()
	//static
	router.StaticFS("/static", http.Dir("./static"))
	router.StaticFS("/device", http.Dir("./device")).Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/x-javascript")
		c.Next()
	})
	router.Use(gin.Logger(), middleware.Recovery(), middleware.AddHeader())
	//api
	rootAPI := router.Group("/api")
	api.InitApi(rootAPI)
	//ws
	room.NewRoomServer(router.Group("/socket"))
	go func() {
		// err := http.ListenAndServeTLS(":443", "ssl/cert.pem", "ssl/key.pem", router)
		err := router.RunTLS(":443", "ssl/cert.pem", "ssl/key.pem")
		if err != nil {
			glog.Error(err)
		}
	}()
	router.Run(":8080")
}
