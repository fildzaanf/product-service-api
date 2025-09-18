package middleware

import (
	"io"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func InitLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		ForceColors:     true,
		DisableColors:   false,
		PadLevelText:    true,
	})

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		logrus.SetOutput(os.Stdout)
		logrus.Warn("failed to save logs to file, using stdout instead.")
	}

	logrus.SetLevel(logrus.InfoLevel)
}

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		req := c.Request()

		reqLog := logrus.WithFields(logrus.Fields{
			"time":       start.Format("2006-01-02T15:04:05.000Z07:00"),
			"remote_ip":  c.RealIP(),
			"user_agent": req.UserAgent(),
			"method":     req.Method,
			"uri":        req.RequestURI,
		})

		switch req.Method {
		case "":
			reqLog.Error("\033[31m[REQUEST] invalid request method\033[0m")
		case "GET", "POST", "PUT", "DELETE":
			reqLog.Info("\033[34m[REQUEST] incoming request\033[0m")
		default:
			reqLog.Warn("\033[33m[REQUEST] unusual request method\033[0m")
		}

		err := next(c)
		latency := time.Since(start)
		res := c.Response()

		resLog := logrus.WithFields(logrus.Fields{
			"time":       time.Now().Format("2006-01-02T15:04:05.000Z07:00"),
			"remote_ip":  c.RealIP(),
			"user_agent": req.UserAgent(),
			"method":     req.Method,
			"uri":        req.RequestURI,
			"status":     res.Status,
			"latency":    latency.String(),
		})

		switch {
		case res.Status >= 500 && latency > 5*time.Second:
			resLog.Fatal("\033[35m[CRITICAL] high-latency server error occurred \033[0m")
		case res.Status >= 500:
			resLog.Error("\033[31m[ERROR] server encountered an error \033[0m")
		case res.Status >= 400:
			resLog.Warn("\033[33m[WARNING] client made an invalid request \033[0m")
		default:
			resLog.Info("\033[32m[RESPONSE] request processed successfully \033[0m")
		}

		return err
	}
}
