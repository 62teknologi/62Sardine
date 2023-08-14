package main

import (
	"fmt"
	"net/http"

	"github.com/62teknologi/62sardine/app/http/controllers"
	"github.com/62teknologi/62sardine/config"

	"github.com/gin-gonic/gin"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("cannot load config: %w", err)
		return
	}

	r := gin.Default()
	r.Static("/storage", "./storage")
	r.SetTrustedProxies(nil)

	apiV1 := r.Group("/api/v1")
	{
		c := &controllers.FileController{}
		r := apiV1

		r.GET("/files", c.FindAll)
		r.POST("/files", c.Upload)
		r.DELETE("/files", c.Delete)
		r.GET("/temporary-url", c.TempUrl)

	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"code":   http.StatusOK,
		})
	})

	err = r.Run(config.HTTPServerAddress)

	if err != nil {
		fmt.Printf("cannot run server: %w", err)
		return
	}
}
