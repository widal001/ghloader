APP=target/ghloader

.PHONY: configure setup install cli web

setup: configure install

configure:
	scripts/setup.sh

install:
	go build -o $(APP) ./cmd

cli: url file
	$(APP) -url $(url) -file $(file)

web-up: install
	@echo "Starting the web application in the background"
	$(APP) -web &
	sleep 2
	open http://localhost:8080
	@echo "View the web app at http://localhost:8080"

web-down:
	@echo "Tearing down the web application"
	pkill -f "$(APP) -web"

# Ensure 'url' and 'file' are passed
url:
	$(if $(url),,$(error You must pass both a `url` and `file` parameter to target `make run`))

file:
	$(if $(file),,$(error You must pass both a `url` and `file` parameter to target `make run`))
