package tokens

import (
	"os"
	"time"

	"github.com/MSPR-PayeTonKawa/auth/types"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = os.Getenv("JWT_KEY")

func CreateToken(userID string) (*types.TokenDetails, error) {
	td := &types.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()   // Access token expires after 15 minutes
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // Refresh token expires after 7 days

	atClaims := &types.Claims{
		UserID: userID,
		Claims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.AtExpires, 0)),
		},
	}

	rtClaims := &types.Claims{
		UserID: userID,
		Claims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.RtExpires, 0)),
		},
	}

	var err error
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString(JwtKey)
	if err != nil {
		return nil, err
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(JwtKey)
	if err != nil {
		return nil, err
	}

	return td, nil
}
