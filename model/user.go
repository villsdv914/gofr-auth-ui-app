package model

import "time"

type User struct {
    ID           string    `json:"id"`
    Email        string    `json:"email"`
    Password     string    `json:"password"`
    Role         string    `json:"role"`
    RefreshToken string    `json:"refresh_token,omitempty"`
    TokenExpiry  time.Time `json:"token_expiry,omitempty"`
}

type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token"`
}