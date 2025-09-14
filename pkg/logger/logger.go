package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.Logger.SetPrefix("[INFO] ")
	l.Logger.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Logger.SetPrefix("[ERROR] ")
	l.Logger.Println(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Logger.SetPrefix("[DEBUG] ")
	l.Logger.Println(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.Logger.SetPrefix("[WARN] ")
	l.Logger.Println(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logger.SetPrefix("[INFO] ")
	l.Logger.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logger.SetPrefix("[ERROR] ")
	l.Logger.Printf(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logger.SetPrefix("[DEBUG] ")
	l.Logger.Printf(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Logger.SetPrefix("[WARN] ")
	l.Logger.Printf(format, v...)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		fmt.Printf("→ %s %s", method, path)
		if raw != "" {
			fmt.Printf("?%s", raw)
		}
		fmt.Println()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		fmt.Printf("← %s %s %d %v\n", method, path, status, duration)
	}
}