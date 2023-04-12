package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// Set Gin to test mode so it doesn't output akward logs
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
