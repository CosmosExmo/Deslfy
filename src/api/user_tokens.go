package api

/* type createUserTokenRequest struct {
	ExpireAt int64 `json:"expire_at" binding:"required"`
}

func (server *Server) createUserToken(ctx *gin.Context) {
	var req createUserTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	expireAt := time.Unix(req.ExpireAt, 0)
	duration := time.Until(expireAt)
	accessToken, _, err := server.tokenMaker.CreateToken(authPayload.Username, duration, token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserTokenParams{
		Owner:    authPayload.Username,
		Token:    accessToken,
		ExpireAt: expireAt,
	}

	token, err := server.store.CreateUserToken(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, token)
}

type getUserTokenRequest struct {
	ID int32 `uri:"id" binding:"required,min=0"`
}

func (server *Server) getUserToken(ctx *gin.Context) {
	var req getUserTokenRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.GetUserTokenParams{
		Owner: authPayload.Username,
		ID: req.ID,
	}
	token, err := server.store.GetUserToken(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, token)
}

func (server *Server) getUserTokens(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	token, err := server.store.GetUserTokens(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, token)
}

type deleteUserTokenRequest struct {
	ID int32 `json:"id" binding:"required,min=0"`
}

var deleteSuccessResponse = "Delete Succesfull"

func (server *Server) deleteUserToken(ctx *gin.Context) {
	var req deleteUserTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.DeleteUserTokenParams{
		ID: req.ID,
		Owner: authPayload.Username,
	}
	err := server.store.DeleteUserToken(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, deleteSuccessResponse)
} */
