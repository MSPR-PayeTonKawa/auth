package tokens

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	userID := "user1"
	td, err := CreateToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, td.AccessToken)
	assert.NotEmpty(t, td.RefreshToken)
	assert.NotZero(t, td.AtExpires)
	assert.NotZero(t, td.RtExpires)
}
