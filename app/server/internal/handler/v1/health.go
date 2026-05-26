package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"boread/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Ping 健康检查
// @Summary  健康检查
// @Tags     system
// @Produce  json
// @Success  200  {object}  response.Response
// @Router   /ping [get]
func (h *HealthHandler) Ping(c *gin.Context) {
	response.Success(c, gin.H{"status": "ok"})
}

func (h *HealthHandler) NoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "route not found"})
}

func (h *HealthHandler) NoMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"code": 405, "message": "method not allowed"})
}
