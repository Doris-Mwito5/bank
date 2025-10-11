package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Doris-Mwito5/simple-bank/internal/token"
	"github.com/gin-gonic/gin"
)

const (
    authorizationHeaderKey  = "Authorization"
    authorizationTypeBearer = "bearer" 
    authorizationPayloadKey = "authorization_payload"
)


func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader) // ✅ use the actual header

		if len(fields) < 2 {
			err := errors.New("invalid authorization header")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type: %s", authorizationType)
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
