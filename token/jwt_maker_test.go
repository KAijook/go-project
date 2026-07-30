package token

import (
	"fmt"
	"my-gin-app/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJwtMaker(t *testing.T) {
	maker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, int64(duration))
	require.NoError(t, err)
	require.NotEmpty(t, token)
	fmt.Println(token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}
