APP=target/ghloader

.PHONY: configure setup install run

setup: configure install

configure:
	scripts/setup.sh

install:
	go build -o $(APP) ./cmd

run: url file
	$(APP) $(url) $(file)

# Ensure 'url' and 'file' are passed
url:
	$(if $(url),,$(error You must pass both a `url` and `file` parameter to target `make run`))

file:
	$(if $(file),,$(error You must pass both a `url` and `file` parameter to target `make run`))
