package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/util"
	//need _ because we are not explicitly using it
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		TokenDuration:     time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}
