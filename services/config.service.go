package services

import (
	database "studenent_manager/models"

	"github.com/spf13/viper"
)

var Config *database.Config

func LoadConfig() {
	v := viper.New()
	v.AutomaticEnv() // set enviroment variable if has
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("MODE", "debug")
	v.SetConfigType("dotenv")
	v.SetConfigName(".env")
	v.AddConfigPath("./")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&Config); err != nil {
		panic(err)
	}

	if err := Config.Validate(); err != nil {
		panic(err)
	}
}

func InitializeRepository() {
	InitializeStudentRepository()
}
