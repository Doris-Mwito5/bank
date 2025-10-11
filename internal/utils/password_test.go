package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassord, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassord)

	err = ValidatePassword(password, hashedPassord)
	require.NoError(t, err)
	
	wrongPassword := RandomString(6)
	err = ValidatePassword(wrongPassword, hashedPassord)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
