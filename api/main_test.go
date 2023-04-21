package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/mativm02/bank_system/db/sqlc"
	"github.com/mativm02/bank_system/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: 15 * time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	// Set Gin to test mode so it doesn't output akward logs
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
