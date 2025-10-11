package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Doris-Mwito5/simple-bank/internal/db/sqlc"
	"github.com/Doris-Mwito5/simple-bank/internal/token"
	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	FromAccountID int64   `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64   `json:"to_account_id" binding:"required,min=1"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Currency      string  `json:"currency" binding:"required"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req TransferRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := s.validateAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesnt belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	_, valid = s.validateAccount(ctx, req.ToAccountID, req.Currency)

	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transfer)
}

func (s *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
