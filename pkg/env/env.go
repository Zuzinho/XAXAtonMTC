package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type envVarName string

const (
	tgBotToken         envVarName = "TG_BOT_TOKEN"
	smsApiID           envVarName = "SMS_API_ID"
	dbConnectionString envVarName = "BD_CONNECTION_STRING"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
}

func MustTokenBot() string {
	token, exists := os.LookupEnv(string(tgBotToken))
	if !exists {
		log.Fatal("TgBot token not found")
	}

	return token
}

func MustSmsApiID() string {
	apiID, exists := os.LookupEnv(string(smsApiID))
	if !exists {
		log.Fatal("Sms api id not found")
	}

	return apiID
}

func MustDBConnectionString() string {
	connString, exists := os.LookupEnv(string(dbConnectionString))
	if !exists {
		log.Fatal("DB connection string not found")
	}

	return connString
}
