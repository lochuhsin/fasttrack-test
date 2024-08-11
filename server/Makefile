.PHONY: build-server
build-server:
	go build -o app cmd/main.go

.PHONY: build-server-race
build-server-race:
	go build -o -race cmd/main.go

.PHONY: run-server
run-server: build-server
	./app
