package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seanchan/dyparser/parser"
)

type HttpResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func main() {
	r := gin.Default()
	r.GET("/hi", func(c *gin.Context) {
		jsonRes := HttpResponse{
			Code: 200,
			Msg:  "解析成功",
			Data: "hi",
		}
		c.JSON(http.StatusOK, jsonRes)
	})
	r.GET("/parse", func(c *gin.Context) {
		result, err := parser.Parse(c)
		if err != nil {
			c.JSON(http.StatusForbidden, HttpResponse{Code: 404, Msg: err.Error()})
		}
		c.JSON(http.StatusOK, HttpResponse{Code: 200, Data: result})
	})

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器 (设置 5 秒的超时时间)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
