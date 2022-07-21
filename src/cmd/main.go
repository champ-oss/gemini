package main

import (
	"github.com/champ-oss/gemini/adapter/github"
	cfg "github.com/champ-oss/gemini/config"
	"github.com/champ-oss/gemini/repository"
	gemini "github.com/champ-oss/gemini/service/common"
)

func main() {
	config := cfg.LoadConfig()

	repo, err := repository.NewRepository(config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	if err != nil {
		panic(err)
	}

	githubClient, err := github.NewAdapter(config.GitHubAppId, config.GitHubInstallationId, config.GitHubPem)
	if err != nil {
		panic(err)
	}

	svc := gemini.NewService(repo, githubClient, config.MinutesBetweenChecks)

	err = svc.StartGemini(config.Repos)
	if err != nil {
		panic(err)
	}
}
