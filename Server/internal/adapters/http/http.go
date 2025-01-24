package http

import (
	"errors"
	"net/http"
	"server/domain"
	"server/internal/core/services"
	"server/log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPHandler struct {
	s      *services.MessengerService
	log    *log.Logger
	metric *prometheus.CounterVec
}

func NewHandlerMessengerServiceRepository(service *services.MessengerService, log *log.Logger, metric *prometheus.CounterVec) *HTTPHandler {
	return &HTTPHandler{s: service, log: log, metric: metric}
}

func (h *HTTPHandler) SaveMessage(ctx *gin.Context) {
	h.log.Info("http: Получен запрос SaveMessage", zap.String("client_ip", ctx.ClientIP()))

	//metric
	h.metric.WithLabelValues(ctx.Request.Method).Inc()
	var message *domain.Message

	if err := ctx.ShouldBindJSON(message); err != nil {
		h.log.Error("http: Ошибка записи JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}
	err := h.s.SaveMessage(message)
	if err != nil {
		h.log.Error("http: Ошибка сохранения message", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": errors.Unwrap(err)})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Message сохранен успешно",
	})
	h.log.Info("http: Message сохранен успешно", zap.String("message", message.Body))

}

func (h *HTTPHandler) ReadMessage(ctx *gin.Context) {
	h.log.Info("http: Получен запрос ReadMessage", zap.String("client_ip", ctx.ClientIP()))

	//metric
	h.metric.WithLabelValues(ctx.Request.Method).Inc()

	id := ctx.Param("id")

	message, err := h.s.ReadMessage(id)
	if err != nil {
		h.log.Error("http: Ошибка чтения message", zap.String("id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": errors.Unwrap(err)})
		return
	}

	ctx.JSON(http.StatusOK, message)
	h.log.Info("http: Message прочитан успешно", zap.String("id", id))

}

func (h *HTTPHandler) ReadMessages(ctx *gin.Context) {
	h.log.Info("http: Получен запрос ReadMessages", zap.String("client_ip", ctx.ClientIP()))

	//metric
	h.metric.WithLabelValues(ctx.Request.Method).Inc()

	slice, err := h.s.ReadMessages()

	if err != nil {
		h.log.Error("http: Ошибка чтения messages", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": errors.Unwrap(err)})
		return
	}

	ctx.JSON(http.StatusOK, slice)
	h.log.Info("http: Messages прочитан успешно", zap.String("number", strconv.Itoa(len(slice))))

}

func (h *HTTPHandler) MetricHandler(ctx *gin.Context) {
	h.log.Info("Запрошен endpoint на метрики", zap.String("client_ip", ctx.ClientIP()))
	promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
}
