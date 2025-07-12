package handler

import (
	"gofr-auth-ui-app/model"
	"gofr-auth-ui-app/service"

	"github.com/golang-jwt/jwt/v5"
	"gofr.dev/pkg/gofr"
)

func SignUp(ctx *gofr.Context) (interface{}, error) {
	var u model.User
	if err := ctx.Bind(&u); err != nil {
		return nil, err
	}
	return service.SignUp(ctx, &u)
}

func Login(ctx *gofr.Context) (interface{}, error) {
	var creds model.User
	if err := ctx.Bind(&creds); err != nil {
		return nil, err
	}
	return service.Login(ctx, &creds)
}

func Me(ctx *gofr.Context) (interface{}, error) {
	// First try to get claims from custom middleware context
	type contextKey string
	const ClaimsKey contextKey = "jwt_claims"
	
	if claimsValue := ctx.Request.Context().Value(ClaimsKey); claimsValue != nil {
		if claims, ok := claimsValue.(jwt.MapClaims); ok {
			return map[string]interface{}{
				"source": "custom_middleware",
				"claims": claims,
			}, nil
		}
	}

	// Fallback to GoFr's built-in auth info
	authInfo := ctx.GetAuthInfo()
	if authInfo != nil {
		claims := authInfo.GetClaims()
		if claims != nil {
			return map[string]interface{}{
				"source": "gofr_auth_info",
				"claims": claims,
			}, nil
		}
	}

	return map[string]interface{}{
		"error": "No claims found",
		"context_keys": "checked custom middleware context",
	}, nil
}

func RefreshToken(ctx *gofr.Context) (interface{}, error) {
	var refreshReq model.RefreshRequest
	if err := ctx.Bind(&refreshReq); err != nil {
		return nil, err
	}
	return service.RefreshToken(ctx, &refreshReq)
}
