package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"XConf/admin-api/config"
)

func CreateApp(c *gin.Context) {
	var req = struct {
		AppName     string `json:"appName"        binding:"required"`
		Description string `json:"description"`
	}{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	app, err := config.CreateApp(req.AppName, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

func DeleteApp(c *gin.Context) {
	var req = struct {
		AppName string `json:"appName"        binding:"required"`
	}{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := config.DeleteApp(req.AppName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func ListApps(c *gin.Context) {
	apps, err := config.ListApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, apps)
}
