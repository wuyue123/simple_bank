package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "pxsemic.com/simplebank/db/sqlc"
	"pxsemic.com/simplebank/util"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)  
	os.Exit(m.Run())
}


func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store,config)
	require.NoError(t, err)

	return server
}
