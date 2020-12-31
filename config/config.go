package config

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ENV string

type PopulateConfig struct {
	Migrate bool
	Init    bool
	File    string
}

func InitConfig(env string, path string) {
	ENV = env

	// load config file
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error while reading config file: %v", err)
	}

	ConfigureLogger()
}

func ConfigureLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	partialPath := fmt.Sprintf("%s.logging.", ENV)
	logOutputType := fmt.Sprintf("%s", viper.Get(partialPath + "outputType"))
	logFilePath := fmt.Sprintf("%s", viper.Get( partialPath + "filePath"))

	if logOutputType == "console" {
		// pass
	} else if logOutputType == "file" {
		logFile, err := os.OpenFile(logFilePath, os.O_RDWR, os.ModeAppend)
		if err != nil {
			log.Fatalf("Error while opening log file: %v", err)
		}
		log.SetOutput(logFile)
	} else {
		log.Fatalf("Invalid logging.outputType")
	}
}

func GetByKey(key string) string {
	return fmt.Sprintf("%s", viper.Get(key))
}

func GetPostgresDSN() string {
	partialPath := fmt.Sprintf("%s.database.", ENV)
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		viper.Get(partialPath+"host"),
		viper.Get(partialPath+"user"),
		viper.Get(partialPath+"password"),
		viper.Get(partialPath+"dbname"),
		viper.Get(partialPath+"port"),
		viper.Get(partialPath+"sslmode"))
}

func GetPopulateConfig() PopulateConfig {
	partialPath := fmt.Sprintf("%s.database.populate.", ENV)
	migrate, _ := strconv.ParseBool(fmt.Sprint(viper.Get(partialPath + "migrate")))
	init, _ := strconv.ParseBool(fmt.Sprint(viper.Get(partialPath + "init")))
	file := fmt.Sprint(viper.Get(partialPath + "file"))

	return PopulateConfig{
		Migrate: migrate,
		Init:    init,
		File:    file}
}
