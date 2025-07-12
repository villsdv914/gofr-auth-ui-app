package migration

import (
	"gofr.dev/pkg/gofr/migration"
)

func AddRefreshTokenColumns(d migration.Datasource) error {
	query := `
	ALTER TABLE users 
	ADD COLUMN IF NOT EXISTS refresh_token VARCHAR(255),
	ADD COLUMN IF NOT EXISTS token_expiry TIMESTAMP;`

	_, err := d.SQL.Exec(query)
	return err
}