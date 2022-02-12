package api

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()

	if code != 0 {
		log.Fatalf("Failed to run")
	}

	os.Exit(code)
}
