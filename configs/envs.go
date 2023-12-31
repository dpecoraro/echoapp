package configs

import (
	"os"
)

func EnvMongoURI() string {
	return os.Getenv("MONGODBURI")
}

func GetDatabaseName() string {
	return os.Getenv("DATABASENAME")
}
