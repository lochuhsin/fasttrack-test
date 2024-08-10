.PHONY: run-server
run-server:
	go build -o app server/cmd/main.go && ./app

.PHONY: build-server
build-server:
	go build -o app server/cmd/main.go


.PHONY: build-server-race
build-server-race:
	go build -race -o app server/cmd/main.go 