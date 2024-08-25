package store

import (
	"context"
	"os"
	"path/filepath"
)

func (p postgres) Migrate(path string) error {
	// Get the absolute path of the migrations directory
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}
	migrationDir := filepath.Join(rootDir, path)

	// Read the migration files
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return err
	}

	// Execute the migrations
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrationSQL, err := os.ReadFile(filepath.Join(migrationDir, file.Name()))
			if err != nil {
				return err
			}

			_, err = p.conn.Exec(context.Background(), string(migrationSQL))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
