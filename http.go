package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/playmixer/home-mix/database"
	"github.com/playmixer/home-mix/tools"
)

func startHttp() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	conn, err := database.Connect()
	if err != nil {
		panic(err)
	}

	result := []database.Ping{}
	conn.Find(&result)

	r.StaticFile("/", "./www/pages/index.html")
	r.StaticFS("/static", http.Dir("./www/static/"))

	r.GET("/api/v0/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, result)
	})

	r.Run(fmt.Sprintf(":%s", tools.Getenv("HTTP_PORT", "8000")))
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
