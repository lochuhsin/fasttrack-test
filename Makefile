.PHONY: build-server
build-server:
	cd server && make build-server

.PHONY: build-server-race
build-server-race:
	cd server && make build-server-race

.PHONY: run-server
run-server:
	cd server && make run-server

.PHONY: install-cli
install-cli:
	cd cli && make install