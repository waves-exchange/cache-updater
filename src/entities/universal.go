package entities;

import (
	"os"

	"github.com/joho/godotenv"
)
type DAappNumberRecord struct {
	Key, Type string
	Value *int
}

type DAppStringRecord struct {
	Key, Type string
	Value *string
}

func StrArrayContains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func unwrapDefaultRegex (rawregex *string, defaultRegex string) string {
	if rawregex == nil {
		return defaultRegex
	} else {
		return *rawregex
    }
}

func GetDBCredentials () (string, string, string) {
	envLoadErr := godotenv.Load(".env")
	if envLoadErr != nil {
		_ = godotenv.Load(".env.example")
	}

	dbuser := os.Getenv("DB_USERNAME");
	dbpass := os.Getenv("DB_PASS");
	dbdatabase := os.Getenv("DB_NAME");

	return dbuser, dbpass, dbdatabase
}