package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"pxsemic.com/simplebank/util"
)

func TestJwtMaker_CreateToken(t *testing.T) {
	maker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration, TokenTypeAccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, TokenTypeAccessToken)
	require.NoError(t, err)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.Equal(t, TokenTypeAccessToken, payload.Type)
	require.WithinDuration(t, issuedAt, payload.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt.Time, time.Second)
}

// TestExpiredToken tests the behavior of the JwtMaker when the token is expired.
func TestExpiredToken(t *testing.T) {
	maker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	token, payload, err := maker.CreateToken(username, role, -duration, TokenTypeAccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, TokenTypeAccessToken)
	require.Error(t, err)
	require.Equal(t, ErrExpiredToken, err)
	require.Nil(t, payload)
}
