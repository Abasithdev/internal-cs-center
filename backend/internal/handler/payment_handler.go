package handler

import (
	common_errors "errors"
	"net/http"

	"abasithdev.github.io/internal-cs-center-backend/internal/errors"
	"abasithdev.github.io/internal-cs-center-backend/internal/service"
	"abasithdev.github.io/internal-cs-center-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(payment *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: payment}
}

// ListPayments godoc
// @Summary List payments
// @Description Get list of payments with filters
// @Tags payments
// @Accept json
// @Produce json
// @Param page query int false "page number" default(1)
// @Param size query int false "page size" default(10)
// @Param status query string false "filter by status"
// @Param search query string false "search term"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /payments [get]
func (paymentHandler *PaymentHandler) ListPayments(context *gin.Context) {
	page := utils.QueryInt(context, "page", 1)
	size := utils.QueryInt(context, "size", 10)
	status := context.Query("status")
	search := context.Query("search")
	sortBy := context.DefaultQuery("sortBy", "date")
	orderBy := context.DefaultQuery("orderBy", "desc")

	params := service.ListRequest{
		Page:    page,
		Size:    size,
		Status:  status,
		Search:  search,
		SortBy:  sortBy,
		OrderBy: orderBy,
	}

	total := paymentHandler.paymentService.GetTotalByFilter(params)
	result := paymentHandler.paymentService.GetList(params)
	completed, process, failed := paymentHandler.paymentService.GetStatusSummary()

	context.JSON(http.StatusOK, gin.H{
		"meta": result,
		"summary": gin.H{
			"total":      total,
			"completed":  completed,
			"processing": process,
			"failed":     failed,
		},
	})
}

// ReviewPayment godoc
// @Summary Review payment
// @Description Review a payment (operational role required)
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "payment id"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /payments/{id}/review [put]
func (paymentHandler *PaymentHandler) ReviewPayment(ctx *gin.Context) {
	role, _ := ctx.Get("role")

	if role.(string) != "operational" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Forbidden"})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must not empty"})
		return
	}

	if err := paymentHandler.paymentService.Review(id); err != nil {
		var notFoundErr *errors.NotFoundError
		if common_errors.As(err, &notFoundErr) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}
		// for any other error, return internal server error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
