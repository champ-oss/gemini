package config

import (
	"encoding/base64"
	"github.com/champ-oss/gemini/model"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Debug                bool
	Repos                []*model.Repo
	GitHubAppId          int64
	GitHubInstallationId int64
	GitHubPem            []byte
	DbHost               string
	DbPort               string
	DbUsername           string
	DbPassword           string
	DbName               string
	MinutesBetweenChecks float64
	DropTables           bool
}

// LoadConfig loads configuration values from environment variables
func LoadConfig() *Config {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.Info("Loading configuration from environment")

	config := Config{
		Debug:                parseBool("DEBUG", true),
		Repos:                parseRepos("REPOS", ",", "/"),
		GitHubAppId:          parseInt64("GITHUB_APP_ID", 0),
		GitHubInstallationId: parseInt64("GITHUB_INSTALLATION_ID", 0),
		GitHubPem:            parseBase64("GITHUB_PEM", []byte("")),
		DbHost:               parseString("DB_HOST", "localhost"),
		DbPort:               parseString("DB_PORT", "3306"),
		DbUsername:           parseString("DB_USERNAME", "root"),
		DbPassword:           parseString("DB_PASSWORD", "secret"),
		DbName:               parseString("DB_NAME", "gemini"),
		MinutesBetweenChecks: parseFloat("MINUTES_BETWEEN_CHECKS", 5),
		DropTables:           parseBool("DROP_TABLES", false),
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debugging mode enabled")
	}

	return &config
}

// parseBool parses an environment variable as a boolean value
func parseBool(key string, fallback bool) bool {
	if value := os.Getenv(key); strings.ToLower(value) == "true" {
		return true
	}
	if value := os.Getenv(key); strings.ToLower(value) == "false" {
		return false
	}
	return fallback
}

// parseRepos parses an environment variable containing a list of repos
func parseRepos(key string, repoSeparator string, ownerSeparator string) []*model.Repo {
	envValue := os.Getenv(key)
	var repos []*model.Repo

	if envValue == "" {
		return repos
	}
	repoStrings := strings.Split(envValue, repoSeparator)

	for _, repoString := range repoStrings {
		if strings.Count(repoString, ownerSeparator) != 1 {
			log.Errorf("Repo definition is invalid: %s", repoString)
			continue
		}

		repoString = strings.Replace(repoString, " ", "", -1)

		repos = append(repos, &model.Repo{
			Owner: strings.Split(repoString, ownerSeparator)[0],
			Name:  strings.Split(repoString, ownerSeparator)[1],
		})
	}
	return repos
}

// parseString parses an environment variable as a string value
func parseString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// parseFloat parses an environment variable as a float value
func parseFloat(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Errorf("Unable to parse value of %s into float: %s", key, value)
		return fallback
	}
	return parsed
}

// parseBase64 parses an environment variable as base64 bytes
func parseBase64(key string, fallback []byte) []byte {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		log.Errorf("Unable to parse value of %s as base64: %s", key, value)
		return fallback
	}
	return decoded
}

// parseInt64 parses an environment variable as int64
func parseInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Errorf("Unable to parse value of %s as int64: %s", key, value)
		return fallback
	}
	return parsed
}
