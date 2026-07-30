package api

import (
	"errors"

	db "my-gin-app/db/sqlc"
	"my-gin-app/token"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Balance  int64  `json:"balance" binding:"required,min=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503", "23505":
				ctx.JSON(403, errorResponse(err))
				return
			}
		}
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(c *gin.Context) {
	var req getAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(404, errorResponse(errors.New("không tìm thấy tài khoản")))
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account does not belong to the authenticated user")
		c.JSON(403, errorResponse(err))
		return
	}

	c.JSON(200, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(c *gin.Context) {
	var req listAccountsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	accounts, err := server.store.ListAccounts(c.Request.Context(), db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		c.JSON(500, errorResponse(err))
		return
	}

	c.JSON(200, accounts)
}
