package usertokens

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nisil-thaing/golang-gingonic-practice/internal/app/users"
)

func GenerateTokens(user users.UserSchema, secretKey string) (*UserTokensPublicInfo, error) {
	currentTime := time.Now().Local()
	accessExpiresAt := currentTime.Add(time.Duration(24) * time.Hour)
	refreshExpiresAt := currentTime.Add(time.Duration(168) * time.Hour)
	accessClaims := JWTSigningClaims{
		UserID:    user.UserID,
		Email:     user.Email,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		UserType:  user.Type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
		},
	}

	signedAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secretKey))

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	refreshClaims := JWTSigningClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
		},
	}

	signedRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	userTokens := UserTokensPublicInfo{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
		ExpiresAt:    refreshExpiresAt,
	}

	return &userTokens, nil
}
