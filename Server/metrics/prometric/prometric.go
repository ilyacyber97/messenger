package prometric

import (
	"messenger/logs/zaplog"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests"},
		[]string{"method"})
)

func MetricHandler(ctx *gin.Context) {
	zaplog.Logger.Info("Запрошен endpoint на метрики", zaplog.LogHandler("client_ip", ctx.ClientIP()))
	promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
}
