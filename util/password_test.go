package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	pasword := RandomString(6)

	hashedPassword1, err := HashPassword(pasword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)
	err = CheckPassword(pasword, hashedPassword1)
	require.NoError(t, err)

	wrongPasword := RandomString(6)
	err = CheckPassword(wrongPasword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// this is somewhat unnecessary
	hashedPassword2, err := HashPassword(pasword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
