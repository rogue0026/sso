build:
	go build -o ./cmd/sso/sso ./cmd/sso/main.go

run:
	TOKEN_SIGNING_KEY="asdfasdfasdf" ./cmd/sso/sso -cfg=./configs/local.yaml

migrator:
	go build -o ./cmd/migrator/migrator ./cmd/migrator/main.go