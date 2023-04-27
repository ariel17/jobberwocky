package configs

import "os"

const (
	databaseNameKey     = "DATABASE_NAME"
	defaultDatabaseName = "production.db"
)

var (
	databaseName = ""
)

func GetDatabaseName() string {
	return databaseName
}

func init() {
	if databaseName = os.Getenv(databaseNameKey); databaseName == "" {
		databaseName = defaultDatabaseName
	}
}