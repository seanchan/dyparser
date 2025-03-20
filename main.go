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
	"github.com/spf13/viper"
)

type HttpResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func main() {
	// 设置viper读取配置文件
	viper.SetConfigName("config") // 配置文件名 (不带扩展名)
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 配置文件路径

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	port := viper.GetString("server.port")
	if port == "" {
		port = "8081" // 默认端口
	}
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
		query := c.Query("query")
		source := c.Query("source")
		if query == "" || source == "" {
			c.AbortWithStatusJSON(http.StatusBadGateway, HttpResponse{Code: http.StatusBadGateway, Msg: "missing query or source"})
			return
		}
		result, err := parser.Parse(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, HttpResponse{Code: http.StatusNotFound, Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, HttpResponse{Code: 200, Data: result})
	})

	srv := &http.Server{
		Addr:    ":" + port,
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
