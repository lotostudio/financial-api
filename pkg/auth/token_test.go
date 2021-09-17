package auth

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func newTestJWTManager(t *testing.T) TokenManager {
	m, err := NewJWTManager("key", time.Duration(1)*time.Hour, 32)

	require.NoError(t, err)

	return m
}

func TestNewJWTManager_emptyKey(t *testing.T) {
	_, err := NewJWTManager("", time.Duration(1)*time.Hour, 32)

	require.Errorf(t, err, "empty signing key")
}

func TestJWTManager_IssueAndDecode(t *testing.T) {
	m := newTestJWTManager(t)
	userId := "1"

	token, err := m.Issue(userId)

	require.NoError(t, err)
	require.NotNil(t, token)

	id, err := m.Decode(token)

	require.NoError(t, err)
	require.Equal(t, userId, id)
}

func TestJWTManager_DecodeErr(t *testing.T) {
	m := newTestJWTManager(t)

	_, err := m.Decode("qwe")

	require.Error(t, err)
}

func TestJWTManager_Refresh(t *testing.T) {
	m := newTestJWTManager(t)

	_, err := m.Random()

	require.NoError(t, err)
}
