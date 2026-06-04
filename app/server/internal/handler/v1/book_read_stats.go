package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

// BookReadStatsHandler 阅读统计处理器 (只读聚合查询)
type BookReadStatsHandler struct {
	svc *service.BookReadStatsService
}

func NewBookReadStatsHandler(svc *service.BookReadStatsService) *BookReadStatsHandler {
	return &BookReadStatsHandler{svc: svc}
}

// GetDailyStats 按日阅读统计
// @Summary   按日阅读统计
// @Tags      book-read-stats
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ReadStatsQuery  true  "查询参数"
// @Success  200  {object}  response.Response{data=[]dto.ReadEventDailyResponse}
// @Router   /api/book/read-stats/daily [post]
func (h *BookReadStatsHandler) GetDailyStats(c *gin.Context) {
	var req dto.ReadStatsQuery
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.GetDailyStats(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetBookStats 按书阅读统计
// @Summary   按书阅读统计
// @Tags      book-read-stats
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ReadStatsQuery  true  "查询参数"
// @Success  200  {object}  response.Response{data=[]dto.ReadEventBookResponse}
// @Router   /api/book/read-stats/books [post]
func (h *BookReadStatsHandler) GetBookStats(c *gin.Context) {
	var req dto.ReadStatsQuery
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.GetBookStats(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetTotalStats 总阅读统计
// @Summary   总阅读统计
// @Tags      book-read-stats
// @Security  BearerAuth
// @Produce   json
// @Success  200  {object}  response.Response{data=dto.ReadEventTotalResponse}
// @Router   /api/book/read-stats/total [get]
func (h *BookReadStatsHandler) GetTotalStats(c *gin.Context) {
	resp, err := h.svc.GetTotalStats(c.Request.Context(), utils.GetUserID(c))
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}
