package router

import (
	"net/http"
	"time"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"west.garden/template/common/config"
	"west.garden/template/handler"
)

func InitApiRouter(cfg *config.Configuration) *gin.Engine {
	r := initDefaultRouter(cfg)

	v1 := r.Group("/api")
	{
		v1.GET("/ping", handler.Pong)
	}
	return r
}

func initDefaultRouter(cfg *config.Configuration) *gin.Engine {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	corsConfig := getCorsConfig()
	corsConfig.AllowAllOrigins = true
	r.Use(cors.New(corsConfig))
	if cfg.Http.LimitConnection > 0 {
		r.Use(limit.MaxAllowed(cfg.Http.LimitConnection))
	}
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"result": false,
			"error":  "Method Not Allowed",
		})
		return
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"result": false,
			"error":  "Endpoint Not Found",
		})
		return
	})
	// 最大运行上传文件大小
	r.MaxMultipartMemory = cfg.Http.MaxMultipartMemory * 1024 * 1024
	return r
}

func getCorsConfig() cors.Config {
	return cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Pragma", "Cache-Control", "Connection", "Content-Type", "FilePath-Length", "FilePath-Type", "Authorization", "X-Forwarded-For", "User-Agent", "Referer"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}
