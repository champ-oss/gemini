package config

import (
	"github.com/champ-oss/gemini/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_ParseBool_true(t *testing.T) {
	_ = os.Setenv("FOO", "true")
	assert.True(t, parseBool("FOO", false))
	_ = os.Setenv("FOO", "True")
	assert.True(t, parseBool("FOO", false))
}

func Test_ParseBool_false(t *testing.T) {
	_ = os.Setenv("FOO", "false")
	assert.False(t, parseBool("FOO", true))
	_ = os.Setenv("FOO", "False")
	assert.False(t, parseBool("FOO", true))
}

func Test_ParseBool_invalid(t *testing.T) {
	_ = os.Setenv("FOO", "nonsense")
	assert.False(t, parseBool("FOO", false))
}

func Test_ParseBool_unset(t *testing.T) {
	os.Clearenv()
	assert.False(t, parseBool("FOO", false))
	os.Clearenv()
	assert.True(t, parseBool("FOO", true))
}

func Test_ParseRepos_success(t *testing.T) {
	_ = os.Setenv("FOO", "owner1/repo1,owner2/repo2")
	expected := []*model.Repo{
		{Owner: "owner1", Name: "repo1"},
		{Owner: "owner2", Name: "repo2"},
	}
	assert.Equal(t, expected, parseRepos("FOO", ",", "/"))
}

func Test_ParseRepos_spaces(t *testing.T) {
	_ = os.Setenv("FOO", "owner1/repo1, owner2/repo2 ")
	expected := []*model.Repo{
		{Owner: "owner1", Name: "repo1"},
		{Owner: "owner2", Name: "repo2"},
	}
	assert.Equal(t, expected, parseRepos("FOO", ",", "/"))
}

func Test_ParseRepos_invalid(t *testing.T) {
	_ = os.Setenv("FOO", "owner1/repo1,owner2/repo2,nonsense")
	expected := []*model.Repo{
		{Owner: "owner1", Name: "repo1"},
		{Owner: "owner2", Name: "repo2"},
	}
	assert.Equal(t, expected, parseRepos("FOO", ",", "/"))
}

func Test_ParseRepos_unset(t *testing.T) {
	os.Clearenv()
	var expected []*model.Repo
	assert.Equal(t, expected, parseRepos("FOO", ",", "/"))
}

func Test_parseFloat_success(t *testing.T) {
	_ = os.Setenv("FOO", "1.10")
	assert.Equal(t, 1.10, parseFloat("FOO", 0))
}

func Test_parseFloat_integer(t *testing.T) {
	_ = os.Setenv("FOO", "5")
	assert.Equal(t, float64(5), parseFloat("FOO", 0))
}

func Test_parseFloat_invalid(t *testing.T) {
	_ = os.Setenv("FOO", "nonsense")
	assert.Equal(t, float64(99), parseFloat("FOO", 99))
}

func Test_parseFloat_unset(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, float64(66), parseFloat("FOO", 66))
}

func Test_ParseString_success(t *testing.T) {
	_ = os.Setenv("FOO", "stuff")
	assert.Equal(t, "stuff", parseString("FOO", ""))
}

func Test_ParseString_empty(t *testing.T) {
	_ = os.Setenv("FOO", "")
	assert.Equal(t, "something", parseString("FOO", "something"))
}

func Test_ParseString_unset(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, "something", parseString("FOO", "something"))
}

func Test_ParseBase64_success(t *testing.T) {
	_ = os.Setenv("FOO", "c29tZXRoaW5n")
	assert.Equal(t, "something", string(parseBase64("FOO", []byte(""))))
}

func Test_ParseBase64_invalid(t *testing.T) {
	_ = os.Setenv("FOO", "111111111")
	assert.Equal(t, []byte("something"), parseBase64("FOO", []byte("something")))
}

func Test_ParseBase64_unset(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, []byte("something"), parseBase64("FOO", []byte("something")))
}

func Test_ParseInt64_success(t *testing.T) {
	_ = os.Setenv("FOO", "1651589627")
	assert.Equal(t, int64(1651589627), parseInt64("FOO", 0))
}

func Test_ParseInt64_invalid(t *testing.T) {
	_ = os.Setenv("FOO", "nonsense")
	assert.Equal(t, int64(1651589627), parseInt64("FOO", 1651589627))
}

func Test_ParseInt64_unset(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, int64(1651589627), parseInt64("FOO", 1651589627))
}

func Test_LoadConfig(t *testing.T) {
	_ = os.Setenv("DEBUG", "true")
	_ = os.Setenv("REPOS", "owner1/repo1,owner2/repo2")
	_ = os.Setenv("GITHUB_APP_ID", "1651589627")
	_ = os.Setenv("GITHUB_INSTALLATION_ID", "1651589627")
	_ = os.Setenv("GITHUB_PEM", "c29tZXRoaW5n")
	_ = os.Setenv("DB_HOST", "dbhost")
	_ = os.Setenv("DB_PORT", "dbport")
	_ = os.Setenv("DB_USERNAME", "dbusername")
	_ = os.Setenv("DB_PASSWORD", "dbpass")
	_ = os.Setenv("DB_NAME", "dbname")
	_ = os.Setenv("MINUTES_BETWEEN_CHECKS", "34634.3")

	config := LoadConfig()

	assert.Equal(t, true, config.Debug)

	assert.Equal(t, []*model.Repo{
		{Owner: "owner1", Name: "repo1"},
		{Owner: "owner2", Name: "repo2"},
	}, config.Repos)

	assert.Equal(t, int64(1651589627), config.GitHubAppId)
	assert.Equal(t, int64(1651589627), config.GitHubInstallationId)
	assert.Equal(t, []byte("something"), config.GitHubPem)
	assert.Equal(t, "dbhost", config.DbHost)
	assert.Equal(t, "dbport", config.DbPort)
	assert.Equal(t, "dbusername", config.DbUsername)
	assert.Equal(t, "dbpass", config.DbPassword)
	assert.Equal(t, "dbname", config.DbName)
	assert.Equal(t, 34634.3, config.MinutesBetweenChecks)
}
