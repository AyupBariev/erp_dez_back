package cleanhandler

import (
	"erp/internal/app/response"
	"erp/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AggregatorPayoutHandler struct {
	uc *application.AggregatorPayoutUseCase
}

func NewAggregatorPayoutHandler(uc *application.AggregatorPayoutUseCase) *AggregatorPayoutHandler {
	return &AggregatorPayoutHandler{uc: uc}
}

func (h *AggregatorPayoutHandler) ListDayPayouts(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	list, err := h.uc.GetDayPayouts(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.FromAggregatorDayPayoutList(list))
}

func (h *AggregatorPayoutHandler) PayDay(c *gin.Context) {
	type req struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channel := c.Param("channel")
	date := c.Param("date") // 2006-01-02
	if err := h.uc.MarkPaid(channel, date, r.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
