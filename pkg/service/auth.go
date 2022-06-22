package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"gitlab.com/trencetech/mypipe-api/db/model"
)

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"give_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

type AppleIDClaims struct{}

type JWTConfig struct {
	Algo      jwt.SigningMethod
	ExpiresIn int
	Key       string
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    string
}

func (s Service) IssueAuthToken(user model.User) (AuthToken, error) {
	accessToken, refreshToken, expiresAt, err := s.generateTokenPair(user)
	if err != nil {
		return AuthToken{}, err
	}
	authTokens := AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
	return authTokens, nil
}
func (s Service) generateTokenPair(user model.User) (accessToken, refreshToken, expiryTime string, err error) {
	atExpiresIn := time.Now().Add(time.Duration(s.JWTConfig.ExpiresIn) * time.Second).Unix()
	rtExpiresIn := time.Now().Add(30 * (24 * time.Hour)).Unix()
	exToTime := time.Now().Add(time.Duration(s.JWTConfig.ExpiresIn) * time.Second)
	expiryTime = exToTime.Format(time.RFC3339Nano)
	at := jwt.NewWithClaims(s.JWTConfig.Algo, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      atExpiresIn,
	})

	rt := jwt.NewWithClaims(s.JWTConfig.Algo, jwt.MapClaims{
		"sub": user.ID,
		"exp": rtExpiresIn,
	})

	accessToken, err = at.SignedString([]byte(s.JWTConfig.Key))
	if err != nil {
		return "", "", "", err
	}
	refreshToken, err = rt.SignedString([]byte(s.JWTConfig.Key))
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, expiryTime, nil
}

func (s Service) ValidateGoogleJWT(tokenString string) (GoogleClaims, error) {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	claimStruct := GoogleClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimStruct,
		retrieveKeyFromPem,
	)

	if err != nil {
		logger.Err(err).Msg("could not execute jwt.ParseWithClaims")
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		logger.Info().Msg("invalid google JWT")
		return GoogleClaims{}, errors.New("invalid google JWT")
	}
	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		logger.Info().Msg("GOOGLE_JWT_ERROR: iss is invalid")
		return GoogleClaims{}, errors.New("iss is invalid")
	}

	if claims.Audience != os.Getenv("GOOGLE_CLIENT_ID") {
		logger.Info().Msg("GOOGLE_JWT_ERROR: aud is invalid")
		return GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		logger.Info().Msg("GOOGLE_JWT_ERROR: jwt expired")
		return GoogleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}

func retrieveKeyFromPem(t *jwt.Token) (interface{}, error) {
	pem, err := getGooglePublicKey(fmt.Sprintf("%s", t.Header["kid"]))
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
	if err != nil {
		return nil, err
	}

	return key, nil
}

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}
