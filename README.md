# gemini

An application which gathers data for the following DevOps-related metrics:
- Deployment frequency
- Lead time for changes
- Time to restore service
- Change failure rate

This data helps to measure and evaluate your organization's efficiency at delivering excellent software and operational capabilities.
The inspiration for this project comes from the [Accelerate State of DevOps Report](https://services.google.com/fh/files/misc/state-of-devops-2021.pdf). Please see the report
for detailed descriptions for each of the metrics.

[![gobuild](https://github.com/champ-oss/gemini/actions/workflows/gobuild.yml/badge.svg)](https://github.com/champ-oss/gemini/actions/workflows/gobuild.yml)
[![gotest](https://github.com/champ-oss/gemini/actions/workflows/gotest.yml/badge.svg)](https://github.com/champ-oss/gemini/actions/workflows/gotest.yml)
[![release](https://github.com/champ-oss/gemini/actions/workflows/release.yml/badge.svg)](https://github.com/champ-oss/gemini/actions/workflows/release.yml)
[![sonar](https://github.com/champ-oss/gemini/actions/workflows/sonar.yml/badge.svg)](https://github.com/champ-oss/gemini/actions/workflows/sonar.yml)
[![tftest](https://github.com/champ-oss/gemini/actions/workflows/tftest.yml/badge.svg)](https://github.com/champ-oss/gemini/actions/workflows/tftest.yml)

[![SonarCloud](https://sonarcloud.io/images/project_badges/sonarcloud-black.svg)](https://sonarcloud.io/summary/new_code?id=champ-oss_gemini)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=champ-oss_gemini&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=champ-oss_gemini)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=champ-oss_gemini&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=champ-oss_gemini)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=champ-oss_gemini&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=champ-oss_gemini)



## Metric Data
This section explains how metric data is gathered. There are many acceptable ways of measuring the key DevOps metrics outlined above. In this project specifically, we chose to gather metrics
from Git repositories, which often contain valuable markers representing an application's lifecycle.

To this end, a few assumptions were made about the Git repositories used as the source of data being collected:
- GitHub is the host of the repositories (although support could be added in the future for other hosts such as BitBucket, GitLab, etc)
- Terraform is being used as the deployment tool
- The application infrastructure is represented as two or more repositories:
  - A Terraform module which defines the application in an environment-agnostic way.
  - At least one Terraform "environment repository" representing a single deployment of the above application Terraform module. Typically, an application would have many environment repositories (ex: dev, qa, prod).
- A GitHub workflow is used to run `terraform apply` in the "environment repository"



Explanation of metric data:
- **Deployment frequency** - Based on the frequency of `terraform apply` running inside a GitHub workflow in the "environment repository".
The assumption is this workflow would only run on the default branch of a repository and that it represents a single deployment event for an application in that environment.

- **Lead time for changes** - Based on the time difference between a commit being made to the application Terraform module and that exact commit subsequently being referenced in the "environment repository". 
In other words, this tracks the time it takes to deploy a single change made to the application to production.

- **Time to restore service** - Based on the time difference between a failed `terraform apply` workflow and the next successful one.

- **Change failure rate** - Based on the percentage of `terraform apply` workflows which fail, compared to those which succeed.



## Usage
This project is set up to be easily deployed into AWS using Terraform. Please see the `terraform/examples/complete` folder for an example of everything needed to deploy the application.

AWS services being used:
- ECS Fargate
- Aurora RDS MySQL



## Go Code Structure
#### `src/adapter`
Contains any code which we need in order to interface with an external service. At present, there is an adapter implementation for interfacing with the GitHub API. 
Additional adapter implementations could be written for Git providers such as BitBucket or GitLab.

#### `src/cmd`
Contains the entry point for the entire application, processes the configuration parameters, and also sets up and injects the dependencies to use.

#### `src/config`
Parses environment variables and configures the application at runtime.

#### `src/mocks`
Contains code for testing which is entirely generated using mockgen, which is part of [gomock](https://github.com/golang/mock). Use the command `make mocks` to update the mock code.

#### `src/model`
Contains all the data models for the application, defined as [gorm](https://gorm.io/index.html) models.

#### `src/repository`
Contains the implementation for the gorm database. The implementation currently uses MySQL.

#### `src/service`
Contains the business logic for gathering each of the DevOps metrics outlined above. This code essentially connects the `adapter` code to the `repository`. 
The `common` service is the main service for the application.


## Terraform Code Structure
#### `terraform/`
Contains the Terraform module code for deploying the application.

#### `terraform/examples/complete`
Contains a working example of deploying gemini to an AWS environment. This example is also used for integration testing the application.

#### `terraform/test`
Contains Go test code for integration testing the application. See the [Integration Testing](#integration-testing) section for more information.



## Grafana

Grafana is the user interface for viewing metric data. All Grafana configurations and dashboards are controlled using
Terraform. If manual changes are made to a dashboard, it is important to copy the JSON model of the dashboard to the
files in the `terraform/` folder. Without this step, changes to dashboards will be overwritten in the
next deployment.

### Grafana Authentication

GitHub OAuth authentication is supported and should be used instead of basic auth. Once OAuth login is tested and
working, basic auth can be disabled. Below is the process to achieve this:

1. Immediately after the initial deployment, have users log in using GitHub OAuth and verify it is working.
2. Log in as the `admin` user and grant specific users the Admin role.
3. Using the Grafana UI, generate an API key for Terraform to use.
4. In the Terraform environment config:
    1. set `grafana_force_oauth` to `true` which will disable basic auth and the login form
    2. KMS encrypt the API key generated in step 3 and add it to the config for `terraform_api_key`
    3. Set `use_terraform_api_key` to true
5. Apply the Terraform changes, and you will see that the login form is disabled and all users are forced to use OAuth.

NOTE: When setting `grafana_force_oauth` to `true` it breaks the ability for Terraform to configure Grafana unless you
switch it to use the API key to authenticate. Therefore, generating an API key and setting `use_terraform_api_key` to
`true` is mandatory.




## Configuration
All configuration is done using environment variables. Typically, you would only override these variables when running
the application locally for testing. Otherwise, Terraform manages many of these values in a production environment and exposes variables
for changing them. See `terraform/variables.tf` for the list of variables which can be tweaked.

| Name | Description | Default |
|------|--------|---------|
| `REPOS` | Comma separated list of repository names to gather data, with the organization name included. (ex: `champ-oss/gemini`) |  |
| `GITHUB_APP_ID` | GitHub Application ID for authentication |  |
| `GITHUB_INSTALLATION_ID` | GitHub Installation ID for authentication |  |
| `GITHUB_PEM` | Private key for the GitHub app for authentication |  |
| `DB_HOST` | Hostname of the MySQL server to use | `localhost` |
| `DB_PORT` | TCP Port for the MySQL server | `3306` |
| `DB_USERNAME` | Username to connect to the MySQL server | `root` |
| `DB_PASSWORD` | Password to connect to the MySQL server | `secret` |
| `DB_NAME` | Name of the database to use | `gemini` |
| `MINUTES_BETWEEN_CHECKS` | Time interval between checking for new data. Note this may need to be adjusted to account for GitHub rate limiting | `5` |
| `DEBUG` | Enable debug logging | `true` |




## Integration Testing
This application is fully tested in a live AWS environment, using the [`tftest` workflow](https://github.com/champ-oss/gemini/blob/main/.github/workflows/tftest.yml). 
This workflow runs on every commit and helps ensure that changes do not break the functionality of the application.

Below is an overview of the integration testing steps:
1. The `tftest` workflow runs the `go test` command which triggers the `TestGemini` function in the `terraform/test/gemini_test.go` file.
2. Go uses the [terratest](https://terratest.gruntwork.io/) library to stand up a fully functional copy of the application in an AWS account used for testing.
3. Once the application is fully deployed and running, a series of tests are run to ensure the application is functioning correctly. For integration testing, the application is configured to pull data from a sample git repository.
4. The MySQL database is queried to ensure that each of the metric data tables are being populated as expected.
5. The Grafana API is queried to ensure that the configured data source is returning data successfully.


## Bugs / Issues / Features
Please use the built-in [GitHub issues tracker](https://github.com/champ-oss/gemini/issues).


## Contributing
We welcome any and all contributions! Please see [GitHub issues tracker](https://github.com/champ-oss/gemini/issues) for current bugs or enhancements, which may be a good place to start if you would like to contribute. 
If you decide to work on an issue, please assign it to yourself. We are happy to review pull requests!


### Setting up a development environment
- Go 1.17 is currently being used for development. 
- You can use `make download` to install all the dependencies.
- Use `make run` to test the application locally. This will start a local mysql database as well.
- Use `make test` to run all the unit tests
- Use `make coverage` to run all the unit tests and check code coverage. The coverage report will be opened in your browser.
- Use `make mocks` to generate/update mock test files. This will be needed when updating the adapter, service, or repository
- Use `make fmt` to properly format the Go and Terraform code.
- Use `make tidy` to run tidy up Go dependencies