package api

import (
	"database/sql"
	db "desly/db/sqlc"
	"desly/token"
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateDeslyParams{
		Redirect: req.Redirect,
		Owner:    authPayload.Username,
	}

	desly, err := server.store.CreateDesly(ctx, arg)

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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.GetDeslyParams{
		Desly: req.Desly,
		Owner: authPayload.Username,
	}

	desly, err := server.store.GetDesly(ctx, arg)

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

	redirect, err := server.store.GetRedirectByDesly(ctx, req.Desly)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, redirect)
}
