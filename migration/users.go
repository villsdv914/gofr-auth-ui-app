package migration

import (
	"gofr.dev/pkg/gofr/migration"
)

func CreateUsersTable(d migration.Datasource) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(50) DEFAULT 'user',
		refresh_token VARCHAR(255),
		token_expiry TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := d.SQL.Exec(query)
	return err
}
