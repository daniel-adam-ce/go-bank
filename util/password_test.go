package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	pasword := RandomString(6)

	hashedPassword, err := HashPassword(pasword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	err = CheckPassword(pasword, hashedPassword)
	require.NoError(t, err)

	wrongPasword := RandomString(6)
	err = CheckPassword(wrongPasword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
