package config

import (
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

type database struct {
	URL string
}

type jwt struct {
	Secret string
	Issuer string
}
type Config struct {
	Database database
	JWT      jwt
}

func LoadEnv(fileName string) {
	reg := regexp.MustCompile(`^(.*` + "twitter-go-graphql" + `)`)
	dir, _ := os.Getwd()
	rootPath := reg.Find([]byte(dir))
	if err := godotenv.Load(string(rootPath) + `/` + fileName); err != nil {
		godotenv.Load()
	}

}
func New() *Config {
	//godotenv.Load()
	return &Config{
		Database: database{
			URL: os.Getenv("DATABASE_URL"),
		},
		JWT: jwt{
			Secret: os.Getenv("JWT_SECRET"),
			Issuer: os.Getenv("DOMAIN"),
		},
	}
}
