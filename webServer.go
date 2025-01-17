package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"smServer/config"
	"smServer/server/web"
	"time"
)

func main() {
	host := config.File.MustValue("web_server", "host", "localhost")
	port := config.File.MustValue("web_server", "port", "8088")

	router := gin.Default()
	web.Init(router)
	s := &http.Server{
		Addr:           host + ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Println("注册服务启动失败", err)
		panic(err)
	}

}
