package middleware

import (
	"log/slog"
	"time"
	"github.com/gin-gonic/gin"
)



func Logger() gin.HandlerFunc {
 return func(ctx *gin.Context) {
	start := time.Now()
	status := ctx.Writer.Status()	

    ctx.Next()
	duration := time.Since(start)
	logArgs := []any{
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"status", status,
			"duration", duration,
			"ip", ctx.ClientIP(),
		}
		if len(ctx.Errors) > 0 {
			logArgs = append(logArgs, "errors: ", ctx.Errors.String())
			return 
		}
		switch {
		case status >= 500 : slog.Error("something wrong on server", logArgs...)
		case status >= 400 : slog.Error("something wrong on the client", logArgs...)
		case status >= 200 : slog.Info("request completed", logArgs...)
		default : slog.Debug("request: ", logArgs...)
		}
 }
}
