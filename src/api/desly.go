package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createDeslyRequest struct {
	Redirect string `json:"redirect" binding:"required,min=15"`
}

func (server *Server) createDesly(ctx *gin.Context) {
	var req createDeslyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	desly, err := server.store.CreateDesly(ctx, req.Redirect)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, desly)
}

type getDeslyRequest struct {
	Desly string `uri:"desly" binding:"required,len=6"`
}

func (server *Server) getDesly(ctx *gin.Context) {
	var req getDeslyRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	desly, err := server.store.GetDesly(ctx, req.Desly)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, desly)
}

type redirectRequest struct {
	Desly string `uri:"desly" binding:"required,len=6"`
}

func (server *Server) redirect(ctx *gin.Context) {
	var req redirectRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	desly, err := server.store.GetDesly(ctx, req.Desly)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, desly.Redirect)
}
