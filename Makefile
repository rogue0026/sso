build:
	go build -o ./cmd/sso/sso ./cmd/sso/main.go

run:
	./cmd/sso/sso -cfg=./configs/local.yaml