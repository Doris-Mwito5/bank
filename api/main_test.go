package api

import (
	"os"
	"testing"
	"time"

	db "github.com/Doris-Mwito5/simple-bank/internal/db/sqlc"
	"github.com/Doris-Mwito5/simple-bank/internal/utils"
	"github.com/Doris-Mwito5/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
    config := util.Config{
        TokenSymmetricKey:   utils.RandomString(32),
        AccessTokenDuration: time.Minute,
    }

    server, err := NewServer(store, config)
    require.NoError(t, err)

    return server
}


func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
