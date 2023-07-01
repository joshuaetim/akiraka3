package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshuaetim/frontdesk/controller"
	"gorm.io/gorm"
)

type DashboardHandler interface {
	GetUsersCount(*gin.Context)
}

func NewDashboardHandler(db *gorm.DB) DashboardHandler {
	return dashboardHandler{
		controller: controller.NewDashboardController(db),
	}
}

type dashboardHandler struct {
	controller controller.DashboardController
}

func (dash dashboardHandler) GetUsersCount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": dash.controller.GetUsersCount()})
}
