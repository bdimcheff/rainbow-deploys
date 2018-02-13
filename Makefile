COLOR := $(shell git rev-parse HEAD | cut -c 1-6)

.PHONY: build
build: rainbow-deploys

rainbow-deploys: main.go
	go build

image:
	@echo Building with color $(COLOR)
	docker build . --build-arg COLOR=$(COLOR)