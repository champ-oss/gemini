
run:
	docker run -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=gemini -d mysql:5.7 || true
	export REPOS=champtitles/tflint-ruleset-champtitles && cd src && go run cmd/main.go

test:
	cd src && go test ./...

coverage:
	cd src && go test -json -coverprofile=cover.out ./... > result.json
	cd src && go tool cover -func cover.out
	cd src && go tool cover -html=cover.out

mocks:
	cd src && go install github.com/golang/mock/mockgen@latest
	cd src && mockgen -source adapter/adapter.go -destination mocks/mock_adapter/mock.go -package mock_adapter
	cd src && mockgen -source repository/repository.go -destination mocks/mock_repository/mock.go -package mock_repository
	cd src && mockgen -source service/commits/commits.go -destination mocks/mock_service_commits/mock.go -package mock_service_commits
	cd src && mockgen -source service/actions/actions.go -destination mocks/mock_service_actions/mock.go -package mock_service_actions
	cd src && mockgen -source service/terraformrefs/terraformrefs.go -destination mocks/mock_service_terraformrefs/mock.go -package mock_service_terraformrefs
	cd src && mockgen -source service/pullrequests/pullrequests.go -destination mocks/mock_service_pullrequests/mock.go -package mock_service_pullrequests

fmt:
	cd src && go fmt ./...
	terraform fmt -recursive -diff

tidy:
	cd src && go mod tidy
	cd terraform/test/src && go mod tidy