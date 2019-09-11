package main

import (
	"blog/pkg/setting"
	"blog/routers"
	"fmt"
	"net/http"
)

func main(){
	//router:=gin.Default()
	//router.GET("/test", func(c *gin.Context) {
	//	c.JSON(200,gin.H{"message":"test"})
	//})
	router:= routers.InitRouter()
	s:=&http.Server{
		Addr: fmt.Sprintf(":%d",setting.HTTPort),
		Handler:router,
		ReadTimeout:setting.ReadTimeout,
		WriteTimeout:setting.WriteTimeout,
		MaxHeaderBytes:1<<20,
	}
	s.ListenAndServe()
}