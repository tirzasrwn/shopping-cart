package configs

import (
	"database/sql"
	"log"

	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port         int
	Environment  string
	Debug        bool
	JWTSecret    string
	DSN          string
	Domain       string
	DB           *sql.DB
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	MongoUrl     string
}

func InitializeAppConfig() {
	viper.SetConfigName("development.env") // allow directly reading from .env file
	viper.SetConfigType("env")
	viper.AddConfigPath("./configs")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	AppConfig.Port = viper.GetInt("PORT")
	AppConfig.Environment = viper.GetString("ENVIRONMENT")
	AppConfig.Debug = viper.GetBool("DEBUG")
	AppConfig.JWTSecret = viper.GetString("JWT_SECRET")
	AppConfig.DSN = viper.GetString("DSN")
	AppConfig.Domain = viper.GetString("DOMAIN")
	AppConfig.JWTIssuer = viper.GetString("JWTISSUER")
	AppConfig.JWTAudience = viper.GetString("JWTAUDIENCE")
	AppConfig.CookieDomain = viper.GetString("COOKIEDOMAIN")
	AppConfig.MongoUrl = viper.GetString("MONGOURL")
	log.Println("[INIT] configuration loaded")
}
