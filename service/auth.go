package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"gofr-auth-ui-app/model"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gofr.dev/pkg/gofr"
	"golang.org/x/crypto/bcrypt"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "my_secret_key"
	}
	return []byte(secret)
}

func getJWTExpiry() time.Duration {
	expiryStr := os.Getenv("JWT_EXPIRY")
	if expiryStr == "" {
		expiryStr = "3600"
	}
	expiry, _ := strconv.Atoi(expiryStr)
	return time.Duration(expiry) * time.Second
}

func getRefreshTokenExpiry() time.Duration {
	expiryStr := os.Getenv("REFRESH_TOKEN_EXPIRY")
	if expiryStr == "" {
		expiryStr = "604800" // 7 days
	}
	expiry, _ := strconv.Atoi(expiryStr)
	return time.Duration(expiry) * time.Second
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func generateAccessToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(getJWTExpiry()).Unix(),
	})
	return token.SignedString(getJWTSecret())
}

func SignUp(ctx *gofr.Context, user *model.User) (interface{}, error) {
	if ctx.SQL == nil {
		return nil, errors.New("database not initialized")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	_, err := ctx.SQL.Exec("INSERT INTO users (email, password, role) VALUES ($1, $2, $3)", user.Email, user.Password, "user")
	if err != nil {
		return nil, err
	}
	return "User registered", nil
}

func Login(ctx *gofr.Context, creds *model.User) (interface{}, error) {
	if ctx.SQL == nil {
		return nil, errors.New("database not initialized")
	}

	var dbUser model.User
	err := ctx.SQL.QueryRow("SELECT id, email, password, role FROM users WHERE email = $1", creds.Email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password, &dbUser.Role)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(creds.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := generateAccessToken(&dbUser)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiryTime := time.Now().Add(getRefreshTokenExpiry())
	_, err = ctx.SQL.Exec("UPDATE users SET refresh_token = $1, token_expiry = $2 WHERE id = $3", refreshToken, expiryTime, dbUser.ID)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func RefreshToken(ctx *gofr.Context, refreshReq *model.RefreshRequest) (interface{}, error) {
	if ctx.SQL == nil {
		return nil, errors.New("database not initialized")
	}

	var dbUser model.User
	err := ctx.SQL.QueryRow("SELECT id, email, role, refresh_token, token_expiry FROM users WHERE refresh_token = $1", refreshReq.RefreshToken).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Role, &dbUser.RefreshToken, &dbUser.TokenExpiry)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(dbUser.TokenExpiry) {
		return nil, errors.New("refresh token expired")
	}

	accessToken, err := generateAccessToken(&dbUser)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiryTime := time.Now().Add(getRefreshTokenExpiry())
	_, err = ctx.SQL.Exec("UPDATE users SET refresh_token = $1, token_expiry = $2 WHERE id = $3", newRefreshToken, expiryTime, dbUser.ID)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
