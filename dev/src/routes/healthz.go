package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthz godoc
// @Summary Health Check
// @Description checks the health of the service
// @Accept plain
// @Produce plain
// @Success 200 "healthy"
// @Router /healthz [get]
func healthz(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}
