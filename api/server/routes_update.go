package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iron-io/functions/api/models"
	"github.com/iron-io/runner/common"
)

func handleRouteUpdate(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	log := common.Logger(ctx)

	var wroute models.RouteWrapper

	err := c.BindJSON(&wroute)
	if err != nil {
		log.WithError(err).Debug(models.ErrInvalidJSON)
		c.JSON(http.StatusBadRequest, simpleError(models.ErrInvalidJSON))
		return
	}

	if wroute.Route == nil {
		log.WithError(err).Error(models.ErrInvalidJSON)
		c.JSON(http.StatusBadRequest, simpleError(models.ErrRoutesMissingNew))
		return
	}

	wroute.Route.AppName = c.Param("app")
	wroute.Route.Path = c.Param("route")

	if err := wroute.Validate(); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, simpleError(err))
		return
	}

	_, err = Api.Datastore.UpdateRoute(wroute.Route)
	if err != nil {
		log.WithError(err).Debug(models.ErrRoutesUpdate)
		c.JSON(http.StatusInternalServerError, simpleError(models.ErrRoutesUpdate))
		return
	}

	c.JSON(http.StatusOK, routeResponse{"Route successfully updated", wroute.Route})
}
