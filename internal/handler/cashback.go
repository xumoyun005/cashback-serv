package handler

import (
	"cashback-serv/internal/service"
	"cashback-serv/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @title Cashback Service API
// @version 1.0
// @description Cashback amount control of users
// @host localhost:8080
// @BasePath /

type CashbackHandler struct {
	service *service.CashbackService
}

func NewCashbackHandler(service *service.CashbackService) *CashbackHandler {
	return &CashbackHandler{service: service}
}

func (h *CashbackHandler) RegisterRoutes(router *gin.Engine) {
	cashback := router.Group("/cashback")
	{
		cashback.POST("/increase", h.IncreaseCashback)
		cashback.POST("/decrease", h.DecreaseCashback)
		cashback.GET("/:turon_user_id", h.GetCashback)
		cashback.GET("/:turon_user_id/history", h.GetCashbackHistory)
	}
}

// @Summary Cashback increase
// @Description Increase cashback of the user
// @Tags cashback
// @Accept json
// @Produce json
// @Param request body models.CashbackRequest true "Cashback increase"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cashback/increase [post]
func (h *CashbackHandler) IncreaseCashback(c *gin.Context) {
	var req models.CashbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.HostIP = c.ClientIP()
	req.Type = c.GetHeader("User-Agent")

	if err := h.service.IncreaseCashback(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cashback successfully increased"})
}

// @Summary Cashback amount decrease
// @Description Cashback amount decrease of the user
// @Tags cashback
// @Accept json
// @Produce json
// @Param request body models.CashbackRequest true "Cashback amount decrease "
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cashback/decrease [post]
func (h *CashbackHandler) DecreaseCashback(c *gin.Context) {
	var req models.CashbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.HostIP = c.ClientIP()
	req.Type = c.GetHeader("User-Agent")

	if err := h.service.DecreaseCashback(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cashback successfully decreased"})
}

// @Summary GET Cashback
// @Description Cashback amount of the user
// @Tags cashback
// @Accept json
// @Produce json
// @Param turon_user_id path int true "Turon User ID"
// @Param from_date query string false "Start date" format(date) example(2024-03-01)
// @Param to_date query string false "End date" format(date) example(2024-03-20)
// @Success 200 {object} models.Cashback
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cashback/{turon_user_id} [get]
func (h *CashbackHandler) GetCashback(c *gin.Context) {
	turonUserID, err := strconv.ParseInt(c.Param("turon_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return
	}

	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	// Validate date formats if provided
	if fromDate != "" {
		if _, err := time.Parse("2006-01-02", fromDate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from_date format. Use YYYY-MM-DD"})
			return
		}
	}
	if toDate != "" {
		if _, err := time.Parse("2006-01-02", toDate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to_date format. Use YYYY-MM-DD"})
			return
		}
	}

	cashback, err := h.service.GetCashbackByUserID(turonUserID, fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cashback data"})
		return
	}

	if cashback == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cashback not found "})
		return
	}

	c.JSON(http.StatusOK, cashback)
}

// @Summary CashbackHistory of the user
// @Description Get cashback history with optional date filtering and pagination
// @Tags cashback
// @Accept json
// @Produce json
// @Param turon_user_id path int true "Turon User ID"
// @Param from_date query string false "Start date" format(date) example(2024-03-01)
// @Param to_date query string false "End date" format(date) example(2024-03-20)
// @Param page query int false "Page number" default(1) minimum(1)
// @Param page_size query int false "Items per page" default(10) minimum(1) maximum(100)
// @Success 200 {object} map[string]interface{} "data: array of cashback history, pagination: pagination info"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cashback/{turon_user_id}/history [get]
func (h *CashbackHandler) GetCashbackHistory(c *gin.Context) {
	turonUserID, err := strconv.ParseInt(c.Param("turon_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return
	}

	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 64)

	pagination := &models.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	history, err := h.service.GetCashbackHistoryByUserID(turonUserID, fromDate, toDate, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       history,
		"pagination": pagination,
	})
}
