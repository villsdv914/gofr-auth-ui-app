package main

import (
	"gofr-auth-ui-app/handler"
	"gofr-auth-ui-app/middleware"
	usermigration "gofr-auth-ui-app/migration"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/migration"
)

func main() {
	app := gofr.New()

	// Run migrations
	app.Migrate(map[int64]migration.Migrate{
		1: {UP: usermigration.CreateUsersTable},
		2: {UP: usermigration.AddRefreshTokenColumns},
	})

	// Apply JWT middleware globally but with route exclusions
	app.UseMiddleware(middleware.JWTAuthMiddleware())
	
	app.AddStaticFiles("/ui", "./ui")
	
	// Public routes (authentication is handled by middleware exclusions)
	app.POST("/signup", handler.SignUp)
	app.POST("/login", handler.Login)
	app.POST("/refresh-token", handler.RefreshToken)
	
	// Protected routes (require authentication)
	app.GET("/me", handler.Me)

	app.Run()
}
