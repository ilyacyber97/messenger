package http

import (
	"errors"
	"server/internal/core/domain"
	"server/internal/core/ports"
	"server/logs/zaplog"
	"server/metrics/prometric"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	s ports.MessengerService
}

func NewHandlerMessengerServiceRepository(MessengerService ports.MessengerService) *HTTPHandler {
	return &HTTPHandler{s: MessengerService}
}

func (h *HTTPHandler) SaveMessage(ctx *gin.Context) {
	zaplog.Logger.Info("http: Получен запрос SaveMessage", zaplog.LogHandler("client_ip", ctx.ClientIP()))

	//metric
	prometric.RequestCounter.WithLabelValues(ctx.Request.Method).Inc()
	var message domain.Message

	if err := ctx.ShouldBindJSON(&message); err != nil {
		zaplog.Logger.Error("http: Ошибка записи JSON", zaplog.LogError(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}
	err := h.s.SaveMessage(message)
	if err != nil {
		zaplog.Logger.Error("http: Ошибка сохранения message", zaplog.LogError(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": errors.Unwrap(err)})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Message сохранен успешно",
	})
	zaplog.Logger.Info("http: Message сохранен успешно", zaplog.LogString("message", message.Body))

}

func (h *HTTPHandler) ReadMessage(ctx *gin.Context) {
	zaplog.Logger.Info("http: Получен запрос ReadMessage", zaplog.LogHandler("client_ip", ctx.ClientIP()))

	//metric
	prometric.RequestCounter.WithLabelValues(ctx.Request.Method).Inc()

	id := ctx.Param("id")

	message, err := h.s.ReadMessage(id)
	if err != nil {
		zaplog.Logger.Error("http: Ошибка чтения message", zaplog.LogString("id", id), zaplog.LogError(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": errors.Unwrap(err)})
		return
	}

	ctx.JSON(http.StatusOK, message)
	zaplog.Logger.Info("http: Message прочитан успешно", zaplog.LogString("id", id))

}

func (h *HTTPHandler) ReadMessages(ctx *gin.Context) {
	zaplog.Logger.Info("http: Получен запрос ReadMessages", zaplog.LogHandler("client_ip", ctx.ClientIP()))

	//metric
	prometric.RequestCounter.WithLabelValues(ctx.Request.Method).Inc()

	slice, err := h.s.ReadMessages()

	if err != nil {
		zaplog.Logger.Error("http: Ошибка чтения messages", zaplog.LogError(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": errors.Unwrap(err)})
		return
	}

	ctx.JSON(http.StatusOK, slice)
	zaplog.Logger.Info("http: Messages прочитан успешно", zaplog.LogString("number", strconv.Itoa(len(slice))))

}
