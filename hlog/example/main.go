package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hwholiday/ghost/hlog"
	"go.uber.org/zap"
	"net/http"
)

func AddTraceId() gin.HandlerFunc {
	return func(g *gin.Context) {
		traceId := g.GetHeader("traceId")
		if traceId == "" {
			traceId = uuid.New().String()
		}
		ctx, log := hlog.GetLogger().AddCtx(g.Request.Context(), zap.Any("traceId", traceId))
		g.Request = g.Request.WithContext(ctx)
		log.Info("AddTraceId success")
		g.Next()
	}
}

// curl http://127.0.0.1:8888/test
func main() {
	hlog.NewLogger(hlog.SetPrintTime(hlog.PrintTimestamp))
	g := gin.New()
	g.Use(AddTraceId())
	g.GET("/test", func(context *gin.Context) {
		log := hlog.GetLogger().GetCtx(context.Request.Context())
		log.Info("test")
		log.Debug("test")
		context.JSON(200, "success")
	})
	hlog.GetLogger().Info("hconf example success")
	http.ListenAndServe(":8888", g)
}

// curl http://127.0.0.1:8888/test
