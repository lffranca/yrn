package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
)

func main() {
	slog.Info("start api")

	engine := gin.Default()

	engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	if err := engine.Run(); err != nil {
		panic(err)
	}
}
