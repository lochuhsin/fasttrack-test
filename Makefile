.PHONY: run
run:
	go build -o app server/cmd/main.go && ./app