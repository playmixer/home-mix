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

	r.StaticFile("/", "./www/pages/index.html")
	r.StaticFS("/static", http.Dir("./www/static/"))

	r.GET("/api/v0/ping", handleGetPing)

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

type tResultGetPing struct {
	IP     string `json:"ip"`
	Mac    string `json:"mac"`
	Owner  string `json:"name"`
	Online bool   `json:"online"`
}

func handleGetPing(ctx *gin.Context) {
	conn, err := database.Connect()
	if err != nil {
		panic(err)
	}
	result := []tResultGetPing{}
	conn.Model(&database.Ping{}).Select("pings.ip, pings.mac, devices.owner, pings.online").Joins("left join devices on devices.mac = pings.mac").Find(&result)
	ctx.JSON(http.StatusOK, result)
}
