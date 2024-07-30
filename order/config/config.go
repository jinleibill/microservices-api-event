package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func GetDataSourceURL() string {
	cfg := fmt.Sprintf("host=%s user=root password=secret dbname=orders port=5432 sslmode=disable TimeZone=Asia/Shanghai", getEnvironmentValue("DATA_SOURCE_URL"))

	return cfg
}

func GetApplicationPort() int {
	portStr := getEnvironmentValue("APPLICATION_PORT")
	port, err := strconv.Atoi(portStr)

	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}

	return port
}

func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}

	return os.Getenv(key)
}
