package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitLogger() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()
		log.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"latency":  time.Since(startTime),
			"clientIP": c.ClientIP(),
		}).Info("request details")
	}
}
